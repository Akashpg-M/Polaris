// package handler

// import (
// 	"context"
// 	"log/slog"
// 	"net/http"
// 	"time"

// 	"github.com/Akashpg-M/polaris/backend/internal/core/actor"
// 	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
// 	"github.com/Akashpg-M/polaris/backend/internal/core/ports"
// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"google.golang.org/protobuf/proto"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// type IngestionHandler struct {
// 	// publisher ports.TelemetryPublisher
// 	// registry  *ConnectionRegistry
// 	actorRegistry *actor.ActorRegistry
// }

// func NewIngestionHandler(pub ports.TelemetryPublisher, reg *ConnectionRegistry) *IngestionHandler {
// 	return &IngestionHandler{actorRegistry: reg}
// }

// func (h *IngestionHandler) HandleIoTConnection(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		slog.Error("[Gateway] WebSocket upgrade failed", "error", err)
// 		return
// 	}
// 	defer conn.Close()

// 	var nodeID string

// 	for {
// 		// 1. Read Raw Binary format packets matching pure Protobuf specs
// 		msgType, data, err := conn.ReadMessage()
// 		if err != nil {
// 			// Clean break on normal client disconnects (tab closures/purges)
// 			break
// 		}

// 		if msgType != websocket.BinaryMessage {
// 			slog.Warn("[Gateway] Security violation: Non-binary payload dropped from client context.")
// 			continue // Drop corrupt frames without destroying the underlying persistent TCP pipe
// 		}

// 		// 2. Ultra-fast Protobuf Unmarshaling (Zero JSON parsing allocations)
// 		var payload pb.SpatialObject
// 		if err := proto.Unmarshal(data, &payload); err != nil {
// 			slog.Warn("[Gateway] Failed to decode protobuf binary wire format", "error", err)
// 			continue
// 		}

// 		// First-time registration handshake check
// 		if nodeID == "" {
// 			nodeID = payload.Id 
// 			if nodeID != "" {
// 				h.registry.Register(nodeID, conn)
// 				slog.Info("[Gateway] Uplink channel mapped successfully", "node_id", nodeID)
// 			}
// 		}

// 		// 3. Inject server-side timestamp for high-resolution transit metrics
// 		payload.Timestamp = time.Now().UnixMilli()

// 		// 4. Forward straight to Redis Streams / Kafka via Port interface
// 		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 		err = h.publisher.Publish(ctx, &payload)
// 		cancel()

// 		if err != nil {
// 			slog.Error("[Gateway] Core event emission failed, dropping payload but keeping socket alive", "node_id", nodeID, "error", err)
// 			continue
// 		}
// 	}

// 	// 5. Clean up tracking states safely when connection drops
// 	if nodeID != "" {
// 		h.registry.Unregister(nodeID)
// 		slog.Info("[Gateway] Uplink disconnected and cleaned from registry", "node_id", nodeID)
// 	}
// }

package handler

import (
	"log/slog"
	"net/http"
	"time"
	"errors"
	"sync/atomic"

	"github.com/Akashpg-M/polaris/backend/internal/core/actor"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
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
}

func NewIngestionHandler(reg *actor.ActorRegistry) *IngestionHandler {
	return &IngestionHandler{
		actorRegistry: reg,
		activeSockets: 0, // High-performance, lock-free counter
	}
}

func (h *IngestionHandler) GetActiveConnectionsCount() int64 {
	return atomic.LoadInt64(&h.activeSockets)
}

func (h *IngestionHandler) HandleIoTConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("[Gateway] WebSocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&h.activeSockets, 1)
	defer atomic.AddInt64(&h.activeSockets, -1)
	
	var nodeID string

	for {
		// 1. Read Raw Binary format matching pure Protobuf specs
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType != websocket.BinaryMessage {
			slog.Warn("[Gateway] Security violation: Non-binary payload dropped.")
			continue
		}

		// 2. Fast Protobuf Unmarshaling (Executed on the edge worker thread)
		var payload pb.SpatialObject
		if err := proto.Unmarshal(data, &payload); err != nil {
			slog.Warn("[Gateway] Failed to decode protobuf binary wire format", "error", err)
			continue
		}

		// Initial connection mapping handshake check
		if nodeID == "" {
			nodeID = payload.Id
			slog.Info("[Gateway] Device mapped to local gateway workspace", "node_id", nodeID)
		}

		// Inject server-side ingestion timestamp for profiling metric tracking
		payload.Timestamp = time.Now().UnixMilli()

		// 3. ROUTING LAYER BOUNDARY: Fetch actor and push payload to its inbox mailbox channel
		assetActor := h.actorRegistry.GetOrCreate(nodeID)
		
		// Enforce backpressure safety loops
		if err := assetActor.Push(actor.TelemetryMsg{Payload: &payload}); err != nil {
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