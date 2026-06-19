package handler

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/ports"
	"github.com/Akashpg-M/polaris/backend/internal/core/domain/pb"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		return
	}
	defer conn.Close()

	var nodeID string

	for {
		// 1. Read Raw Binary instead of JSON
		msgType, data, err := conn.ReadMessage()
		if err != nil || msgType != websocket.BinaryMessage {
			log.Printf("[Gateway] Connection lost or non-binary message sent.")
			break
		}

		// 2. Ultra-fast Protobuf Unmarshaling
		var payload pb.SpatialObject
		if err := proto.Unmarshal(data, &payload); err != nil {
			slog.Warn("Failed to decode protobuf payload")
			continue
		}

		if nodeID == "" {
			nodeID = payload.Id // Protobuf generates uppercase field names
			h.registry.Register(nodeID, conn)
		}

		// Inject server-side timestamp using Protobuf's specific time format
		payload.Timestamp = timestamppb.Now()

		// Push to Redis 
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = h.publisher.Publish(ctx, &payload)
		cancel()

		if err != nil {
			continue
		}
	}

	if nodeID != "" {
		h.registry.Unregister(nodeID)
	}
}