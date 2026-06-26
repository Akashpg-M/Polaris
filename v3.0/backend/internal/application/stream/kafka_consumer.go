package stream

import (
	"context"
	"log/slog"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type KafkaConsumer struct {
	reader *kafka.Reader
	engine *spatial.Engine
}

func NewKafkaConsumer(brokerURL string, engine *spatial.Engine) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerURL},
		Topic:   KafkaTelemetryTopic, 
		GroupID: "polaris_engine_group",
	})

	return &KafkaConsumer{
		reader: reader,
		engine: engine,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context, workerID string) {
	slog.Info("Kafka Geohash Consumer Worker Started", "worker_id", workerID)

	var batch []*pb.SpatialObject
	batchSize := 1000

	for {
		select {
		case <-ctx.Done():
			c.reader.Close()
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				continue
			}

			var payload pb.SpatialObject
			if err := proto.Unmarshal(msg.Value, &payload); err == nil {
				batch = append(batch, &payload)
			}

			if len(batch) >= batchSize {
				c.engine.BatchUpdate(batch)
				batch = nil
			}
		}
	}
}