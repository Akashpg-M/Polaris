package actor

import (
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
)

// TelemetryMsg encapsulates incoming hardware device frames
type TelemetryMsg struct {
	Payload  *pb.SpatialObject
	Envelope *events.TelemetryEnvelope
}

// CommandMsg represents remote overriding directives sent from operators
type CommandMsg struct {
	Directive string
	Payload   map[string]interface{}
}
