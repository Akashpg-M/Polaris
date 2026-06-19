package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Akashpg-M/polaris/backend/internal/core/domain/pb"
	"github.com/segmentio/kafka-go"
	"github.com/uber/h3-go/v4"
	"google.golang.org/protobuf/proto"
)

const TelemetryTopic = "telemetry.ingress"
const H3Resolution = 7 // Resolution 7 gives roughly 1.2km wide hexagons

type KafkaStreamAdapter struct {
	writer *kafka.Writer
}

func NewKafkaStreamAdapter(brokerURL string) *KafkaStreamAdapter {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    TelemetryTopic,
		// Hash balancer ensures the H3 Key determines the partition
		Balancer: &kafka.Hash{}, 
		// Async mode for extreme throughput (acts like a shock absorber)
		Async:    true,
	}

	return &KafkaStreamAdapter{writer: writer}
}

// Publish implements ports.TelemetryPublisher
func (k *KafkaStreamAdapter) Publish(ctx context.Context, payload *pb.SpatialObject) error {
	// 1. Calculate the Spatial Partition Key (Uber H3)
	latLng := h3.NewLatLng(payload.Lat, payload.Lon)
	h3Cell := h3.LatLngToCell(latLng, H3Resolution)
	h3Key := h3Cell.String()

	// 2. Serialize to raw Protocol Buffers
	data, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("protobuf marshal failed: %w", err)
	}

	// 3. Write to Kafka using the H3 Hex as the routing Key
	err = k.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(h3Key),
		Value: data,
	})

	if err != nil {
		slog.Error("Kafka write failed", "error", err)
		return err
	}

	return nil
}

func (k *KafkaStreamAdapter) Close() error {
	return k.writer.Close()
}