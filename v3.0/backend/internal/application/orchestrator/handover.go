package orchestrator

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Akashpg-M/polaris/backend/algo_/geo"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/segmentio/kafka-go"
	"github.com/uber/h3-go/v4"
)

const (
	HandoverTopic    = "telemetry.system.handovers"
	H3Resolution     = 7
	// At Resolution 7, a hex radius is roughly 1.2km. 
	// We trigger handover when they pass the 1.0km threshold.
	HandoverThresholdKm = 1.0 
)

// HandoverState is the memory context transferred between Engine Nodes
type HandoverState struct {
	NodeID          string      `json:"node_id"`
	AssetClass      pb.NodeType `json:"asset_class"`
	TargetH3Hex     string      `json:"target_h3_hex"`
	ActiveDirective string      `json:"active_directive"` 
	TargetLat       float64     `json:"target_lat"`
	TargetLon       float64     `json:"target_lon"`
}

type HandoverManager struct {
	engine *spatial.Engine
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewHandoverManager(brokerURL string, engine *spatial.Engine) *HandoverManager {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    HandoverTopic,
		Balancer: &kafka.Hash{}, // Ensures target hex routes to the correct node
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerURL},
		Topic:   HandoverTopic,
		GroupID: "polaris_engine_group", // Shares the same consumer group load balancing
	})

	return &HandoverManager{
		engine: engine,
		writer: writer,
		reader: reader,
	}
}

// EvaluateBoundary is called on every telemetry ping by the Engine
func (h *HandoverManager) EvaluateBoundary(ctx context.Context, payload *pb.SpatialObject, activeCmd *CommandPayload) {
	// 1. Find the center of the current H3 cell
	latLng := h3.NewLatLng(payload.Lat, payload.Lon)
	currentCell, err := h3.LatLngToCell(latLng, H3Resolution)
	if err != nil {
		slog.Error("Failed to compute H3 cell", "error", err)
		return
	}
	cellCenter, err := currentCell.LatLng()
	if err != nil {
		slog.Error("Failed to compute H3 cell center", "error", err)
		return
	}
	// 2. Check distance from the drone to the absolute center of the hexagon
	distToCenter := geo.Haversine(payload.Lat, payload.Lon, cellCenter.Lat, cellCenter.Lng)
	
	// 3. If near the edge, project the next cell and initiate Handover Broadcast
	if distToCenter >= HandoverThresholdKm {
		// Project slightly forward based on heading to figure out WHICH neighbor they are entering
		projectedLat, projectedLon := geo.ProjectCoordinate(payload.Lat, payload.Lon, float64(payload.HeadingDeg), 0.2) // Project 200 meters forward		
		
		targetCell, err := h3.LatLngToCell(h3.NewLatLng(projectedLat, projectedLon), H3Resolution)
		if err != nil {
			slog.Error("Failed to compute target H3 cell", "error", err)
			return
		}		
		if targetCell != currentCell {
			h.emitHandover(ctx, payload, targetCell.String(), activeCmd)
		}
	}
}

// emitHandover pushes the memory context to the Kafka Event Fabric
func (h *HandoverManager) emitHandover(ctx context.Context, payload *pb.SpatialObject, targetHex string, activeCmd *CommandPayload) {
	state := HandoverState{
		NodeID:      payload.Id,
		AssetClass:  payload.Type,
		TargetH3Hex: targetHex,
	}

	if activeCmd != nil {
		state.ActiveDirective = activeCmd.Directive
		state.TargetLat = activeCmd.TargetLat
		state.TargetLon = activeCmd.TargetLon
	}

	data, _ := json.Marshal(state)

	_ = h.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(targetHex), // Using the NEXT hex as the key ensures it routes to the correct destination node
		Value: data,
	})

	slog.Info("Edge Handover Protocol Triggered", "node", payload.Id, "target_hex", targetHex)
}


// ListenForHandovers runs in a background goroutine on all Engine Nodes
func (h *HandoverManager) ListenForHandovers(ctx context.Context) {
	slog.Info("Distributed Edge Handover Protocol Active")

	for {
		select {
		case <-ctx.Done():
			h.reader.Close()
			h.writer.Close()
			return
		default:
			msg, err := h.reader.ReadMessage(ctx)
			if err != nil {
				continue
			}

			var state HandoverState
			if err := json.Unmarshal(msg.Value, &state); err != nil {
				continue
			}

			// Pre-warm the Engine's QuadTree and RAM with the incoming asset
			// Even if the asset hasn't physically connected to our Gateway yet!
			dummyPayload := &pb.SpatialObject{
				Id:     state.NodeID,
				Type:   state.AssetClass,
				Status: pb.NodeStatus_NODE_STATUS_ACTIVE,
			}
			
			h.engine.BatchUpdate([]*pb.SpatialObject{dummyPayload})

			// Re-apply any active orchestrator commands to the local memory
			if state.ActiveDirective != "" {
				slog.Info("Handover state received. Context synchronized.", 
					"node", state.NodeID, 
					"restored_directive", state.ActiveDirective)
				
				// In a full implementation, you would map this restored command back into 
				// your orchestrator's local tracking state here.
			}
		}
	}
}