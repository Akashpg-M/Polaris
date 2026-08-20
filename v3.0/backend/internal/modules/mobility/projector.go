package mobility

import (
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	legacy "github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	mobilityspatial "github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

type TelemetryProjector struct{ Manager *mobilityspatial.Manager }

func profile(t pb.NodeType) model.MobilityProfile {
	switch t {
	case pb.NodeType_NODE_TYPE_DRONE:
		return model.MobilityAerialDrone
	case pb.NodeType_NODE_TYPE_ROBOT:
		return model.MobilityGroundRobot
	case pb.NodeType_NODE_TYPE_STATIC_SENSOR:
		return model.MobilityStatic
	default:
		return model.MobilityRoadVehicle
	}
}
func (p *TelemetryProjector) ApplyEnvelope(e *events.TelemetryEnvelope) legacy.Classification {
	speed := e.Payload.VelocityMps
	heading := e.Payload.HeadingDeg
	s := model.SpatialState{TenantID: e.TenantID, DeviceID: e.DeviceID, ReportedPosition: model.Position{Latitude: e.Payload.Lat, Longitude: e.Payload.Lon}, HeadingDegrees: &heading, SpeedMPS: &speed, MobilityProfile: profile(e.Payload.Type), ObservedAt: time.UnixMilli(e.ObservedAt).UTC(), BootID: e.DeviceBootID, BootStartedAt: time.UnixMilli(e.BootStartedAt).UTC(), SequenceNumber: e.SequenceNumber}
	err := p.Manager.Upsert(s)
	if err == nil {
		return legacy.Accepted
	}
	if err == mobilityspatial.ErrStaleVersion {
		return legacy.OutOfOrder
	}
	return legacy.OutOfOrder
}
