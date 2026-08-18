package outbox

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/segmentio/kafka-go"
)

const (
	LifecycleTopic     = "device.lifecycle.v1"
	TaskTopic          = "task.lifecycle.v1"
	CommandTopic       = "device.command.v1"
	CommandAckTopic    = "device.command.ack.v1"
	CommandResultTopic = "device.command.result.v1"
)

type Relay struct {
	store    *repository.RegistryStore
	writer   *kafka.Writer
	batch    int
	interval time.Duration
	done     chan struct{}
}

func New(store *repository.RegistryStore, broker string, batch int, interval time.Duration) *Relay {
	if batch < 1 {
		batch = 100
	}
	return &Relay{store: store, writer: &kafka.Writer{Addr: kafka.TCP(broker), Balancer: &kafka.Hash{}}, batch: batch, interval: interval, done: make(chan struct{})}
}
func (r *Relay) Start(ctx context.Context) {
	defer close(r.done)
	defer r.writer.Close()
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.flush(ctx)
		case <-ctx.Done():
			shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			r.flush(shutdown)
			cancel()
			return
		}
	}
}
func (r *Relay) flush(ctx context.Context) {
	events, err := r.store.ClaimOutbox(ctx, r.batch)
	if err != nil {
		slog.Error("outbox claim failed", "error", err)
		return
	}
	for _, e := range events {
		topic, key, value := route(e)
		if err = r.writer.WriteMessages(ctx, kafka.Message{Topic: topic, Key: []byte(key), Value: value}); err != nil {
			_ = r.store.MarkOutboxFailed(ctx, e.OutboxID, err.Error())
			continue
		}
		if err = r.store.MarkOutboxPublished(ctx, e.OutboxID); err != nil {
			slog.Error("outbox publish marker failed; event will replay", "event_id", e.EventID, "error", err)
		}
	}
}

func route(e repository.OutboxEvent) (string, string, []byte) {
	key := e.TenantID + ":" + e.EventID
	value, _ := json.Marshal(map[string]interface{}{"event_id": e.EventID, "event_type": e.EventType, "schema_version": 1, "tenant_id": e.TenantID, "payload": json.RawMessage(e.Payload)})
	if e.EventType == "command.created.v1" || e.EventType == "command.retry.requested.v1" {
		var envelope struct {
			DeviceID string `json:"device_id"`
		}
		if json.Unmarshal(e.Payload, &envelope) == nil && envelope.DeviceID != "" {
			key = e.TenantID + ":" + envelope.DeviceID
		}
		return CommandTopic, key, e.Payload
	}
	if e.EventType == "command.acknowledged.v1" {
		return CommandAckTopic, key, value
	}
	if e.EventType == "command.result.v1" {
		return CommandResultTopic, key, value
	}
	if len(e.EventType) >= 5 && e.EventType[:5] == "task." {
		return TaskTopic, key, value
	}
	return LifecycleTopic, key, value
}
