package graph

import (
	
	"container/heap"
	"errors"
	"math"
)

// PathResult encapsulates the final calculated route
type PathResult struct {
	SourceID    int64
	TargetID    int64
	TotalCost   float64
	TotalDistKm float64
	RouteNodes  []int64 // Ordered list of intersection IDs to drive through
}

// FindShortestPath calculates the optimal route factoring in dynamic congestion
func (rn *RoadNetwork) FindShortestPath(sourceID, targetID int64) (*PathResult, error) {
	// RLock allows infinite concurrent routing reads, only blocking if Kafka is updating weights
	rn.mu.RLock()
	defer rn.mu.RUnlock()

	if _, ok := rn.nodes[sourceID]; !ok {
		return nil, errors.New("source node not found in graph")
	}
	if _, ok := rn.nodes[targetID]; !ok {
		return nil, errors.New("target node not found in graph")
	}

	// Tracking structures
	costs := make(map[int64]float64)       // Lowest known cost to reach a node
	distances := make(map[int64]float64)   // Actual physical distance in km
	previous := make(map[int64]int64)      // Breadcrumbs to reconstruct the path
	
	for id := range rn.nodes {
		costs[id] = math.Inf(1)
	}
	costs[sourceID] = 0
	distances[sourceID] = 0

	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &pqItem{nodeID: sourceID, cost: 0})

	for pq.Len() > 0 {
		// 1. Pop the most optimal pending node
		current := heap.Pop(pq).(*pqItem)

		// 2. Target reached. Break early (Dijkstra guarantees this is the shortest path)
		if current.nodeID == targetID {
			break
		}

		// 3. Skip if we've already found a better way to this node since queueing it
		if current.cost > costs[current.nodeID] {
			continue
		}

		// 4. Explore all outward roads from the current intersection
		for _, edge := range rn.edges[current.nodeID] {
			// The true magic: Cost is Physical Distance * Real-Time Traffic Congestion
			traversalCost := edge.BaseDistanceKm * edge.CongestionMult
			newCost := costs[current.nodeID] + traversalCost

			// If we found a faster route to the neighbor, update it and queue it
			if newCost < costs[edge.TargetNodeID] {
				costs[edge.TargetNodeID] = newCost
				distances[edge.TargetNodeID] = distances[current.nodeID] + edge.BaseDistanceKm
				previous[edge.TargetNodeID] = current.nodeID
				
				heap.Push(pq, &pqItem{
					nodeID: edge.TargetNodeID,
					cost:   newCost,
				})
			}
		}
	}

	// 5. Reconstruct the final path by walking backward from the target
	if math.IsInf(costs[targetID], 1) {
		return nil, errors.New("no valid path exists between these nodes")
	}

	var path []int64
	curr := targetID
	for {
		path = append([]int64{curr}, path...) // Prepend
		if curr == sourceID {
			break
		}
		curr = previous[curr]
	}

	return &PathResult{
		SourceID:    sourceID,
		TargetID:    targetID,
		TotalCost:   costs[targetID],
		TotalDistKm: distances[targetID],
		RouteNodes:  path,
	}, nil
}