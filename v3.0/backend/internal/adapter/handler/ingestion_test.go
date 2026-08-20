package handler

import (
	"testing"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

func TestRegisteredDeviceTypeConstrainsTelemetryProfile(t *testing.T) {
	tests := []struct {
		deviceType string
		nodeType   pb.NodeType
		allowed    bool
	}{
		{"delivery_drone", pb.NodeType_NODE_TYPE_DRONE, true},
		{"delivery_drone", pb.NodeType_NODE_TYPE_SEDAN, false},
		{"connected_vehicle", pb.NodeType_NODE_TYPE_BIKE, true},
		{"connected_vehicle", pb.NodeType_NODE_TYPE_SUV, true},
		{"static_camera", pb.NodeType_NODE_TYPE_STATIC_SENSOR, true},
		{"compute_node", pb.NodeType_NODE_TYPE_STATIC_SENSOR, false},
		{"unknown", pb.NodeType_NODE_TYPE_DRONE, false},
	}
	for _, test := range tests {
		if got := telemetryTypeAllowed(test.deviceType, test.nodeType); got != test.allowed {
			t.Errorf("%s/%s allowed=%v want %v", test.deviceType, test.nodeType, got, test.allowed)
		}
	}
}
