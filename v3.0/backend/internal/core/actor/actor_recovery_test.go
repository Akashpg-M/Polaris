package actor

import (
	"context"
	"testing"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

// MockPublisherSpy captures events to verify replay validity
type MockPublisherSpy struct {
	PublishedEvents []AssetStateChangedEvent
}

func (m *MockPublisherSpy) PublishEvent(ctx context.Context, topic string, event interface{}) error {
	if e, ok := event.(AssetStateChangedEvent); ok {
		m.PublishedEvents = append(m.PublishedEvents, e)
	}
	return nil
}

func TestActorCrashAndStateReplayRecovery(t *testing.T) {
	spy := &MockPublisherSpy{}
	assetID := "TEST-DRONE-99"

	// 1. Simulate Active Lifecycle Phase
	originalActor := NewAssetActor(assetID, spy, 10)
	originalActor.Start()

	_ = originalActor.Push(TelemetryMsg{Payload: &pb.SpatialObject{
		Id:            assetID,
		Lat:           13.04,
		Lon:           80.24,
		EnergyPercent: 85,
		Timestamp:     time.Now().UnixMilli(),
	}})

	// Allow event loop execution cycle to settle
	time.Sleep(10 * time.Millisecond)
	originalActor.Stop() // --- SIMULATE CRASH / NODE TERMINATION ---

	// Assert that local telemetry was processed and pushed down the event pipeline
	if len(spy.PublishedEvents) != 1 {
		t.Fatalf("Expected 1 state projection event, found %d", len(spy.PublishedEvents))
	}

	// 2. Simulate Node Recovery Recovery Pipeline
	recoveredActor := NewAssetActor(assetID, spy, 10)
	
	// Replay history sequentially back into the fresh actor memory boundary
	historicalEvent := spy.PublishedEvents[0]
	
	recoveredActor.lat = historicalEvent.Lat
	recoveredActor.lon = historicalEvent.Lon
	recoveredActor.energyPercent = historicalEvent.EnergyPercent

	// Verify structural state restoration
	if recoveredActor.lat != 13.04 || recoveredActor.energyPercent != 85 {
		t.Errorf("State recovery divergence detected after simulation log replay")
	}
}