package routing

import (
	"errors"
	"sync/atomic"
	"time"
)

type RoutingCostSnapshot struct {
	Version     uint64    `json:"version"`
	GeneratedAt time.Time `json:"generated_at"`
	EdgeCosts   []float64 `json:"edge_costs"`
}
type SnapshotStore struct {
	current atomic.Pointer[RoutingCostSnapshot]
}

func NewSnapshotStore(g *RoadGraph) *SnapshotStore {
	s := &SnapshotStore{}
	cost := make([]float64, len(g.edges))
	for i, e := range g.edges {
		cost[i] = e.BaseTravelTime.Seconds()
	}
	s.current.Store(&RoutingCostSnapshot{Version: 1, GeneratedAt: time.Now().UTC(), EdgeCosts: cost})
	return s
}
func (s *SnapshotStore) Load() *RoutingCostSnapshot { return s.current.Load() }
func (s *SnapshotStore) Swap(next RoutingCostSnapshot) error {
	old := s.current.Load()
	if old != nil && next.Version <= old.Version {
		return errors.New("snapshot version must increase")
	}
	if old != nil && len(next.EdgeCosts) != len(old.EdgeCosts) {
		return errors.New("snapshot edge count mismatch")
	}
	copyCosts := append([]float64(nil), next.EdgeCosts...)
	next.EdgeCosts = copyCosts
	if next.GeneratedAt.IsZero() {
		next.GeneratedAt = time.Now().UTC()
	}
	s.current.Store(&next)
	return nil
}
