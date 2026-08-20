package extension

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
)

type DefaultTaskPlanner struct{}

func (DefaultTaskPlanner) Name() string                { return "core.default/v1" }
func (DefaultTaskPlanner) Supports(taskcore.Task) bool { return true }
func (DefaultTaskPlanner) Plan(_ context.Context, req PlanningRequest) (ExecutionPlan, error) {
	var requirements taskcore.Requirements
	if len(req.Task.Requirements) > 0 && json.Unmarshal(req.Task.Requirements, &requirements) != nil {
		return ExecutionPlan{}, errors.New("invalid task requirements")
	}
	if requirements.PlanningMode == taskcore.PlanningPolarisRequired {
		return ExecutionPlan{}, ErrPlanningRequired
	}
	if len(req.Task.Target) == 0 {
		return ExecutionPlan{}, errors.New("task target is empty")
	}
	now := time.Now().UTC()
	validUntil := req.Task.ExpiresAt
	return ExecutionPlan{PlannerName: "core.default/v1", SchemaVersion: 1, CommandType: req.Task.TaskType, Payload: req.Task.Target, GeneratedAt: now, ValidUntil: &validUntil}, nil
}
