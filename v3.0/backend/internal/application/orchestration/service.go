package orchestration

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	twincore "github.com/Akashpg-M/polaris/backend/internal/core/twin"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	store       *repository.RegistryStore
	redis       *redis.Client
	maxAttempts int
	metrics     *Metrics
	extensions  *extension.Registry
}

type CreateTaskInput struct {
	ProjectID     *string
	TaskType      string
	Priority      string
	Requirements  taskcore.Requirements
	Target        json.RawMessage
	ExpiresAt     time.Time
	CorrelationID string
}

type CreateResult struct {
	Task    taskcore.Task   `json:"task"`
	Command *command.Record `json:"command,omitempty"`
	Timing  *TaskPathTiming `json:"timing,omitempty"`
}

type TaskPathTiming struct {
	CandidateSelectionDurationUS int64 `json:"candidate_selection_duration_us"`
	RoutingDurationUS            int64 `json:"routing_duration_us"`
	PersistenceDurationUS        int64 `json:"persistence_duration_us"`
	TotalDurationUS              int64 `json:"total_duration_us"`
}

func (t *TaskPathTiming) add(other TaskPathTiming) {
	t.CandidateSelectionDurationUS += other.CandidateSelectionDurationUS
	t.RoutingDurationUS += other.RoutingDurationUS
	t.PersistenceDurationUS += other.PersistenceDurationUS
}

func NewService(store *repository.RegistryStore, redisClient *redis.Client, maxAttempts int, metrics *Metrics) *Service {
	registry := extension.NewRegistry()
	registry.RegisterTaskPlanner(extension.DefaultTaskPlanner{})
	return NewServiceWithRegistry(store, redisClient, maxAttempts, metrics, registry)
}

func NewServiceWithRegistry(store *repository.RegistryStore, redisClient *redis.Client, maxAttempts int, metrics *Metrics, registry *extension.Registry) *Service {
	if maxAttempts < 1 {
		maxAttempts = 5
	}
	if metrics == nil {
		metrics = NewMetrics()
	}
	if registry == nil {
		registry = extension.NewRegistry()
		registry.RegisterTaskPlanner(extension.DefaultTaskPlanner{})
	}
	return &Service{store: store, redis: redisClient, maxAttempts: maxAttempts, metrics: metrics, extensions: registry}
}

func (s *Service) CreateTask(ctx context.Context, tenant string, principal auth.OperatorPrincipal, requestID string, in CreateTaskInput) (result CreateResult, err error) {
	started := time.Now()
	timing := TaskPathTiming{}
	defer func() {
		timing.TotalDurationUS = time.Since(started).Microseconds()
		result.Timing = &timing
	}()
	if tenant == "" || in.TaskType == "" || !validPriority(in.Priority) || in.ExpiresAt.Before(time.Now()) {
		return CreateResult{}, ErrInvalidTask
	}
	if len(in.Target) == 0 || !json.Valid(in.Target) {
		return CreateResult{}, ErrInvalidTask
	}
	if capability := command.RequiredCapability(in.TaskType); capability == "" {
		return CreateResult{}, ErrUnsupportedCommand
	} else if !contains(in.Requirements.RequiredCapabilities, capability) {
		in.Requirements.RequiredCapabilities = append(in.Requirements.RequiredCapabilities, capability)
	}
	if in.Requirements.MinimumBattery < 0 || in.Requirements.MinimumBattery > 100 {
		return CreateResult{}, ErrInvalidTask
	}
	if in.TaskType == "NAVIGATE" || in.TaskType == "RELOCATE" {
		if in.Requirements.PlanningMode == "" {
			in.Requirements.PlanningMode = taskcore.PlanningDeviceLocal
		}
		if in.Requirements.PlanningMode != taskcore.PlanningDeviceLocal && in.Requirements.PlanningMode != taskcore.PlanningPolarisRequired {
			return CreateResult{}, ErrInvalidTask
		}
	} else if in.Requirements.PlanningMode != "" {
		return CreateResult{}, ErrInvalidTask
	}
	if in.ProjectID != nil && in.Requirements.ProjectID == "" {
		in.Requirements.ProjectID = *in.ProjectID
	}
	requirements, _ := json.Marshal(in.Requirements)
	correlation := in.CorrelationID
	if correlation == "" {
		correlation = auth.NewID()
	}
	v := taskcore.Task{TaskID: auth.NewID(), TenantID: tenant, ProjectID: in.ProjectID, TaskType: in.TaskType, Status: string(taskcore.Pending), Priority: in.Priority, Requirements: requirements, Target: in.Target, CorrelationID: correlation, CreatedBy: principal.APIKeyID, ExpiresAt: in.ExpiresAt}
	persistenceStarted := time.Now()
	if err := s.store.CreateTask(ctx, v, principal.APIKeyID, requestID); err != nil {
		timing.PersistenceDurationUS += time.Since(persistenceStarted).Microseconds()
		return CreateResult{}, err
	}
	timing.PersistenceDurationUS += time.Since(persistenceStarted).Microseconds()
	s.metrics.TasksCreated.Add(1)
	created, err := s.store.GetTask(ctx, tenant, v.TaskID)
	if err != nil {
		return CreateResult{}, err
	}
	cmd, assignTiming, assignErr := s.assignTimed(ctx, created, principal.APIKeyID, requestID)
	timing.add(assignTiming)
	if assignErr != nil && !errors.Is(assignErr, ErrNoEligibleDevice) && !errors.Is(assignErr, repository.ErrConflict) {
		return CreateResult{Task: created}, assignErr
	}
	created, _ = s.store.GetTask(ctx, tenant, v.TaskID)
	if assignErr == nil {
		s.metrics.CommandsCreated.Add(1)
		return CreateResult{Task: created, Command: &cmd}, nil
	}
	return CreateResult{Task: created}, nil
}

func (s *Service) Assign(ctx context.Context, v taskcore.Task, actor, requestID string) (command.Record, error) {
	record, _, err := s.assignTimed(ctx, v, actor, requestID)
	return record, err
}

func (s *Service) assignTimed(ctx context.Context, v taskcore.Task, actor, requestID string) (command.Record, TaskPathTiming, error) {
	timing := TaskPathTiming{}
	selectionStarted := time.Now()
	var requirements taskcore.Requirements
	if err := json.Unmarshal(v.Requirements, &requirements); err != nil {
		return command.Record{}, timing, ErrInvalidTask
	}
	candidates, err := s.store.EligibleDevices(ctx, v.TenantID, requirements)
	if err != nil {
		return command.Record{}, timing, err
	}
	eligible := make(map[string]repository.DeviceCandidate, len(candidates))
	eligibleIDs := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		eligible[candidate.DeviceID] = candidate
		eligibleIDs = append(eligibleIDs, candidate.DeviceID)
	}
	candidateTiming := &extension.CandidateTiming{}
	request := extension.CandidateRequest{TenantID: v.TenantID, EligibleDeviceIDs: eligibleIDs, RequiredCapabilities: requirements.RequiredCapabilities, DeviceTypeIDs: requirements.AllowedDeviceTypes, Limit: 50, Context: map[string]any{"maximum_distance_meters": requirements.MaximumDistanceM}, Timing: candidateTiming}
	if lat, lon, ok := targetCoordinates(v.Target); ok {
		request.Context["target_latitude"], request.Context["target_longitude"] = lat, lon
	}
	proposals := make([]extension.Candidate, 0, len(candidates))
	domainRanked := false
	if provider, providerErr := s.extensions.CandidateProvider(request); providerErr == nil {
		proposals, err = provider.Candidates(ctx, request)
		if err != nil {
			timing.RoutingDurationUS += candidateTiming.RoutingDuration.Microseconds()
			timing.CandidateSelectionDurationUS += (time.Since(selectionStarted) - candidateTiming.RoutingDuration).Microseconds()
			return command.Record{}, timing, err
		}
		// A domain provider ranks the candidates it understands; it must not
		// accidentally become an exclusive capability filter. Preserve every
		// core-eligible device as an unscored fallback (important for generic
		// non-spatial tasks and for spatial raw-result limits).
		proposals = includeUnrankedEligible(proposals, candidates)
		domainRanked = true
	} else {
		for _, candidate := range candidates {
			proposals = append(proposals, extension.Candidate{DeviceID: candidate.DeviceID})
		}
	}
	// Candidate evaluation is deliberately bounded. Providers may rank a
	// domain-specific subset and Core appends deterministic eligible fallbacks,
	// but a single task must not fan out Redis reads across an entire tenant.
	if request.Limit > 0 && len(proposals) > request.Limit {
		proposals = proposals[:request.Limit]
	}
	type candidateState struct {
		twin, connection *redis.MapStringStringCmd
	}
	states := make(map[string]candidateState, len(eligible))
	pipe := s.redis.Pipeline()
	for _, proposal := range proposals {
		if _, allowed := eligible[proposal.DeviceID]; !allowed {
			continue
		}
		if _, exists := states[proposal.DeviceID]; exists {
			continue
		}
		states[proposal.DeviceID] = candidateState{
			twin:       pipe.HGetAll(ctx, "polaris:twin:"+v.TenantID+":"+proposal.DeviceID),
			connection: pipe.HGetAll(ctx, "polaris:connection:"+v.TenantID+":"+proposal.DeviceID),
		}
	}
	if _, pipeErr := pipe.Exec(ctx); pipeErr != nil && !errors.Is(pipeErr, redis.Nil) {
		timing.RoutingDurationUS += candidateTiming.RoutingDuration.Microseconds()
		timing.CandidateSelectionDurationUS += (time.Since(selectionStarted) - candidateTiming.RoutingDuration).Microseconds()
		return command.Record{}, timing, pipeErr
	}
	ranked := make([]rankedCandidate, 0, len(proposals))
	for order, proposal := range proposals {
		candidate, allowed := eligible[proposal.DeviceID]
		if !allowed {
			continue
		}
		loaded := states[candidate.DeviceID]
		state, stateErr := loaded.twin.Result()
		connection, connectionErr := loaded.connection.Result()
		if stateErr != nil || connectionErr != nil || (state["connectivity_status"] != "ONLINE" && !activeConnectionState(connection)) {
			continue
		}
		var reported struct {
			Lat           float64 `json:"lat"`
			Lon           float64 `json:"lon"`
			EnergyPercent int32   `json:"energy_percent"`
		}
		hasReported := json.Unmarshal([]byte(state["reported_state"]), &reported) == nil
		_, _, hasSpatialTarget := targetCoordinates(v.Target)
		if (!hasReported && (requirements.MinimumBattery > 0 || hasSpatialTarget)) || (hasReported && reported.EnergyPercent < requirements.MinimumBattery) {
			continue
		}
		distance := 0.0
		if targetLat, targetLon, ok := targetCoordinates(v.Target); ok {
			distance = haversineMeters(reported.Lat, reported.Lon, targetLat, targetLon)
			if requirements.MaximumDistanceM > 0 && distance > requirements.MaximumDistanceM {
				continue
			}
		}
		ranked = append(ranked, rankedCandidate{id: candidate.DeviceID, distance: distance, battery: reported.EnergyPercent, domainScore: proposal.DomainScore, proposalOrder: order, twin: twinFromState(v.TenantID, candidate.DeviceID, state, connection)})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if domainRanked {
			if ranked[i].domainScore != nil && ranked[j].domainScore != nil && *ranked[i].domainScore != *ranked[j].domainScore {
				return *ranked[i].domainScore < *ranked[j].domainScore
			}
			if (ranked[i].domainScore != nil) != (ranked[j].domainScore != nil) {
				return ranked[i].domainScore != nil
			}
			if ranked[i].proposalOrder != ranked[j].proposalOrder {
				return ranked[i].proposalOrder < ranked[j].proposalOrder
			}
		}
		if ranked[i].distance != ranked[j].distance {
			return ranked[i].distance < ranked[j].distance
		}
		if ranked[i].battery != ranked[j].battery {
			return ranked[i].battery > ranked[j].battery
		}
		return ranked[i].id < ranked[j].id
	})
	timing.RoutingDurationUS += candidateTiming.RoutingDuration.Microseconds()
	timing.CandidateSelectionDurationUS += (time.Since(selectionStarted) - candidateTiming.RoutingDuration).Microseconds()
	selectionStarted = time.Now()
	fresh, recheckErr := s.store.EligibleDevices(ctx, v.TenantID, requirements)
	if recheckErr != nil {
		timing.CandidateSelectionDurationUS += time.Since(selectionStarted).Microseconds()
		return command.Record{}, timing, recheckErr
	}
	stillEligible := make(map[string]struct{}, len(fresh))
	for _, candidate := range fresh {
		stillEligible[candidate.DeviceID] = struct{}{}
	}
	timing.CandidateSelectionDurationUS += time.Since(selectionStarted).Microseconds()
	for _, candidate := range ranked {
		if _, ok := stillEligible[candidate.id]; !ok {
			continue
		}
		twin := candidate.twin
		if twin.Connectivity != "ONLINE" {
			continue
		}
		plannedTask := v
		plannedTask.AssignedDeviceID = &candidate.id
		planners := s.extensions.TaskPlanners(plannedTask)
		if len(planners) == 0 {
			return command.Record{}, timing, ErrUnsupportedCommand
		}
		var plan extension.ExecutionPlan
		var planErr error
		planningStarted := time.Now()
		for _, planner := range planners {
			plan, planErr = planner.Plan(ctx, extension.PlanningRequest{Task: plannedTask, DeviceTwin: twin})
			if errors.Is(planErr, extension.ErrPlanningUnsupported) {
				continue
			}
			break
		}
		timing.RoutingDurationUS += time.Since(planningStarted).Microseconds()
		if planErr != nil {
			return command.Record{}, timing, planErr
		}
		persistenceStarted := time.Now()
		cmd, err := s.store.AssignTaskWithPlan(ctx, v, candidate.id, actor, requestID, s.maxAttempts, plan)
		timing.PersistenceDurationUS += time.Since(persistenceStarted).Microseconds()
		if err == nil {
			return cmd, timing, nil
		}
		if !errors.Is(err, repository.ErrConflict) {
			return command.Record{}, timing, err
		}
	}
	return command.Record{}, timing, ErrNoEligibleDevice
}

func includeUnrankedEligible(proposals []extension.Candidate, eligible []repository.DeviceCandidate) []extension.Candidate {
	seen := make(map[string]struct{}, len(proposals))
	for _, proposal := range proposals {
		seen[proposal.DeviceID] = struct{}{}
	}
	for _, candidate := range eligible {
		if _, ok := seen[candidate.DeviceID]; ok {
			continue
		}
		proposals = append(proposals, extension.Candidate{DeviceID: candidate.DeviceID})
	}
	return proposals
}

func (s *Service) RetryTask(ctx context.Context, tenant, id string, principal auth.OperatorPrincipal, requestID string, ttl time.Duration) (CreateResult, error) {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	if err := s.store.RetryTask(ctx, tenant, id, principal.APIKeyID, requestID, time.Now().Add(ttl)); err != nil {
		return CreateResult{}, err
	}
	v, err := s.store.GetTask(ctx, tenant, id)
	if err != nil {
		return CreateResult{}, err
	}
	cmd, assignErr := s.Assign(ctx, v, principal.APIKeyID, requestID)
	v, _ = s.store.GetTask(ctx, tenant, id)
	if assignErr == nil {
		s.metrics.CommandsCreated.Add(1)
		return CreateResult{Task: v, Command: &cmd}, nil
	}
	if errors.Is(assignErr, ErrNoEligibleDevice) {
		return CreateResult{Task: v}, nil
	}
	return CreateResult{Task: v}, assignErr
}

type rankedCandidate struct {
	id            string
	distance      float64
	battery       int32
	domainScore   *float64
	proposalOrder int
	twin          twincore.DeviceTwin
}

func (s *Service) loadTwin(ctx context.Context, tenant, device string) (twincore.DeviceTwin, error) {
	state, err := s.redis.HGetAll(ctx, "polaris:twin:"+tenant+":"+device).Result()
	if err != nil {
		return twincore.DeviceTwin{}, err
	}
	connection, connectionErr := s.redis.HGetAll(ctx, "polaris:connection:"+tenant+":"+device).Result()
	if connectionErr != nil && !errors.Is(connectionErr, redis.Nil) {
		return twincore.DeviceTwin{}, connectionErr
	}
	return twinFromState(tenant, device, state, connection), nil
}

func twinFromState(tenant, device string, state, connection map[string]string) twincore.DeviceTwin {
	twin := twincore.DeviceTwin{TenantID: tenant, DeviceID: device, Connectivity: state["connectivity_status"], Components: map[string]twincore.ComponentEnvelope{}}
	if twin.Connectivity != "ONLINE" && activeConnectionState(connection) {
		twin.Connectivity = "ONLINE"
	}
	for field, raw := range state {
		if strings.HasPrefix(field, "component:") {
			var c twincore.ComponentEnvelope
			if json.Unmarshal([]byte(raw), &c) == nil {
				twin.Components[strings.TrimPrefix(field, "component:")] = c
			}
		}
	}
	return twin
}

func activeConnectionState(state map[string]string) bool {
	if state["gateway_id"] == "" {
		return false
	}
	expiresAt, err := strconv.ParseInt(state["lease_expires_at"], 10, 64)
	return err == nil && expiresAt > time.Now().UnixMilli()
}

func targetCoordinates(raw json.RawMessage) (float64, float64, bool) {
	var value map[string]interface{}
	if json.Unmarshal(raw, &value) != nil {
		return 0, 0, false
	}
	lat, lok := number(value["lat"])
	lon, ook := number(value["lon"])
	return lat, lon, lok && ook && lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

func number(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case string:
		f, err := strconv.ParseFloat(n, 64)
		return f, err == nil
	}
	return 0, false
}

func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earth = 6371000.0
	toRad := math.Pi / 180
	dLat, dLon := (lat2-lat1)*toRad, (lon2-lon1)*toRad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*toRad)*math.Cos(lat2*toRad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * earth * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func validPriority(value string) bool {
	return value == "LOW" || value == "NORMAL" || value == "HIGH" || value == "CRITICAL"
}

var (
	ErrInvalidTask        = errors.New("invalid task")
	ErrUnsupportedCommand = errors.New("unsupported command type")
	ErrNoEligibleDevice   = errors.New("no eligible device")
)
