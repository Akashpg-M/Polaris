package model

import "time"

type MobilityProfile string

const (
	MobilityRoadVehicle MobilityProfile = "ROAD_VEHICLE"
	MobilityGroundRobot MobilityProfile = "GROUND_ROBOT"
	MobilityAerialDrone MobilityProfile = "AERIAL_DRONE"
	MobilityStatic      MobilityProfile = "STATIC"
)

type Position struct {
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	AltitudeM *float64 `json:"altitude_meters,omitempty"`
}

type SpatialState struct {
	TenantID         string          `json:"tenant_id"`
	DeviceID         string          `json:"device_id"`
	Position         Position        `json:"position"`
	ReportedPosition Position        `json:"reported_position"`
	HeadingDegrees   *float64        `json:"heading_degrees,omitempty"`
	SpeedMPS         *float64        `json:"speed_mps,omitempty"`
	MobilityProfile  MobilityProfile `json:"mobility_profile"`
	H3Cell           uint64          `json:"h3_cell"`
	ObservedAt       time.Time       `json:"observed_at"`
	IndexedAt        time.Time       `json:"indexed_at"`
	BootID           string          `json:"boot_id"`
	BootStartedAt    time.Time       `json:"boot_started_at"`
	SequenceNumber   uint64          `json:"sequence_number"`
	Quality          MobilityQuality `json:"quality"`
}

type MobilityQuality struct {
	Valid      bool     `json:"valid"`
	Confidence float64  `json:"confidence"`
	Anomalies  []string `json:"anomalies,omitempty"`
}

func (s SpatialState) NewerThan(current SpatialState) bool {
	if current.DeviceID == "" {
		return true
	}
	if s.BootID == current.BootID {
		return s.SequenceNumber > current.SequenceNumber
	}
	return s.BootStartedAt.After(current.BootStartedAt)
}
