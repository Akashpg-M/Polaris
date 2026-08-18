package spatial

import (
	"testing"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
)

func testEnvelope(boot string, bootStarted int64, sequence uint64, lat float64) *events.TelemetryEnvelope {
	now := time.Now().UTC()
	p := &pb.SpatialObject{Id: "device-1", TenantId: "tenant-1", DeviceBootId: boot, SequenceNumber: sequence,
		BootStartedAt: bootStarted, ObservedAt: now.UnixMilli(), SchemaVersion: 1, Type: pb.NodeType_NODE_TYPE_DRONE,
		Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: lat, Lon: 80, EnergyPercent: 80}
	return events.NewTelemetryEnvelope(p, now, "", "", "")
}

func TestLatestStateClassifications(t *testing.T) {
	e := NewEngine()
	started := time.Now().Add(-time.Minute).UnixMilli()
	cases := []struct {
		name     string
		envelope *events.TelemetryEnvelope
		want     Classification
	}{
		{"first", testEnvelope("boot-a", started, 1, 13.0), Accepted},
		{"duplicate", testEnvelope("boot-a", started, 1, 99.0), Duplicate},
		{"newer sequence", testEnvelope("boot-a", started, 3, 13.2), Accepted},
		{"out of order", testEnvelope("boot-a", started, 2, 99.0), OutOfOrder},
		{"boot conflict", testEnvelope("boot-conflict", started, 1, 99.0), BootConflict},
		{"new boot", testEnvelope("boot-b", started+1000, 1, 13.3), NewBoot},
		{"retired boot", testEnvelope("boot-a", started, 4, 99.0), RetiredBoot},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := e.ApplyEnvelope(tc.envelope); got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
		})
	}
	results := e.FindNearest("tenant-1", 13.3, 80, 1, pb.NodeType_NODE_TYPE_DRONE)
	if len(results) != 1 || results[0].NodeID != "device-1" {
		t.Fatalf("latest accepted state missing: %#v", results)
	}
}
