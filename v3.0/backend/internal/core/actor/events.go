package actor

// AssetStateChangedEvent tracks snapshot drifts without global memory locking
type AssetStateChangedEvent struct {
	AssetID        string  `json:"asset_id"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	EnergyPercent  int32   `json:"energy_percent"`
	VelocityMps    float64 `json:"velocity_mps"`
	HeadingDeg     float64 `json:"heading_deg"`
	TimestampMilli int64   `json:"timestamp_milli"`
}