package actor

import (
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

// TelemetryMsg encapsulates incoming hardware device frames
type TelemetryMsg struct {
	Payload *pb.SpatialObject
}

// CommandMsg represents remote overriding directives sent from operators
type CommandMsg struct {
	Directive string
	Payload   map[string]interface{}
}