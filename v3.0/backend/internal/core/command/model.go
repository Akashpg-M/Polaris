package command

import (
	"encoding/json"
	"time"
)

const SchemaVersion = 1

type Status string

const (
	Pending      Status = "PENDING"
	Delivered    Status = "DELIVERED"
	Acknowledged Status = "ACKNOWLEDGED"
	Completed    Status = "COMPLETED"
	Failed       Status = "FAILED"
	Expired      Status = "EXPIRED"
	Cancelled    Status = "CANCELLED"
)

type Envelope struct {
	FrameType      string          `json:"frame_type"`
	CommandID      string          `json:"command_id"`
	CommandType    string          `json:"command_type"`
	SchemaVersion  int             `json:"schema_version"`
	TenantID       string          `json:"tenant_id"`
	DeviceID       string          `json:"device_id"`
	TaskID         string          `json:"task_id"`
	SequenceNumber int64           `json:"sequence_number"`
	CreatedAt      time.Time       `json:"created_at"`
	ExpiresAt      time.Time       `json:"expires_at"`
	CorrelationID  string          `json:"correlation_id"`
	CausationID    string          `json:"causation_id"`
	Payload        json.RawMessage `json:"payload"`
	// DeliveryObservation is volatile timing evidence added after the durable
	// command decision is read from the outbox. It is never persisted as part
	// of command identity and may differ across at-least-once delivery attempts.
	DeliveryObservation *DeliveryObservation `json:"delivery_observation,omitempty"`
}

type DeliveryObservation struct {
	RelayPublishedAt  time.Time `json:"relay_published_at"`
	GatewayReceivedAt time.Time `json:"gateway_received_at,omitempty"`
}

type Record struct {
	CommandID      string           `db:"command_id" json:"command_id"`
	TenantID       string           `db:"tenant_id" json:"tenant_id"`
	DeviceID       string           `db:"device_id" json:"device_id"`
	TaskID         string           `db:"task_id" json:"task_id"`
	CommandType    string           `db:"command_type" json:"command_type"`
	Payload        json.RawMessage  `db:"payload" json:"payload"`
	Status         string           `db:"status" json:"status"`
	SequenceNumber int64            `db:"sequence_number" json:"sequence_number"`
	CorrelationID  string           `db:"correlation_id" json:"correlation_id"`
	CausationID    string           `db:"causation_id" json:"causation_id"`
	AttemptCount   int              `db:"attempt_count" json:"attempt_count"`
	MaxAttempts    int              `db:"max_attempts" json:"max_attempts"`
	Version        int64            `db:"version" json:"version"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	AvailableAt    time.Time        `db:"available_at" json:"available_at"`
	SentAt         *time.Time       `db:"sent_at" json:"sent_at,omitempty"`
	AcknowledgedAt *time.Time       `db:"acknowledged_at" json:"acknowledged_at,omitempty"`
	CompletedAt    *time.Time       `db:"completed_at" json:"completed_at,omitempty"`
	ExpiresAt      time.Time        `db:"expires_at" json:"expires_at"`
	AckStatus      *string          `db:"ack_status" json:"ack_status,omitempty"`
	Result         *json.RawMessage `db:"result" json:"result,omitempty"`
	LastError      *string          `db:"last_error" json:"last_error,omitempty"`
}

func (r Record) Envelope() Envelope {
	return Envelope{FrameType: "COMMAND", CommandID: r.CommandID, CommandType: r.CommandType, SchemaVersion: SchemaVersion, TenantID: r.TenantID, DeviceID: r.DeviceID, TaskID: r.TaskID, SequenceNumber: r.SequenceNumber, CreatedAt: r.CreatedAt, ExpiresAt: r.ExpiresAt, CorrelationID: r.CorrelationID, CausationID: r.CausationID, Payload: r.Payload}
}

func (e Envelope) PartitionKey() string { return e.TenantID + ":" + e.DeviceID }

type Ack struct {
	FrameType      string    `json:"frame_type"`
	CommandID      string    `json:"command_id"`
	SequenceNumber int64     `json:"sequence_number"`
	Status         string    `json:"status"`
	ReceivedAt     time.Time `json:"received_at"`
	Reason         string    `json:"reason,omitempty"`
}

type Result struct {
	FrameType      string          `json:"frame_type"`
	CommandID      string          `json:"command_id"`
	SequenceNumber int64           `json:"sequence_number"`
	Status         string          `json:"status"`
	CompletedAt    time.Time       `json:"completed_at"`
	Result         json.RawMessage `json:"result,omitempty"`
	Reason         string          `json:"reason,omitempty"`
}

func ValidTransition(from, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case Pending:
		return to == Delivered || to == Cancelled || to == Expired || to == Failed
	case Delivered:
		return to == Acknowledged || to == Pending || to == Expired || to == Failed
	case Acknowledged:
		return to == Completed || to == Failed
	}
	return false
}

func IsTerminal(status Status) bool {
	return status == Completed || status == Failed || status == Expired || status == Cancelled
}

func RequiredCapability(commandType string) string {
	switch commandType {
	case "RELOCATE":
		return "receive_relocation_command"
	case "NAVIGATE", "RETURN_HOME":
		return "navigate"
	case "CAPTURE_IMAGE":
		return "capture_image"
	case "RUN_MODEL":
		return "run_model"
	case "THERMAL_SCAN", "START_SCAN":
		return "thermal_scan"
	case "STOP":
		return "receive_relocation_command"
	default:
		return ""
	}
}
