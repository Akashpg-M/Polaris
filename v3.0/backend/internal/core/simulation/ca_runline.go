package simulation

import (
	"context"
	"math/rand"
	"time"
)

type CellularAutomataRunline struct {
	maxRiskThreshold float64
	computeDeadline  time.Duration
}

func NewCellularAutomataRunline(maxRisk float64, deadline time.Duration) *CellularAutomataRunline {
	return &CellularAutomataRunline{
		maxRiskThreshold: maxRisk,
		computeDeadline:  deadline,
	}
}

// ValidateAsync spins up an isolated forward-projection slice without blocking the routing path
func (ca *CellularAutomataRunline) ValidateAsync(ctx context.Context, req ValidationRequest) <-chan SimulationResult {
	out := make(chan SimulationResult, 1)

	go func() {
		defer close(out)
		startTime := time.Now()

		// Local timeout safety circuit bound to our real-time deadline parameter
		simCtx, cancel := context.WithTimeout(ctx, ca.computeDeadline)
		defer cancel()

		// Channel to catch the localized worker calculation results
		resultChan := make(chan float64, 1)

		go func() {
			// This simulates a discrete cellular automaton state matrix transition
			// In production, you map the target edge into a bitmapped cell sequence
			// representing vehicle slots (Nagel-Schreckenberg model matrix).
			
			// We execute a brief micro-step iteration simulating 20x real-world speed
			time.Sleep(3 * time.Millisecond) 
			
			// Generate a simulated density/risk value
			// In production, this calculates actual emergent bottlenecks or deadlocks
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			resultChan <- r.Float64()
		}()

		select {
		case <-simCtx.Done():
			// The simulation runline breached its real-time compute threshold constraint!
			// Fail-safe default strategy: Reject the risk proposal to preserve stability.
			out <- SimulationResult{
				RouteID:    req.RouteID,
				AllowRoute: false,
				RiskScore:  1.0,
				ComputeTime: time.Since(startTime),
			}
		case calculatedRisk := <-resultChan:
			isAllowed := calculatedRisk <= ca.maxRiskThreshold

			out <- SimulationResult{
				RouteID:    req.RouteID,
				AllowRoute: isAllowed,
				RiskScore:  calculatedRisk,
				ComputeTime: time.Since(startTime),
			}
		}
	}()

	return out
}