package main

import (
	"bytes"
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
	mode := env("IDENTITY_CHECK_MODE", "basic")
	token := mustEnv("DEVICE_TOKEN")
	device := mustEnv("SMOKE_DEVICE_ID")
	switch mode {
	case "basic":
		expectRejected("pol_dev_invalid.invalid")
		connect(token).Close()
		spoof(token, device)
	case "rejected":
		expectRejected(token)
	case "send":
		c := connect(token)
		send(c, device, "alpha_logistics", 1)
		c.Close()
	case "ticket":
		ticketCheck(device)
	case "revoke-session":
		c := connect(token)
		revoke()
		send(c, device, "alpha_logistics", 2)
		c.SetReadDeadline(time.Now().Add(3 * time.Second))
		if _, _, err := c.ReadMessage(); err == nil {
			panic("revoked active session remained open")
		}
	default:
		panic("unknown mode")
	}
	fmt.Println("PASS: " + mode)
}
func connect(token string) *websocket.Conn {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+token)
	c, r, err := websocket.DefaultDialer.Dial(env("GATEWAY_URL", "ws://localhost:6080")+"/ws/telemetry", h)
	if err != nil {
		if r != nil {
			panic(fmt.Sprintf("WebSocket rejected: %d", r.StatusCode))
		}
		panic(err)
	}
	return c
}
func expectRejected(token string) {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+token)
	c, r, err := websocket.DefaultDialer.Dial(env("GATEWAY_URL", "ws://localhost:6080")+"/ws/telemetry", h)
	if c != nil {
		c.Close()
		panic("invalid/revoked credential connected")
	}
	if err == nil || r == nil || r.StatusCode != http.StatusUnauthorized {
		panic("expected HTTP 401 before WebSocket upgrade")
	}
}
func spoof(token, device string) {
	c := connect(token)
	defer c.Close()
	send(c, "spoofed-device", "other_tenant", 1)
	c.SetReadDeadline(time.Now().Add(3 * time.Second))
	if _, _, err := c.ReadMessage(); err == nil {
		panic("spoofed identity was not disconnected")
	}
}
func send(c *websocket.Conn, device, tenant string, seq uint64) {
	now := time.Now().UTC()
	if raw := os.Getenv("TELEMETRY_SEQUENCE"); raw != "" {
		_, _ = fmt.Sscan(raw, &seq)
	}
	boot := env("DEVICE_BOOT_ID", "identity-check-boot")
	p := &pb.SpatialObject{Id: device, TenantId: tenant, Type: pb.NodeType_NODE_TYPE_DRONE, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13.0067, Lon: 80.2206, VelocityMps: 10, EnergyPercent: 90, DeviceBootId: boot, SequenceNumber: seq, BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1}
	b, err := proto.Marshal(p)
	if err != nil {
		panic(err)
	}
	if err = c.WriteMessage(websocket.BinaryMessage, b); err != nil {
		panic(err)
	}
}
func revoke() {
	body := bytes.NewBufferString(`{}`)
	url := env("ENGINE_URL", "http://localhost:6081") + "/api/v1/devices/" + mustEnv("SMOKE_DEVICE_ID") + "/credentials/" + mustEnv("DEVICE_CREDENTIAL_ID") + "/revoke"
	req, _ := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Authorization", "Bearer "+mustEnv("OPERATOR_TOKEN"))
	req.Header.Set("X-Tenant-ID", "alpha_logistics")
	req.Header.Set("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		var v interface{}
		_ = json.NewDecoder(r.Body).Decode(&v)
		panic(fmt.Sprintf("revoke failed %d %#v", r.StatusCode, v))
	}
}
func ticketCheck(device string) {
	url := env("ENGINE_URL", "http://localhost:6081") + "/api/v1/devices/" + device + "/connection-ticket"
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer "+mustEnv("OPERATOR_TOKEN"))
	req.Header.Set("X-Tenant-ID", "alpha_logistics")
	req.Header.Set("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	var body struct {
		Data struct {
			Ticket string `json:"ticket"`
		} `json:"data"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	r.Body.Close()
	if r.StatusCode != 201 || body.Data.Ticket == "" {
		panic("ticket issue failed")
	}
	target := env("GATEWAY_URL", "ws://localhost:6080") + "/ws/telemetry?ticket=" + body.Data.Ticket
	c, _, err := websocket.DefaultDialer.Dial(target, nil)
	if err != nil {
		panic(err)
	}
	c.Close()
	c, response, err := websocket.DefaultDialer.Dial(target, nil)
	if c != nil {
		c.Close()
		panic("one-time ticket was reused")
	}
	if err == nil || response == nil || response.StatusCode != 401 {
		panic("consumed ticket did not return 401")
	}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(k + " required")
	}
	return v
}
