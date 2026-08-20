package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
	"github.com/gin-gonic/gin"
)

type MobilityAPI struct {
	spatial  *spatial.Manager
	routing  routing.RoutingEngine
	maxLimit int
}

func NewMobilityAPI(s *spatial.Manager, r routing.RoutingEngine, maxLimit int) *MobilityAPI {
	return &MobilityAPI{spatial: s, routing: r, maxLimit: maxLimit}
}
func (a *MobilityAPI) Register(r *gin.RouterGroup) {
	r.GET("/spatial/devices/nearby", a.nearby)
	r.POST("/routes", a.route)
	r.GET("/routes/calculate", a.legacyRoute)
}

func (a *MobilityAPI) legacyRoute(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	srcLat, e1 := strconv.ParseFloat(c.Query("src_lat"), 64)
	srcLon, e2 := strconv.ParseFloat(c.Query("src_lon"), 64)
	dstLat, e3 := strconv.ParseFloat(c.Query("tgt_lat"), 64)
	dstLon, e4 := strconv.ParseFloat(c.Query("tgt_lon"), 64)
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		apiError(c, 400, "INVALID_ROUTE_REQUEST", "Valid source and target coordinates are required")
		return
	}
	result, err := a.routing.Route(c, routing.RouteRequest{TenantID: tenant, MobilityProfile: model.MobilityRoadVehicle, Origin: model.Position{Latitude: srcLat, Longitude: srcLon}, Destination: model.Position{Latitude: dstLat, Longitude: dstLon}, Policy: routing.RouteFastest})
	if err != nil {
		routingError(c, err)
		return
	}
	route := make([]gin.H, len(result.Waypoints))
	for i, p := range result.Waypoints {
		route[i] = gin.H{"lat": p.Latitude, "lon": p.Longitude}
	}
	c.JSON(200, gin.H{"status": "success", "total_dist_km": result.DistanceMeters / 1000, "estimated_duration_ms": result.EstimatedTime.Milliseconds(), "road_graph_version": result.GraphVersion, "routing_snapshot_version": result.SnapshotVersion, "route": route})
}
func (a *MobilityAPI) nearby(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	lat, e1 := strconv.ParseFloat(c.Query("lat"), 64)
	lon, e2 := strconv.ParseFloat(c.Query("lon"), 64)
	radius, e3 := strconv.ParseFloat(c.DefaultQuery("radius_meters", "5000"), 64)
	limit, e4 := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil || limit < 1 || limit > a.maxLimit {
		apiError(c, 400, "INVALID_SPATIAL_QUERY", "Coordinates, radius, or limit are invalid")
		return
	}
	items, err := a.spatial.Nearby(tenant, model.Position{Latitude: lat, Longitude: lon}, radius, limit)
	if err != nil {
		apiError(c, 400, "INVALID_SPATIAL_QUERY", err.Error())
		return
	}
	apiData(c, 200, gin.H{"count": len(items), "devices": items})
}
func (a *MobilityAPI) route(c *gin.Context) {
	tenant, ok := tenantFor(c)
	if !ok {
		apiError(c, 400, "TENANT_REQUIRED", "Tenant scope is required")
		return
	}
	var req routing.RouteRequest
	if c.ShouldBindJSON(&req) != nil {
		apiError(c, 400, "INVALID_ROUTE_REQUEST", "A valid route request is required")
		return
	}
	req.TenantID = tenant
	if req.MobilityProfile == "" {
		req.MobilityProfile = model.MobilityRoadVehicle
	}
	if req.Policy == "" {
		req.Policy = routing.RouteFastest
	}
	result, err := a.routing.Route(c, req)
	if err != nil {
		routingError(c, err)
		return
	}
	apiData(c, 200, result)
}
func routingError(c *gin.Context, err error) {
	code := http.StatusUnprocessableEntity
	name := err.Error()
	switch {
	case errors.Is(err, routing.ErrBusy):
		code = http.StatusTooManyRequests
	case errors.Is(err, routing.ErrTimeout):
		code = http.StatusGatewayTimeout
	case errors.Is(err, routing.ErrUnavailable):
		code = http.StatusServiceUnavailable
	case errors.Is(err, routing.ErrNoRoute), errors.Is(err, routing.ErrNoRoadNode), errors.Is(err, routing.ErrOutsideRegion), errors.Is(err, routing.ErrUnsupportedProfile):
		code = http.StatusUnprocessableEntity
	default:
		code = http.StatusInternalServerError
		name = "ROUTING_ERROR"
	}
	apiError(c, code, name, err.Error())
}
