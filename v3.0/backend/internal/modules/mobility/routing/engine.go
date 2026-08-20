package routing

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type routeJob struct {
	ctx    context.Context
	req    RouteRequest
	result chan routeResponse
}
type routeResponse struct {
	route RouteResult
	err   error
}
type EngineConfig struct {
	Workers, QueueCapacity, MaxExpansions, MaxConcurrentPerTenant int
	Timeout                                                       time.Duration
}
type Engine struct {
	graph     *RoadGraph
	snapshots *SnapshotStore
	cfg       EngineConfig
	jobs      chan routeJob
	closed    chan struct{}
	wg        sync.WaitGroup
	started   atomic.Bool
	tenantMu  sync.Mutex
	tenant    map[string]*tenantLimiter
	requests  atomic.Uint64
	busy      atomic.Uint64
}

type tenantLimiter struct {
	slots chan struct{}
	refs  int
}

type EngineStats struct {
	Requests      uint64 `json:"requests"`
	Busy          uint64 `json:"routing_busy"`
	QueueDepth    int    `json:"queue_depth"`
	QueueCapacity int    `json:"queue_capacity"`
	ActiveTenants int    `json:"active_tenants"`
}

func NewEngine(g *RoadGraph, s *SnapshotStore, c EngineConfig) *Engine {
	return &Engine{graph: g, snapshots: s, cfg: c, jobs: make(chan routeJob, c.QueueCapacity), closed: make(chan struct{}), tenant: map[string]*tenantLimiter{}}
}
func (e *Engine) Start() {
	if !e.started.CompareAndSwap(false, true) {
		return
	}
	for i := 0; i < e.cfg.Workers; i++ {
		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			for {
				select {
				case <-e.closed:
					return
				case j := <-e.jobs:
					r, err := e.route(j.ctx, j.req)
					select {
					case j.result <- routeResponse{r, err}:
					case <-j.ctx.Done():
					}
				}
			}
		}()
	}
}
func (e *Engine) acquireTenant(id string) (*tenantLimiter, bool) {
	e.tenantMu.Lock()
	defer e.tenantMu.Unlock()
	limiter := e.tenant[id]
	if limiter == nil {
		limiter = &tenantLimiter{slots: make(chan struct{}, e.cfg.MaxConcurrentPerTenant)}
		e.tenant[id] = limiter
	}
	select {
	case limiter.slots <- struct{}{}:
		limiter.refs++
		return limiter, true
	default:
		return nil, false
	}
}
func (e *Engine) releaseTenant(id string, limiter *tenantLimiter) {
	e.tenantMu.Lock()
	defer e.tenantMu.Unlock()
	<-limiter.slots
	limiter.refs--
	if limiter.refs == 0 && e.tenant[id] == limiter {
		delete(e.tenant, id)
	}
}
func (e *Engine) Route(ctx context.Context, req RouteRequest) (RouteResult, error) {
	e.requests.Add(1)
	if !e.started.Load() {
		return RouteResult{}, ErrUnavailable
	}
	if req.MobilityProfile != model.MobilityRoadVehicle {
		return RouteResult{}, ErrUnsupportedProfile
	}
	limiter, acquired := e.acquireTenant(req.TenantID)
	if !acquired {
		e.busy.Add(1)
		return RouteResult{}, ErrBusy
	}
	defer e.releaseTenant(req.TenantID, limiter)
	timeout := e.cfg.Timeout
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < timeout {
		timeout = time.Until(deadline)
	}
	routeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	result := make(chan routeResponse, 1)
	select {
	case e.jobs <- routeJob{routeCtx, req, result}:
	case <-routeCtx.Done():
		return RouteResult{}, classifyContext(routeCtx.Err())
	default:
		e.busy.Add(1)
		return RouteResult{}, ErrBusy
	}
	select {
	case r := <-result:
		return r.route, r.err
	case <-routeCtx.Done():
		return RouteResult{}, classifyContext(routeCtx.Err())
	}
}
func classifyContext(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeout
	}
	return err
}
func (e *Engine) route(ctx context.Context, req RouteRequest) (RouteResult, error) {
	if err := spatial.ValidatePosition(req.Origin); err != nil {
		return RouteResult{}, ErrOutsideRegion
	}
	if err := spatial.ValidatePosition(req.Destination); err != nil {
		return RouteResult{}, ErrOutsideRegion
	}
	from, err := e.graph.nodeIndex.Nearest(ctx, req.Origin)
	if err != nil {
		return RouteResult{}, ErrNoRoadNode
	}
	to, err := e.graph.nodeIndex.Nearest(ctx, req.Destination)
	if err != nil {
		return RouteResult{}, ErrNoRoadNode
	}
	snapshot := e.snapshots.Load()
	found, err := AStar(ctx, e.graph, snapshot, from.ID, to.ID, req.Policy, e.cfg.MaxExpansions)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return RouteResult{}, ErrTimeout
		}
		return RouteResult{}, err
	}
	out := routeResult(e.graph, snapshot, req.Policy, found)
	out.RouteID = auth.NewID()
	if len(out.Waypoints) == 0 {
		out.Waypoints = []model.Position{req.Origin}
	}
	return out, nil
}
func (e *Engine) SnapshotStore() *SnapshotStore { return e.snapshots }
func (e *Engine) Stats() EngineStats {
	e.tenantMu.Lock()
	activeTenants := len(e.tenant)
	e.tenantMu.Unlock()
	return EngineStats{Requests: e.requests.Load(), Busy: e.busy.Load(), QueueDepth: len(e.jobs), QueueCapacity: cap(e.jobs), ActiveTenants: activeTenants}
}
func (e *Engine) Close() {
	if e.started.CompareAndSwap(true, false) {
		close(e.closed)
		e.wg.Wait()
	}
}
