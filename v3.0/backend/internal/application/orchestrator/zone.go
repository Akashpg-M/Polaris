package orchestrator

import pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"

// Zone is a read-only predicted demand view retained for the operator API.
// Predictions never issue commands directly; a future advisor may translate
// them into ordinary durable Polaris tasks.
type Zone struct {
	ID             string      `json:"id"`
	Lat            float64     `json:"lat"`
	Lon            float64     `json:"lon"`
	RadiusKm       float64     `json:"radius_km"`
	RequiredAssets int         `json:"required_assets"`
	TargetClass    pb.NodeType `json:"target_class"`
	TenantID       string      `json:"tenant_id"`
}
