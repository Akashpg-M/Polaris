package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type result struct {
	Index               string  `json:"index"`
	Devices             int     `json:"devices"`
	IndexedDevices      int     `json:"indexed_devices"`
	Queries             int     `json:"queries"`
	P50US               int64   `json:"p50_us"`
	P95US               int64   `json:"p95_us"`
	P99US               int64   `json:"p99_us"`
	OperationsPerSecond float64 `json:"operations_per_second"`
}

func main() {
	devices := flag.Int("devices", 1000, "active device count")
	queries := flag.Int("queries", 1000, "nearest-query count")
	radius := flag.Float64("search-radius", 5000, "search radius in meters")
	limit := flag.Int("candidate-limit", 10, "candidate limit")
	moving := flag.Float64("moving-percent", 40, "percentage of devices updated before queries")
	duration := flag.Duration("duration", 0, "optional minimum benchmark duration")
	flag.Parse()
	if *devices < 1 || *queries < 1 {
		panic("devices and queries must be positive")
	}
	rng := rand.New(rand.NewSource(42))
	states := make([]model.SpatialState, 0, *devices*85/100)
	movingStates := []model.SpatialState{}
	for i := 0; i < *devices; i++ {
		// 40% road vehicles, 20% robots, 25% static spatial devices, and
		// 15% deliberately non-spatial compute/other devices.
		if i >= *devices*85/100 {
			continue
		}
		profile := model.MobilityStatic
		if i < *devices*40/100 {
			profile = model.MobilityRoadVehicle
		} else if i < *devices*60/100 {
			profile = model.MobilityGroundRobot
		}
		s := model.SpatialState{TenantID: "benchmark", DeviceID: fmt.Sprintf("device-%06d", i), MobilityProfile: profile, Position: model.Position{Latitude: 13.0 + rng.Float64()*.2, Longitude: 80.1 + rng.Float64()*.2}}
		states = append(states, s)
		if profile != model.MobilityStatic {
			movingStates = append(movingStates, s)
		}
	}
	run := func(name string, index spatial.SpatialIndex) result {
		for _, s := range states {
			_ = index.Upsert(s)
		}
		moveCount := int(float64(len(movingStates)) * *moving / 100)
		for i := 0; i < moveCount; i++ {
			s := movingStates[i]
			s.Position.Latitude += .0001
			_ = index.Upsert(s)
		}
		latencies := make([]time.Duration, 0, *queries)
		started := time.Now()
		count := 0
		for count < *queries || (*duration > 0 && time.Since(started) < *duration) {
			q := model.Position{Latitude: 13.0 + rng.Float64()*.2, Longitude: 80.1 + rng.Float64()*.2}
			at := time.Now()
			_, _ = index.Nearest(q, *limit, *radius)
			latencies = append(latencies, time.Since(at))
			count++
		}
		elapsed := time.Since(started)
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		percentile := func(p float64) int64 { return latencies[int(float64(len(latencies)-1)*p)].Microseconds() }
		return result{name, *devices, len(states), len(latencies), percentile(.50), percentile(.95), percentile(.99), float64(len(latencies)) / elapsed.Seconds()}
	}
	results := []result{run("rtree", spatial.NewRTreeSpatialIndex()), run("linear", spatial.NewLinearSpatialIndex())}
	data, _ := json.MarshalIndent(map[string]any{"measured_at": time.Now().UTC(), "parameters": map[string]any{"devices": *devices, "queries": *queries, "radius_meters": *radius, "candidate_limit": *limit, "moving_percent_of_mobile": *moving, "distribution": "40% road, 20% robot, 25% static spatial, 15% non-spatial"}, "results": results}, "", "  ")
	fmt.Println(string(data))
}
