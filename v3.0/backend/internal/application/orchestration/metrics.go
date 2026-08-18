package orchestration

import "sync/atomic"

type Metrics struct {
	TasksCreated, TasksCompleted, TasksFailed                                 atomic.Int64
	CommandsCreated, CommandsDelivered, CommandsAcknowledged, CommandsRetried atomic.Int64
	CommandsExpired, CommandsFailed, ActiveConnections, PendingCommands       atomic.Int64
}

func NewMetrics() *Metrics { return &Metrics{} }

func (m *Metrics) Snapshot() map[string]int64 {
	return map[string]int64{
		"tasks_created_total": m.TasksCreated.Load(), "tasks_completed_total": m.TasksCompleted.Load(), "tasks_failed_total": m.TasksFailed.Load(),
		"commands_created_total": m.CommandsCreated.Load(), "commands_delivered_total": m.CommandsDelivered.Load(), "commands_acknowledged_total": m.CommandsAcknowledged.Load(),
		"commands_retried_total": m.CommandsRetried.Load(), "commands_expired_total": m.CommandsExpired.Load(), "commands_failed_total": m.CommandsFailed.Load(),
		"active_device_connections": m.ActiveConnections.Load(), "pending_commands": m.PendingCommands.Load(),
	}
}
