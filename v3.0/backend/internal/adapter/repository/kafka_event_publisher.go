package repository

import (
	"context"
	"fmt"

	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/segmentio/kafka-go"
)

// KafkaEventPublisher writes the canonical versioned envelope and partitions by
// stable tenant+device identity so all boots and sequences stay ordered.
type KafkaEventPublisher struct {
	writer *kafka.Writer
}

func NewKafkaEventPublisher(brokerURL string) *KafkaEventPublisher {
	return &KafkaEventPublisher{writer: &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    TelemetryTopic,
		Balancer: &kafka.Hash{},
		Async:    false,
	}}
}

func (p *KafkaEventPublisher) PublishEvent(ctx context.Context, _ string, event interface{}) error {
	envelope, ok := event.(*events.TelemetryEnvelope)
	if !ok || envelope == nil {
		return fmt.Errorf("unsupported kafka event type %T", event)
	}
	data, err := envelope.Marshal()
	if err != nil {
		return fmt.Errorf("marshal telemetry envelope: %w", err)
	}
	partitionKey := envelope.PartitionKey()
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(partitionKey), Value: data})
}

func (p *KafkaEventPublisher) Ready(ctx context.Context, brokerURL string) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokerURL)
	if err != nil {
		return err
	}
	return conn.Close()
}

func (p *KafkaEventPublisher) Close() error { return p.writer.Close() }
