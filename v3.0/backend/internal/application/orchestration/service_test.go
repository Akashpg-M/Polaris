package orchestration

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
)

func TestDeterministicHelpers(t *testing.T) {
	if d := haversineMeters(13.0067, 80.2206, 13.0067, 80.2206); math.Abs(d) > 0.001 {
		t.Fatalf("same point distance=%f", d)
	}
	if _, _, ok := targetCoordinates(json.RawMessage(`{"lat":91,"lon":80}`)); ok {
		t.Fatal("invalid latitude accepted")
	}
	if !validPriority("HIGH") || validPriority("URGENT") {
		t.Fatal("priority validation incorrect")
	}
}

func TestActiveConnectionStateRequiresLiveLease(t *testing.T) {
	if activeConnectionState(map[string]string{"gateway_id": "gateway-1", "lease_expires_at": "1"}) {
		t.Fatal("expired connection accepted")
	}
	future := time.Now().Add(time.Minute).UnixMilli()
	if !activeConnectionState(map[string]string{"gateway_id": "gateway-1", "lease_expires_at": fmt.Sprint(future)}) {
		t.Fatal("live connection rejected")
	}
}

func TestDomainProviderCannotDropCoreEligibleCandidates(t *testing.T) {
	score := 1.0
	got := includeUnrankedEligible(
		[]extension.Candidate{{DeviceID: "ranked-camera", DomainScore: &score}},
		[]repository.DeviceCandidate{{DeviceID: "vehicle"}, {DeviceID: "compute"}},
	)
	if len(got) != 3 || got[0].DeviceID != "ranked-camera" || got[1].DeviceID != "vehicle" || got[2].DeviceID != "compute" {
		t.Fatalf("eligible fallback candidates lost: %#v", got)
	}
	got = includeUnrankedEligible(got, []repository.DeviceCandidate{{DeviceID: "vehicle"}})
	if len(got) != 3 {
		t.Fatalf("duplicate fallback candidate added: %#v", got)
	}
}
