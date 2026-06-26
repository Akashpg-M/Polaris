package orchestrator

import (
	"context"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

// StaticZoneStrategy simulates a database table of logistics hubs or smart-city sectors
type StaticZoneStrategy struct {}

func (s *StaticZoneStrategy) GetTargetZones(ctx context.Context) []Zone {
	return []Zone{
		{
			ID:             "Hub-Guindy",
			Lat:            13.0067,
			Lon:            80.2206,
			RadiusKm:       5.0,
			RequiredAssets: 3,
			TargetClass:    pb.NodeType_NODE_TYPE_DRONE,
      TenantID:       "alpha_logistics",
		},
		{
			ID:             "Hub-Adyar",
			Lat:            13.0012,
			Lon:            80.2565,
			RadiusKm:       3.0,
			RequiredAssets: 2,
			TargetClass:    pb.NodeType_NODE_TYPE_DRONE,
      TenantID:       "alpha_logistics",
		},
	}
}