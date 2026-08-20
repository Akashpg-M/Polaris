package mobility

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/segmentio/kafka-go"
)

type TrafficConsumer struct {
	reader          *kafka.Reader
	dlq             *kafka.Writer
	traffic         *routing.TrafficManager
	refreshInterval time.Duration
	done            chan struct{}
}

func NewTrafficConsumer(broker string, traffic *routing.TrafficManager, refreshInterval time.Duration) *TrafficConsumer {
	return &TrafficConsumer{reader: kafka.NewReader(kafka.ReaderConfig{Brokers: []string{broker}, Topic: "telemetry.ingress", GroupID: "polaris_traffic_group", CommitInterval: 0}), dlq: &kafka.Writer{Addr: kafka.TCP(broker), Topic: "telemetry.dead-letter.v1", Balancer: &kafka.Hash{}}, traffic: traffic, refreshInterval: refreshInterval, done: make(chan struct{})}
}
func (c *TrafficConsumer) Start(ctx context.Context) {
	defer close(c.done)
	defer c.reader.Close()
	defer c.dlq.Close()
	refreshCtx, stopRefresh := context.WithCancel(ctx)
	defer stopRefresh()
	go func() {
		ticker := time.NewTicker(c.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case now := <-ticker.C:
				_ = c.traffic.Refresh(now.UTC())
			case <-refreshCtx.Done():
				return
			}
		}
	}()
	slog.Info("Mobility map-matched traffic consumer started", "refresh_interval", c.refreshInterval, "traffic_scope", "SHARED_TRUSTED")
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			continue
		}
		envelope, parseErr := events.Unmarshal(msg.Value)
		if parseErr != nil {
			for {
				err = c.dlq.WriteMessages(ctx, kafka.Message{Key: msg.Key, Value: msg.Value, Headers: []kafka.Header{{Key: "error_reason", Value: []byte(parseErr.Error())}, {Key: "consumer", Value: []byte("polaris_traffic_group")}, {Key: "source_partition", Value: []byte(fmt.Sprint(msg.Partition))}, {Key: "source_offset", Value: []byte(fmt.Sprint(msg.Offset))}}})
				if err == nil {
					break
				}
				select {
				case <-time.After(250 * time.Millisecond):
				case <-ctx.Done():
					return
				}
			}
		} else if roadTelemetry(envelope.Payload.Type) {
			heading := envelope.Payload.HeadingDeg
			_ = c.traffic.Observe(ctx, routing.TrafficObservation{Position: model.Position{Latitude: envelope.Payload.Lat, Longitude: envelope.Payload.Lon}, HeadingDegrees: &heading, SpeedMPS: envelope.Payload.VelocityMps, ObservedAt: time.UnixMilli(envelope.ObservedAt).UTC()})
		}
		for {
			if err = c.reader.CommitMessages(ctx, msg); err == nil {
				break
			}
			select {
			case <-time.After(250 * time.Millisecond):
			case <-ctx.Done():
				return
			}
		}
	}
}
func roadTelemetry(t pb.NodeType) bool {
	return t == pb.NodeType_NODE_TYPE_BIKE || t == pb.NodeType_NODE_TYPE_AUTO || t == pb.NodeType_NODE_TYPE_SEDAN || t == pb.NodeType_NODE_TYPE_SUV
}
func (c *TrafficConsumer) Wait(ctx context.Context) error {
	select {
	case <-c.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
