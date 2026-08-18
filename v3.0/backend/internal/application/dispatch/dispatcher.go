package dispatch

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/outbox"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type Dispatcher struct {
	reader *kafka.Reader
	redis  *redis.Client
	owners *repository.ConnectionOwnershipStore
	health atomic.Int64
	done   chan struct{}
}

func New(broker string, redisClient *redis.Client, owners *repository.ConnectionOwnershipStore) *Dispatcher {
	reader := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{broker}, Topic: outbox.CommandTopic, GroupID: "polaris-command-dispatcher", MinBytes: 1, MaxBytes: 10e6, CommitInterval: 0})
	d := &Dispatcher{reader: reader, redis: redisClient, owners: owners, done: make(chan struct{})}
	d.health.Store(time.Now().UnixMilli())
	return d
}

func (d *Dispatcher) Start(ctx context.Context) {
	defer close(d.done)
	defer d.reader.Close()
	for {
		message, err := d.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("command dispatcher fetch failed", "error", err)
			continue
		}
		var envelope command.Envelope
		if json.Unmarshal(message.Value, &envelope) != nil || envelope.FrameType != "COMMAND" || envelope.PartitionKey() != string(message.Key) {
			slog.Error("invalid durable command envelope", "partition", message.Partition, "offset", message.Offset)
			if err = d.reader.CommitMessages(ctx, message); err != nil {
				slog.Error("invalid command offset commit failed", "error", err)
			}
			continue
		}
		ownership, lookupErr := d.owners.Get(ctx, envelope.TenantID, envelope.DeviceID)
		if lookupErr == nil && ownership.LeaseExpiresAt.After(time.Now()) {
			if err = d.redis.Publish(ctx, repository.GatewayCommandChannel(ownership.GatewayID), message.Value).Err(); err != nil {
				slog.Error("command routing notification failed", "command_id", envelope.CommandID, "error", err)
				continue
			}
		} else if lookupErr != nil && !errors.Is(lookupErr, repository.ErrNotFound) {
			slog.Error("command ownership lookup failed", "command_id", envelope.CommandID, "error", lookupErr)
			continue
		}
		if err = d.reader.CommitMessages(ctx, message); err != nil {
			slog.Error("command dispatcher commit failed; notification may replay", "command_id", envelope.CommandID, "error", err)
			continue
		}
		d.health.Store(time.Now().UnixMilli())
	}
}

func (d *Dispatcher) Healthy() bool {
	return d.health.Load() > 0
}

func (d *Dispatcher) Wait(ctx context.Context) error {
	select {
	case <-d.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
