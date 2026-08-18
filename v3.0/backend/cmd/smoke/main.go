package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	gateway := env("GATEWAY_URL", "ws://localhost:6080")
	engine := env("ENGINE_URL", "http://localhost:6081")
	id := env("SMOKE_DEVICE_ID", fmt.Sprintf("SMOKE-%d", time.Now().UnixNano()))
	deviceToken := os.Getenv("DEVICE_TOKEN")
	operatorToken := os.Getenv("OPERATOR_TOKEN")
	if deviceToken == "" || operatorToken == "" {
		panic("DEVICE_TOKEN and OPERATOR_TOKEN are required")
	}
	dashboardHeaders := http.Header{}
	dashboardHeaders.Set("Authorization", "Bearer "+operatorToken)
	dashboardHeaders.Set("X-Tenant-ID", "alpha_logistics")
	dashboard, _, err := websocket.DefaultDialer.Dial(gateway+"/ws/dashboard", dashboardHeaders)
	must(err)
	defer dashboard.Close()
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+deviceToken)
	telemetry, _, err := websocket.DefaultDialer.Dial(gateway+"/ws/telemetry", headers)
	must(err)
	defer telemetry.Close()
	now := time.Now().UTC()
	payload := &pb.SpatialObject{Id: id, TenantId: "alpha_logistics", Type: pb.NodeType_NODE_TYPE_DRONE, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13.0067, Lon: 80.2206, VelocityMps: 12.5, EnergyPercent: 91,
		DeviceBootId: "boot-" + id, SequenceNumber: 1, BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1}
	data, err := proto.Marshal(payload)
	must(err)
	started := time.Now()
	must(telemetry.WriteMessage(websocket.BinaryMessage, data))

	dashboard.SetReadDeadline(time.Now().Add(15 * time.Second))
	for {
		_, frame, readErr := dashboard.ReadMessage()
		must(readErr)
		var event struct {
			ID string `json:"id"`
		}
		if json.Unmarshal(frame, &event) == nil && event.ID == id {
			break
		}
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/api/v1/nodes/match?tenant_id=alpha_logistics&lat=13.0067&lon=80.2206&radius_km=1&class=5", engine)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+operatorToken)
		req.Header.Set("X-Tenant-ID", "alpha_logistics")
		response, getErr := http.DefaultClient.Do(req)
		if getErr == nil {
			var body struct {
				Count int `json:"count"`
			}
			_ = json.NewDecoder(response.Body).Decode(&body)
			response.Body.Close()
			if body.Count > 0 {
				result, _ := json.Marshal(map[string]interface{}{"id": id, "lat": payload.Lat, "lon": payload.Lon, "end_to_end_latency_ms": time.Since(started).Milliseconds()})
				fmt.Println(string(result))
				return
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	panic("engine did not expose smoke-test node")
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
