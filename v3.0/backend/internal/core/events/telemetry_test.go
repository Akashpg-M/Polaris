package events

import (
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"testing"
	"time"
)

func validFrame(now time.Time) *pb.SpatialObject {
	return &pb.SpatialObject{Id: "device-1", TenantId: "tenant-1", DeviceBootId: "boot-1", SequenceNumber: 1,
		BootStartedAt: now.Add(-time.Minute).UnixMilli(), ObservedAt: now.UnixMilli(), SchemaVersion: 1,
		Type: pb.NodeType_NODE_TYPE_DRONE, Status: pb.NodeStatus_NODE_STATUS_ACTIVE, Lat: 13, Lon: 80, EnergyPercent: 50}
}

func TestEnvelopeRoundTripAndStableIdentity(t *testing.T) {
	now := time.Now().UTC()
	p := validFrame(now)
	a := NewTelemetryEnvelope(p, now, "", "", "")
	b := NewTelemetryEnvelope(p, now.Add(time.Second), "", "", "")
	if a.EventID != b.EventID {
		t.Fatal("device tuple must produce a stable replay identity")
	}
	data, err := a.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.PartitionKey() != "tenant-1:device-1" {
		t.Fatalf("unexpected key %q", decoded.PartitionKey())
	}
}

func TestFrameValidationRejectsPermanentFailures(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name   string
		mutate func(*pb.SpatialObject)
	}{
		{"unsupported schema", func(p *pb.SpatialObject) { p.SchemaVersion = 99 }},
		{"invalid coordinate", func(p *pb.SpatialObject) { p.Lat = 91 }},
		{"missing identity", func(p *pb.SpatialObject) { p.Id = "" }},
		{"invalid battery", func(p *pb.SpatialObject) { p.EnergyPercent = 101 }},
		{"invalid velocity", func(p *pb.SpatialObject) { p.VelocityMps = -1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := validFrame(now)
			tc.mutate(p)
			if ValidateFrame(p, now) == nil {
				t.Fatal("expected rejection")
			}
		})
	}
}
