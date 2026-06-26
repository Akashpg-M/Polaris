package simulation

import (
	"time"
)

// ValidationRequest specifies the isolated spatial tracks to simulate
type ValidationRequest struct {
	RouteID    string    `json:"route_id"`
	TargetEdge string    `json:"target_edge"`
	H3Sectors  []string  `json:"h3_sectors"` // Limits snapshot boundaries
	Timestamp  time.Time `json:"timestamp"`
}

// SimulationResult returns the deterministic risk evaluation
type SimulationResult struct {
	RouteID     string  `json:"route_id"`
	AllowRoute  bool    `json:"allow_route"`  // Strict binary gate
	RiskScore   float64 `json:"risk_score"`   // 0.0 (Safe) to 1.0 (Gridlock)
	ComputeTime time.Duration `json:"compute_time"`
}