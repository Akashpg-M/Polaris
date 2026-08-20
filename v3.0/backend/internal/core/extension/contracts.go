package extension

import (
	"context"
	"encoding/json"
	"time"

	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	twincore "github.com/Akashpg-M/polaris/backend/internal/core/twin"
)

type ModuleState string

const (
	ModuleStarting ModuleState = "STARTING"
	ModuleReady    ModuleState = "READY"
	ModuleDegraded ModuleState = "DEGRADED"
	ModuleFailed   ModuleState = "FAILED"
	ModuleStopped  ModuleState = "STOPPED"
)

type ModuleStatus struct {
	State      ModuleState             `json:"state"`
	Message    string                  `json:"message,omitempty"`
	Components map[string]ModuleStatus `json:"components,omitempty"`
	Details    map[string]any          `json:"details,omitempty"`
}

type Module interface {
	Name() string
	Start(context.Context) error
	Ready(context.Context) ModuleStatus
	Close(context.Context) error
}

type CandidateRequest struct {
	TenantID             string
	EligibleDeviceIDs    []string
	RequiredCapabilities []string
	DeviceTypeIDs        []string
	ProjectIDs           []string
	Limit                int
	Context              map[string]any
	Timing               *CandidateTiming
}

// CandidateTiming is optional request-scoped instrumentation. Providers add
// only time spent in domain routing so Core can distinguish lookup/ranking
// from routing without coupling itself to a concrete capability module.
type CandidateTiming struct{ RoutingDuration time.Duration }

type Candidate struct {
	DeviceID    string         `json:"device_id"`
	DomainScore *float64       `json:"domain_score,omitempty"`
	Attributes  map[string]any `json:"attributes,omitempty"`
}

type CandidateProvider interface {
	Name() string
	Supports(CandidateRequest) bool
	Candidates(context.Context, CandidateRequest) ([]Candidate, error)
}

type PlanningRequest struct {
	Task       taskcore.Task
	DeviceTwin twincore.DeviceTwin
}

type ExecutionPlan struct {
	PlannerName   string          `json:"planner_name"`
	SchemaVersion uint32          `json:"schema_version"`
	CommandType   string          `json:"command_type"`
	Payload       json.RawMessage `json:"payload"`
	GeneratedAt   time.Time       `json:"generated_at"`
	ValidUntil    *time.Time      `json:"valid_until,omitempty"`
	Metadata      map[string]any  `json:"metadata,omitempty"`
}

type TaskPlanner interface {
	Name() string
	Supports(taskcore.Task) bool
	Plan(context.Context, PlanningRequest) (ExecutionPlan, error)
}
