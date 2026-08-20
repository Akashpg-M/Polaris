package task

import (
	"encoding/json"
	"time"
)

type Status string
type PlanningMode string

const (
	Pending    Status = "PENDING"
	Assigning  Status = "ASSIGNING"
	Assigned   Status = "ASSIGNED"
	InProgress Status = "IN_PROGRESS"
	Completed  Status = "COMPLETED"
	Failed     Status = "FAILED"
	Cancelled  Status = "CANCELLED"
	Expired    Status = "EXPIRED"
)

const (
	PlanningDeviceLocal     PlanningMode = "DEVICE_LOCAL"
	PlanningPolarisRequired PlanningMode = "POLARIS_REQUIRED"
)

type Requirements struct {
	RequiredCapabilities []string     `json:"required_capabilities,omitempty"`
	MinimumBattery       int32        `json:"minimum_battery,omitempty"`
	AllowedDeviceTypes   []string     `json:"allowed_device_types,omitempty"`
	MaximumDistanceM     float64      `json:"max_distance_meters,omitempty"`
	ProjectID            string       `json:"project_id,omitempty"`
	PlanningMode         PlanningMode `json:"planning_mode,omitempty"`
	Custom               any          `json:"custom_constraints,omitempty"`
}

type Target struct {
	Latitude  *float64        `json:"lat,omitempty"`
	Longitude *float64        `json:"lon,omitempty"`
	H3Cell    string          `json:"h3_cell,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

type Task struct {
	TaskID           string          `db:"task_id" json:"task_id"`
	TenantID         string          `db:"tenant_id" json:"tenant_id"`
	ProjectID        *string         `db:"project_id" json:"project_id,omitempty"`
	TaskType         string          `db:"task_type" json:"task_type"`
	Status           string          `db:"status" json:"status"`
	Priority         string          `db:"priority" json:"priority"`
	Requirements     json.RawMessage `db:"requirements" json:"requirements"`
	Target           json.RawMessage `db:"target" json:"target"`
	AssignedDeviceID *string         `db:"assigned_device_id" json:"assigned_device_id,omitempty"`
	CorrelationID    string          `db:"correlation_id" json:"correlation_id"`
	CreatedBy        string          `db:"created_by" json:"created_by"`
	Version          int64           `db:"version" json:"version"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at" json:"updated_at"`
	AssignedAt       *time.Time      `db:"assigned_at" json:"assigned_at,omitempty"`
	StartedAt        *time.Time      `db:"started_at" json:"started_at,omitempty"`
	CompletedAt      *time.Time      `db:"completed_at" json:"completed_at,omitempty"`
	FailedAt         *time.Time      `db:"failed_at" json:"failed_at,omitempty"`
	ExpiresAt        time.Time       `db:"expires_at" json:"expires_at"`
	FailureReason    *string         `db:"failure_reason" json:"failure_reason,omitempty"`
}

func ValidTransition(from, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case Pending:
		return to == Assigning || to == Cancelled || to == Expired || to == Failed
	case Assigning:
		return to == Assigned || to == Pending || to == Cancelled || to == Expired || to == Failed
	case Assigned:
		return to == InProgress || to == Cancelled || to == Expired || to == Failed
	case InProgress:
		return to == Completed || to == Failed || to == Expired
	}
	return false
}

func IsTerminal(status Status) bool {
	return status == Completed || status == Failed || status == Cancelled || status == Expired
}
