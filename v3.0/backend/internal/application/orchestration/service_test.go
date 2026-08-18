package orchestration

import (
	"encoding/json"
	"math"
	"testing"
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
