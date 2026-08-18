package stream

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/graph"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/segmentio/kafka-go"
)

type TrafficAnalyzer struct {
	reader  *kafka.Reader
	dlq     *kafka.Writer
	network *graph.RoadNetwork
}

func NewTrafficAnalyzer(brokerURL string, network *graph.RoadNetwork) *TrafficAnalyzer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerURL},
		Topic:          KafkaTelemetryTopic,     // Using the constant we defined in archiver.go
		GroupID:        "polaris_traffic_group", // A distinct group so it gets its own copy of the data
		CommitInterval: 0,
	})

	return &TrafficAnalyzer{
		reader:  reader,
		dlq:     &kafka.Writer{Addr: kafka.TCP(brokerURL), Topic: DeadLetterTopic, Balancer: &kafka.Hash{}},
		network: network,
	}
}

func (t *TrafficAnalyzer) Start(ctx context.Context) {
	slog.Info("Dynamic Traffic Analyzer Online. Monitoring Kafka stream for congestion events...")

	defer t.reader.Close()
	defer t.dlq.Close()
	for {
		msg, err := t.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}
		envelope, parseErr := events.Unmarshal(msg.Value)
		if parseErr != nil {
			for {
				err = t.dlq.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: msg.Value, Headers: []kafka.Header{
					{Key: "error_reason", Value: []byte(parseErr.Error())}, {Key: "consumer", Value: []byte("polaris_traffic_group")},
					{Key: "source_partition", Value: []byte(fmt.Sprint(msg.Partition))}, {Key: "source_offset", Value: []byte(fmt.Sprint(msg.Offset))},
				}})
				if err == nil {
					break
				}
				select {
				case <-time.After(250 * time.Millisecond):
				case <-ctx.Done():
					return
				}
			}
		} else if envelope.Payload.Type != pb.NodeType_NODE_TYPE_STATIC_SENSOR {
			t.processCongestion(envelope.Payload)
		}
		for {
			if err = t.reader.CommitMessages(ctx, msg); err == nil {
				break
			}
			slog.Error("traffic Kafka commit failed; retrying same offset", "partition", msg.Partition, "offset", msg.Offset, "error", err)
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				return
			}
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
