package repository

import (
	"context"
	"fmt"
	"github.com/Akashpg-M/polaris/backend/internal/core/domain/pb"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

const TelemetryStreamKey = "telemetry:ingress"

type RedisStreamAdapter struct {
	client *redis.Client
}

func NewRedisStreamAdapter(redisURL string) (*RedisStreamAdapter, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}
	
	client := redis.NewClient(opts)
	
	// Ping to ensure connection is alive
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
  
	return &RedisStreamAdapter{client: client}, nil
}

func (r *RedisStreamAdapter) Publish(ctx context.Context, payload *pb.SpatialObject) error {
	// 1. Serialize to ultra-compact binary byte array
	data, err := proto.Marshal(payload)
	if err != nil {
		return fmt.Errorf("protobuf marshal failed: %w", err)
	}

	err = r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: TelemetryStreamKey,
		MaxLen: 100000, 
		Approx: true,
		Values: map[string]interface{}{
			"data": data, // Store raw bytes in Redis
		},
	}).Err()

	return err
}