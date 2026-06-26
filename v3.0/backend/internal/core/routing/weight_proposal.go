package routing

import (
	"time"
)

// WeightProposal represents an isolated, non-binding prediction output
type WeightProposal struct {
	EdgeID          string    `json:"edge_id"`
	ProposedWeight  float64   `json:"proposed_weight"`
	ConfidenceScore float64   `json:"confidence_score"` // Explicit Gating Range: 0.0 to 1.0
	Timestamp       time.Time `json:"timestamp"`
}