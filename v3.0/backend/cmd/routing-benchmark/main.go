package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
)

type routeClass struct {
	Name        string
	Origin      model.Position
	Destination model.Position
}

type measured struct {
	Class       string                  `json:"class"`
	Algorithm   routing.SearchAlgorithm `json:"algorithm"`
	LatencyUS   int64                   `json:"latency_us"`
	Allocations float64                 `json:"allocations_per_run"`
	Metrics     routing.SearchMetrics   `json:"metrics"`
	Error       string                  `json:"error,omitempty"`
}

func main() {
	graphPath := flag.String("graph", "data/chennai-metro.osm.pbf", "OSM PBF road graph")
	version := flag.String("version", "chennai-v1", "road graph version")
	repetitions := flag.Int("repetitions", 3, "timed repetitions per algorithm and route")
	flag.Parse()
	ctx := context.Background()
	graph, err := routing.LoadOSMPBF(ctx, *graphPath, *version)
	if err != nil {
		panic(err)
	}
	snapshot := routing.NewSnapshotStore(graph).Load()
	classes := []routeClass{
		{"short-1-3km", model.Position{Latitude: 13.0604, Longitude: 80.2496}, model.Position{Latitude: 13.0740, Longitude: 80.2600}},
		{"medium-5-15km", model.Position{Latitude: 13.0067, Longitude: 80.2206}, model.Position{Latitude: 13.0827, Longitude: 80.2707}},
		{"long-cross-city", model.Position{Latitude: 12.9010, Longitude: 80.2279}, model.Position{Latitude: 13.1600, Longitude: 80.3000}},
		{"dense-city", model.Position{Latitude: 13.0400, Longitude: 80.2100}, model.Position{Latitude: 13.0900, Longitude: 80.2850}},
		{"edge-of-graph", model.Position{Latitude: 12.8000, Longitude: 80.0500}, model.Position{Latitude: 13.2400, Longitude: 80.3400}},
	}
	results := []measured{}
	for _, routeClass := range classes {
		from, fromErr := graph.NodeIndex().Nearest(ctx, routeClass.Origin)
		to, toErr := graph.NodeIndex().Nearest(ctx, routeClass.Destination)
		if fromErr != nil || toErr != nil {
			results = append(results, measured{Class: routeClass.Name, Error: fmt.Sprintf("snap origin=%v destination=%v", fromErr, toErr)})
			continue
		}
		var oracleCost float64
		for _, algorithm := range []routing.SearchAlgorithm{routing.AlgorithmAStar, routing.AlgorithmDijkstra} {
			started := time.Now()
			var metrics routing.SearchMetrics
			for n := 0; n < *repetitions; n++ {
				metrics, err = routing.MeasureSearch(ctx, graph, snapshot, from.ID, to.ID, routing.RouteFastest, graph.NodeCount()*2, algorithm)
				if err != nil {
					break
				}
			}
			entry := measured{Class: routeClass.Name, Algorithm: algorithm, LatencyUS: time.Since(started).Microseconds() / int64(max(*repetitions, 1)), Metrics: metrics}
			if err != nil {
				entry.Error = err.Error()
			} else {
				entry.Allocations = testing.AllocsPerRun(1, func() {
					_, _ = routing.MeasureSearch(ctx, graph, snapshot, from.ID, to.ID, routing.RouteFastest, graph.NodeCount()*2, algorithm)
				})
				if algorithm == routing.AlgorithmAStar {
					oracleCost = metrics.Cost
				} else if math.Abs(metrics.Cost-oracleCost) > 1e-6 {
					entry.Error = "cost differs from A* oracle"
				}
			}
			results = append(results, entry)
		}
	}
	out, _ := json.MarshalIndent(map[string]any{"measured_at": time.Now().UTC(), "road_graph_version": graph.Version(), "nodes": graph.NodeCount(), "edges": graph.EdgeCount(), "results": results}, "", "  ")
	fmt.Println(string(out))
}
