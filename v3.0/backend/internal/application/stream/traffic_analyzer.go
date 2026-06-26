package stream

import (
	"context"
	"log/slog"
	"math"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/algo_/graph"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type TrafficAnalyzer struct {
	reader  *kafka.Reader
	network *graph.RoadNetwork
}

func NewTrafficAnalyzer(brokerURL string, network *graph.RoadNetwork) *TrafficAnalyzer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerURL},
		Topic:   KafkaTelemetryTopic, // Using the constant we defined in archiver.go
		GroupID: "polaris_traffic_group", // A distinct group so it gets its own copy of the data
	})

	return &TrafficAnalyzer{
		reader:  reader,
		network: network,
	}
}

func (t *TrafficAnalyzer) Start(ctx context.Context) {
	slog.Info("Dynamic Traffic Analyzer Online. Monitoring Kafka stream for congestion events...")

	for {
		select {
		case <-ctx.Done():
			t.reader.Close()
			return
		default:
			msg, err := t.reader.ReadMessage(ctx)
			if err != nil {
				continue
			}

			var payload pb.SpatialObject
			if err := proto.Unmarshal(msg.Value, &payload); err != nil {
				continue
			}

			// We only care about moving assets for traffic analysis
			if payload.Type == pb.NodeType_NODE_TYPE_STATIC_SENSOR {
				continue
			}

			t.processCongestion(&payload)
		}
	}
}

func (t *TrafficAnalyzer) processCongestion(payload *pb.SpatialObject) {
	// 1. Snap the drone/vehicle to the nearest intersection
	nearestNode, err := t.network.GetNearestIntersection(payload.Lat, payload.Lon)
	if err != nil {
		return
	}

	// 2. Calculate the Congestion Multiplier
	// Assuming a baseline clear-road speed of ~15 m/s (54 km/h).
	// If a vehicle is doing 3 m/s, the multiplier becomes 5.0 (5x cost to travel).
	baselineSpeed := 15.0
	currentSpeed := math.Max(1.0, float64(payload.VelocityMps)) // Prevent division by zero
	
	multiplier := math.Max(1.0, baselineSpeed/currentSpeed)

	// Cap extreme multipliers so we don't sever the graph entirely
	if multiplier > 10.0 {
		multiplier = 10.0
	}

	// 3. Apply the dynamic weight to the road network
	// Note: If multiplier > 1.5, we log it as a traffic event
	if multiplier > 1.5 {
		slog.Debug("Traffic congestion detected", 
			"node", nearestNode, 
			"velocity_mps", payload.VelocityMps, 
			"new_weight", multiplier)
	}

	_ = t.network.UpdateSegmentCongestion(nearestNode, multiplier)
}