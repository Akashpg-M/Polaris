package stream

import (
	"context"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/internal/core/domain/pb"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/proto"
)

const DeadLetterTopic = "telemetry.dead_letters"

type KafkaPostgresArchiver struct {
	reader *kafka.Reader
	writer *kafka.Writer // Used to emit bad payloads to the DLQ topic
	db     *sqlx.DB
}

func NewKafkaPostgresArchiver(brokerURL, postgresURL string) (*KafkaPostgresArchiver, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerURL},
		Topic:    repository.TelemetryTopic,
		GroupID:  "polaris_archive_group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    DeadLetterTopic,
		Balancer: &kafka.Hash{},
	}

	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		return nil, err
	}

	return &KafkaPostgresArchiver{reader: reader, writer: writer, db: db}, nil
}

func (a *KafkaPostgresArchiver) Start(ctx context.Context) {
	slog.Info("Fault-Tolerant Kafka Time-Series Archiver Worker Active")

	for {
		select {
		case <-ctx.Done():
			a.reader.Close()
			a.writer.Close()
			return
		default:
			// Fetch the raw binary data package from Kafka
			msg, err := a.reader.ReadMessage(ctx)
			if err != nil {
				continue
			}

			// 1. Attempt binary parsing
			var payload pb.SpatialObject
			if err := proto.Unmarshal(msg.Value, &payload); err != nil {
				slog.Warn("Failed parsing binary stream packet. Shifting to DLQ.")
				a.sendToDLQ(ctx, msg.Key, msg.Value, "protobuf_unmarshal_failed")
				continue
			}

			// 2. Attempt relational long-term persistence execution
			_, dbErr := a.db.ExecContext(ctx, `
				INSERT INTO telemetry_history (tenant_id, node_id, asset_class, lat, lon, status, battery, recorded_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
				payload.TenantId, payload.Id, int(payload.Type), payload.Lat, payload.Lon, string(payload.Status), payload.EnergyPercent, payload.Timestamp.AsTime(),
			)

			// 3. If database constraints reject it, isolate and continue the group pipeline
			if dbErr != nil {
				slog.Error("Database constraint failure. Dropping packet to DLQ.", "node_id", payload.Id, "err", dbErr)
				a.sendToDLQ(ctx, msg.Key, msg.Value, dbErr.Error())
				continue
			}
		}
	}
}

func (a *KafkaPostgresArchiver) sendToDLQ(ctx context.Context, key []byte, value []byte, reason string) {
	_ = a.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Headers: []kafka.Header{
			{Key: "error_reason", Value: []byte(reason)},
			{Key: "failed_at", Value: []byte(time.Now().UTC().String())},
		},
	})
}