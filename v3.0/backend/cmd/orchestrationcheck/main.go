package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var (
	engineURL  = env("ENGINE_URL", "http://localhost:6081")
	gatewayURL = env("GATEWAY_URL", "ws://localhost:6080")
	tenantID   = env("TENANT_ID", "alpha_logistics")
)

func main() {
	mode := env("ORCHESTRATION_CHECK_MODE", "complete")
	switch mode {
	case "complete":
		complete(false)
	case "duplicate":
		complete(true)
	case "offline":
		offline()
	case "fencing":
		fencing()
	case "wrong-ack":
		wrongAck()
	case "capability-mismatch":
		capabilityMismatch()
	case "receive-no-ack":
		receiveNoAck()
	case "resume":
		resume()
	default:
		panic("unknown orchestration check mode: " + mode)
	}
	fmt.Println("PASS: " + mode)
}

func complete(dropFirstAck bool) {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 1)
	waitOnline(device)
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	first := readCommand(conn, 15*time.Second)
	if first.TaskID != taskID || first.DeviceID != device {
		panic("command was not bound to task/device")
	}
	if dropFirstAck {
		second := readCommand(conn, 15*time.Second)
		if second.CommandID != first.CommandID || second.SequenceNumber != first.SequenceNumber {
			panic("retry did not preserve command identity and sequence")
		}
		ack(conn, second, "DUPLICATE")
		result(conn, second, 1)
		waitTask(taskID, "COMPLETED")
		cmd := getCommand(first.CommandID)
		if cmd.AttemptCount < 2 {
			panic("lost ACK did not produce a bounded retry")
		}
		return
	}
	ack(conn, first, "ACCEPTED")
	result(conn, first, 1)
	waitTask(taskID, "COMPLETED")
}

func offline() {
	device := mustEnv("SMOKE_DEVICE_ID")
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	if status := getTask(taskID).Status; status != "PENDING" {
		panic("offline task did not remain pending: " + status)
	}
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 1)
	commandFrame := readCommand(conn, 15*time.Second)
	if commandFrame.TaskID != taskID {
		panic("reconnect reconciliation delivered wrong task")
	}
	ack(conn, commandFrame, "ACCEPTED")
	result(conn, commandFrame, 1)
	waitTask(taskID, "COMPLETED")
}

func fencing() {
	device := mustEnv("SMOKE_DEVICE_ID")
	first := connect(mustEnv("DEVICE_TOKEN"))
	sendTelemetry(first, device, 1)
	second := connect(mustEnv("DEVICE_TOKEN"))
	defer second.Close()
	sendTelemetry(second, device, 2)
	first.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := first.ReadMessage(); err == nil {
		panic("stale ownership socket remained usable")
	}
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	frame := readCommand(second, 15*time.Second)
	ack(second, frame, "ACCEPTED")
	result(second, frame, 1)
	waitTask(taskID, "COMPLETED")
}

func wrongAck() {
	deviceA := mustEnv("SMOKE_DEVICE_ID")
	a := connect(mustEnv("DEVICE_TOKEN"))
	defer a.Close()
	sendTelemetry(a, deviceA, 1)
	waitOnline(deviceA)
	taskID := createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	frame := readCommand(a, 15*time.Second)
	b := connect(mustEnv("DEVICE_TOKEN_B"))
	defer b.Close()
	sendTelemetry(b, mustEnv("DEVICE_ID_B"), 1)
	ack(b, frame, "ACCEPTED")
	b.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := b.ReadMessage(); err == nil {
		panic("wrong-device ACK did not close the connection")
	}
	if status := getCommand(frame.CommandID).Status; status != "DELIVERED" {
		panic("wrong-device ACK mutated command: " + status)
	}
	ack(a, frame, "ACCEPTED")
	result(a, frame, 1)
	waitTask(taskID, "COMPLETED")
}

func capabilityMismatch() {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 1)
	waitOnline(device)
	taskID := createTask("CAPTURE_IMAGE", []string{"capture_image"}, 0, time.Now().Add(2*time.Second))
	if status := getTask(taskID).Status; status != "PENDING" {
		panic("capability mismatch was assigned")
	}
	time.Sleep(3 * time.Second)
	waitTask(taskID, "EXPIRED")
}

func receiveNoAck() {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	sendTelemetry(conn, device, 1)
	waitOnline(device)
	_ = createTask("RELOCATE", []string{"receive_relocation_command"}, 30, time.Now().Add(time.Minute))
	frame := readCommand(conn, 15*time.Second)
	fmt.Println("COMMAND_ID=" + frame.CommandID)
	conn.Close()
}

func resume() {
	device := mustEnv("SMOKE_DEVICE_ID")
	conn := connect(mustEnv("DEVICE_TOKEN"))
	defer conn.Close()
	sendTelemetry(conn, device, 2)
	frame := readCommand(conn, 15*time.Second)
	if expected := mustEnv("EXPECTED_COMMAND_ID"); frame.CommandID != expected {
		panic("gateway recovery changed command identity")
	}
	ack(conn, frame, "DUPLICATE")
	result(conn, frame, 1)
	waitTask(frame.TaskID, "COMPLETED")
}

func connect(token string) *websocket.Conn {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	conn, response, err := websocket.DefaultDialer.Dial(gatewayURL+"/ws/telemetry", headers)
	if err != nil {
		if response != nil {
			panic(fmt.Sprintf("connect rejected: %d", response.StatusCode))
		}
		panic(err)
	}
	return conn
}

func sendTelemetry(conn *websocket.Conn, device string, sequence uint64) {
	now := time.Now().UTC()
	frame := &pb.SpatialObject{Id: device, TenantId: tenantID, Type: pb.NodeType_NODE_TYPE_DRONE, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13.0067, Lon: 80.2206, VelocityMps: 10, EnergyPercent: 90, DeviceBootId: "phase3-" + device, SequenceNumber: sequence, BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1}
	payload, err := proto.Marshal(frame)
	must(err)
	must(conn.WriteMessage(websocket.BinaryMessage, payload))
}

func createTask(commandType string, capabilities []string, battery int, expires time.Time) string {
	body := map[string]interface{}{"task_type": commandType, "priority": "HIGH", "requirements": map[string]interface{}{"required_capabilities": capabilities, "minimum_battery": battery, "max_distance_meters": 10000}, "target": map[string]interface{}{"lat": 13.0068, "lon": 80.2207}, "expires_at": expires.UTC()}
	if project := os.Getenv("TASK_PROJECT_ID"); project != "" {
		body["project_id"] = project
		body["requirements"].(map[string]interface{})["project_id"] = project
	}
	var response struct {
		Data struct {
			Task struct {
				TaskID string `json:"task_id"`
			} `json:"task"`
		} `json:"data"`
	}
	api(http.MethodPost, "/api/v1/tasks", body, &response)
	if response.Data.Task.TaskID == "" {
		panic("task API returned no task ID")
	}
	return response.Data.Task.TaskID
}

type taskView struct {
	Status string `json:"status"`
}

func getTask(id string) taskView {
	var response struct {
		Data struct {
			Task taskView `json:"task"`
		} `json:"data"`
	}
	api(http.MethodGet, "/api/v1/tasks/"+id, nil, &response)
	return response.Data.Task
}
func getCommand(id string) command.Record {
	var response struct {
		Data command.Record `json:"data"`
	}
	api(http.MethodGet, "/api/v1/commands/"+id, nil, &response)
	return response.Data
}
func waitTask(id, expected string) {
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if getTask(id).Status == expected {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	panic("task did not reach " + expected + "; current=" + getTask(id).Status)
}
func waitOnline(device string) {
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		var response struct {
			Data struct {
				Connectivity struct {
					Status string `json:"status"`
				} `json:"connectivity"`
			} `json:"data"`
		}
		api(http.MethodGet, "/api/v1/devices/"+device+"/twin", nil, &response)
		if response.Data.Connectivity.Status == "ONLINE" {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	panic("device did not become online")
}

func readCommand(conn *websocket.Conn, timeout time.Duration) command.Envelope {
	conn.SetReadDeadline(time.Now().Add(timeout))
	for {
		messageType, data, err := conn.ReadMessage()
		must(err)
		if messageType != websocket.TextMessage {
			continue
		}
		var frame command.Envelope
		if json.Unmarshal(data, &frame) == nil && frame.FrameType == "COMMAND" {
			return frame
		}
	}
}
func ack(conn *websocket.Conn, frame command.Envelope, status string) {
	must(conn.WriteJSON(command.Ack{FrameType: "COMMAND_ACK", CommandID: frame.CommandID, SequenceNumber: frame.SequenceNumber, Status: status, ReceivedAt: time.Now().UTC()}))
}
func result(conn *websocket.Conn, frame command.Envelope, executionCount int) {
	payload, _ := json.Marshal(map[string]int{"execution_count": executionCount})
	must(conn.WriteJSON(command.Result{FrameType: "COMMAND_RESULT", CommandID: frame.CommandID, SequenceNumber: frame.SequenceNumber, Status: "SUCCEEDED", CompletedAt: time.Now().UTC(), Result: payload}))
}

func api(method, path string, body interface{}, target interface{}) {
	var content []byte
	if body != nil {
		content, _ = json.Marshal(body)
	}
	req, _ := http.NewRequest(method, engineURL+path, bytes.NewReader(content))
	req.Header.Set("Authorization", "Bearer "+mustEnv("OPERATOR_TOKEN"))
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	must(err)
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var failure interface{}
		_ = json.NewDecoder(response.Body).Decode(&failure)
		panic(fmt.Sprintf("API %s %s failed: %d %#v", method, path, response.StatusCode, failure))
	}
	must(json.NewDecoder(response.Body).Decode(target))
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(key + " is required")
	}
	return value
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
