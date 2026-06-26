// package handler

// import (
// 	"context"
// 	"log"
// 	"log/slog"
// 	"net/http"
// 	"time"

// 	"github.com/Akashpg-M/polaris/backend/internal/core/ports"
// 	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
// 	"github.com/gin-gonic/gin"
// 	"github.com/gorilla/websocket"
// 	"google.golang.org/protobuf/proto"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// type IngestionHandler struct {
// 	publisher ports.TelemetryPublisher
// 	registry  *ConnectionRegistry
// }

// func NewIngestionHandler(pub ports.TelemetryPublisher, reg *ConnectionRegistry) *IngestionHandler {
// 	return &IngestionHandler{publisher: pub, registry: reg}
// }

// func (h *IngestionHandler) HandleIoTConnection(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	var nodeID string

// 	for {
// 		// 1. Read Raw Binary instead of JSON
// 		msgType, data, err := conn.ReadMessage()
// 		if err != nil || msgType != websocket.BinaryMessage {
// 			log.Printf("[Gateway] Connection lost or non-binary message sent.")
// 			break
// 		}

// 		// 2. Ultra-fast Protobuf Unmarshaling
// 		var payload pb.SpatialObject
// 		if err := proto.Unmarshal(data, &payload); err != nil {
// 			slog.Warn("Failed to decode protobuf payload")
// 			continue
// 		}

// 		if nodeID == "" {
// 			nodeID = payload.Id // Protobuf generates uppercase field names
// 			h.registry.Register(nodeID, conn)
// 		}

// 		// Inject server-side timestamp using Protobuf's specific time format
// 		payload.Timestamp = time.Now().UnixMilli()

// 		// Push to Redis 
// 		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 		err = h.publisher.Publish(ctx, &payload)
// 		cancel()

// 		if err != nil {
// 			continue
// 		}
// 	}

// 	if nodeID != "" {
// 		h.registry.Unregister(nodeID)
// 	}
// }

package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type IngestionHandler struct {
	publisher ports.TelemetryPublisher
	registry  *ConnectionRegistry
}

func NewIngestionHandler(pub ports.TelemetryPublisher, reg *ConnectionRegistry) *IngestionHandler {
	return &IngestionHandler{publisher: pub, registry: reg}
}

func (h *IngestionHandler) HandleIoTConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("[Gateway] WebSocket upgrade failed", "error", err)
		return
	}
	defer conn.Close()

	var nodeID string

	for {
		// 1. Read Raw Binary format packets matching pure Protobuf specs
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			// Clean break on normal client disconnects (tab closures/purges)
			break
		}

		if msgType != websocket.BinaryMessage {
			slog.Warn("[Gateway] Security violation: Non-binary payload dropped from client context.")
			continue // Drop corrupt frames without destroying the underlying persistent TCP pipe
		}

		// 2. Ultra-fast Protobuf Unmarshaling (Zero JSON parsing allocations)
		var payload pb.SpatialObject
		if err := proto.Unmarshal(data, &payload); err != nil {
			slog.Warn("[Gateway] Failed to decode protobuf binary wire format", "error", err)
			continue
		}

		// First-time registration handshake check
		if nodeID == "" {
			nodeID = payload.Id 
			if nodeID != "" {
				h.registry.Register(nodeID, conn)
				slog.Info("[Gateway] Uplink channel mapped successfully", "node_id", nodeID)
			}
		}

		// 3. Inject server-side timestamp for high-resolution transit metrics
		payload.Timestamp = time.Now().UnixMilli()

		// 4. Forward straight to Redis Streams / Kafka via Port interface
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = h.publisher.Publish(ctx, &payload)
		cancel()

		if err != nil {
			slog.Error("[Gateway] Core event emission failed, dropping payload but keeping socket alive", "node_id", nodeID, "error", err)
			continue
		}
	}

	// 5. Clean up tracking states safely when connection drops
	if nodeID != "" {
		h.registry.Unregister(nodeID)
		slog.Info("[Gateway] Uplink disconnected and cleaned from registry", "node_id", nodeID)
	}
}