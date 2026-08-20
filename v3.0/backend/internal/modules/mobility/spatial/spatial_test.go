package spatial

import (
	"errors"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

func state(tenant, id string, lat, lon float64, seq uint64) model.SpatialState {
	now := time.Unix(1_700_000_000+int64(seq), 0).UTC()
	return model.SpatialState{TenantID: tenant, DeviceID: id, Position: model.Position{Latitude: lat, Longitude: lon}, ReportedPosition: model.Position{Latitude: lat, Longitude: lon}, MobilityProfile: model.MobilityRoadVehicle, ObservedAt: now, BootID: "boot-1", BootStartedAt: time.Unix(1_699_999_000, 0).UTC(), SequenceNumber: seq}
}
func ids(v []SpatialCandidate) []string {
	out := make([]string, len(v))
	for i, c := range v {
		out[i] = c.State.DeviceID
	}
	return out
}
func TestGeodesicEdgeCases(t *testing.T) {
	cases := []struct {
		a, b     model.Position
		min, max float64
	}{{model.Position{Latitude: 0, Longitude: 179.9}, model.Position{Latitude: 0, Longitude: -179.9}, 20_000, 25_000}, {model.Position{Latitude: 89, Longitude: 0}, model.Position{Latitude: 89, Longitude: 90}, 150_000, 160_000}, {model.Position{Latitude: 13, Longitude: 80}, model.Position{Latitude: 13, Longitude: 80}, 0, .001}}
	for _, tc := range cases {
		got := DistanceMeters(tc.a, tc.b)
		if got < tc.min || got > tc.max {
			t.Fatalf("distance %f outside [%f,%f]", got, tc.min, tc.max)
		}
	}
}
func TestRTreeMatchesLinearOracleRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	tree := NewRTreeSpatialIndex()
	linear := NewLinearSpatialIndex()
	for n := 0; n < 1000; n++ {
		s := state("tenant-a", string(rune(1000+n)), rng.Float64()*170-85, rng.Float64()*360-180, 1)
		if err := tree.Upsert(s); err != nil {
			t.Fatal(err)
		}
		_ = linear.Upsert(s)
	}
	queries := []model.Position{{Latitude: 13, Longitude: 80}, {Latitude: 0, Longitude: 179.95}, {Latitude: 75, Longitude: -40}}
	for _, q := range queries {
		a, _ := tree.Nearest(q, 10, 0)
		b, _ := linear.Nearest(q, 10, 0)
		if !reflect.DeepEqual(ids(a), ids(b)) {
			t.Fatalf("nearest mismatch %v != %v", ids(a), ids(b))
		}
	}
}
func TestRTreeMutationAndAntimeridian(t *testing.T) {
	tree := NewRTreeSpatialIndex()
	_ = tree.Upsert(state("t", "west", 0, 179.95, 1))
	_ = tree.Upsert(state("t", "east", 0, -179.95, 1))
	got, err := tree.WithinRadius(model.Position{Latitude: 0, Longitude: 180}, 20_000)
	if err != nil || len(got) != 2 {
		t.Fatalf("antimeridian search got=%v err=%v", ids(got), err)
	}
	moved := state("t", "west", 20, 20, 2)
	_ = tree.Upsert(moved)
	_ = tree.Remove("t", "east")
	got, _ = tree.WithinRadius(model.Position{Latitude: 0, Longitude: 180}, 20_000)
	if len(got) != 0 {
		t.Fatalf("removed/moved device returned: %v", ids(got))
	}
}
func TestManagerVersionTenantAndReplayInvariants(t *testing.T) {
	cfg := ManagerConfig{H3Resolution: 8, ShardResolution: 6, MinMoveMeters: 5, MaxIndexAge: time.Minute, MaxH3Rings: 12, MaxRadiusMeters: 10_000, MaxDevicesPerTenant: 100}
	m := NewManager(cfg)
	s := state("a", "d1", 13.0067, 80.2206, 2)
	if err := m.Upsert(s); err != nil {
		t.Fatal(err)
	}
	old := state("a", "d1", 14, 81, 1)
	if err := m.Upsert(old); !errors.Is(err, ErrStaleVersion) {
		t.Fatalf("expected stale version, got %v", err)
	}
	_ = m.Upsert(state("b", "d2", 13.0067, 80.2206, 1))
	got, _ := m.Nearby("a", model.Position{Latitude: 13.0067, Longitude: 80.2206}, 1000, 10)
	if !reflect.DeepEqual(ids(got), []string{"d1"}) {
		t.Fatalf("tenant leak: %v", ids(got))
	}
	first := m.Snapshot()
	replay := NewManager(cfg)
	shuffled := append([]model.SpatialState(nil), first...)
	sort.Slice(shuffled, func(i, j int) bool { return shuffled[i].DeviceID > shuffled[j].DeviceID })
	if err := replay.Rebuild(shuffled); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, replay.Snapshot()) {
		t.Fatal("replay was not deterministic")
	}
	_ = m.EvictInactive("a", "d1", "ACTIVE", "OFFLINE")
	got, _ = m.Nearby("a", model.Position{Latitude: 13.0067, Longitude: 80.2206}, 1000, 10)
	if len(got) != 0 {
		t.Fatal("offline device remained indexed")
	}
}
func TestMovementThresholdPreservesReportedState(t *testing.T) {
	m := NewManager(ManagerConfig{H3Resolution: 8, ShardResolution: 6, MinMoveMeters: 100, MaxIndexAge: time.Hour, MaxH3Rings: 12, MaxRadiusMeters: 10_000, MaxDevicesPerTenant: 10})
	a := state("t", "d", 13, 80, 1)
	_ = m.Upsert(a)
	b := state("t", "d", 13.00001, 80.00001, 2)
	_ = m.Upsert(b)
	got, _ := m.Get("t", "d")
	if got.ReportedPosition == got.Position {
		t.Fatal("reported and indexed positions should be independent below threshold")
	}
	if got.SequenceNumber != 2 {
		t.Fatal("latest source version not retained")
	}
}

func BenchmarkRTreeNearest(b *testing.B)  { benchmarkNearest(b, NewRTreeSpatialIndex()) }
func BenchmarkLinearNearest(b *testing.B) { benchmarkNearest(b, NewLinearSpatialIndex()) }
func benchmarkNearest(b *testing.B, index SpatialIndex) {
	for i := 0; i < 5000; i++ {
		_ = index.Upsert(state("bench", string(rune(i+1)), 13+float64(i%100)*.001, 80+float64(i/100)*.001, 1))
	}
	q := model.Position{Latitude: 13.04, Longitude: 80.02}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = index.Nearest(q, 10, 5000)
	}
}
