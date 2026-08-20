package routing

import (
	"container/heap"
	"context"
	"math"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type queueItem struct {
	node           int64
	priority, cost float64
	index          int
}
type priorityQueue []*queueItem

func (q priorityQueue) Len() int { return len(q) }
func (q priorityQueue) Less(i, j int) bool {
	if q[i].priority != q[j].priority {
		return q[i].priority < q[j].priority
	}
	return q[i].node < q[j].node
}
func (q priorityQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i]; q[i].index = i; q[j].index = j }
func (q *priorityQueue) Push(v any)   { n := v.(*queueItem); n.index = len(*q); *q = append(*q, n) }
func (q *priorityQueue) Pop() any     { old := *q; n := old[len(old)-1]; *q = old[:len(old)-1]; return n }

type searchResult struct {
	edgeIndexes             []int
	cost, distance, seconds float64
	expanded                int
	heapOperations          int
}

type SearchAlgorithm string

const (
	AlgorithmAStar    SearchAlgorithm = "ASTAR"
	AlgorithmDijkstra SearchAlgorithm = "DIJKSTRA"
)

type SearchMetrics struct {
	Algorithm      SearchAlgorithm `json:"algorithm"`
	Cost           float64         `json:"cost"`
	DistanceMeters float64         `json:"distance_meters"`
	ExpandedNodes  int             `json:"expanded_nodes"`
	HeapOperations int             `json:"heap_operations"`
	EdgeCount      int             `json:"edge_count"`
}

func edgeCost(e RoadEdge, index int, policy RouteCostPolicy, s *RoutingCostSnapshot) float64 {
	if policy == RouteShortest {
		return e.DistanceM
	}
	if s != nil && index < len(s.EdgeCosts) && s.EdgeCosts[index] > 0 {
		return s.EdgeCosts[index]
	}
	return e.BaseTravelTime.Seconds()
}
func heuristic(g *RoadGraph, node, target int64, policy RouteCostPolicy) float64 {
	a := g.nodes[node].Position
	b := g.nodes[target].Position
	d := spatial.DistanceMeters(a, b)
	if policy == RouteFastest {
		return d / g.maxSpeedMPS
	}
	return d
}
func runSearch(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, policy RouteCostPolicy, maxExpansions int, useHeuristic bool) (searchResult, error) {
	if _, ok := g.nodes[from]; !ok {
		return searchResult{}, ErrNoRoadNode
	}
	if _, ok := g.nodes[to]; !ok {
		return searchResult{}, ErrNoRoadNode
	}
	if from == to {
		return searchResult{}, nil
	}
	cost := map[int64]float64{from: 0}
	previous := map[int64]int{}
	visited := map[int64]bool{}
	q := &priorityQueue{}
	heap.Init(q)
	heap.Push(q, &queueItem{node: from})
	expanded, heapOperations := 0, 1
	for q.Len() > 0 {
		if ctx.Err() != nil {
			return searchResult{}, ctx.Err()
		}
		cur := heap.Pop(q).(*queueItem)
		heapOperations++
		if visited[cur.node] {
			continue
		}
		visited[cur.node] = true
		expanded++
		if expanded > maxExpansions {
			return searchResult{}, ErrTimeout
		}
		if cur.node == to {
			break
		}
		for _, ei := range g.adjacency[cur.node] {
			e := g.edges[ei]
			next := cost[cur.node] + edgeCost(e, ei, policy, s)
			old, ok := cost[e.ToID]
			if !ok || next < old {
				cost[e.ToID] = next
				previous[e.ToID] = ei
				h := 0.0
				if useHeuristic {
					h = heuristic(g, e.ToID, to, policy)
				}
				heap.Push(q, &queueItem{node: e.ToID, cost: next, priority: next + h})
				heapOperations++
			}
		}
	}
	if _, ok := cost[to]; !ok {
		return searchResult{}, ErrNoRoute
	}
	rev := []int{}
	at := to
	for at != from {
		ei, ok := previous[at]
		if !ok {
			return searchResult{}, ErrNoRoute
		}
		rev = append(rev, ei)
		at = g.edges[ei].FromID
	}
	edges := make([]int, len(rev))
	var distance, seconds float64
	for i := range rev {
		ei := rev[len(rev)-1-i]
		edges[i] = ei
		e := g.edges[ei]
		distance += e.DistanceM
		if s != nil && ei < len(s.EdgeCosts) {
			seconds += s.EdgeCosts[ei]
		} else {
			seconds += e.BaseTravelTime.Seconds()
		}
	}
	return searchResult{edgeIndexes: edges, cost: cost[to], distance: distance, seconds: seconds, expanded: expanded, heapOperations: heapOperations}, nil
}

func AStar(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, p RouteCostPolicy, max int) (searchResult, error) {
	return runSearch(ctx, g, s, from, to, p, max, true)
}
func Dijkstra(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, p RouteCostPolicy, max int) (searchResult, error) {
	return runSearch(ctx, g, s, from, to, p, max, false)
}

// MeasureSearch exposes algorithm work without exposing mutable graph details.
// It is used by the reproducible real-road benchmark and correctness oracle.
func MeasureSearch(ctx context.Context, g *RoadGraph, s *RoutingCostSnapshot, from, to int64, p RouteCostPolicy, max int, algorithm SearchAlgorithm) (SearchMetrics, error) {
	result, err := runSearch(ctx, g, s, from, to, p, max, algorithm == AlgorithmAStar)
	if err != nil {
		return SearchMetrics{}, err
	}
	return SearchMetrics{Algorithm: algorithm, Cost: result.cost, DistanceMeters: result.distance, ExpandedNodes: result.expanded, HeapOperations: result.heapOperations, EdgeCount: len(result.edgeIndexes)}, nil
}

func routeResult(g *RoadGraph, s *RoutingCostSnapshot, p RouteCostPolicy, r searchResult) RouteResult {
	waypoints := []model.Position{}
	ids := make([]EdgeID, len(r.edgeIndexes))
	if len(r.edgeIndexes) > 0 {
		waypoints = append(waypoints, g.nodes[g.edges[r.edgeIndexes[0]].FromID].Position)
	}
	for i, ei := range r.edgeIndexes {
		e := g.edges[ei]
		ids[i] = e.ID
		waypoints = append(waypoints, g.nodes[e.ToID].Position)
	}
	version := uint64(0)
	if s != nil {
		version = s.Version
	}
	return RouteResult{GraphVersion: g.Version(), SnapshotVersion: version, Policy: p, DistanceMeters: r.distance, EstimatedTime: timeDuration(r.seconds), Waypoints: waypoints, EdgeIDs: ids, ExpandedNodes: r.expanded}
}
func timeDuration(seconds float64) time.Duration {
	return time.Duration(math.Round(seconds * float64(time.Second)))
}
