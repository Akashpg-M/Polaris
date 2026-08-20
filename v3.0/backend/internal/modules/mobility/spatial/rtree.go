package spatial

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

const rtreeFanout = 16

type box struct{ minX, minY, maxX, maxY float64 }

func pointBox(p model.Position) box { return box{p.Longitude, p.Latitude, p.Longitude, p.Latitude} }
func (a box) intersects(b box) bool {
	return a.minX <= b.maxX && a.maxX >= b.minX && a.minY <= b.maxY && a.maxY >= b.minY
}
func union(a, b box) box {
	return box{math.Min(a.minX, b.minX), math.Min(a.minY, b.minY), math.Max(a.maxX, b.maxX), math.Max(a.maxY, b.maxY)}
}

type rtreeEntry struct {
	state  model.SpatialState
	bounds box
}
type rtreeNode struct {
	bounds   box
	leaf     bool
	entries  []rtreeEntry
	children []*rtreeNode
}

// RTreeSpatialIndex is a packed STR R-tree. Mutations are O(1) and mark the
// hierarchy dirty; the next read atomically rebuilds it. This suits telemetry
// bursts and gives exact queries after geodesic post-filtering.
type RTreeSpatialIndex struct {
	mu               sync.RWMutex
	items            map[string]model.SpatialState
	root             *rtreeNode
	dirty            bool
	mutations        atomic.Uint64
	rebuilds         atomic.Uint64
	rebuildNanos     atomic.Uint64
	rebuildWaitNanos atomic.Uint64
}

type RTreeStats struct {
	Mutations        uint64        `json:"mutations"`
	Rebuilds         uint64        `json:"rebuilds"`
	TotalRebuildTime time.Duration `json:"total_rebuild_time"`
	TotalLockWait    time.Duration `json:"total_lock_wait"`
}

func NewRTreeSpatialIndex() *RTreeSpatialIndex {
	return &RTreeSpatialIndex{items: map[string]model.SpatialState{}, dirty: true}
}
func (i *RTreeSpatialIndex) Upsert(s model.SpatialState) error {
	if e := ValidatePosition(s.Position); e != nil {
		return e
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.items[indexKey(s.TenantID, s.DeviceID)] = s
	i.dirty = true
	i.mutations.Add(1)
	return nil
}
func (i *RTreeSpatialIndex) Remove(t, d string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.items, indexKey(t, d))
	i.dirty = true
	i.mutations.Add(1)
	return nil
}

func (i *RTreeSpatialIndex) ensure() {
	i.mu.RLock()
	dirty := i.dirty
	i.mu.RUnlock()
	if !dirty {
		return
	}
	waitStarted := time.Now()
	i.mu.Lock()
	i.rebuildWaitNanos.Add(uint64(time.Since(waitStarted)))
	defer i.mu.Unlock()
	if !i.dirty {
		return
	}
	rebuildStarted := time.Now()
	entries := make([]rtreeEntry, 0, len(i.items))
	for _, s := range i.items {
		entries = append(entries, rtreeEntry{s, pointBox(s.Position)})
	}
	i.root = buildSTR(entries)
	i.dirty = false
	i.rebuilds.Add(1)
	i.rebuildNanos.Add(uint64(time.Since(rebuildStarted)))
}

func (i *RTreeSpatialIndex) Stats() RTreeStats {
	return RTreeStats{Mutations: i.mutations.Load(), Rebuilds: i.rebuilds.Load(), TotalRebuildTime: time.Duration(i.rebuildNanos.Load()), TotalLockWait: time.Duration(i.rebuildWaitNanos.Load())}
}

func buildSTR(entries []rtreeEntry) *rtreeNode {
	if len(entries) == 0 {
		return nil
	}
	sort.Slice(entries, func(a, b int) bool { return entries[a].state.Position.Longitude < entries[b].state.Position.Longitude })
	leaves := make([]*rtreeNode, 0, (len(entries)+rtreeFanout-1)/rtreeFanout)
	for start := 0; start < len(entries); start += rtreeFanout {
		end := start + rtreeFanout
		if end > len(entries) {
			end = len(entries)
		}
		part := append([]rtreeEntry(nil), entries[start:end]...)
		sort.Slice(part, func(a, b int) bool { return part[a].state.Position.Latitude < part[b].state.Position.Latitude })
		n := &rtreeNode{leaf: true, entries: part, bounds: part[0].bounds}
		for _, e := range part[1:] {
			n.bounds = union(n.bounds, e.bounds)
		}
		leaves = append(leaves, n)
	}
	return buildParents(leaves)
}
func buildParents(nodes []*rtreeNode) *rtreeNode {
	if len(nodes) == 1 {
		return nodes[0]
	}
	sort.Slice(nodes, func(a, b int) bool { return nodes[a].bounds.minX < nodes[b].bounds.minX })
	next := make([]*rtreeNode, 0, (len(nodes)+rtreeFanout-1)/rtreeFanout)
	for s := 0; s < len(nodes); s += rtreeFanout {
		e := s + rtreeFanout
		if e > len(nodes) {
			e = len(nodes)
		}
		n := &rtreeNode{children: append([]*rtreeNode(nil), nodes[s:e]...), bounds: nodes[s].bounds}
		for _, c := range nodes[s+1 : e] {
			n.bounds = union(n.bounds, c.bounds)
		}
		next = append(next, n)
	}
	return buildParents(next)
}
func searchNode(n *rtreeNode, q box, out map[string]model.SpatialState) {
	if n == nil || !n.bounds.intersects(q) {
		return
	}
	if n.leaf {
		for _, e := range n.entries {
			if e.bounds.intersects(q) {
				out[indexKey(e.state.TenantID, e.state.DeviceID)] = e.state
			}
		}
		return
	}
	for _, c := range n.children {
		searchNode(c, q, out)
	}
}

func queryBoxes(p model.Position, radius float64) []box {
	if radius <= 0 || radius >= math.Pi*EarthRadiusMeters {
		return []box{{-180, -90, 180, 90}}
	}
	latDelta := radius / EarthRadiusMeters * 180 / math.Pi
	cos := math.Cos(p.Latitude * math.Pi / 180)
	lonDelta := 180.0
	if math.Abs(cos) > .000001 {
		lonDelta = math.Min(180, latDelta/math.Abs(cos))
	}
	minY, maxY := math.Max(-90, p.Latitude-latDelta), math.Min(90, p.Latitude+latDelta)
	minX, maxX := p.Longitude-lonDelta, p.Longitude+lonDelta
	if minX < -180 {
		return []box{{-180, minY, maxX, maxY}, {minX + 360, minY, 180, maxY}}
	}
	if maxX > 180 {
		return []box{{minX, minY, 180, maxY}, {-180, minY, maxX - 360, maxY}}
	}
	return []box{{minX, minY, maxX, maxY}}
}
func (i *RTreeSpatialIndex) WithinRadius(p model.Position, radius float64) ([]SpatialCandidate, error) {
	if e := ValidatePosition(p); e != nil {
		return nil, e
	}
	i.ensure()
	i.mu.RLock()
	defer i.mu.RUnlock()
	found := map[string]model.SpatialState{}
	for _, q := range queryBoxes(p, radius) {
		searchNode(i.root, q, found)
	}
	out := make([]SpatialCandidate, 0, len(found))
	for _, s := range found {
		d := DistanceMeters(p, s.Position)
		if radius <= 0 || d <= radius {
			out = append(out, SpatialCandidate{s, d})
		}
	}
	sortCandidates(out)
	return out, nil
}
func (i *RTreeSpatialIndex) Nearest(p model.Position, limit int, radius float64) ([]SpatialCandidate, error) {
	out, e := i.WithinRadius(p, radius)
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, e
}
