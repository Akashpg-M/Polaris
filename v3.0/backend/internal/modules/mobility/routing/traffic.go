package routing

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type EdgeTrafficState struct {
	SampleCount    uint64    `json:"sample_count"`
	EWMASpeedMPS   float64   `json:"ewma_speed_mps"`
	LastSpeedMPS   float64   `json:"last_speed_mps"`
	LastObservedAt time.Time `json:"last_observed_at"`
	Confidence     float64   `json:"confidence"`
}

type TrafficObservation struct {
	Position       model.Position
	HeadingDegrees *float64
	SpeedMPS       float64
	ObservedAt     time.Time
}

type TrafficManager struct {
	mu        sync.RWMutex
	graph     *RoadGraph
	snapshots *SnapshotStore
	states    map[int]EdgeTrafficState
	maxAge    time.Duration
}

func NewTrafficManager(g *RoadGraph, s *SnapshotStore, maxAge time.Duration) *TrafficManager {
	return &TrafficManager{graph: g, snapshots: s, states: map[int]EdgeTrafficState{}, maxAge: maxAge}
}
func (t *TrafficManager) Observe(ctx context.Context, o TrafficObservation) error {
	if o.SpeedMPS < 0 || time.Since(o.ObservedAt) > t.maxAge {
		return errors.New("traffic observation rejected")
	}
	ei, confidence, err := t.match(ctx, o)
	if err != nil {
		return err
	}
	t.mu.Lock()
	s := t.states[ei]
	alpha := .3
	if s.SampleCount == 0 {
		s.EWMASpeedMPS = o.SpeedMPS
	} else {
		s.EWMASpeedMPS = alpha*o.SpeedMPS + (1-alpha)*s.EWMASpeedMPS
	}
	s.SampleCount++
	s.LastSpeedMPS = o.SpeedMPS
	s.LastObservedAt = o.ObservedAt
	s.Confidence = confidence
	t.states[ei] = s
	t.mu.Unlock()
	return nil
}
func (t *TrafficManager) match(ctx context.Context, o TrafficObservation) (int, float64, error) {
	node, err := t.graph.nodeIndex.Nearest(ctx, o.Position)
	if err != nil {
		return 0, 0, err
	}
	best, bestScore := -1, math.Inf(1)
	for _, ei := range t.graph.incident[node.ID] {
		e := t.graph.edges[ei]
		a, b := t.graph.nodes[e.FromID].Position, t.graph.nodes[e.ToID].Position
		distance := segmentDistance(o.Position, a, b)
		score := distance
		if o.HeadingDegrees != nil {
			bearing := initialBearing(a, b)
			delta := math.Abs(*o.HeadingDegrees - bearing)
			if delta > 180 {
				delta = 360 - delta
			}
			score += delta * 2
		}
		if score < bestScore {
			best, bestScore = ei, score
		}
	}
	if best < 0 || bestScore > 100 {
		return 0, 0, ErrNoRoadNode
	}
	return best, math.Max(0, 1-bestScore/100), nil
}
func segmentDistance(p, a, b model.Position) float64 {
	latScale := 111320.0
	lonScale := latScale * math.Cos(p.Latitude*math.Pi/180)
	ax, ay := (a.Longitude-p.Longitude)*lonScale, (a.Latitude-p.Latitude)*latScale
	bx, by := (b.Longitude-p.Longitude)*lonScale, (b.Latitude-p.Latitude)*latScale
	dx, dy := bx-ax, by-ay
	if dx == 0 && dy == 0 {
		return math.Hypot(ax, ay)
	}
	u := -(ax*dx + ay*dy) / (dx*dx + dy*dy)
	u = math.Max(0, math.Min(1, u))
	return math.Hypot(ax+u*dx, ay+u*dy)
}
func initialBearing(a, b model.Position) float64 {
	lat1, lat2 := a.Latitude*math.Pi/180, b.Latitude*math.Pi/180
	dlon := (b.Longitude - a.Longitude) * math.Pi / 180
	y := math.Sin(dlon) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(dlon)
	v := math.Atan2(y, x) * 180 / math.Pi
	if v < 0 {
		v += 360
	}
	return v
}
func (t *TrafficManager) Refresh(now time.Time) error {
	old := t.snapshots.Load()
	if old == nil {
		return ErrUnavailable
	}
	costs := make([]float64, len(t.graph.edges))
	t.mu.RLock()
	states := make(map[int]EdgeTrafficState, len(t.states))
	for edge, state := range t.states {
		states[edge] = state
	}
	t.mu.RUnlock()
	for i, e := range t.graph.edges {
		base := e.BaseTravelTime.Seconds()
		state, ok := states[i]
		if !ok {
			costs[i] = base
			continue
		}
		age := now.Sub(state.LastObservedAt)
		confidence := state.Confidence * math.Exp(-float64(age)/float64(t.maxAge))
		baseSpeed := e.DistanceM / base
		observed := math.Max(.5, state.EWMASpeedMPS)
		multiplier := math.Max(1, baseSpeed/observed)
		costs[i] = base * (1 + confidence*(multiplier-1))
	}
	// Expired observations no longer affect costs and must not accumulate for
	// the lifetime of the process. Compare timestamps before deletion so a
	// concurrent fresh observation for the same edge is never removed.
	t.mu.Lock()
	for edge, state := range states {
		if now.Sub(state.LastObservedAt) > 5*t.maxAge && t.states[edge].LastObservedAt.Equal(state.LastObservedAt) {
			delete(t.states, edge)
		}
	}
	t.mu.Unlock()
	return t.snapshots.Swap(RoutingCostSnapshot{Version: old.Version + 1, GeneratedAt: now.UTC(), EdgeCosts: costs})
}

func (t *TrafficManager) StateCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.states)
}

func (t *TrafficManager) OverlayBytes() int64 {
	return int64(len(t.graph.edges)) * 8
}
