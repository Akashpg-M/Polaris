package twin

import (
	"encoding/json"
	"time"
)

// ComponentEnvelope is the storage contract for independently versioned twin
// components. Core code treats Payload as opaque; capability modules own it.
type ComponentEnvelope struct {
	Type           string          `json:"type"`
	SchemaVersion  uint32          `json:"schema_version"`
	ObservedAt     time.Time       `json:"observed_at"`
	BootID         string          `json:"boot_id"`
	SequenceNumber uint64          `json:"sequence_number"`
	Payload        json.RawMessage `json:"payload"`
}

type DeviceTwin struct {
	TenantID     string                       `json:"tenant_id"`
	DeviceID     string                       `json:"device_id"`
	Connectivity string                       `json:"connectivity"`
	Components   map[string]ComponentEnvelope `json:"components"`
}
