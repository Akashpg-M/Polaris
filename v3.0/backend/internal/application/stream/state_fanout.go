package stream

import (
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
)

type StateFanout struct {
	Primary     stateApplier
	Projections []stateApplier
}

func (f *StateFanout) ApplyEnvelope(e *events.TelemetryEnvelope) spatial.Classification {
	result := f.Primary.ApplyEnvelope(e)
	for _, p := range f.Projections {
		p.ApplyEnvelope(e)
	}
	return result
}
