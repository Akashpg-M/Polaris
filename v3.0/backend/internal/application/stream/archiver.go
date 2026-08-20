package stream

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

const KafkaTelemetryTopic = "telemetry.ingress"
const DeadLetterTopic = "telemetry.dead-letter.v1"

type KafkaPostgresArchiver struct {
	reader       *kafka.Reader
	writer       *kafka.Writer
	db           *sqlx.DB
	done         chan struct{}
	maxRetries   int
	lastProgress atomic.Int64
}

func NewKafkaPostgresArchiver(brokerURL, postgresURL string) (*KafkaPostgresArchiver, error) {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)
	a := &KafkaPostgresArchiver{
		reader: kafka.NewReader(kafka.ReaderConfig{Brokers: []string{brokerURL}, Topic: KafkaTelemetryTopic, GroupID: "polaris_archive_group", CommitInterval: 0}),
		writer: &kafka.Writer{Addr: kafka.TCP(brokerURL), Topic: DeadLetterTopic, Balancer: &kafka.Hash{}},
		db:     db, done: make(chan struct{}), maxRetries: 5,
	}
	a.lastProgress.Store(time.Now().UnixMilli())
	return a, nil
}

func (a *KafkaPostgresArchiver) archive(ctx context.Context, e *events.TelemetryEnvelope) error {
	p := e.Payload
	_, err := a.db.ExecContext(ctx, `
		INSERT INTO telemetry_history
		(event_id, tenant_id, device_id, device_boot_id, sequence_number, asset_type,
		 lat, lon, geom, status, velocity_mps, heading_deg, battery, observed_at,
		 ingested_at, schema_version, correlation_id, recorded_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,ST_SetSRID(ST_MakePoint($8,$7),4326),$9,$10,$11,$12,$13::timestamptz,$14::timestamptz,$15,$16,($13::timestamptz AT TIME ZONE 'UTC'))
		ON CONFLICT DO NOTHING`,
		e.EventID, e.TenantID, e.DeviceID, e.DeviceBootID, e.SequenceNumber, int(p.Type),
		p.Lat, p.Lon, int(p.Status), p.VelocityMps, p.HeadingDeg, p.EnergyPercent,
		time.UnixMilli(e.ObservedAt), time.UnixMilli(e.IngestedAt), e.SchemaVersion, e.CorrelationID)
	return err
}

func (a *KafkaPostgresArchiver) sendToDLQ(ctx context.Context, msg kafka.Message, reason string) error {
	return a.writer.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: msg.Value, Headers: []kafka.Header{
		{Key: "error_reason", Value: []byte(reason)}, {Key: "source_topic", Value: []byte(msg.Topic)},
		{Key: "source_partition", Value: []byte(fmt.Sprint(msg.Partition))}, {Key: "source_offset", Value: []byte(fmt.Sprint(msg.Offset))},
		{Key: "failed_at", Value: []byte(time.Now().UTC().Format(time.RFC3339Nano))},
	}})
}

func (a *KafkaPostgresArchiver) process(ctx context.Context, msg kafka.Message) bool {
	e, err := events.Unmarshal(msg.Value)
	if err != nil {
		if dlqErr := a.sendToDLQ(ctx, msg, err.Error()); dlqErr != nil {
			slog.Error("archive poison event DLQ failed", "error", dlqErr)
			return false
		}
		return true
	}
	var lastErr error
	for attempt := 1; attempt <= a.maxRetries; attempt++ {
		if err := a.archive(ctx, e); err == nil {
			return true
		} else {
			lastErr = err
		}
		slog.Warn("transient PostgreSQL archive failure", "event_id", e.EventID, "attempt", attempt, "error", lastErr)
		if attempt < a.maxRetries {
			select {
			case <-time.After(time.Duration(attempt*50) * time.Millisecond):
			case <-ctx.Done():
				return false
			}
		}
	}
	if err := a.sendToDLQ(ctx, msg, "retry_exhausted: "+lastErr.Error()); err != nil {
		slog.Error("archive retry-exhausted DLQ failed", "error", err)
		return false
	}
	return true
}

func (a *KafkaPostgresArchiver) Start(ctx context.Context) {
	defer close(a.done)
	defer a.reader.Close()
	defer a.writer.Close()
	defer a.db.Close()
	slog.Info("Idempotent Kafka PostgreSQL archiver active")
	for {
		msg, err := a.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("archive consumer shutdown complete")
				return
			}
			slog.Error("archive fetch failed", "error", err)
			continue
		}
		for !a.process(ctx, msg) {
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				slog.Error("archive shutdown with uncommitted message", "partition", msg.Partition, "offset", msg.Offset)
				return
			}
		}
		for {
			if err := a.reader.CommitMessages(ctx, msg); err == nil {
				break
			} else {
				slog.Error("archive Kafka commit failed; retrying same offset", "partition", msg.Partition, "offset", msg.Offset, "error", err)
			}
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				return
			}
		}
		a.lastProgress.Store(time.Now().UnixMilli())
	}
}

func (a *KafkaPostgresArchiver) Ready(ctx context.Context) error { return a.db.PingContext(ctx) }
func (a *KafkaPostgresArchiver) Healthy() bool {
	select {
	case <-a.done:
		return false
	default:
		return true
	}
}
func (a *KafkaPostgresArchiver) Wait(ctx context.Context) error {
	select {
	case <-a.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
