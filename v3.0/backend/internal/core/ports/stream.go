package ports

import (
	"context"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

type TelemetryPublisher interface {
	Publish(ctx context.Context, payload *pb.SpatialObject) error
}