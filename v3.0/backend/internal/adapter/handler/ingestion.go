package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/actor"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type IngestionHandler struct {
	// Replacing legacy publisher interface with our new Actor Partition Registry
	actorRegistry *actor.ActorRegistry
	activeSockets int64
	authenticator DeviceAuthenticator
	connections   *DeviceConnectionManager
}

type DeviceAuthenticator interface {
	ResolveDevice(context.Context, string) (auth.DevicePrincipal, error)
	ConsumeTicket(context.Context, string) (auth.DevicePrincipal, error)
	RevalidateDevice(context.Context, auth.DevicePrincipal) error
}

func NewIngestionHandler(reg *actor.ActorRegistry, authenticator DeviceAuthenticator, connections *DeviceConnectionManager) *IngestionHandler {
	return &IngestionHandler{
		actorRegistry: reg,
		activeSockets: 0, // High-performance, lock-free counter
		authenticator: authenticator,
		connections:   connections,
	}
}

func (h *IngestionHandler) GetActiveConnectionsCount() int64 {
	return atomic.LoadInt64(&h.activeSockets)
}

func (h *IngestionHandler) HandleIoTConnection(c *gin.Context) {
	principal, err := h.authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "DEVICE_AUTHENTICATION_FAILED", "message": "Device credential is invalid or inactive"}})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("[Gateway] WebSocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()
	conn.SetReadLimit(events.MaxFrameBytes)
	session, err := h.connections.Register(c.Request.Context(), conn, principal)
	if err != nil {
		rejectFrame(conn, "connection ownership could not be established")
		return
	}
	defer h.connections.Unregister(context.Background(), session)

	atomic.AddInt64(&h.activeSockets, 1)
	defer atomic.AddInt64(&h.activeSockets, -1)

	nodeID, tenantID := principal.DeviceID, principal.TenantID
	var bootID string

	for {
		// 1. Read Raw Binary format matching pure Protobuf specs
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if err := h.authenticator.RevalidateDevice(c.Request.Context(), principal); err != nil {
			session.close(websocket.ClosePolicyViolation, "credential or device was revoked")
			return
		}
		if msgType == websocket.TextMessage {
			if err := h.handleDeviceControl(c.Request.Context(), principal, data); err != nil {
				session.close(websocket.ClosePolicyViolation, "command response rejected")
				return
			}
			continue
		}
		if msgType != websocket.BinaryMessage {
			slog.Warn("[Gateway] Security violation: unsupported WebSocket frame dropped.")
			continue
		}

		// 2. Fast Protobuf Unmarshaling (Executed on the edge worker thread)
		var payload pb.SpatialObject
		if err := proto.Unmarshal(data, &payload); err != nil {
			rejectFrame(conn, "malformed protobuf frame")
			return
		}
		if err := events.ValidateFrame(&payload, time.Now()); err != nil {
			slog.Warn("[Gateway] Rejected invalid telemetry before Kafka", "error", err)
			rejectFrame(conn, err.Error())
			return
		}
		if payload.Id != principal.DeviceID || payload.TenantId != principal.TenantID {
			rejectFrame(conn, "payload identity does not match authenticated device")
			return
		}
		// The principal, not the untrusted frame, owns platform identity.
		payload.Id = principal.DeviceID
		payload.TenantId = principal.TenantID

		// Initial connection mapping handshake check
		if bootID == "" {
			bootID = payload.DeviceBootId
			slog.Info("[Gateway] Device mapped to local gateway workspace", "node_id", nodeID)
		} else if payload.DeviceBootId != bootID {
			slog.Warn("[Gateway] Rejected device identity change within connection", "expected_device", nodeID, "actual_device", payload.Id)
			rejectFrame(conn, "device identity changed within connection")
			return
		}

		ingestedAt := time.Now().UTC()
		if payload.Timestamp == 0 {
			payload.Timestamp = payload.ObservedAt
		}
		envelope := events.NewTelemetryEnvelope(&payload, ingestedAt,
			c.GetHeader("X-Correlation-ID"), c.GetHeader("X-Causation-ID"), c.GetHeader("traceparent"))

		// 3. ROUTING LAYER BOUNDARY: Fetch actor and push payload to its inbox mailbox channel
		assetActor := h.actorRegistry.GetOrCreate(tenantID + ":" + nodeID)

		// Enforce backpressure safety loops
		if err := assetActor.Push(actor.TelemetryMsg{Payload: &payload, Envelope: envelope}); err != nil {
			if errors.Is(err, actor.ErrMailboxSaturated) {
				slog.Error("[Gateway] System backpressure hit, dropping frame to preserve stability", "node_id", nodeID)
				// Optional: In a critical system, you can signal the edge drone to back off here
			}
			continue
		}
	}

	// Clean up runtime structures safely if the persistent socket drops
	if nodeID != "" {
		slog.Info("[Gateway] Telemetry channel closed at edge boundary", "node_id", nodeID)
		// NOTE: In an event-sourced distributed control system, we do NOT destroy the actor immediately
		// because the actor might still have pending messages to clear in its channel mailbox queue!
	}
}

func (h *IngestionHandler) handleDeviceControl(ctx context.Context, principal auth.DevicePrincipal, data []byte) error {
	var header struct {
		FrameType string `json:"frame_type"`
	}
	if len(data) > 64*1024 || json.Unmarshal(data, &header) != nil {
		return errors.New("invalid device control envelope")
	}
	switch header.FrameType {
	case "COMMAND_ACK":
		var ack command.Ack
		if json.Unmarshal(data, &ack) != nil || ack.CommandID == "" || ack.SequenceNumber < 1 {
			return errors.New("invalid command acknowledgement")
		}
		if ack.Status != "ACCEPTED" && ack.Status != "REJECTED" && ack.Status != "DUPLICATE" && ack.Status != "EXPIRED" && ack.Status != "UNSUPPORTED" {
			return errors.New("unsupported acknowledgement status")
		}
		if err := h.connections.store.ApplyCommandAck(ctx, principal.TenantID, principal.DeviceID, ack); err != nil {
			return err
		}
		h.connections.metrics.CommandsAcknowledged.Add(1)
		return nil
	case "COMMAND_RESULT":
		var result command.Result
		if json.Unmarshal(data, &result) != nil || result.CommandID == "" || result.SequenceNumber < 1 || (result.Status != "SUCCEEDED" && result.Status != "COMPLETED" && result.Status != "FAILED") {
			return errors.New("invalid command result")
		}
		if err := h.connections.store.ApplyCommandResult(ctx, principal.TenantID, principal.DeviceID, result); err != nil {
			return err
		}
		if result.Status == "SUCCEEDED" || result.Status == "COMPLETED" {
			h.connections.metrics.TasksCompleted.Add(1)
		} else {
			h.connections.metrics.CommandsFailed.Add(1)
			h.connections.metrics.TasksFailed.Add(1)
		}
		return nil
	default:
		return errors.New("unsupported device frame type")
	}
}

func (h *IngestionHandler) authenticate(c *gin.Context) (auth.DevicePrincipal, error) {
	if h.authenticator == nil {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	if ticket := c.Query("ticket"); ticket != "" {
		return h.authenticator.ConsumeTicket(c.Request.Context(), ticket)
	}
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	return h.authenticator.ResolveDevice(c.Request.Context(), strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
}

func rejectFrame(conn *websocket.Conn, reason string) {
	message := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, fmt.Sprintf("telemetry rejected: %s", reason))
	_ = conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
}
