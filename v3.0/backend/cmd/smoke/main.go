package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {
	gateway := env("GATEWAY_URL", "ws://127.0.0.1:6080")
	engine := env("ENGINE_URL", "http://127.0.0.1:6081")
	id := env("SMOKE_DEVICE_ID", fmt.Sprintf("SMOKE-%d", time.Now().UnixNano()))
	deviceToken := os.Getenv("DEVICE_TOKEN")
	operatorToken := os.Getenv("OPERATOR_TOKEN")
	tenant := env("SMOKE_TENANT_ID", "alpha_logistics")
	if deviceToken == "" || operatorToken == "" {
		panic("DEVICE_TOKEN and OPERATOR_TOKEN are required")
	}
	dashboardHeaders := http.Header{}
	dashboardHeaders.Set("Authorization", "Bearer "+operatorToken)
	dashboardHeaders.Set("X-Tenant-ID", tenant)
	dashboard, _, err := websocket.DefaultDialer.Dial(gateway+"/ws/dashboard", dashboardHeaders)
	must(err)
	defer dashboard.Close()
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+deviceToken)
	telemetry, _, err := websocket.DefaultDialer.Dial(gateway+"/ws/telemetry", headers)
	must(err)
	defer telemetry.Close()
	now := time.Now().UTC()
	lat := envFloat("SMOKE_LAT", 13.0067)
	lon := envFloat("SMOKE_LON", 80.2206)
	nodeType := int64(pb.NodeType_NODE_TYPE_DRONE)
	if raw := os.Getenv("SMOKE_NODE_TYPE"); raw != "" {
		if parsed, parseErr := strconv.ParseInt(raw, 10, 32); parseErr == nil {
			nodeType = parsed
		}
	}
	payload := &pb.SpatialObject{Id: id, TenantId: tenant, Type: pb.NodeType(nodeType), Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: lat, Lon: lon, VelocityMps: 12.5, EnergyPercent: 91,
		DeviceBootId: env("SMOKE_BOOT_ID", "boot-"+id), SequenceNumber: envUint64("SMOKE_SEQUENCE", 1), BootStartedAt: envInt64("SMOKE_BOOT_STARTED_AT", now.Add(-time.Minute).UnixMilli()), ObservedAt: now.UnixMilli(), SchemaVersion: 1}
	data, err := proto.Marshal(payload)
	must(err)
	started := time.Now()
	must(telemetry.WriteMessage(websocket.BinaryMessage, data))
	if !envBool("SMOKE_WAIT_FOR_PROJECTION", true) {
		time.Sleep(500 * time.Millisecond)
		result, _ := json.Marshal(map[string]interface{}{"id": id, "lat": payload.Lat, "lon": payload.Lon, "sequence_number": payload.SequenceNumber, "projection_waited": false})
		fmt.Println(string(result))
		return
	}

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
	if !envBool("SMOKE_WAIT_FOR_MATCH", true) {
		result, _ := json.Marshal(map[string]interface{}{"id": id, "lat": payload.Lat, "lon": payload.Lon, "sequence_number": payload.SequenceNumber, "dashboard_observed": true})
		fmt.Println(string(result))
		return
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		url := fmt.Sprintf("%s/api/v1/nodes/match?tenant_id=%s&lat=%f&lon=%f&radius_km=1&class=%d", engine, tenant, lat, lon, nodeType)
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

func envFloat(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envUint64(key string, fallback uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseUint(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envInt64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}
