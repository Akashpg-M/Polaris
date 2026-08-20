package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	gateway := flag.String("gateway", "ws://127.0.0.1:6080/ws/telemetry", "telemetry WebSocket")
	engine := flag.String("engine", "http://127.0.0.1:6081", "engine origin")
	admin := flag.String("admin-token", "", "platform admin token")
	tenant := flag.String("tenant", "", "tenant")
	project := flag.String("project", "", "project")
	device := flag.String("device", "", "device")
	token := flag.String("device-token", "", "device credential")
	requests := flag.Int("requests", 80, "concurrent route requests")
	flag.Parse()
	if *admin == "" || *tenant == "" || *project == "" || *device == "" || *token == "" || *requests < 2 {
		panic("admin-token, tenant, project, device, device-token, and requests are required")
	}
	headers := map[string]string{"Authorization": "Bearer " + *admin, "X-Tenant-ID": *tenant}
	baseline := route(context.Background(), *engine, headers)
	if baseline != http.StatusOK {
		panic(fmt.Sprintf("baseline route failed with HTTP %d", baseline))
	}

	wsHeaders := http.Header{"Authorization": []string{"Bearer " + *token}}
	conn, _, err := websocket.DefaultDialer.Dial(*gateway, wsHeaders)
	must(err)
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var writeMu sync.Mutex
	var telemetry atomic.Int64
	commandReceived := make(chan struct{}, 1)
	go func() {
		for {
			_, payload, readErr := conn.ReadMessage()
			if readErr != nil {
				return
			}
			var envelope command.Envelope
			if json.Unmarshal(payload, &envelope) != nil || envelope.FrameType != "COMMAND" {
				continue
			}
			now := time.Now().UTC()
			writeMu.Lock()
			_ = conn.WriteJSON(command.Ack{FrameType: "COMMAND_ACK", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: "ACCEPTED", ReceivedAt: now})
			_ = conn.WriteJSON(command.Result{FrameType: "COMMAND_RESULT", CommandID: envelope.CommandID, SequenceNumber: envelope.SequenceNumber, Status: "SUCCEEDED", CompletedAt: now, Result: []byte(`{"overload_probe":true}`)})
			writeMu.Unlock()
			select {
			case commandReceived <- struct{}{}:
			default:
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		bootStarted := time.Now().Add(-time.Minute).UnixMilli()
		var sequence uint64
		for {
			select {
			case <-ctx.Done():
				return
			case observed := <-ticker.C:
				sequence++
				frame, _ := proto.Marshal(&pb.SpatialObject{TenantId: *tenant, Id: *device, Type: pb.NodeType_NODE_TYPE_SEDAN, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13.0067, Lon: 80.2206, VelocityMps: 8, EnergyPercent: 90, DeviceBootId: "overload-" + *device, SequenceNumber: sequence, BootStartedAt: bootStarted, ObservedAt: observed.UnixMilli(), SchemaVersion: 1})
				writeMu.Lock()
				err := conn.WriteMessage(websocket.BinaryMessage, frame)
				writeMu.Unlock()
				if err != nil {
					return
				}
				telemetry.Add(1)
			}
		}
	}()
	time.Sleep(time.Second)

	start := make(chan struct{})
	var wg sync.WaitGroup
	var busy, timedOut, succeeded, unexpected atomic.Int64
	for i := 0; i < *requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			status := route(context.Background(), *engine, headers)
			switch status {
			case http.StatusOK:
				succeeded.Add(1)
			case http.StatusTooManyRequests:
				busy.Add(1)
			case http.StatusGatewayTimeout:
				timedOut.Add(1)
			default:
				unexpected.Add(1)
			}
		}()
	}
	close(start)
	time.Sleep(50 * time.Millisecond)
	taskBody := map[string]any{"project_id": *project, "task_type": "RUN_MODEL", "priority": "HIGH", "requirements": map[string]any{"required_capabilities": []string{"run_model"}, "project_id": *project}, "target": map[string]any{"model": "overload-proof"}, "expires_at": time.Now().Add(2 * time.Minute).UTC()}
	taskStatus, _ := api(context.Background(), http.MethodPost, *engine+"/api/v1/tasks", taskBody, headers)
	wg.Wait()
	commandOK := false
	select {
	case <-commandReceived:
		commandOK = true
	case <-time.After(15 * time.Second):
	}
	postStatus := 0
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		postStatus = route(context.Background(), *engine, headers)
		if postStatus == http.StatusOK {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	readyStatus, readyBody := api(context.Background(), http.MethodGet, *engine+"/readyz", nil, nil)
	result := map[string]any{"baseline_status": baseline, "flood_requests": *requests, "routing_busy": busy.Load(), "routing_timeout": timedOut.Load(), "routing_success": succeeded.Load(), "unexpected": unexpected.Load(), "telemetry_sent_during_overload": telemetry.Load(), "generic_task_status": taskStatus, "generic_command_completed": commandOK, "post_overload_route_status": postStatus, "ready_status": readyStatus, "ready": json.RawMessage(readyBody)}
	encoded, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(encoded))
	if busy.Load() == 0 || unexpected.Load() != 0 || telemetry.Load() == 0 || taskStatus != http.StatusCreated || !commandOK || postStatus != http.StatusOK || readyStatus != http.StatusOK {
		os.Exit(1)
	}
}

func route(ctx context.Context, engine string, headers map[string]string) int {
	body := map[string]any{"mobility_profile": "ROAD_VEHICLE", "origin": map[string]float64{"latitude": 13.0067, "longitude": 80.2206}, "destination": map[string]float64{"latitude": 13.18, "longitude": 80.30}, "policy": "FASTEST"}
	status, _ := api(ctx, http.MethodPost, engine+"/api/v1/routes", body, headers)
	return status
}

func api(ctx context.Context, method, url string, body any, headers map[string]string) (int, []byte) {
	var reader io.Reader
	if body != nil {
		encoded, _ := json.Marshal(body)
		reader = bytes.NewReader(encoded)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	must(err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, []byte(err.Error())
	}
	defer response.Body.Close()
	payload, _ := io.ReadAll(response.Body)
	return response.StatusCode, payload
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
