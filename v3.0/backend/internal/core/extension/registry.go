package extension

import (
	"context"
	"errors"
	"fmt"
	"sync"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
)

var (
	ErrNoCandidateProvider = errors.New("no candidate provider supports request")
	ErrNoTaskPlanner       = errors.New("no task planner supports task")
	ErrPlanningUnsupported = errors.New("planner does not support selected device")
	ErrPlanningRequired    = errors.New("compatible Polaris planner required")
)

// Registry is populated explicitly by application composition. It deliberately
// has no init hooks or dynamic plugin loading.
type Registry struct {
	mu        sync.RWMutex
	modules   []Module
	providers []CandidateProvider
	planners  []TaskPlanner
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) RegisterModule(m Module) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.modules = append(r.modules, m)
}
func (r *Registry) RegisterCandidateProvider(p CandidateProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = append(r.providers, p)
}
func (r *Registry) RegisterTaskPlanner(p TaskPlanner) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.planners = append(r.planners, p)
}

func (r *Registry) CandidateProvider(req CandidateRequest) (CandidateProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.providers {
		if p.Supports(req) {
			return p, nil
		}
	}
	return nil, ErrNoCandidateProvider
}

func (r *Registry) TaskPlanner(v taskcore.Task) (TaskPlanner, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.planners {
		if p.Supports(v) {
			return p, nil
		}
	}
	return nil, ErrNoTaskPlanner
}

func (r *Registry) TaskPlanners(v taskcore.Task) []TaskPlanner {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []TaskPlanner{}
	for _, p := range r.planners {
		if p.Supports(v) {
			out = append(out, p)
		}
	}
	return out
}

func (r *Registry) Start(ctx context.Context) error {
	r.mu.RLock()
	modules := append([]Module(nil), r.modules...)
	r.mu.RUnlock()
	started := make([]Module, 0, len(modules))
	for _, m := range modules {
		if err := m.Start(ctx); err != nil {
			for i := len(started) - 1; i >= 0; i-- {
				_ = started[i].Close(context.Background())
			}
			return fmt.Errorf("start module %s: %w", m.Name(), err)
		}
		started = append(started, m)
	}
	return nil
}

func (r *Registry) Close(ctx context.Context) error {
	r.mu.RLock()
	modules := append([]Module(nil), r.modules...)
	r.mu.RUnlock()
	var joined error
	for i := len(modules) - 1; i >= 0; i-- {
		joined = errors.Join(joined, modules[i].Close(ctx))
	}
	return joined
}

func (r *Registry) Status(ctx context.Context) map[string]ModuleStatus {
	r.mu.RLock()
	modules := append([]Module(nil), r.modules...)
	r.mu.RUnlock()
	out := make(map[string]ModuleStatus, len(modules))
	for _, m := range modules {
		out[m.Name()] = m.Ready(ctx)
	}
	return out
}
