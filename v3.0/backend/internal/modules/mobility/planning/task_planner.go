package planning

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
)

type Planner struct {
	SpatialState func(tenant, device string) (model.SpatialState, bool)
	Routing      routing.RoutingEngine
	MaxPlanAge   time.Duration
}

func (p *Planner) Name() string { return "mobility.task-planner/v1" }
func (p *Planner) Supports(t taskcore.Task) bool {
	if t.TaskType != "NAVIGATE" && t.TaskType != "RELOCATE" {
		return false
	}
	var requirements taskcore.Requirements
	return json.Unmarshal(t.Requirements, &requirements) == nil && requirements.PlanningMode == taskcore.PlanningPolarisRequired
}

type target struct {
	Lat       *float64                `json:"lat"`
	Lon       *float64                `json:"lon"`
	Latitude  *float64                `json:"latitude"`
	Longitude *float64                `json:"longitude"`
	Policy    routing.RouteCostPolicy `json:"policy"`
}
type routePayload struct {
	RouteID                string                  `json:"route_id"`
	RouteSchemaVersion     uint32                  `json:"route_schema_version"`
	RoadGraphVersion       string                  `json:"road_graph_version"`
	RoutingSnapshotVersion uint64                  `json:"routing_snapshot_version"`
	GeneratedAt            time.Time               `json:"generated_at"`
	ValidUntil             time.Time               `json:"valid_until"`
	Origin                 model.Position          `json:"origin"`
	Destination            model.Position          `json:"destination"`
	Waypoints              []model.Position        `json:"waypoints"`
	DistanceMeters         float64                 `json:"distance_meters"`
	EstimatedDurationMS    int64                   `json:"estimated_duration_ms"`
	Policy                 routing.RouteCostPolicy `json:"policy"`
}

func (p *Planner) Plan(ctx context.Context, req extension.PlanningRequest) (extension.ExecutionPlan, error) {
	if p.Routing == nil {
		return extension.ExecutionPlan{}, routing.ErrUnavailable
	}
	if req.Task.AssignedDeviceID == nil {
		return extension.ExecutionPlan{}, errors.New("planning requires selected device")
	}
	state, ok := p.SpatialState(req.Task.TenantID, *req.Task.AssignedDeviceID)
	if !ok {
		return extension.ExecutionPlan{}, errors.New("selected device has no active spatial state")
	}
	var t target
	if json.Unmarshal(req.Task.Target, &t) != nil {
		return extension.ExecutionPlan{}, errors.New("invalid mobility target")
	}
	lat, lon := t.Lat, t.Lon
	if lat == nil {
		lat = t.Latitude
	}
	if lon == nil {
		lon = t.Longitude
	}
	if lat == nil || lon == nil {
		return extension.ExecutionPlan{}, errors.New("mobility target coordinates required")
	}
	policy := t.Policy
	if policy == "" {
		policy = routing.RouteFastest
	}
	destination := model.Position{Latitude: *lat, Longitude: *lon}
	route, err := p.Routing.Route(ctx, routing.RouteRequest{TenantID: req.Task.TenantID, MobilityProfile: state.MobilityProfile, Origin: state.ReportedPosition, Destination: destination, Policy: policy})
	if err != nil {
		if errors.Is(err, routing.ErrUnsupportedProfile) {
			return extension.ExecutionPlan{}, extension.ErrPlanningUnsupported
		}
		return extension.ExecutionPlan{}, err
	}
	now := time.Now().UTC()
	valid := now.Add(p.MaxPlanAge)
	if req.Task.ExpiresAt.Before(valid) {
		valid = req.Task.ExpiresAt
	}
	payload, err := json.Marshal(routePayload{route.RouteID, 1, route.GraphVersion, route.SnapshotVersion, now, valid, state.ReportedPosition, destination, route.Waypoints, route.DistanceMeters, route.EstimatedTime.Milliseconds(), route.Policy})
	if err != nil {
		return extension.ExecutionPlan{}, err
	}
	return extension.ExecutionPlan{PlannerName: p.Name(), SchemaVersion: 1, CommandType: req.Task.TaskType, Payload: payload, GeneratedAt: now, ValidUntil: &valid, Metadata: map[string]any{"route_id": route.RouteID, "road_graph_version": route.GraphVersion, "routing_snapshot_version": route.SnapshotVersion}}, nil
}
