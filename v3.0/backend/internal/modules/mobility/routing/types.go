package routing

import (
	"context"
	"errors"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type EdgeID int64
type NodeType uint8

const (
	NodeUnknown NodeType = iota
	NodeRoad
	NodeIntersection
	NodeChargingStation
)

type RoadClass string

const (
	RoadMotorway     RoadClass = "MOTORWAY"
	RoadPrimary      RoadClass = "PRIMARY"
	RoadSecondary    RoadClass = "SECONDARY"
	RoadTertiary     RoadClass = "TERTIARY"
	RoadResidential  RoadClass = "RESIDENTIAL"
	RoadUnclassified RoadClass = "UNCLASSIFIED"
)

type RoadNode struct {
	ID       int64
	Position model.Position
	Type     NodeType
}
type RoadEdge struct {
	ID             EdgeID
	FromID, ToID   int64
	DistanceM      float64
	BaseTravelTime time.Duration
	RoadClass      RoadClass
}
type RouteCostPolicy string

const (
	RouteShortest RouteCostPolicy = "SHORTEST"
	RouteFastest  RouteCostPolicy = "FASTEST"
)

type RouteRequest struct {
	TenantID        string                `json:"tenant_id"`
	MobilityProfile model.MobilityProfile `json:"mobility_profile"`
	Origin          model.Position        `json:"origin"`
	Destination     model.Position        `json:"destination"`
	Policy          RouteCostPolicy       `json:"policy"`
}
type RouteResult struct {
	RouteID         string           `json:"route_id"`
	GraphVersion    string           `json:"road_graph_version"`
	SnapshotVersion uint64           `json:"snapshot_version"`
	Policy          RouteCostPolicy  `json:"policy"`
	DistanceMeters  float64          `json:"distance_meters"`
	EstimatedTime   time.Duration    `json:"estimated_time"`
	Waypoints       []model.Position `json:"waypoints"`
	EdgeIDs         []EdgeID         `json:"edge_ids"`
	ExpandedNodes   int              `json:"expanded_nodes"`
}
type RoutingEngine interface {
	Route(context.Context, RouteRequest) (RouteResult, error)
}
type RoadNodeIndex interface {
	Nearest(context.Context, model.Position) (RoadNode, error)
}

var (
	ErrNoRoute            = errors.New("NO_ROUTE")
	ErrNoRoadNode         = errors.New("NO_ROAD_NODE")
	ErrOutsideRegion      = errors.New("OUTSIDE_ROUTING_REGION")
	ErrUnsupportedProfile = errors.New("UNSUPPORTED_MOBILITY_PROFILE")
	ErrUnavailable        = errors.New("ROUTING_UNAVAILABLE")
	ErrBusy               = errors.New("ROUTING_BUSY")
	ErrTimeout            = errors.New("ROUTING_TIMEOUT")
)
