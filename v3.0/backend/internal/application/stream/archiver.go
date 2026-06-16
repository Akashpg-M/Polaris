package stream

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/internal/core/domain"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"
)

const DeadLetterStreamKey = "telemetry:dead_letters"

type PostgresArchiver struct {
	redisClient *redis.Client
	db          *sqlx.DB
}

func NewPostgresArchiver(redisURL, postgresURL string) (*PostgresArchiver, error) {
	rOpt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(rOpt)

	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		return nil, err
	}

	redisClient.XGroupCreateMkStream(context.Background(), StreamName, "polaris_archive_group", "$")

	return &PostgresArchiver{redisClient: redisClient, db: db}, nil
}

func (a *PostgresArchiver) Start(ctx context.Context) {
	slog.Info("Fault-Tolerant Time-Series Archiver Worker Started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			streams, err := a.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    "polaris_archive_group",
				Consumer: "archiver-node-1",
				Streams:  []string{StreamName, ">"},
				Count:    500,
				Block:    5 * time.Second,
			}).Result()

			if err != nil || len(streams) == 0 {
				continue
			}

			// Individual message processing safeguards healthy records from batch crashes
			for _, msg := range streams[0].Messages {
				var p domain.TelemetryPayload
				rawData := msg.Values["data"].(string)

				// 1. Unmarshal check
				if err := json.Unmarshal([]byte(rawData), &p); err != nil {
					slog.Warn("Malformed JSON caught. Shifting to DLQ.", "msg_id", msg.ID)
					a.sendToDLQ(ctx, msg.ID, rawData, "json_unmarshal_failed")
					a.redisClient.XAck(ctx, StreamName, "polaris_archive_group", msg.ID)
					continue
				}

				// 2. Individual Database Insert
				_, dbErr := a.db.ExecContext(ctx, `
					INSERT INTO telemetry_history (tenant_id, node_id, asset_class, lat, lon, status, battery, recorded_at) 
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
					p.TenantID, p.NodeID, p.Class, p.Lat, p.Lon, p.Status, p.Battery, p.Timestamp,
				)

				// 3. Database constraint handling
				if dbErr != nil {
					slog.Error("Database constraint failure. Purging row to DLQ.", "node", p.NodeID, "error", dbErr)
					a.sendToDLQ(ctx, msg.ID, rawData, dbErr.Error())
					a.redisClient.XAck(ctx, StreamName, "polaris_archive_group", msg.ID)
					continue
				}

				// 4. Confirm successful extraction to Redis pipeline
				a.redisClient.XAck(ctx, StreamName, "polaris_archive_group", msg.ID)
			}
		}
	}
}

// sendToDLQ pushes failed events into a separate Redis stream for later debugging
func (a *PostgresArchiver) sendToDLQ(ctx context.Context, sourceID string, payload string, reason string) {
	_ = a.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: DeadLetterStreamKey,
		MaxLen: 10000,
		Approx: true,
		Values: map[string]interface{}{
			"original_message_id": sourceID,
			"payload":             payload,
			"failure_reason":      reason,
			"failed_at":           time.Now().UTC().String(),
		},
	}).Err()
}