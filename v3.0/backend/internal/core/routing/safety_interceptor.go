package routing

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/simulation"
)

var ErrRouteRejectedBySimulation = errors.New("route_cancelled_by_cellular_automata_risk_safety_gate")

// TargetAssetMailbox interface allows the interceptor to securely pass validated paths back into individual actor threads
type TargetAssetMailbox interface {
	Push(msg interface{}) error
}

type SafetyInterceptor struct {
	caRunline *simulation.CellularAutomataRunline
}

func NewSafetyInterceptor(ca *simulation.CellularAutomataRunline) *SafetyInterceptor {
	return &SafetyInterceptor{caRunline: ca}
}

// InterceptAndEnforce routes calculated coordinates through an asynchronous risk assessment gate before execution
func (si *SafetyInterceptor) InterceptAndEnforce(ctx context.Context, routeID string, targetEdge string, asset ActorRegistryInterface, commandMsg interface{}) error {
	slog.Info("[Safety Interceptor] Initializing real-time forward-projection simulation pass", "route_id", routeID, "target_edge", targetEdge)

	// 1. Construct localized spatial tracking parameters (scoping snapshot to target region)
	req := simulation.ValidationRequest{
		RouteID:    routeID,
		TargetEdge: targetEdge,
		H3Sectors:  []string{"chennai_core_grid_sector"}, // Region-limited allocation constraint
		Timestamp:  time.Now(),
	}

	// 2. Dispatch non-blocking Cellular Automata simulation thread sweep
	simFuture := si.caRunline.ValidateAsync(ctx, req)

	// Block only the current dispatch routine until the real-time execution deadline finishes
	result, ok := <-simFuture
	if !ok {
		return errors.New("simulation_channel_severed_unexpectedly")
	}

	slog.Info("[Safety Interceptor] Simulation iteration completed", 
		"route_id", routeID, 
		"risk_score", result.RiskScore, 
		"approved", result.AllowRoute,
		"compute_duration", result.ComputeTime,
	)

	// 3. Evaluate Binary Decision Gate Matrix
	if !result.AllowRoute {
		slog.Warn("[Safety Interceptor] ROUTE OVERRIDE ABORTED: Fluid simulation detects downstream gridlock risk thresholds breached", "route_id", routeID)
		return ErrRouteRejectedBySimulation
	}

	// 4. Safe Verification Clearance: Forward command cleanly down to the stateful Actor Mailbox
	slog.Info("[Safety Interceptor] Route safety verified. Emitting command packet down to asset boundary", "route_id", routeID)
	return asset.Push(commandMsg)
}

// Minimal Interface wrapper to maintain loose coupling across core layers
type ActorRegistryInterface interface {
	Push(msg interface{}) error
}