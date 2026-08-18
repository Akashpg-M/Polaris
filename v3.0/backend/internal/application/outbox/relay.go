package outbox

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/segmentio/kafka-go"
)

const LifecycleTopic = "device.lifecycle.v1"

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
	return &Relay{store: store, writer: &kafka.Writer{Addr: kafka.TCP(broker), Topic: LifecycleTopic, Balancer: &kafka.Hash{}}, batch: batch, interval: interval, done: make(chan struct{})}
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
		value, _ := json.Marshal(map[string]interface{}{"event_id": e.EventID, "event_type": e.EventType, "schema_version": 1, "tenant_id": e.TenantID, "payload": json.RawMessage(e.Payload)})
		if err = r.writer.WriteMessages(ctx, kafka.Message{Key: []byte(e.TenantID + ":" + e.EventID), Value: value}); err != nil {
			_ = r.store.MarkOutboxFailed(ctx, e.OutboxID, err.Error())
			continue
		}
		if err = r.store.MarkOutboxPublished(ctx, e.OutboxID); err != nil {
			slog.Error("outbox publish marker failed; event will replay", "event_id", e.EventID, "error", err)
		}
	}
}
