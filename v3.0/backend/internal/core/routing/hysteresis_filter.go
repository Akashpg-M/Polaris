package routing

import (
	"sync"
)

type HysteresisFilter struct {
	mu                sync.RWMutex
	confidenceGateway float64            // Minimum score required to bypass the filter (e.g., 0.75)
	smoothingFactor   float64            // Exponential Smoothing Alpha parameter (e.g., 0.3)
	stepBufferPercent float64            // Suppression boundary ratio (e.g., 0.04 for +/-4%)
	stabilizedWeights map[string]float64 // In-memory telemetry cache for tracking the smoothed line
}

// NewHysteresisFilter configures a mathematically isolated dampening matrix
func NewHysteresisFilter(minConfidence, alpha, stepBuffer float64) *HysteresisFilter {
	return &HysteresisFilter{
		confidenceGateway: minConfidence,
		smoothingFactor:   alpha,
		stepBufferPercent: stepBuffer,
		stabilizedWeights: make(map[string]float64),
	}
}

// Evaluate Engine Interceptor passes proposals through a moving window filter to avoid feedback loops
func (hf *HysteresisFilter) Evaluate(proposal WeightProposal, currentBaseWeight float64) (float64, bool) {
	// Gate 1: Confidence score validation boundary
	if proposal.ConfidenceScore < hf.confidenceGateway {
		return currentBaseWeight, false // Reject due to low statistical certainty
	}

	hf.mu.Lock()
	defer hf.mu.Unlock()

	lastStableWeight, exists := hf.stabilizedWeights[proposal.EdgeID]
	if !exists {
		lastStableWeight = currentBaseWeight
	}

	// Gate 2: Compute Exponential Moving Average (EMA) to smooth out traffic spikes
	// W_smoothed = (α * W_proposed) + ((1 - α) * W_historical)
	smoothedWeight := (hf.smoothingFactor * proposal.ProposedWeight) + ((1.0 - hf.smoothingFactor) * lastStableWeight)

	// Gate 3: Enforce strict step buffer change threshold to eliminate micro-vibrations
	deltaPercent := (smoothedWeight - lastStableWeight) / lastStableWeight
	
	if deltaPercent < hf.stepBufferPercent && deltaPercent > -hf.stepBufferPercent {
		// The calculated drift is within the safe dampening boundaries. Suppress graph mutation.
		return lastStableWeight, false
	}

	// Mutation threshold passed. Cache and commit the newly stabilized value.
	hf.stabilizedWeights[proposal.EdgeID] = smoothedWeight
	return smoothedWeight, true
}

// ClearCache drops tracked state metrics on command during system teardowns
func (hf *HysteresisFilter) ClearCache() {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	hf.stabilizedWeights = make(map[string]float64)
}