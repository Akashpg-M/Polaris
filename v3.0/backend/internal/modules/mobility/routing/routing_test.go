package routing

import (
	"context"
	"errors"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

func testGraph(t *testing.T) *RoadGraph {
	b := NewGraphBuilder("test-v1")
	nodes := []RoadNode{{1, model.Position{Latitude: 0, Longitude: 0}, NodeRoad}, {2, model.Position{Latitude: 0, Longitude: .01}, NodeRoad}, {3, model.Position{Latitude: .01, Longitude: .01}, NodeRoad}, {4, model.Position{Latitude: .01, Longitude: 0}, NodeRoad}}
	for _, n := range nodes {
		if err := b.AddNode(n); err != nil {
			t.Fatal(err)
		}
	}
	edges := []RoadEdge{{1, 1, 2, 1000, time.Minute, RoadResidential}, {2, 2, 3, 1000, time.Minute, RoadResidential}, {3, 1, 4, 1200, time.Minute, RoadResidential}, {4, 4, 3, 1200, time.Minute, RoadResidential}}
	for _, e := range edges {
		if err := b.AddEdge(e); err != nil {
			t.Fatal(err)
		}
	}
	g, err := b.Build()
	if err != nil {
		t.Fatal(err)
	}
	return g
}
func TestAStarEqualsDijkstraAndHonorsDirection(t *testing.T) {
	g := testGraph(t)
	s := NewSnapshotStore(g).Load()
	for _, policy := range []RouteCostPolicy{RouteShortest, RouteFastest} {
		a, ea := AStar(context.Background(), g, s, 1, 3, policy, 100)
		d, ed := Dijkstra(context.Background(), g, s, 1, 3, policy, 100)
		if ea != nil || ed != nil || math.Abs(a.cost-d.cost) > 1e-9 {
			t.Fatalf("policy %s A*=%v/%v Dijkstra=%v/%v", policy, a.cost, ea, d.cost, ed)
		}
	}
	if _, err := AStar(context.Background(), g, s, 3, 1, RouteShortest, 100); !errors.Is(err, ErrNoRoute) {
		t.Fatalf("one-way route should fail, got %v", err)
	}
}
func TestSnapshotChangesCostWithoutMutatingTopology(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	costs := append([]float64(nil), store.Load().EdgeCosts...)
	costs[0], costs[2] = 1000, 1000
	if err := store.Swap(RoutingCostSnapshot{Version: 2, EdgeCosts: costs}); err != nil {
		t.Fatal(err)
	}
	result, err := AStar(context.Background(), g, store.Load(), 1, 3, RouteFastest, 100)
	if err != nil {
		t.Fatal(err)
	}
	if result.edgeIndexes[0] != 1 {
		t.Fatalf("dynamic fastest route did not change: %v", result.edgeIndexes)
	}
	if g.EdgeCount() != 4 {
		t.Fatal("topology mutated")
	}
}
func TestSnapshotConcurrentReadersUseOneVersion(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	var wg sync.WaitGroup
	versions := make(chan uint64, 200)
	for n := 0; n < 100; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snapshot := store.Load()
			_, _ = AStar(context.Background(), g, snapshot, 1, 3, RouteFastest, 100)
			versions <- snapshot.Version
		}()
	}
	for version := uint64(2); version <= 5; version++ {
		old := store.Load()
		_ = store.Swap(RoutingCostSnapshot{Version: version, EdgeCosts: append([]float64(nil), old.EdgeCosts...)})
	}
	wg.Wait()
	close(versions)
	for version := range versions {
		if version < 1 || version > 5 {
			t.Fatalf("mixed/invalid snapshot version %d", version)
		}
	}
}
func TestKDTreeAndNodeTypeSemantics(t *testing.T) {
	g := testGraph(t)
	node, err := g.NodeIndex().Nearest(context.Background(), model.Position{Latitude: 0, Longitude: .0099})
	if err != nil || node.ID != 2 {
		t.Fatalf("nearest node=%v err=%v", node.ID, err)
	}
	if NodeRoad == NodeIntersection || NodeIntersection == NodeChargingStation {
		t.Fatal("NodeType enum values overlap")
	}
}
func TestBoundedEngineIsolationAndUnsupportedProfile(t *testing.T) {
	g := testGraph(t)
	e := NewEngine(g, NewSnapshotStore(g), EngineConfig{Workers: 1, QueueCapacity: 1, MaxExpansions: 100, MaxConcurrentPerTenant: 1, Timeout: time.Second})
	e.Start()
	defer e.Close()
	_, err := e.Route(context.Background(), RouteRequest{TenantID: "t", MobilityProfile: model.MobilityAerialDrone, Origin: g.nodes[1].Position, Destination: g.nodes[3].Position, Policy: RouteFastest})
	if !errors.Is(err, ErrUnsupportedProfile) {
		t.Fatalf("expected unsupported profile, got %v", err)
	}
}

func TestRoutingQueueSaturationIsExplicit(t *testing.T) {
	g := testGraph(t)
	engine := NewEngine(g, NewSnapshotStore(g), EngineConfig{Workers: 0, QueueCapacity: 1, MaxExpansions: 100, MaxConcurrentPerTenant: 1, Timeout: 50 * time.Millisecond})
	engine.Start()
	defer engine.Close()
	req := RouteRequest{TenantID: "tenant-a", MobilityProfile: model.MobilityRoadVehicle, Origin: g.nodes[1].Position, Destination: g.nodes[3].Position, Policy: RouteFastest}
	firstReq := req
	done := make(chan error, 1)
	go func() { _, err := engine.Route(context.Background(), firstReq); done <- err }()
	time.Sleep(5 * time.Millisecond)
	req.TenantID = "tenant-b"
	if _, err := engine.Route(context.Background(), req); !errors.Is(err, ErrBusy) {
		t.Fatalf("expected ROUTING_BUSY, got %v", err)
	}
	if err := <-done; !errors.Is(err, ErrTimeout) {
		t.Fatalf("queued request should time out, got %v", err)
	}
}
func TestInvalidSnapshotRejected(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	if err := store.Swap(RoutingCostSnapshot{Version: 2, EdgeCosts: []float64{1}}); err == nil {
		t.Fatal("partial snapshot accepted")
	}
}

func TestTrafficObservationUpdatesExactlyOneDirectedEdge(t *testing.T) {
	g := testGraph(t)
	store := NewSnapshotStore(g)
	traffic := NewTrafficManager(g, store, time.Minute)
	heading := 90.0
	err := traffic.Observe(context.Background(), TrafficObservation{Position: model.Position{Latitude: 0, Longitude: .005}, HeadingDegrees: &heading, SpeedMPS: 1, ObservedAt: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	before := append([]float64(nil), store.Load().EdgeCosts...)
	if err = traffic.Refresh(time.Now().UTC()); err != nil {
		t.Fatal(err)
	}
	changed := 0
	for i, cost := range store.Load().EdgeCosts {
		if math.Abs(cost-before[i]) > 1e-9 {
			changed++
		}
	}
	if changed != 1 || traffic.StateCount() != 1 {
		t.Fatalf("one observation changed %d edges and retained %d states", changed, traffic.StateCount())
	}
}

func BenchmarkAStar(b *testing.B)    { benchmarkSearch(b, true) }
func BenchmarkDijkstra(b *testing.B) { benchmarkSearch(b, false) }
func benchmarkSearch(b *testing.B, astar bool) {
	const side = 30
	builder := NewGraphBuilder("benchmark-grid")
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			id := int64(y*side + x + 1)
			_ = builder.AddNode(RoadNode{ID: id, Position: model.Position{Latitude: 13 + float64(y)*.001, Longitude: 80 + float64(x)*.001}, Type: NodeRoad})
		}
	}
	edgeID := EdgeID(1)
	for y := 0; y < side; y++ {
		for x := 0; x < side; x++ {
			from := int64(y*side + x + 1)
			if x+1 < side {
				to := from + 1
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: from, ToID: to, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: to, ToID: from, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
			}
			if y+1 < side {
				to := from + side
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: from, ToID: to, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
				_ = builder.AddEdge(RoadEdge{ID: edgeID, FromID: to, ToID: from, DistanceM: 110, BaseTravelTime: 10 * time.Second, RoadClass: RoadResidential})
				edgeID++
			}
		}
	}
	g, _ := builder.Build()
	snapshot := NewSnapshotStore(g).Load()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if astar {
			_, _ = AStar(context.Background(), g, snapshot, 1, int64(side*side), RouteFastest, side*side*2)
		} else {
			_, _ = Dijkstra(context.Background(), g, snapshot, 1, int64(side*side), RouteFastest, side*side*2)
		}
	}
}
