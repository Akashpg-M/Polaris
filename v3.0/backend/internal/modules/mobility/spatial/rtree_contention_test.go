package spatial

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

func populatedRTree(size int) (*RTreeSpatialIndex, []model.SpatialState) {
	rng := rand.New(rand.NewSource(41))
	index := NewRTreeSpatialIndex()
	states := make([]model.SpatialState, size)
	for n := range states {
		states[n] = model.SpatialState{TenantID: "tenant", DeviceID: fmt.Sprintf("device-%06d", n), Position: model.Position{Latitude: 12.9 + rng.Float64()*.25, Longitude: 80.1 + rng.Float64()*.25}}
		_ = index.Upsert(states[n])
	}
	return index, states
}

func TestRTreeDirtyReadRebuildIsBoundedAndExact(t *testing.T) {
	index, states := populatedRTree(5000)
	started := time.Now()
	result, err := index.WithinRadius(states[0].Position, 10)
	if err != nil || len(result) == 0 || result[0].State.DeviceID != states[0].DeviceID {
		t.Fatalf("dirty rebuild query was not exact: count=%d err=%v", len(result), err)
	}
	stats := index.Stats()
	if stats.Rebuilds != 1 || stats.Mutations != 5000 {
		t.Fatalf("unexpected rebuild accounting: %#v", stats)
	}
	if elapsed := time.Since(started); elapsed > 500*time.Millisecond {
		t.Fatalf("5k dirty-tree rebuild exceeded local safety budget: %v", elapsed)
	}
	_, _ = index.WithinRadius(states[1].Position, 10)
	if index.Stats().Rebuilds != 1 {
		t.Fatal("clean read rebuilt an unchanged tree")
	}
}

func BenchmarkRTreeMovementQueryMix(b *testing.B) {
	for _, size := range []int{1000, 5000, 10000} {
		for _, writes := range []int{80, 50, 20} {
			name := fmt.Sprintf("devices=%d/writes=%d", size, writes)
			b.Run(name, func(b *testing.B) {
				index, states := populatedRTree(size)
				_, _ = index.WithinRadius(states[0].Position, 5000)
				b.ResetTimer()
				for n := 0; n < b.N; n++ {
					state := states[n%len(states)]
					if n%100 < writes {
						state.Position.Latitude += float64(n%7) * .000001
						_ = index.Upsert(state)
					} else {
						_, _ = index.Nearest(state.Position, 10, 5000)
					}
				}
			})
		}
	}
}
