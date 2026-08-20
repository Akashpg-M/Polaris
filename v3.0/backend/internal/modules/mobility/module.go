package mobility

import (
	"context"
	"sync"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type RebuildLoader func(context.Context) ([]model.SpatialState, error)
type Module struct {
	cfg     Config
	Spatial *spatial.Manager
	loader  RebuildLoader
	mu      sync.RWMutex
	state   extension.ModuleState
	message string
	router  *routing.Engine
	graph   *routing.RoadGraph
	traffic *routing.TrafficManager
}

func New(cfg Config, loader RebuildLoader) *Module {
	manager := spatial.NewManager(spatial.ManagerConfig{H3Resolution: cfg.H3Resolution, ShardResolution: cfg.H3ShardResolution, MinMoveMeters: cfg.IndexMinMoveMeters, MaxIndexAge: cfg.IndexMaxAge, MaxH3Rings: cfg.MaxH3Rings, MaxRadiusMeters: cfg.MaxSearchRadiusMeters, MaxDevicesPerTenant: cfg.MaxActiveDevicesPerTenant})
	return &Module{cfg: cfg, Spatial: manager, loader: loader, state: extension.ModuleStarting}
}
func (m *Module) Name() string { return "mobility" }
func (m *Module) Start(ctx context.Context) error {
	m.mu.Lock()
	m.state = extension.ModuleStarting
	m.mu.Unlock()
	if m.cfg.SpatialEnabled && m.loader != nil {
		states, err := m.loader(ctx)
		if err != nil {
			m.set(extension.ModuleDegraded, "spatial rebuild failed: "+err.Error())
		} else if err = m.Spatial.Rebuild(states); err != nil {
			m.set(extension.ModuleDegraded, "spatial rebuild failed: "+err.Error())
		}
	}
	if m.cfg.RoutingEnabled {
		graph, err := routing.LoadOSMPBF(ctx, m.cfg.RoadGraphPath, m.cfg.RoadGraphVersion)
		if err != nil {
			m.set(extension.ModuleDegraded, "road graph unavailable: "+err.Error())
			if m.cfg.Required {
				return err
			}
		} else {
			snapshots := routing.NewSnapshotStore(graph)
			engine := routing.NewEngine(graph, snapshots, routing.EngineConfig{Workers: m.cfg.RoutingWorkers, QueueCapacity: m.cfg.RoutingQueueCapacity, MaxExpansions: m.cfg.MaxRouteExpansions, MaxConcurrentPerTenant: m.cfg.MaxConcurrentRoutesTenant, Timeout: m.cfg.RoutingTimeout})
			engine.Start()
			m.mu.Lock()
			m.graph, m.router = graph, engine
			m.traffic = routing.NewTrafficManager(graph, snapshots, m.cfg.MaxTrafficObservationAge)
			m.mu.Unlock()
		}
	}
	m.mu.Lock()
	if m.state == extension.ModuleStarting {
		m.state = extension.ModuleReady
		m.message = "spatial and configured routing components started"
	}
	m.mu.Unlock()
	return nil
}
func (m *Module) set(s extension.ModuleState, msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state, m.message = s, msg
}
func (m *Module) Ready(context.Context) extension.ModuleStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	spatialState := extension.ModuleReady
	if !m.cfg.SpatialEnabled {
		spatialState = extension.ModuleStopped
	}
	routingState := extension.ModuleReady
	if !m.cfg.RoutingEnabled {
		routingState = extension.ModuleStopped
	} else if m.router == nil {
		routingState = extension.ModuleFailed
	}
	details := map[string]any{}
	if m.graph != nil {
		details["road_graph_version"] = m.graph.Version()
		details["road_nodes"] = m.graph.NodeCount()
		details["road_edges"] = m.graph.EdgeCount()
		details["routing_snapshot_version"] = m.router.SnapshotStore().Load().Version
		details["routing_runtime"] = m.router.Stats()
		details["traffic_scope"] = m.cfg.TrafficScope
		details["traffic_refresh_interval"] = m.cfg.TrafficRefreshInterval.String()
		if m.traffic != nil {
			details["traffic_edge_states"] = m.traffic.StateCount()
			details["traffic_overlay_bytes"] = m.traffic.OverlayBytes()
		}
	}
	return extension.ModuleStatus{State: m.state, Message: m.message, Components: map[string]extension.ModuleStatus{"spatial": {State: spatialState}, "routing": {State: routingState}}, Details: details}
}
func (m *Module) Route(ctx context.Context, req routing.RouteRequest) (routing.RouteResult, error) {
	m.mu.RLock()
	r := m.router
	m.mu.RUnlock()
	if r == nil {
		return routing.RouteResult{}, routing.ErrUnavailable
	}
	return r.Route(ctx, req)
}
func (m *Module) Close(context.Context) error {
	m.mu.Lock()
	if m.router != nil {
		m.router.Close()
		m.router = nil
	}
	m.state = extension.ModuleStopped
	m.message = "stopped"
	m.mu.Unlock()
	return nil
}
func (m *Module) Traffic() *routing.TrafficManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.traffic
}
