package spatial

import (
	"hash/fnv"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/geo"
	"github.com/Akashpg-M/polaris/backend/algo_/quadtree"
	"github.com/Akashpg-M/polaris/backend/internal/core/domain"
)

// ShardCount dictates how many independent memory partitions exist.
// 32 is a standard default to minimize lock contention across highly concurrent goroutines.
const ShardCount = 32

type EngineShard struct {
	mu    sync.RWMutex
	nodes map[string]*domain.TelemetryPayload
}

type Engine struct {
	shards []*EngineShard
	qt     *quadtree.SafeQuadTree
}

// MatchResult is the DTO sent back to the dispatcher
type MatchResult struct {
	NodeID     string  `json:"node_id"`
	Class      uint16  `json:"asset_class"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	DistanceKm float64 `json:"distance_km"`
	ETASec     int     `json:"eta_seconds"`
	RouteType  string  `json:"route_type"`
}

func NewEngine() *Engine {
	bounds := quadtree.Bounds{X: 12.0, Y: 79.5, Width: 2.0, Height: 1.5}
	shards := make([]*EngineShard, ShardCount)
	for i := 0; i < ShardCount; i++ {
		shards[i] = &EngineShard{nodes: make(map[string]*pb.SpatialObject)}
	}
	return &Engine{shards: shards, qt: quadtree.NewSafeQuadTree(bounds)}
}

// getShard picks the correct memory partition using an FNV-1a hash of the NodeID
func (e *Engine) getShard(nodeID string) *EngineShard {
	h := fnv.New32a()
	h.Write([]byte(nodeID))
	return e.shards[h.Sum32()%ShardCount]
}


func (e *Engine) BatchUpdate(payloads []*pb.SpatialObject) {
	if len(payloads) == 0 { return }

	for _, p := range payloads {
		shard := e.getShard(p.Id)

		shard.mu.Lock()
		if _, exists := shard.nodes[p.Id]; exists {
			e.qt.Remove(p.Id)
		}
		shard.nodes[p.Id] = p
		shard.mu.Unlock()

		e.qt.Insert(quadtree.Point{
			Lat:      p.Lat,
			Lon:      p.Lon,
			ID:       p.Id,
			Class:    uint16(p.Type), // Map the Proto Enum to the QuadTree Class
			TenantID: p.TenantId,
		})
	}
}

func (e *Engine) FindNearest(tenantID string, lat, lon, radiusKm float64, reqType pb.NodeType) []MatchResult {
	x, y, w, h := geo.BoundingBox(lat, lon, radiusKm)
	searchBounds := quadtree.Bounds{X: x, Y: y, Width: w, Height: h}

	candidates := e.qt.Search(searchBounds, uint16(reqType))
	var results []MatchResult

	for _, c := range candidates {
		if c.TenantID != tenantID { continue }

		dist := geo.Haversine(lat, lon, c.Lat, c.Lon)
		if dist <= radiusKm {
			
			eta := int((dist / 40.0) * 3600) // simplified ETA

			results = append(results, MatchResult{
				NodeID:     c.ID,
				Type:       pb.NodeType(c.Class),
				Lat:        c.Lat,
				Lon:        c.Lon,
				DistanceKm: dist,
				ETASec:     eta,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].ETASec < results[j].ETASec })
	if len(results) > 500 { return results[:500] }
	return results
}