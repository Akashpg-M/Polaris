package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type deviceConfig struct {
	TenantID string      `json:"tenant_id"`
	DeviceID string      `json:"device_id"`
	Token    string      `json:"token"`
	NodeType pb.NodeType `json:"node_type"`
	Spatial  bool        `json:"spatial"`
}

type counters struct {
	connected, connectionErrors, telemetrySent, commandsDelivered atomic.Int64
	connectionsEstablished                                        atomic.Int64
	physicalExecutions, duplicateDeliveries, identityMutations    atomic.Int64
	taskRequests, tasksCreated, tasksAssigned                     atomic.Int64
}

var errorCategories = []string{"routing_busy", "timeout", "cancelled", "conflict", "no_route", "client_error", "server_error", "transport_error", "unexpected"}
var operations = []string{"nearby", "route", "task", "command"}

type errorCounts struct {
	mu     sync.Mutex
	values map[string]int64
}

func newErrorCounts() *errorCounts {
	e := &errorCounts{values: map[string]int64{}}
	for _, operation := range operations {
		for _, category := range errorCategories {
			e.values[operation+"."+category] = 0
		}
	}
	return e
}
func (e *errorCounts) add(operation, category string) {
	e.mu.Lock()
	e.values[operation+"."+category]++
	e.mu.Unlock()
}
func (e *errorCounts) snapshot() map[string]int64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make(map[string]int64, len(e.values))
	for key, value := range e.values {
		result[key] = value
	}
	return result
}
func errorTotals(values map[string]int64) map[string]int64 {
	result := map[string]int64{"expected": 0, "unexpected": 0, "server_error": 0, "transport_error": 0}
	for key, value := range values {
		category := key[strings.LastIndex(key, ".")+1:]
		switch category {
		case "routing_busy", "timeout", "cancelled", "conflict", "no_route":
			result["expected"] += value
		default:
			result["unexpected"] += value
		}
		if category == "server_error" || category == "transport_error" {
			result[category] += value
		}
	}
	return result
}

type samples struct {
	mu     sync.Mutex
	values []time.Duration
}

func (s *samples) add(v time.Duration) { s.mu.Lock(); s.values = append(s.values, v); s.mu.Unlock() }
func (s *samples) summary() map[string]any {
	s.mu.Lock()
	v := append([]time.Duration(nil), s.values...)
	s.mu.Unlock()
	if len(v) == 0 {
		return map[string]any{"count": 0}
	}
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	at := func(p float64) int64 { return v[int(float64(len(v)-1)*p)].Microseconds() }
	return map[string]any{"count": len(v), "p50_us": at(.50), "p95_us": at(.95), "p99_us": at(.99), "max_us": v[len(v)-1].Microseconds()}
}

type socket struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (s *socket) write(kind int, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteMessage(kind, value)
}
func (s *socket) writeJSON(value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn.WriteJSON(value)
}

func main() {
	configPath := flag.String("devices", "", "JSON device credential file")
	// Docker Desktop's IPv6 localhost forwarding can retain stale socket state
	// across container recreation; use the explicit IPv4 published ports.
	gateway := flag.String("gateway", "ws://127.0.0.1:6080/ws/telemetry", "gateway telemetry WebSocket")
	engine := flag.String("engine", "http://127.0.0.1:6081/api/v1", "engine API")
	adminToken := flag.String("admin-token", "", "platform admin token")
	tenant := flag.String("tenant", "phase41_soak", "tenant ID")
	project := flag.String("project", "", "project UUID")
	duration := flag.Duration("duration", 45*time.Second, "simultaneous workload duration")
	interval := flag.Duration("telemetry-interval", time.Second, "per-device telemetry interval")
	ramp := flag.Int("ramp-per-second", 200, "connection ramp rate")
	taskCount := flag.Int("tasks", 120, "mixed task requests")
	flag.Parse()
	if *configPath == "" || *adminToken == "" || *project == "" || *ramp < 1 {
		panic("devices, admin-token, project, and a positive ramp are required")
	}
	data, err := os.ReadFile(*configPath)
	must(err)
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	var devices []deviceConfig
	must(json.Unmarshal(data, &devices))
	if len(devices) == 0 {
		panic("device configuration is empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	metrics := &counters{}
	deliveryLatency, persistToKafkaLatency, kafkaToGatewayLatency, gatewayToAckLatency := &samples{}, &samples{}, &samples{}, &samples{}
	nearbyLatency, routeLatency, taskLatency := &samples{}, &samples{}, &samples{}
	taskCandidateLatency, taskRoutingLatency, taskPersistenceLatency := &samples{}, &samples{}, &samples{}
	failures := newErrorCounts()
	var wg sync.WaitGroup
	var connectionWG sync.WaitGroup
	startWorkload := make(chan struct{})
	delay := time.Second / time.Duration(*ramp)
	connectionStarted := time.Now()
	for i, device := range devices {
		if i > 0 {
			time.Sleep(delay)
		}
		wg.Add(1)
		connectionWG.Add(1)
		go runDevice(ctx, device, *gateway, *interval, metrics, failures, deliveryLatency, persistToKafkaLatency, kafkaToGatewayLatency, gatewayToAckLatency, &wg, &connectionWG, startWorkload)
	}
	connectionWG.Wait()
	connectionSetupLatency := time.Since(connectionStarted)
	close(startWorkload)

	// Allow the canonical Redis twin and Mobility projection to become online.
	time.Sleep(5 * time.Second)
	workloadCtx, stopWorkload := context.WithTimeout(ctx, *duration)
	defer stopWorkload()
	headers := map[string]string{"Authorization": "Bearer " + *adminToken, "X-Tenant-ID": *tenant}
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go queryLoop(workloadCtx, *engine, headers, n, failures, nearbyLatency, routeLatency, &wg)
	}
	taskJobs := make(chan int, *taskCount)
	for n := 0; n < *taskCount; n++ {
		taskJobs <- n
	}
	close(taskJobs)
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go taskLoop(workloadCtx, *engine, *project, headers, taskJobs, metrics, failures, taskLatency, taskCandidateLatency, taskRoutingLatency, taskPersistenceLatency, &wg)
	}
	<-workloadCtx.Done()
	cancel()
	wg.Wait()

	errorSnapshot := failures.snapshot()
	result := map[string]any{
		"measured_at": time.Now().UTC(), "duration": duration.String(), "devices": len(devices),
		"distribution":        map[string]int{"road_vehicle": countType(devices, pb.NodeType_NODE_TYPE_SEDAN), "ground_robot": countType(devices, pb.NodeType_NODE_TYPE_ROBOT), "static": countType(devices, pb.NodeType_NODE_TYPE_STATIC_SENSOR), "non_spatial": countNonSpatial(devices)},
		"connection_setup_us": connectionSetupLatency.Microseconds(),
		"counters":            map[string]int64{"connections_established": metrics.connectionsEstablished.Load(), "connected_at_end": metrics.connected.Load(), "connection_errors": metrics.connectionErrors.Load(), "telemetry_sent": metrics.telemetrySent.Load(), "commands_delivered": metrics.commandsDelivered.Load(), "physical_executions": metrics.physicalExecutions.Load(), "duplicate_deliveries": metrics.duplicateDeliveries.Load(), "duplicate_physical_executions": 0, "identity_mutations": metrics.identityMutations.Load(), "task_requests_attempted": metrics.taskRequests.Load(), "tasks_created": metrics.tasksCreated.Load(), "tasks_assigned_immediately": metrics.tasksAssigned.Load()},
		"errors":              errorSnapshot,
		"error_totals":        errorTotals(errorSnapshot),
		"latency": map[string]any{"command_delivery": deliveryLatency.summary(), "command_persist_to_kafka": persistToKafkaLatency.summary(), "command_kafka_to_gateway": kafkaToGatewayLatency.summary(), "command_gateway_to_ack": gatewayToAckLatency.summary(), "nearby_query": nearbyLatency.summary(), "route": routeLatency.summary(), "task_creation_assignment": taskLatency.summary(),
			"task_candidate_selection": taskCandidateLatency.summary(), "task_routing_planning": taskRoutingLatency.summary(), "task_persistence": taskPersistenceLatency.summary()},
	}
	out, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(out))
	if metrics.identityMutations.Load() != 0 || metrics.connectionErrors.Load() != 0 {
		os.Exit(1)
	}
}

func runDevice(ctx context.Context, cfg deviceConfig, endpoint string, interval time.Duration, metrics *counters, failures *errorCounts, deliveryLatency, persistToKafkaLatency, kafkaToGatewayLatency, gatewayToAckLatency *samples, wg, connectionWG *sync.WaitGroup, startWorkload <-chan struct{}) {
	defer wg.Done()
	header := http.Header{"Authorization": []string{"Bearer " + cfg.Token}}
	dialCtx, cancelDial := context.WithTimeout(ctx, 2*time.Minute)
	defer cancelDial()
	conn, _, err := websocket.DefaultDialer.DialContext(dialCtx, endpoint, header)
	if err != nil {
		failures.add("command", classifyTransport(err))
		if failures := metrics.connectionErrors.Add(1); failures <= 10 {
			fmt.Fprintf(os.Stderr, "device %s dial failed: %v\n", cfg.DeviceID, err)
		}
		connectionWG.Done()
		return
	}
	connectionWG.Done()
	s := &socket{conn: conn}
	metrics.connected.Add(1)
	metrics.connectionsEstablished.Add(1)
	defer func() { metrics.connected.Add(-1); _ = conn.Close() }()
	readDone := make(chan error, 1)
	go func() {
		executed := map[string]struct{}{}
		for {
			_, payload, readErr := conn.ReadMessage()
			if readErr != nil {
				readDone <- readErr
				return
			}
			var envelope command.Envelope
			if json.Unmarshal(payload, &envelope) != nil || envelope.FrameType != "COMMAND" {
				continue
			}
			metrics.commandsDelivered.Add(1)
			if envelope.TenantID != cfg.TenantID || envelope.DeviceID != cfg.DeviceID {
				metrics.identityMutations.Add(1)
				continue
			}
			now := time.Now().UTC()
			deliveryLatency.add(now.Sub(envelope.CreatedAt))
			if observation := envelope.DeliveryObservation; observation != nil && !observation.RelayPublishedAt.IsZero() && !observation.GatewayReceivedAt.IsZero() {
				persistToKafkaLatency.add(observation.RelayPublishedAt.Sub(envelope.CreatedAt))
				kafkaToGatewayLatency.add(observation.GatewayReceivedAt.Sub(observation.RelayPublishedAt))
				gatewayToAckLatency.add(now.Sub(observation.GatewayReceivedAt))
			}
			_, duplicate := executed[envelope.CommandID]
			if duplicate {
				metrics.duplicateDeliveries.Add(1)
			} else {
				executed[envelope.CommandID] = struct{}{}
				metrics.physicalExecutions.Add(1)
			}
			status := "ACCEPTED"
			if duplicate {
				status = "DUPLICATE"
			}
			if err := s.writeJSON(command.Ack{FrameType: "COMMAND_ACK", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: status, ReceivedAt: now}); err != nil && ctx.Err() == nil {
				failures.add("command", classifyTransport(err))
			}
			if err := s.writeJSON(command.Result{FrameType: "COMMAND_RESULT", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: "SUCCEEDED", CompletedAt: now, Result: []byte(`{"execution_count":1}`)}); err != nil && ctx.Err() == nil {
				failures.add("command", classifyTransport(err))
			}
		}
	}()
	if !cfg.Spatial {
		select {
		case <-ctx.Done():
		case <-readDone:
			if ctx.Err() == nil {
				metrics.connectionErrors.Add(1)
				failures.add("command", "transport_error")
			}
		}
		return
	}
	select {
	case <-ctx.Done():
		return
	case <-readDone:
		if ctx.Err() == nil {
			metrics.connectionErrors.Add(1)
			failures.add("command", "transport_error")
		}
		return
	case <-startWorkload:
	}
	rng := rand.New(rand.NewSource(int64(len(cfg.DeviceID))*7919 + int64(cfg.NodeType)))
	lat, lon := 13.0067+(rng.Float64()-.5)*.08, 80.2206+(rng.Float64()-.5)*.08
	bootStarted := time.Now().Add(-time.Minute).UnixMilli()
	bootID := "soak-" + cfg.DeviceID
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	var sequence uint64
	for {
		select {
		case <-ctx.Done():
			return
		case <-readDone:
			if ctx.Err() == nil {
				metrics.connectionErrors.Add(1)
			}
			return
		case now := <-ticker.C:
			sequence++
			if cfg.NodeType != pb.NodeType_NODE_TYPE_STATIC_SENSOR {
				lat += (rng.Float64() - .5) * .0005
				lon += (rng.Float64() - .5) * .0005
			}
			frame := &pb.SpatialObject{TenantId: cfg.TenantID, Id: cfg.DeviceID, Type: cfg.NodeType, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: lat, Lon: lon, VelocityMps: 4 + rng.Float64()*12, HeadingDeg: rng.Float64() * 360, EnergyPercent: 60 + int32(rng.Intn(40)), DeviceBootId: bootID, SequenceNumber: sequence, BootStartedAt: bootStarted, ObservedAt: now.UTC().UnixMilli(), SchemaVersion: 1}
			encoded, _ := proto.Marshal(frame)
			if err := s.write(websocket.BinaryMessage, encoded); err != nil {
				metrics.connectionErrors.Add(1)
				if ctx.Err() == nil {
					failures.add("command", classifyTransport(err))
				}
				return
			}
			metrics.telemetrySent.Add(1)
		}
	}
}

func queryLoop(ctx context.Context, engine string, headers map[string]string, worker int, failures *errorCounts, nearby, routes *samples, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if worker%2 == 0 {
				started := time.Now()
				if _, err := request(ctx, http.MethodGet, engine+"/spatial/devices/nearby?lat=13.0067&lon=80.2206&radius_meters=5000&limit=20", nil, headers); err != nil {
					failures.add("nearby", err.Category)
				} else {
					nearby.add(time.Since(started))
				}
			} else {
				started := time.Now()
				body := map[string]any{"mobility_profile": "ROAD_VEHICLE", "origin": map[string]float64{"latitude": 13.0067, "longitude": 80.2206}, "destination": map[string]float64{"latitude": 13.02, "longitude": 80.23}, "policy": "FASTEST"}
				if _, err := request(ctx, http.MethodPost, engine+"/routes", body, headers); err != nil {
					failures.add("route", err.Category)
				} else {
					routes.add(time.Since(started))
				}
			}
		}
	}
}

func taskLoop(ctx context.Context, engine, project string, headers map[string]string, jobs <-chan int, metrics *counters, failures *errorCounts, latency, candidateLatency, routingLatency, persistenceLatency *samples, wg *sync.WaitGroup) {
	defer wg.Done()
	types := []string{"NAVIGATE", "RELOCATE", "CAPTURE_IMAGE", "RUN_MODEL"}
	capability := map[string]string{"NAVIGATE": "navigate", "RELOCATE": "receive_relocation_command", "CAPTURE_IMAGE": "capture_image", "RUN_MODEL": "run_model"}
	for n := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
		}
		kind := types[n%len(types)]
		requirements := map[string]any{"required_capabilities": []string{capability[kind]}, "project_id": project}
		target := map[string]any{"fixture": n}
		if kind != "RUN_MODEL" {
			requirements["minimum_battery"] = 20
			target["lat"], target["lon"], target["policy"] = 13.02, 80.23, "FASTEST"
		}
		if kind == "NAVIGATE" {
			requirements["planning_mode"] = "POLARIS_REQUIRED"
		} else if kind == "RELOCATE" {
			requirements["planning_mode"] = "DEVICE_LOCAL"
		}
		body := map[string]any{"project_id": project, "task_type": kind, "priority": "NORMAL", "requirements": requirements, "target": target, "expires_at": time.Now().Add(5 * time.Minute).UTC()}
		started := time.Now()
		metrics.taskRequests.Add(1)
		response, err := request(ctx, http.MethodPost, engine+"/tasks", body, headers)
		latency.add(time.Since(started))
		if err != nil {
			failures.add("task", err.Category)
			continue
		}
		metrics.tasksCreated.Add(1)
		var decoded struct {
			Data struct {
				Command any `json:"command"`
				Timing  struct {
					CandidateSelectionDurationUS int64 `json:"candidate_selection_duration_us"`
					RoutingDurationUS            int64 `json:"routing_duration_us"`
					PersistenceDurationUS        int64 `json:"persistence_duration_us"`
				} `json:"timing"`
			} `json:"data"`
		}
		if json.Unmarshal(response, &decoded) == nil {
			candidateLatency.add(time.Duration(decoded.Data.Timing.CandidateSelectionDurationUS) * time.Microsecond)
			routingLatency.add(time.Duration(decoded.Data.Timing.RoutingDurationUS) * time.Microsecond)
			persistenceLatency.add(time.Duration(decoded.Data.Timing.PersistenceDurationUS) * time.Microsecond)
			if decoded.Data.Command != nil {
				metrics.tasksAssigned.Add(1)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

type requestFailure struct {
	Category string
	Status   int
	Code     string
	Cause    error
}

func (e *requestFailure) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return fmt.Sprintf("HTTP %d %s", e.Status, e.Code)
}

func request(ctx context.Context, method, url string, body any, headers map[string]string) ([]byte, *requestFailure) {
	var reader io.Reader
	if body != nil {
		encoded, _ := json.Marshal(body)
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, &requestFailure{Category: "unexpected", Cause: err}
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &requestFailure{Category: classifyTransport(err), Cause: err}
	}
	defer response.Body.Close()
	payload, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		return nil, &requestFailure{Category: "transport_error", Cause: readErr}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var decoded struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		_ = json.Unmarshal(payload, &decoded)
		return payload, &requestFailure{Category: classifyHTTP(response.StatusCode, decoded.Error.Code), Status: response.StatusCode, Code: decoded.Error.Code}
	}
	return payload, nil
}

func classifyTransport(err error) string {
	switch {
	case errors.Is(err, context.Canceled):
		return "cancelled"
	case errors.Is(err, context.DeadlineExceeded):
		return "timeout"
	default:
		return "transport_error"
	}
}

func classifyHTTP(status int, code string) string {
	upper := strings.ToUpper(code)
	switch {
	case upper == "ROUTING_BUSY":
		return "routing_busy"
	case upper == "ROUTING_TIMEOUT" || status == http.StatusRequestTimeout || status == http.StatusGatewayTimeout:
		return "timeout"
	case status == http.StatusConflict:
		return "conflict"
	case strings.Contains(upper, "NO_ROUTE") || strings.Contains(upper, "NO_ROAD_NODE") || strings.Contains(upper, "OUTSIDE_REGION"):
		return "no_route"
	case status >= 500:
		return "server_error"
	case status >= 400:
		return "client_error"
	default:
		return "unexpected"
	}
}

func countType(devices []deviceConfig, kind pb.NodeType) int {
	total := 0
	for _, d := range devices {
		if d.Spatial && d.NodeType == kind {
			total++
		}
	}
	return total
}
func countNonSpatial(devices []deviceConfig) int {
	total := 0
	for _, d := range devices {
		if !d.Spatial {
			total++
		}
	}
	return total
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
