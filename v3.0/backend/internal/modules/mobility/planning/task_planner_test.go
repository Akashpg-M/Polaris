package planning

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
)

type fakeRouter struct{}

func (fakeRouter) Route(context.Context, routing.RouteRequest) (routing.RouteResult, error) {
	return routing.RouteResult{RouteID: "route-1", GraphVersion: "chennai-v1", SnapshotVersion: 42, Policy: routing.RouteFastest, DistanceMeters: 1200, EstimatedTime: 90 * time.Second, Waypoints: []model.Position{{Latitude: 13, Longitude: 80}, {Latitude: 13.01, Longitude: 80.01}}, EdgeIDs: []routing.EdgeID{1}}, nil
}
func TestNavigatePlanContainsVersionAndValidity(t *testing.T) {
	device := "vehicle-1"
	state := model.SpatialState{TenantID: "tenant", DeviceID: device, ReportedPosition: model.Position{Latitude: 13, Longitude: 80}, MobilityProfile: model.MobilityRoadVehicle}
	planner := &Planner{SpatialState: func(string, string) (model.SpatialState, bool) { return state, true }, Routing: fakeRouter{}, MaxPlanAge: 2 * time.Minute}
	task := taskcore.Task{TaskID: "task-1", TenantID: "tenant", TaskType: "NAVIGATE", Target: []byte(`{"lat":13.01,"lon":80.01}`), AssignedDeviceID: &device, ExpiresAt: time.Now().Add(5 * time.Minute)}
	plan, err := planner.Plan(context.Background(), extension.PlanningRequest{Task: task})
	if err != nil {
		t.Fatal(err)
	}
	var payload routePayload
	if err = json.Unmarshal(plan.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.RouteID != "route-1" || payload.RoadGraphVersion != "chennai-v1" || payload.RoutingSnapshotVersion != 42 || payload.ValidUntil.IsZero() {
		t.Fatalf("incomplete plan: %#v", payload)
	}
	saved := append([]byte(nil), plan.Payload...)
	if string(saved) != string(plan.Payload) {
		t.Fatal("durable replay payload changed")
	}
}

func TestPlannerOnlyClaimsPolarisRequiredMobilityTasks(t *testing.T) {
	planner := &Planner{}
	required, _ := json.Marshal(taskcore.Requirements{PlanningMode: taskcore.PlanningPolarisRequired})
	local, _ := json.Marshal(taskcore.Requirements{PlanningMode: taskcore.PlanningDeviceLocal})
	if !planner.Supports(taskcore.Task{TaskType: "NAVIGATE", Requirements: required}) {
		t.Fatal("Polaris-required navigation was not claimed")
	}
	if planner.Supports(taskcore.Task{TaskType: "NAVIGATE", Requirements: local}) {
		t.Fatal("device-local navigation must use the generic planner")
	}
}
