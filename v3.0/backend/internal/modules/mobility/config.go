package mobility

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Enabled                   bool
	Required                  bool
	SpatialEnabled            bool
	RoutingEnabled            bool
	H3Resolution              int
	H3ShardResolution         int
	IndexMinMoveMeters        float64
	IndexMaxAge               time.Duration
	MaxH3Rings                int
	MaxSearchRadiusMeters     float64
	MaxRawCandidates          int
	MaxRoutedCandidates       int
	MaxActiveDevicesPerTenant int
	RoutingWorkers            int
	RoutingQueueCapacity      int
	RoutingTimeout            time.Duration
	MaxRouteExpansions        int
	MaxConcurrentRoutesTenant int
	MaxTrafficObservationAge  time.Duration
	TrafficRefreshInterval    time.Duration
	TrafficScope              string
	RoadGraphPath             string
	RoadGraphVersion          string
}

func DefaultConfig() Config {
	return Config{Enabled: true, SpatialEnabled: true, RoutingEnabled: true, H3Resolution: 8, H3ShardResolution: 6,
		IndexMinMoveMeters: 5, IndexMaxAge: 30 * time.Second, MaxH3Rings: 12, MaxSearchRadiusMeters: 10_000,
		MaxRawCandidates: 50, MaxRoutedCandidates: 8, MaxActiveDevicesPerTenant: 10_000,
		RoutingWorkers: 4, RoutingQueueCapacity: 64, RoutingTimeout: 2 * time.Second, MaxRouteExpansions: 250_000,
		MaxConcurrentRoutesTenant: 2, MaxTrafficObservationAge: 10 * time.Minute, TrafficRefreshInterval: 15 * time.Second, TrafficScope: "SHARED_TRUSTED",
		RoadGraphPath: "data/chennai-metro.osm.pbf", RoadGraphVersion: "chennai-v1"}
}

func LoadConfig() (Config, error) {
	c := DefaultConfig()
	c.Enabled = envBool("POLARIS_MODULE_MOBILITY_ENABLED", c.Enabled)
	c.Required = envBool("POLARIS_MODULE_MOBILITY_REQUIRED", c.Required)
	c.SpatialEnabled = envBool("MOBILITY_SPATIAL_ENABLED", c.SpatialEnabled)
	c.RoutingEnabled = envBool("MOBILITY_ROUTING_ENABLED", c.RoutingEnabled)
	c.H3Resolution = envInt("MOBILITY_H3_RESOLUTION", c.H3Resolution)
	c.H3ShardResolution = envInt("MOBILITY_H3_SHARD_RESOLUTION", c.H3ShardResolution)
	c.IndexMinMoveMeters = envFloat("MOBILITY_INDEX_MIN_MOVE_METERS", c.IndexMinMoveMeters)
	c.IndexMaxAge = envDuration("MOBILITY_INDEX_MAX_AGE", c.IndexMaxAge)
	c.MaxH3Rings = envInt("MOBILITY_MAX_H3_RINGS", c.MaxH3Rings)
	c.MaxSearchRadiusMeters = envFloat("MOBILITY_MAX_SEARCH_RADIUS_METERS", c.MaxSearchRadiusMeters)
	c.MaxRawCandidates = envInt("MOBILITY_MAX_RAW_CANDIDATES", c.MaxRawCandidates)
	c.MaxRoutedCandidates = envInt("MOBILITY_MAX_ROUTED_CANDIDATES", c.MaxRoutedCandidates)
	c.MaxActiveDevicesPerTenant = envInt("MOBILITY_MAX_ACTIVE_DEVICES_PER_TENANT", c.MaxActiveDevicesPerTenant)
	c.RoutingWorkers = envInt("MOBILITY_ROUTING_WORKERS", c.RoutingWorkers)
	c.RoutingQueueCapacity = envInt("MOBILITY_ROUTING_QUEUE_CAPACITY", c.RoutingQueueCapacity)
	c.RoutingTimeout = envDuration("MOBILITY_ROUTING_TIMEOUT", c.RoutingTimeout)
	c.MaxRouteExpansions = envInt("MOBILITY_MAX_ROUTE_EXPANSIONS", c.MaxRouteExpansions)
	c.MaxConcurrentRoutesTenant = envInt("MOBILITY_MAX_CONCURRENT_ROUTES_PER_TENANT", c.MaxConcurrentRoutesTenant)
	c.MaxTrafficObservationAge = envDuration("MOBILITY_MAX_TRAFFIC_OBSERVATION_AGE", c.MaxTrafficObservationAge)
	c.TrafficRefreshInterval = envDuration("MOBILITY_TRAFFIC_REFRESH_INTERVAL", c.TrafficRefreshInterval)
	if v := os.Getenv("MOBILITY_TRAFFIC_SCOPE"); v != "" {
		c.TrafficScope = v
	}
	if v := os.Getenv("MOBILITY_ROAD_GRAPH_PATH"); v != "" {
		c.RoadGraphPath = v
	}
	if v := os.Getenv("MOBILITY_ROAD_GRAPH_VERSION"); v != "" {
		c.RoadGraphVersion = v
	}
	return c, c.Validate()
}

func (c Config) Validate() error {
	if c.H3Resolution < 0 || c.H3Resolution > 15 || c.H3ShardResolution < 0 || c.H3ShardResolution > c.H3Resolution {
		return errors.New("invalid H3 resolutions")
	}
	if c.MaxH3Rings < 0 || c.MaxSearchRadiusMeters <= 0 || c.MaxRawCandidates < 1 || c.MaxRoutedCandidates < 1 || c.MaxRoutedCandidates > c.MaxRawCandidates {
		return errors.New("invalid candidate limits")
	}
	if c.RoutingWorkers < 1 || c.RoutingQueueCapacity < 1 || c.RoutingTimeout <= 0 || c.MaxRouteExpansions < 1 || c.MaxConcurrentRoutesTenant < 1 {
		return errors.New("invalid routing limits")
	}
	if c.MaxActiveDevicesPerTenant < 1 || c.IndexMinMoveMeters < 0 || c.IndexMaxAge <= 0 || c.MaxTrafficObservationAge <= 0 || c.TrafficRefreshInterval <= 0 {
		return errors.New("invalid spatial limits")
	}
	if c.TrafficScope != "SHARED_TRUSTED" {
		return errors.New("unsupported traffic scope; only SHARED_TRUSTED is currently implemented")
	}
	return nil
}

func envBool(k string, d bool) bool {
	if v, e := strconv.ParseBool(os.Getenv(k)); e == nil {
		return v
	}
	return d
}
func envInt(k string, d int) int {
	if v, e := strconv.Atoi(os.Getenv(k)); e == nil {
		return v
	}
	return d
}
func envFloat(k string, d float64) float64 {
	if v, e := strconv.ParseFloat(os.Getenv(k), 64); e == nil {
		return v
	}
	return d
}
func envDuration(k string, d time.Duration) time.Duration {
	if v, e := time.ParseDuration(os.Getenv(k)); e == nil {
		return v
	}
	return d
}
