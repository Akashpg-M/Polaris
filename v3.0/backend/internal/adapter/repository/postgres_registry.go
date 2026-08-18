package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/registry"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrNotFound          = errors.New("registry resource not found")
	ErrConflict          = errors.New("registry resource already exists")
	ErrForbidden         = errors.New("registry operation forbidden")
	ErrInvalidTransition = errors.New("invalid lifecycle transition")
)

type RegistryStore struct{ DB *sqlx.DB }

func NewRegistryStore(postgresURL string) (*RegistryStore, error) {
	db, err := sqlx.Connect("postgres", postgresURL)
	if err != nil {
		return nil, err
	}
	return &RegistryStore{DB: db}, nil
}
func (s *RegistryStore) Close() error { return s.DB.Close() }
func jsonOrEmpty(v json.RawMessage) []byte {
	if len(v) == 0 {
		return []byte(`{}`)
	}
	return v
}
func nullableJSON(v []byte) interface{} {
	if len(v) == 0 {
		return nil
	}
	return v
}
func mapPQ(err error) error {
	var p *pq.Error
	if errors.As(err, &p) && p.Code.Class() == "23" {
		if p.Code == "23505" {
			return ErrConflict
		}
		return fmt.Errorf("registry constraint: %w", err)
	}
	return err
}

func insertAuditOutbox(ctx context.Context, tx *sqlx.Tx, tenant, actorID, action, resourceType, resourceID, requestID, eventType string, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	eventID := auth.NewID()
	if _, err = tx.ExecContext(ctx, `INSERT INTO audit_events(audit_id,tenant_id,actor_type,actor_id,action,resource_type,resource_id,request_id,outcome) VALUES($1,NULLIF($2,''),'OPERATOR',$3,$4,$5,$6,$7,'SUCCESS')`, auth.NewID(), tenant, actorID, action, resourceType, resourceID, requestID); err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO outbox_events(outbox_id,aggregate_type,aggregate_id,tenant_id,event_id,event_type,schema_version,payload,status) VALUES($1,$2,$3,$4,$5,$6,1,$7,'PENDING')`, auth.NewID(), resourceType, resourceID, tenant, eventID, eventType, b)
	return err
}

func (s *RegistryStore) BootstrapPlatformAdmin(ctx context.Context, raw string) error {
	if raw == "" {
		return nil
	}
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO operator_api_keys(api_key_id,name,token_prefix,token_hash,role,status) VALUES($1,'development bootstrap',$2,$3,'PLATFORM_ADMIN','ACTIVE') ON CONFLICT(token_prefix) DO UPDATE SET token_hash=EXCLUDED.token_hash,status='ACTIVE',revoked_at=NULL`, auth.NewID(), prefix, auth.Hash(raw))
	return err
}

func (s *RegistryStore) ResolveOperator(ctx context.Context, raw string) (auth.OperatorPrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	var row struct {
		ID      string         `db:"api_key_id"`
		Tenant  sql.NullString `db:"tenant_id"`
		Role    string         `db:"role"`
		Hash    []byte         `db:"token_hash"`
		Status  string         `db:"status"`
		Expires sql.NullTime   `db:"expires_at"`
	}
	err = s.DB.GetContext(ctx, &row, `SELECT api_key_id,tenant_id,role,token_hash,status,expires_at FROM operator_api_keys WHERE token_prefix=$1`, prefix)
	if err != nil || row.Status != "ACTIVE" || (row.Expires.Valid && time.Now().After(row.Expires.Time)) || !auth.Verify(raw, row.Hash) {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	_, _ = s.DB.ExecContext(ctx, `UPDATE operator_api_keys SET last_used_at=NOW() WHERE api_key_id=$1`, row.ID)
	return auth.OperatorPrincipal{APIKeyID: row.ID, TenantID: row.Tenant.String, Role: auth.Role(row.Role)}, nil
}

func (s *RegistryStore) ResolveDevice(ctx context.Context, raw string) (auth.DevicePrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	var row struct {
		CredentialID, TenantID, DeviceID, DeviceType, ProjectID, CredentialStatus, TenantStatus, DeviceStatus string
		Hash                                                                                                  []byte
		Expires                                                                                               sql.NullTime
	}
	err = s.DB.QueryRowxContext(ctx, `SELECT c.credential_id,c.tenant_id,c.device_id,d.device_type_id,COALESCE(d.project_id::text,''),c.status,c.token_hash,c.expires_at,t.status,d.lifecycle_status FROM device_credentials c JOIN devices d USING(tenant_id,device_id) JOIN tenants t USING(tenant_id) WHERE c.token_prefix=$1`, prefix).Scan(&row.CredentialID, &row.TenantID, &row.DeviceID, &row.DeviceType, &row.ProjectID, &row.CredentialStatus, &row.Hash, &row.Expires, &row.TenantStatus, &row.DeviceStatus)
	if err != nil || row.CredentialStatus != "ACTIVE" || row.TenantStatus != "ACTIVE" || row.DeviceStatus != "ACTIVE" || (row.Expires.Valid && time.Now().After(row.Expires.Time)) || !auth.Verify(raw, row.Hash) {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	_, _ = s.DB.ExecContext(ctx, `UPDATE device_credentials SET last_used_at=NOW() WHERE credential_id=$1`, row.CredentialID)
	return auth.DevicePrincipal{TenantID: row.TenantID, DeviceID: row.DeviceID, CredentialID: row.CredentialID, DeviceType: row.DeviceType, ProjectID: row.ProjectID}, nil
}
func (s *RegistryStore) RevalidateDevice(ctx context.Context, p auth.DevicePrincipal) error {
	var ok bool
	err := s.DB.GetContext(ctx, &ok, `SELECT c.status='ACTIVE' AND (c.expires_at IS NULL OR c.expires_at>NOW()) AND d.lifecycle_status='ACTIVE' AND t.status='ACTIVE' FROM device_credentials c JOIN devices d USING(tenant_id,device_id) JOIN tenants t USING(tenant_id) WHERE c.credential_id=$1 AND c.tenant_id=$2 AND c.device_id=$3`, p.CredentialID, p.TenantID, p.DeviceID)
	if err != nil || !ok {
		return auth.ErrInvalidCredential
	}
	return nil
}

func (s *RegistryStore) CreateTenant(ctx context.Context, t registry.Tenant, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO tenants(tenant_id,display_name,status,metadata) VALUES($1,$2,$3,$4)`, t.TenantID, t.DisplayName, t.Status, t.Metadata)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, t.TenantID, actor, "TENANT_CREATED", "tenant", t.TenantID, requestID, "tenant.registered.v1", t); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) GetTenant(ctx context.Context, id string) (registry.Tenant, error) {
	var v registry.Tenant
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM tenants WHERE tenant_id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}
func (s *RegistryStore) SetTenantStatus(ctx context.Context, id, status, actor, requestID string) error {
	if status != "ACTIVE" && status != "SUSPENDED" && status != "DEACTIVATED" {
		return ErrInvalidTransition
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE tenants SET status=$2,updated_at=NOW() WHERE tenant_id=$1`, id, status)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, id, actor, "TENANT_"+status, "tenant", id, requestID, "tenant.lifecycle.changed.v1", map[string]string{"tenant_id": id, "status": status}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) CreateProject(ctx context.Context, p registry.Project, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO projects(project_id,tenant_id,name,description,status,metadata) VALUES($1,$2,$3,$4,$5,$6)`, p.ProjectID, p.TenantID, p.Name, p.Description, p.Status, p.Metadata)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, p.TenantID, actor, "PROJECT_CREATED", "project", p.ProjectID, requestID, "project.registered.v1", p); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) ListProjects(ctx context.Context, tenant string) ([]registry.Project, error) {
	v := []registry.Project{}
	err := s.DB.SelectContext(ctx, &v, `SELECT * FROM projects WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT 100`, tenant)
	return v, err
}
func (s *RegistryStore) GetProject(ctx context.Context, tenant, id string) (registry.Project, error) {
	var v registry.Project
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM projects WHERE tenant_id=$1 AND project_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}
func (s *RegistryStore) UpdateProject(ctx context.Context, tenant, id string, name, description, status *string, metadata []byte, actor, requestID string) error {
	if status != nil && *status != "ACTIVE" && *status != "ARCHIVED" {
		return ErrInvalidTransition
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE projects SET name=COALESCE($3,name),description=COALESCE($4,description),status=COALESCE($5,status),metadata=CASE WHEN $6::jsonb IS NULL THEN metadata ELSE $6::jsonb END,updated_at=NOW() WHERE tenant_id=$1 AND project_id=$2`, tenant, id, name, description, status, nullableJSON(metadata))
	if err != nil {
		return mapPQ(err)
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "PROJECT_UPDATED", "project", id, requestID, "project.updated.v1", map[string]string{"project_id": id}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) CreateDevice(ctx context.Context, d registry.Device, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO devices(tenant_id,device_id,project_id,device_type_id,display_name,lifecycle_status,firmware_version,software_version,model_version,metadata) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, d.TenantID, d.DeviceID, d.ProjectID, d.DeviceTypeID, d.DisplayName, d.LifecycleStatus, d.FirmwareVersion, d.SoftwareVersion, d.ModelVersion, d.Metadata)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, d.TenantID, actor, "DEVICE_REGISTERED", "device", d.DeviceID, requestID, "device.registered.v1", d); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) GetDevice(ctx context.Context, tenant, id string) (registry.Device, error) {
	var v registry.Device
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM devices WHERE tenant_id=$1 AND device_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}
func (s *RegistryStore) UpdateDevice(ctx context.Context, tenant, id string, displayName, firmware, software, model *string, metadata []byte, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE devices SET display_name=COALESCE($3,display_name),firmware_version=COALESCE($4,firmware_version),software_version=COALESCE($5,software_version),model_version=COALESCE($6,model_version),metadata=CASE WHEN $7::jsonb IS NULL THEN metadata ELSE $7::jsonb END,updated_at=NOW() WHERE tenant_id=$1 AND device_id=$2 AND lifecycle_status<>'DECOMMISSIONED'`, tenant, id, displayName, firmware, software, model, nullableJSON(metadata))
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "DEVICE_UPDATED", "device", id, requestID, "device.updated.v1", map[string]string{"device_id": id}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) ListDevices(ctx context.Context, tenant string, limit int, cursor, projectID, deviceType, lifecycle, capability string) ([]registry.Device, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	args := []interface{}{tenant, limit}
	q := `SELECT d.* FROM devices d WHERE d.tenant_id=$1`
	add := func(clause string, v interface{}) { args = append(args, v); q += fmt.Sprintf(clause, len(args)) }
	if cursor != "" {
		add(` AND d.device_id>$%d`, cursor)
	}
	if projectID != "" {
		add(` AND d.project_id=$%d`, projectID)
	}
	if deviceType != "" {
		add(` AND d.device_type_id=$%d`, deviceType)
	}
	if lifecycle != "" {
		add(` AND d.lifecycle_status=$%d`, lifecycle)
	}
	if capability != "" {
		add(` AND EXISTS(SELECT 1 FROM device_capabilities dc WHERE dc.tenant_id=d.tenant_id AND dc.device_id=d.device_id AND dc.capability_id=$%d AND dc.enabled)`, capability)
	}
	q += ` ORDER BY d.device_id LIMIT $2`
	v := []registry.Device{}
	err := s.DB.SelectContext(ctx, &v, q, args...)
	return v, err
}

func (s *RegistryStore) SetLifecycle(ctx context.Context, tenant, id, next, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var current string
	if err = tx.GetContext(ctx, &current, `SELECT lifecycle_status FROM devices WHERE tenant_id=$1 AND device_id=$2 FOR UPDATE`, tenant, id); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if !registry.ValidTransition(current, next) {
		return ErrInvalidTransition
	}
	result, err := tx.ExecContext(ctx, `UPDATE devices SET lifecycle_status=$3,updated_at=NOW(),deactivated_at=CASE WHEN $3='DECOMMISSIONED' THEN NOW() ELSE deactivated_at END WHERE tenant_id=$1 AND device_id=$2`, tenant, id, next)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	payload := map[string]string{"device_id": id, "previous_status": current, "status": next}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "DEVICE_"+next, "device", id, requestID, "device.lifecycle.changed.v1", payload); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) ListCapabilities(ctx context.Context, tenant, id string) ([]registry.Capability, error) {
	v := []registry.Capability{}
	err := s.DB.SelectContext(ctx, &v, `SELECT c.capability_id,c.display_name,c.description,dc.configuration,dc.enabled FROM device_capabilities dc JOIN capabilities c USING(capability_id) WHERE dc.tenant_id=$1 AND dc.device_id=$2 ORDER BY c.capability_id`, tenant, id)
	return v, err
}
func (s *RegistryStore) AllCapabilities(ctx context.Context) ([]registry.Capability, error) {
	v := []registry.Capability{}
	err := s.DB.SelectContext(ctx, &v, `SELECT capability_id,display_name,description,'{}'::jsonb configuration,true enabled FROM capabilities ORDER BY capability_id`)
	return v, err
}
func (s *RegistryStore) PutCapability(ctx context.Context, tenant, id, capability string, configuration []byte, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO device_capabilities(tenant_id,device_id,capability_id,configuration,enabled) VALUES($1,$2,$3,$4,true) ON CONFLICT(tenant_id,device_id,capability_id) DO UPDATE SET configuration=EXCLUDED.configuration,enabled=true`, tenant, id, capability, configuration)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CAPABILITY_ADDED", "device", id, requestID, "device.capabilities.changed.v1", map[string]string{"device_id": id, "capability_id": capability}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) RemoveCapability(ctx context.Context, tenant, id, capability, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `DELETE FROM device_capabilities WHERE tenant_id=$1 AND device_id=$2 AND capability_id=$3`, tenant, id, capability)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CAPABILITY_REMOVED", "device", id, requestID, "device.capabilities.changed.v1", map[string]string{"device_id": id, "capability_id": capability, "enabled": "false"}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) IssueCredential(ctx context.Context, tenant, id, actor, requestID string, expires *time.Time) (registry.CredentialMetadata, string, error) {
	raw, prefix, hash, err := auth.GenerateToken("dev")
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	credentialID := auth.NewID()
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `INSERT INTO device_credentials(credential_id,tenant_id,device_id,token_prefix,token_hash,status,expires_at) SELECT $1,$2,$3,$4,$5,'ACTIVE',$6 WHERE EXISTS(SELECT 1 FROM devices WHERE tenant_id=$2 AND device_id=$3 AND lifecycle_status<>'DECOMMISSIONED')`, credentialID, tenant, id, prefix, hash, expires)
	if err != nil {
		return registry.CredentialMetadata{}, "", mapPQ(err)
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return registry.CredentialMetadata{}, "", ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CREDENTIAL_CREATED", "device", id, requestID, "device.credential.created.v1", map[string]string{"device_id": id, "credential_id": credentialID}); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if err = tx.Commit(); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	return registry.CredentialMetadata{CredentialID: credentialID, TokenPrefix: prefix, Status: "ACTIVE", IssuedAt: time.Now().UTC(), ExpiresAt: expires}, raw, nil
}
func (s *RegistryStore) ListCredentials(ctx context.Context, tenant, id string) ([]registry.CredentialMetadata, error) {
	v := []registry.CredentialMetadata{}
	err := s.DB.SelectContext(ctx, &v, `SELECT credential_id,token_prefix,status,issued_at,expires_at,last_used_at,revoked_at FROM device_credentials WHERE tenant_id=$1 AND device_id=$2 ORDER BY issued_at DESC`, tenant, id)
	return v, err
}
func (s *RegistryStore) RevokeCredential(ctx context.Context, tenant, id, credentialID, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE device_credentials SET status='REVOKED',revoked_at=NOW() WHERE credential_id=$1 AND tenant_id=$2 AND device_id=$3 AND status='ACTIVE'`, credentialID, tenant, id)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrNotFound
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CREDENTIAL_REVOKED", "device", id, requestID, "device.credential.revoked.v1", map[string]string{"device_id": id, "credential_id": credentialID}); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *RegistryStore) RotateCredential(ctx context.Context, tenant, id, oldID, actor, requestID string) (registry.CredentialMetadata, string, error) {
	raw, prefix, hash, err := auth.GenerateToken("dev")
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	newID := auth.NewID()
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE device_credentials SET status='REVOKED',revoked_at=NOW() WHERE credential_id=$1 AND tenant_id=$2 AND device_id=$3 AND status='ACTIVE'`, oldID, tenant, id)
	if err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return registry.CredentialMetadata{}, "", ErrNotFound
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO device_credentials(credential_id,tenant_id,device_id,token_prefix,token_hash,status) VALUES($1,$2,$3,$4,$5,'ACTIVE')`, newID, tenant, id, prefix, hash); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "CREDENTIAL_ROTATED", "device", id, requestID, "device.credential.rotated.v1", map[string]string{"device_id": id, "old_credential_id": oldID, "credential_id": newID}); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	if err = tx.Commit(); err != nil {
		return registry.CredentialMetadata{}, "", err
	}
	return registry.CredentialMetadata{CredentialID: newID, TokenPrefix: prefix, Status: "ACTIVE", IssuedAt: time.Now().UTC()}, raw, nil
}

func (s *RegistryStore) CreateTicket(ctx context.Context, tenant, id, credentialID string, ttl time.Duration) (string, error) {
	raw, prefix, hash, err := auth.GenerateToken("ticket")
	if err != nil {
		return "", err
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO connection_tickets(ticket_prefix,ticket_hash,tenant_id,device_id,credential_id,expires_at) SELECT $1,$2,$3,$4,$5,NOW()+$6::interval WHERE EXISTS(SELECT 1 FROM devices WHERE tenant_id=$3 AND device_id=$4 AND lifecycle_status='ACTIVE')`, prefix, hash, tenant, id, credentialID, ttl.String())
	return raw, err
}
func (s *RegistryStore) ConsumeTicket(ctx context.Context, raw string) (auth.DevicePrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return auth.DevicePrincipal{}, err
	}
	defer tx.Rollback()
	var p auth.DevicePrincipal
	var hash []byte
	var status string
	err = tx.QueryRowxContext(ctx, `SELECT x.tenant_id,x.device_id,x.credential_id,d.device_type_id,COALESCE(d.project_id::text,''),x.ticket_hash,d.lifecycle_status FROM connection_tickets x JOIN devices d USING(tenant_id,device_id) JOIN tenants t USING(tenant_id) WHERE x.ticket_prefix=$1 AND x.consumed_at IS NULL AND x.expires_at>NOW() AND t.status='ACTIVE' FOR UPDATE OF x`, prefix).Scan(&p.TenantID, &p.DeviceID, &p.CredentialID, &p.DeviceType, &p.ProjectID, &hash, &status)
	if err != nil || status != "ACTIVE" || !auth.Verify(raw, hash) {
		return auth.DevicePrincipal{}, auth.ErrInvalidCredential
	}
	if _, err = tx.ExecContext(ctx, `UPDATE connection_tickets SET consumed_at=NOW() WHERE ticket_prefix=$1`, prefix); err != nil {
		return auth.DevicePrincipal{}, err
	}
	if err = tx.Commit(); err != nil {
		return auth.DevicePrincipal{}, err
	}
	return p, nil
}
func (s *RegistryStore) CreateOperatorTicket(ctx context.Context, p auth.OperatorPrincipal, tenant string, ttl time.Duration) (string, error) {
	if p.Role != auth.PlatformAdmin {
		tenant = p.TenantID
	}
	if tenant == "" {
		return "", ErrForbidden
	}
	raw, prefix, hash, err := auth.GenerateToken("ticket")
	if err != nil {
		return "", err
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO operator_connection_tickets(ticket_prefix,ticket_hash,api_key_id,tenant_id,role,expires_at) VALUES($1,$2,$3,$4,$5,NOW()+$6::interval)`, prefix, hash, p.APIKeyID, tenant, p.Role, ttl.String())
	return raw, err
}
func (s *RegistryStore) ConsumeOperatorTicket(ctx context.Context, raw string) (auth.OperatorPrincipal, error) {
	prefix, err := auth.TokenPrefix(raw)
	if err != nil {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return auth.OperatorPrincipal{}, err
	}
	defer tx.Rollback()
	var p auth.OperatorPrincipal
	var hash []byte
	err = tx.QueryRowxContext(ctx, `SELECT api_key_id,tenant_id,role,ticket_hash FROM operator_connection_tickets WHERE ticket_prefix=$1 AND consumed_at IS NULL AND expires_at>NOW() FOR UPDATE`, prefix).Scan(&p.APIKeyID, &p.TenantID, &p.Role, &hash)
	if err != nil || !auth.Verify(raw, hash) {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	if _, err = tx.ExecContext(ctx, `UPDATE operator_connection_tickets SET consumed_at=NOW() WHERE ticket_prefix=$1`, prefix); err != nil {
		return auth.OperatorPrincipal{}, err
	}
	if err = tx.Commit(); err != nil {
		return auth.OperatorPrincipal{}, err
	}
	return p, nil
}

func (s *RegistryStore) Audit(ctx context.Context, tenant, actorID, action, resource, id, requestID, outcome string) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO audit_events(audit_id,tenant_id,actor_type,actor_id,action,resource_type,resource_id,request_id,outcome) VALUES($1,NULLIF($2,''),'OPERATOR',$3,$4,$5,$6,$7,$8)`, auth.NewID(), tenant, actorID, action, resource, id, requestID, outcome)
	return err
}
func (s *RegistryStore) ListAudit(ctx context.Context, tenant string) ([]map[string]interface{}, error) {
	rows, err := s.DB.QueryxContext(ctx, `SELECT audit_id,tenant_id,actor_type,actor_id,action,resource_type,resource_id,request_id,outcome,created_at FROM audit_events WHERE tenant_id=$1 ORDER BY created_at DESC LIMIT 100`, tenant)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []map[string]interface{}{}
	for rows.Next() {
		m := map[string]interface{}{}
		if err = rows.MapScan(m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

type OutboxEvent struct {
	OutboxID, EventID, EventType, TenantID string
	Payload                                []byte
}

func (s *RegistryStore) ClaimOutbox(ctx context.Context, limit int) ([]OutboxEvent, error) {
	if limit < 1 {
		limit = 100
	}
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryxContext(ctx, `SELECT outbox_id,event_id,event_type,tenant_id,payload FROM outbox_events WHERE status IN('PENDING','RETRY_PENDING') AND next_attempt_at<=NOW() ORDER BY created_at LIMIT $1 FOR UPDATE SKIP LOCKED`, limit)
	if err != nil {
		return nil, err
	}
	events := []OutboxEvent{}
	for rows.Next() {
		var e OutboxEvent
		if err = rows.Scan(&e.OutboxID, &e.EventID, &e.EventType, &e.TenantID, &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	rows.Close()
	if len(events) > 0 {
		ids := make([]string, len(events))
		for i, e := range events {
			ids[i] = e.OutboxID
		}
		_, err = tx.ExecContext(ctx, `UPDATE outbox_events SET status='RETRY_PENDING',attempt_count=attempt_count+1,next_attempt_at=NOW()+INTERVAL '5 seconds' WHERE outbox_id=ANY($1)`, pq.Array(ids))
		if err != nil {
			return nil, err
		}
	}
	return events, tx.Commit()
}
func (s *RegistryStore) MarkOutboxPublished(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE outbox_events SET status='PUBLISHED',published_at=NOW(),last_error=NULL WHERE outbox_id=$1`, id)
	return err
}
func (s *RegistryStore) MarkOutboxFailed(ctx context.Context, id string, errText string) error {
	if len(errText) > 500 {
		errText = errText[:500]
	}
	_, err := s.DB.ExecContext(ctx, `UPDATE outbox_events SET status=CASE WHEN attempt_count>=10 THEN 'FAILED' ELSE 'RETRY_PENDING' END,last_error=$2,next_attempt_at=NOW()+INTERVAL '5 seconds' WHERE outbox_id=$1`, id, strings.ReplaceAll(errText, "\n", " "))
	return err
}
