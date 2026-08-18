package orchestration

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	store       *repository.RegistryStore
	redis       *redis.Client
	maxAttempts int
	metrics     *Metrics
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
}

func NewService(store *repository.RegistryStore, redisClient *redis.Client, maxAttempts int, metrics *Metrics) *Service {
	if maxAttempts < 1 {
		maxAttempts = 5
	}
	if metrics == nil {
		metrics = NewMetrics()
	}
	return &Service{store: store, redis: redisClient, maxAttempts: maxAttempts, metrics: metrics}
}

func (s *Service) CreateTask(ctx context.Context, tenant string, principal auth.OperatorPrincipal, requestID string, in CreateTaskInput) (CreateResult, error) {
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
	if in.ProjectID != nil && in.Requirements.ProjectID == "" {
		in.Requirements.ProjectID = *in.ProjectID
	}
	requirements, _ := json.Marshal(in.Requirements)
	correlation := in.CorrelationID
	if correlation == "" {
		correlation = auth.NewID()
	}
	v := taskcore.Task{TaskID: auth.NewID(), TenantID: tenant, ProjectID: in.ProjectID, TaskType: in.TaskType, Status: string(taskcore.Pending), Priority: in.Priority, Requirements: requirements, Target: in.Target, CorrelationID: correlation, CreatedBy: principal.APIKeyID, ExpiresAt: in.ExpiresAt}
	if err := s.store.CreateTask(ctx, v, principal.APIKeyID, requestID); err != nil {
		return CreateResult{}, err
	}
	s.metrics.TasksCreated.Add(1)
	created, err := s.store.GetTask(ctx, tenant, v.TaskID)
	if err != nil {
		return CreateResult{}, err
	}
	cmd, assignErr := s.Assign(ctx, created, principal.APIKeyID, requestID)
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
	var requirements taskcore.Requirements
	if err := json.Unmarshal(v.Requirements, &requirements); err != nil {
		return command.Record{}, ErrInvalidTask
	}
	candidates, err := s.store.EligibleDevices(ctx, v.TenantID, requirements)
	if err != nil {
		return command.Record{}, err
	}
	ranked := make([]rankedCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		state, err := s.redis.HGetAll(ctx, "polaris:twin:"+v.TenantID+":"+candidate.DeviceID).Result()
		if err != nil || state["connectivity_status"] != "ONLINE" {
			continue
		}
		var reported struct {
			Lat           float64 `json:"lat"`
			Lon           float64 `json:"lon"`
			EnergyPercent int32   `json:"energy_percent"`
		}
		if json.Unmarshal([]byte(state["reported_state"]), &reported) != nil || reported.EnergyPercent < requirements.MinimumBattery {
			continue
		}
		distance := 0.0
		if targetLat, targetLon, ok := targetCoordinates(v.Target); ok {
			distance = haversineMeters(reported.Lat, reported.Lon, targetLat, targetLon)
			if requirements.MaximumDistanceM > 0 && distance > requirements.MaximumDistanceM {
				continue
			}
		}
		ranked = append(ranked, rankedCandidate{id: candidate.DeviceID, distance: distance, battery: reported.EnergyPercent})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].distance != ranked[j].distance {
			return ranked[i].distance < ranked[j].distance
		}
		if ranked[i].battery != ranked[j].battery {
			return ranked[i].battery > ranked[j].battery
		}
		return ranked[i].id < ranked[j].id
	})
	for _, candidate := range ranked {
		cmd, err := s.store.AssignTask(ctx, v, candidate.id, actor, requestID, s.maxAttempts)
		if err == nil {
			return cmd, nil
		}
		if !errors.Is(err, repository.ErrConflict) {
			return command.Record{}, err
		}
	}
	return command.Record{}, ErrNoEligibleDevice
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
	id       string
	distance float64
	battery  int32
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
