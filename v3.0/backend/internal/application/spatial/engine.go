package spatial

import (
	"hash/fnv"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/algo_/geo"
	"github.com/Akashpg-M/polaris/algo_/quadtree"
	"github.com/Akashpg-M/polaris/internal/core/domain"
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
	bounds := quadtree.Bounds{
		X:      12.0, // Min Latitude
		Y:      79.5, // Min Longitude
		Width:  2.0,
		Height: 1.5,
	}

	// Initialize the 32 independent shards
	shards := make([]*EngineShard, ShardCount)
	for i := 0; i < ShardCount; i++ {
		shards[i] = &EngineShard{
			nodes: make(map[string]*domain.TelemetryPayload),
		}
	}

	return &Engine{
		shards: shards,
		qt:     quadtree.NewSafeQuadTree(bounds),
	}
}

// getShard picks the correct memory partition using an FNV-1a hash of the NodeID
func (e *Engine) getShard(nodeID string) *EngineShard {
	h := fnv.New32a()
	h.Write([]byte(nodeID))
	return e.shards[h.Sum32()%ShardCount]
}

// BatchUpdate processes thousands of pings concurrently without a global lock
func (e *Engine) BatchUpdate(payloads []domain.TelemetryPayload) {
	if len(payloads) == 0 {
		return
	}

	start := time.Now()

	for _, payload := range payloads {
		p := payload
		shard := e.getShard(p.NodeID)

		// 1. Lock ONLY the specific shard for this specific node
		shard.mu.Lock()
		if _, exists := shard.nodes[p.NodeID]; exists {
			// Wipe old tree index entry to prevent ghosts
			e.qt.Remove(p.NodeID)
		}
		
		// 2. Update RAM Map
		shard.nodes[p.NodeID] = &p
		shard.mu.Unlock() // Release the lock immediately

		// 3. Insert into the thread-safe QuadTree
		e.qt.Insert(quadtree.Point{
			Lat:      p.Lat,
			Lon:      p.Lon,
			ID:       p.NodeID,
			Class:    uint16(p.Class),
			TenantID: p.TenantID,
		})
	}

	slog.Debug("Concurrent sharded batch update complete", "processed", len(payloads), "duration_ms", time.Since(start).Milliseconds())
}

// FindNearest queries the QuadTree, filters by exact distance, and applies context-aware routing
func (e *Engine) FindNearest(tenantID string, lat, lon, radiusKm float64, reqClass uint16) []MatchResult {
	x, y, w, h := geo.BoundingBox(lat, lon, radiusKm)
	searchBounds := quadtree.Bounds{X: x, Y: y, Width: w, Height: h}

	candidates := e.qt.Search(searchBounds, reqClass)
	var results []MatchResult

	for _, c := range candidates {
		if c.TenantID != tenantID {
			continue
		}

		dist := geo.Haversine(lat, lon, c.Lat, c.Lon)
		if dist <= radiusKm {
			var eta int
			var routeType string

			if (c.Class & uint16(domain.ClassDrone)) > 0 {
				eta = int((dist / 60.0) * 3600)
				routeType = "euclidean_air"
			} else {
				streetDist := dist * 1.4 // Rough street network multiplier
				eta = int((streetDist / 40.0) * 3600)
				routeType = "osrm_street"
			}

			results = append(results, MatchResult{
				NodeID:     c.ID,
				Class:      c.Class,
				Lat:        c.Lat,
				Lon:        c.Lon,
				DistanceKm: dist,
				ETASec:     eta,
				RouteType:  routeType,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].ETASec < results[j].ETASec
	})

	if len(results) > 500 {
		return results[:500]
	}
	return results
}