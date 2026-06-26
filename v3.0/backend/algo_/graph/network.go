package graph

import (
	"fmt"
	"sync"
	"math"
	"github.com/Akashpg-M/polaris/backend/algo_/geo"
)

// Intersection represents a physical point on the map where roads meet
type Intersection struct {
	ID  int64 // OSM uses int64 for node IDs
	Lat float64
	Lon float64
}

// RoadSegment is a directional edge connecting two Intersections
type RoadSegment struct {
	WayID          int64   // The OSM Way ID (the street itself)
	TargetNodeID   int64   // Where this segment leads
	BaseDistanceKm float64 // Physical static length
	
	// Dynamic Weighting Attributes
	CurrentSpeedMPS float64 // Real-time velocity reported by Kafka telemetry
	CongestionMult  float64 // 1.0 = clear, 5.0 = heavy traffic
}

// RoadNetwork is the thread-safe adjacency list holding the city's topology
type RoadNetwork struct {
	mu    sync.RWMutex
	nodes map[int64]*Intersection
	edges map[int64][]*RoadSegment
}

func NewRoadNetwork() *RoadNetwork {
	return &RoadNetwork{
		nodes: make(map[int64]*Intersection),
		edges: make(map[int64][]*RoadSegment),
	}
}

// AddIntersection registers a node in the graph
func (rn *RoadNetwork) AddIntersection(id int64, lat, lon float64) {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	rn.nodes[id] = &Intersection{ID: id, Lat: lat, Lon: lon}
}

// AddRoadSegment creates a directional link between two nodes
func (rn *RoadNetwork) AddRoadSegment(sourceID, targetID, wayID int64, distKm float64, isOneWay bool) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Add forward direction
	rn.edges[sourceID] = append(rn.edges[sourceID], &RoadSegment{
		WayID:          wayID,
		TargetNodeID:   targetID,
		BaseDistanceKm: distKm,
		CongestionMult: 1.0, // Default to clear roads
	})

	// Add reverse direction if it's a two-way street
	if !isOneWay {
		rn.edges[targetID] = append(rn.edges[targetID], &RoadSegment{
			WayID:          wayID,
			TargetNodeID:   sourceID,
			BaseDistanceKm: distKm,
			CongestionMult: 1.0,
		})
	}
}

// UpdateSegmentCongestion allows the Kafka stream to dynamically alter route weights
func (rn *RoadNetwork) UpdateSegmentCongestion(sourceID int64, multiplier float64) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	segments, exists := rn.edges[sourceID]
	if !exists {
		return fmt.Errorf("node not found in graph")
	}

	for _, seg := range segments {
		seg.CongestionMult = multiplier
	}
	return nil
}

// GetNodesCount is a helper for metrics
func (rn *RoadNetwork) GetNodesCount() int {
	rn.mu.RLock()
	defer rn.mu.RUnlock()
	return len(rn.nodes)
}

func (rn *RoadNetwork) GetNearestIntersection(lat, lon float64) (int64, error) {
	rn.mu.RLock()
	defer rn.mu.RUnlock()

	if len(rn.nodes) == 0 {
		return 0, fmt.Errorf("road network is empty")
	}

	var bestID int64
	minDist := math.MaxFloat64

	// Note: In an ultra-optimized state, we would index these nodes into your QuadTree.
	// For now, a brute-force memory scan is perfectly fast enough for typical city sizes.
	for id, node := range rn.nodes {
		dist := geo.Haversine(lat, lon, node.Lat, node.Lon)
		if dist < minDist {
			minDist = dist
			bestID = id
		}
	}

	return bestID, nil
}

func (rn *RoadNetwork) GetNodeCoordinates(id int64) (float64, float64, bool) {
	rn.mu.RLock()
	defer rn.mu.RUnlock()
	
	if node, exists := rn.nodes[id]; exists {
		return node.Lat, node.Lon, true
	}
	return 0, 0, false
}

