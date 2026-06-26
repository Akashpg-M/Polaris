package repository

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/actor"
	"github.com/Akashpg-M/polaris/backend/internal/core/routing"
)

// GlobalGraphState interface isolates our projection writer from the Dijkstra read-path
type GlobalGraphState interface {
	GetEdgeBaseWeight(edgeID string) (float64, bool)
	UpdateEdgeDynamicWeight(edgeID string, newWeight float64)
}

type SpatialProjector struct {
	mu             sync.RWMutex
	graphState     GlobalGraphState
	filter         *routing.HysteresisFilter
	edgeDensities  map[string]int     // Tracking device counts on active street lines in RAM
	edgeLastUpdate map[string]time.Time
}

func NewSpatialProjector(graph GlobalGraphState, hf *routing.HysteresisFilter) *SpatialProjector {
	return &SpatialProjector{
		graphState:     graph,
		filter:         hf,
		edgeDensities:  make(map[string]int),
		edgeLastUpdate: make(map[string]time.Time),
	}
}

// ProjectActorEvent processes incoming stream data out of the event mesh loop
func (sp *SpatialProjector) ProjectActorEvent(event actor.AssetStateChangedEvent) {
	// 1. Snapping Logic Hook
	// In production, you match event.Lat/event.Lon to the closest OpenStreetMap EdgeID.
	// For this prototype loop integration, we derive a deterministic target edge string:
	targetEdgeID := "edge_osm_sector_chennai_main"

	sp.mu.Lock()
	sp.edgeDensities[targetEdgeID]++
	density := sp.edgeDensities[targetEdgeID]
	sp.mu.Unlock()

	// 2. Mock Predictive Inference Evaluation Block
	// Scale proposed weight proportionally based on active congestion asset count
	baseWeight, exists := sp.graphState.GetEdgeBaseWeight(targetEdgeID)
	if !exists {
		baseWeight = 10.0 // Default routing unit cost standard
	}
	
	simulatedCongestionMultiplier := 1.0 + (float64(density) * 0.15)
	
	proposal := routing.WeightProposal{
		EdgeID:          targetEdgeID,
		ProposedWeight:  baseWeight * simulatedCongestionMultiplier,
		ConfidenceScore: 0.88, // Passes our strict confidence threshold barrier
		Timestamp:       time.Now(),
	}

	// 3. Pass proposal directly through the Hysteresis Filtering Layer
	if stabilizedWeight, approved := sp.filter.Evaluate(proposal, baseWeight); approved {
		// Event-Driven Mutation: Update computed projection layer only when filter passes
		sp.graphState.UpdateEdgeDynamicWeight(targetEdgeID, stabilizedWeight)
		slog.Info("[Feedback Stabilizer] Weight mutation approved and applied", 
			"edge_id", targetEdgeID, 
			"stabilized_weight", stabilizedWeight,
			"current_density", density,
		)
	} else {
		// Suppressed minor vibration or low confidence
		slog.Debug("[Feedback Stabilizer] Proposal suppressed by hysteresis dampening thresholds", "edge_id", targetEdgeID)
	}
}