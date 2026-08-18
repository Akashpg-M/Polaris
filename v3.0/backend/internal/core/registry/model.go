package registry

import (
	"encoding/json"
	"time"
)

type Tenant struct {
	TenantID    string          `db:"tenant_id" json:"tenant_id"`
	DisplayName string          `db:"display_name" json:"display_name"`
	Status      string          `db:"status" json:"status"`
	Metadata    json.RawMessage `db:"metadata" json:"metadata"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}
type Project struct {
	ProjectID   string          `db:"project_id" json:"project_id"`
	TenantID    string          `db:"tenant_id" json:"tenant_id"`
	Name        string          `db:"name" json:"name"`
	Description *string         `db:"description" json:"description,omitempty"`
	Status      string          `db:"status" json:"status"`
	Metadata    json.RawMessage `db:"metadata" json:"metadata"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}
type Device struct {
	TenantID        string          `db:"tenant_id" json:"tenant_id"`
	DeviceID        string          `db:"device_id" json:"device_id"`
	ProjectID       *string         `db:"project_id" json:"project_id,omitempty"`
	DeviceTypeID    string          `db:"device_type_id" json:"device_type_id"`
	DisplayName     string          `db:"display_name" json:"display_name"`
	LifecycleStatus string          `db:"lifecycle_status" json:"lifecycle_status"`
	FirmwareVersion *string         `db:"firmware_version" json:"firmware_version,omitempty"`
	SoftwareVersion *string         `db:"software_version" json:"software_version,omitempty"`
	ModelVersion    *string         `db:"model_version" json:"model_version,omitempty"`
	Metadata        json.RawMessage `db:"metadata" json:"metadata"`
	RegisteredAt    time.Time       `db:"registered_at" json:"registered_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
	DeactivatedAt   *time.Time      `db:"deactivated_at" json:"deactivated_at,omitempty"`
}
type Capability struct {
	CapabilityID  string          `db:"capability_id" json:"capability_id"`
	DisplayName   string          `db:"display_name" json:"display_name"`
	Description   *string         `db:"description" json:"description,omitempty"`
	Configuration json.RawMessage `db:"configuration" json:"configuration"`
	Enabled       bool            `db:"enabled" json:"enabled"`
}
type CredentialMetadata struct {
	CredentialID string     `db:"credential_id" json:"credential_id"`
	TokenPrefix  string     `db:"token_prefix" json:"token_prefix"`
	Status       string     `db:"status" json:"status"`
	IssuedAt     time.Time  `db:"issued_at" json:"issued_at"`
	ExpiresAt    *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	LastUsedAt   *time.Time `db:"last_used_at" json:"last_used_at,omitempty"`
	RevokedAt    *time.Time `db:"revoked_at" json:"revoked_at,omitempty"`
}

func ValidTransition(from, to string) bool {
	switch from + ">" + to {
	case "REGISTERED>ACTIVE", "REGISTERED>DECOMMISSIONED", "ACTIVE>SUSPENDED", "ACTIVE>DECOMMISSIONED", "SUSPENDED>ACTIVE", "SUSPENDED>DECOMMISSIONED":
		return true
	}
	return false
}
