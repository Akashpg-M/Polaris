package actor

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

var ErrMailboxSaturated = errors.New("actor_mailbox_saturation_backpressure")

// EventPublisher abstracts our asynchronous append-only ledger (Kafka/NATS)
type EventPublisher interface {
	PublishEvent(ctx context.Context, topic string, event interface{}) error
}

// AssetActor owns isolated local state for a singular tracked entity
type AssetActor struct {
	id             string
	mailbox        chan interface{}
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
	eventPublisher EventPublisher

	// State Boundary (Zero Shared Access)
	lat            float64
	lon            float64
	energyPercent  int32
	currentTask    string
	lastPingMilli  int64
}

// NewAssetActor configures a state-isolated entity with strict backpressure limits
func NewAssetActor(id string, pub EventPublisher, capacity int) *AssetActor {
	ctx, cancel := context.WithCancel(context.Background())
	return &AssetActor{
		id:             id,
		mailbox:        make(chan interface{}, capacity),
		ctx:            ctx,
		cancel:         cancel,
		eventPublisher: pub,
	}
}

// Start mounts the actor's single-threaded event execution loop
func (a *AssetActor) Start() {
	go a.eventLoop()
}

// Stop safely tears down the execution routine
func (a *AssetActor) Stop() {
	a.cancel()
}

// Push Enqueues a message or flags immediate backpressure safety drops
func (a *AssetActor) Push(msg interface{}) error {
	select {
	case a.mailbox <- msg:
		return nil
	default:
		slog.Warn("[Actor System] Mailbox full, applying backpressure", "actor_id", a.id)
		return ErrMailboxSaturated
	}
}

func (a *AssetActor) eventLoop() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case msg := <-a.mailbox:
			a.process(msg)
		}
	}
}

func (a *AssetActor) process(msg interface{}) {
	switch m := msg.(type) {
	case TelemetryMsg:
		a.handleTelemetry(m.Payload)
	case CommandMsg:
		a.handleCommand(m)
	default:
		slog.Error("[Actor System] Dropped unknown structure frame payload", "actor_id", a.id)
	}
}

func (a *AssetActor) handleTelemetry(p *pb.SpatialObject) {
	a.mu.Lock()
	a.lat = p.Lat
	a.lon = p.Lon
	a.energyPercent = p.EnergyPercent
	a.lastPingMilli = p.Timestamp
	a.mu.Unlock()

	// Emit immutable mutation record down to the cluster event stream
	event := AssetStateChangedEvent{
		AssetID:        a.id,
		Lat:            p.Lat,
		Lon:            p.Lon,
		EnergyPercent:  p.EnergyPercent,
		VelocityMps:    p.VelocityMps,
		HeadingDeg:     p.HeadingDeg,
		TimestampMilli: time.Now().UnixMilli(),
	}

	if err := a.eventPublisher.PublishEvent(a.ctx, "spatial:updates", event); err != nil {
		slog.Error("[Actor System] Failed to project state change", "actor_id", a.id, "error", err)
	}
}

func (a *AssetActor) handleCommand(c CommandMsg) {
	a.mu.Lock()
	a.currentTask = c.Directive
	a.mu.Unlock()
	slog.Info("[Actor System] Local directive updated", "actor_id", a.id, "task", c.Directive)
}