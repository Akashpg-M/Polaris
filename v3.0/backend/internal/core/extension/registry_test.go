package extension

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
)

type fixtureModule struct{ started, closed bool }

func (m *fixtureModule) Name() string                { return "fixture" }
func (m *fixtureModule) Start(context.Context) error { m.started = true; return nil }
func (m *fixtureModule) Ready(context.Context) ModuleStatus {
	if m.started && !m.closed {
		return ModuleStatus{State: ModuleReady}
	}
	return ModuleStatus{State: ModuleStopped}
}

func TestDefaultPlannerRejectsPolarisRequiredTask(t *testing.T) {
	requirements, _ := json.Marshal(taskcore.Requirements{PlanningMode: taskcore.PlanningPolarisRequired})
	_, err := (DefaultTaskPlanner{}).Plan(context.Background(), PlanningRequest{Task: taskcore.Task{TaskType: "NAVIGATE", Requirements: requirements, Target: []byte(`{"lat":13,"lon":80}`), ExpiresAt: time.Now().Add(time.Minute)}})
	if !errors.Is(err, ErrPlanningRequired) {
		t.Fatalf("expected planner requirement, got %v", err)
	}
}
func (m *fixtureModule) Close(context.Context) error { m.closed = true; return nil }
func TestExplicitModuleLifecycle(t *testing.T) {
	r := NewRegistry()
	m := &fixtureModule{}
	r.RegisterModule(m)
	if err := r.Start(context.Background()); err != nil || !m.started {
		t.Fatalf("module did not start: %v", err)
	}
	if r.Status(context.Background())["fixture"].State != ModuleReady {
		t.Fatal("readiness not exposed")
	}
	_ = r.Close(context.Background())
	if !m.closed {
		t.Fatal("module not closed")
	}
}
func TestMobilityDisabledLeavesGenericPlanning(t *testing.T) {
	r := NewRegistry()
	r.RegisterTaskPlanner(DefaultTaskPlanner{})
	for _, kind := range []string{"CAPTURE_IMAGE", "RUN_MODEL"} {
		task := taskcore.Task{TaskType: kind, Target: []byte(`{"fixture":true}`), ExpiresAt: time.Now().Add(time.Minute)}
		planner, err := r.TaskPlanner(task)
		if err != nil {
			t.Fatalf("%s unavailable with Mobility absent: %v", kind, err)
		}
		plan, err := planner.Plan(context.Background(), PlanningRequest{Task: task})
		if err != nil || plan.CommandType != kind {
			t.Fatalf("generic planning failed: %#v %v", plan, err)
		}
	}
	if _, err := r.TaskPlanner(taskcore.Task{TaskType: "NAVIGATE"}); err != nil {
		t.Fatalf("generic high-level NAVIGATE must remain available when Mobility is disabled: %v", err)
	}
}
