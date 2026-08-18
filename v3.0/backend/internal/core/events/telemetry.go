package events

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
)

const (
	TelemetryEventType   = "polaris.telemetry.observed"
	CurrentSchemaVersion = uint32(1)
	GatewayProducer      = "polaris-gateway"
	MaxFrameBytes        = int64(64 * 1024)
)

// ':' is reserved as the unambiguous tenant/device partition-key separator.
var identityPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

// TelemetryEnvelope is the canonical Kafka value. Device-owned facts stay in
// Payload; trusted platform metadata is added at the gateway boundary.
type TelemetryEnvelope struct {
	EventID        string            `json:"event_id"`
	EventType      string            `json:"event_type"`
	SchemaVersion  uint32            `json:"schema_version"`
	TenantID       string            `json:"tenant_id"`
	DeviceID       string            `json:"device_id"`
	DeviceBootID   string            `json:"device_boot_id"`
	SequenceNumber uint64            `json:"sequence_number"`
	BootStartedAt  int64             `json:"boot_started_at"`
	ObservedAt     int64             `json:"observed_at"`
	IngestedAt     int64             `json:"ingested_at"`
	CorrelationID  string            `json:"correlation_id"`
	CausationID    string            `json:"causation_id,omitempty"`
	Producer       string            `json:"producer"`
	Traceparent    string            `json:"traceparent,omitempty"`
	Payload        *pb.SpatialObject `json:"payload"`
}

func NewTelemetryEnvelope(p *pb.SpatialObject, ingestedAt time.Time, correlationID, causationID, traceparent string) *TelemetryEnvelope {
	observedAt := p.ObservedAt
	if observedAt == 0 {
		observedAt = p.Timestamp
	}
	identity := fmt.Sprintf("%s\x00%s\x00%s\x00%d", p.TenantId, p.Id, p.DeviceBootId, p.SequenceNumber)
	sum := sha256.Sum256([]byte(identity))
	eventID := hex.EncodeToString(sum[:])
	if correlationID == "" {
		correlationID = eventID
	}
	return &TelemetryEnvelope{
		EventID: eventID, EventType: TelemetryEventType, SchemaVersion: p.SchemaVersion,
		TenantID: p.TenantId, DeviceID: p.Id, DeviceBootID: p.DeviceBootId,
		SequenceNumber: p.SequenceNumber, BootStartedAt: p.BootStartedAt,
		ObservedAt: observedAt, IngestedAt: ingestedAt.UnixMilli(),
		CorrelationID: correlationID, CausationID: causationID,
		Producer: GatewayProducer, Traceparent: traceparent, Payload: p,
	}
}

func (e *TelemetryEnvelope) PartitionKey() string     { return e.TenantID + ":" + e.DeviceID }
func (e *TelemetryEnvelope) Marshal() ([]byte, error) { return json.Marshal(e) }
func Unmarshal(data []byte) (*TelemetryEnvelope, error) {
	var e TelemetryEnvelope
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("malformed envelope: %w", err)
	}
	if err := ValidateEnvelope(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func ValidateFrame(p *pb.SpatialObject, now time.Time) error {
	if p == nil {
		return errors.New("missing payload")
	}
	if !identityPattern.MatchString(p.TenantId) {
		return errors.New("invalid tenant_id")
	}
	if !identityPattern.MatchString(p.Id) {
		return errors.New("invalid device_id")
	}
	if !identityPattern.MatchString(p.DeviceBootId) {
		return errors.New("invalid device_boot_id")
	}
	if p.SequenceNumber == 0 {
		return errors.New("sequence_number must be positive")
	}
	if p.SequenceNumber > math.MaxInt64 {
		return errors.New("sequence_number exceeds supported range")
	}
	if p.SchemaVersion != CurrentSchemaVersion {
		return fmt.Errorf("unsupported schema_version %d", p.SchemaVersion)
	}
	if math.IsNaN(p.Lat) || math.IsInf(p.Lat, 0) || p.Lat < -90 || p.Lat > 90 {
		return errors.New("invalid latitude")
	}
	if math.IsNaN(p.Lon) || math.IsInf(p.Lon, 0) || p.Lon < -180 || p.Lon > 180 {
		return errors.New("invalid longitude")
	}
	if p.EnergyPercent < 0 || p.EnergyPercent > 100 {
		return errors.New("battery outside 0..100")
	}
	if math.IsNaN(p.VelocityMps) || math.IsInf(p.VelocityMps, 0) || p.VelocityMps < 0 || p.VelocityMps > 250 {
		return errors.New("invalid velocity")
	}
	if p.Type <= pb.NodeType_NODE_TYPE_UNKNOWN || p.Type > pb.NodeType_NODE_TYPE_STATIC_SENSOR {
		return errors.New("invalid device type")
	}
	observed := p.ObservedAt
	if observed == 0 {
		observed = p.Timestamp
	}
	if observed <= 0 {
		return errors.New("missing observed_at")
	}
	if p.BootStartedAt <= 0 || p.BootStartedAt > observed {
		return errors.New("invalid boot_started_at")
	}
	observedTime := time.UnixMilli(observed)
	if observedTime.Before(now.Add(-24*time.Hour)) || observedTime.After(now.Add(5*time.Minute)) {
		return errors.New("observation timestamp outside allowed window")
	}
	return nil
}

func ValidateEnvelope(e *TelemetryEnvelope) error {
	if e == nil || e.Payload == nil {
		return errors.New("missing envelope payload")
	}
	if e.EventID == "" || e.EventType != TelemetryEventType || e.Producer == "" || e.IngestedAt <= 0 || e.CorrelationID == "" {
		return errors.New("missing platform envelope metadata")
	}
	if e.SchemaVersion != CurrentSchemaVersion {
		return fmt.Errorf("unsupported schema_version %d", e.SchemaVersion)
	}
	if err := ValidateFrame(e.Payload, time.UnixMilli(e.IngestedAt)); err != nil {
		return err
	}
	if e.TenantID != e.Payload.TenantId || e.DeviceID != e.Payload.Id || e.DeviceBootID != e.Payload.DeviceBootId || e.SequenceNumber != e.Payload.SequenceNumber {
		return errors.New("envelope/payload identity mismatch")
	}
	if e.BootStartedAt != e.Payload.BootStartedAt {
		return errors.New("envelope/payload boot timestamp mismatch")
	}
	payloadObservedAt := e.Payload.ObservedAt
	if payloadObservedAt == 0 {
		payloadObservedAt = e.Payload.Timestamp
	}
	if e.ObservedAt != payloadObservedAt || e.SchemaVersion != e.Payload.SchemaVersion {
		return errors.New("envelope/payload schema or observation mismatch")
	}
	if strings.TrimSpace(e.Traceparent) != "" && !strings.HasPrefix(strings.ToLower(e.Traceparent), "00-") {
		return errors.New("invalid traceparent")
	}
	return nil
}
