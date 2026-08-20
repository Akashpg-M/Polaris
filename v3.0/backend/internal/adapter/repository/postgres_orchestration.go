package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	taskcore "github.com/Akashpg-M/polaris/backend/internal/core/task"
	"github.com/lib/pq"
)

type DeviceCandidate struct {
	DeviceID     string `db:"device_id"`
	DeviceTypeID string `db:"device_type_id"`
}

func (s *RegistryStore) CreateTask(ctx context.Context, v taskcore.Task, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO tasks(task_id,tenant_id,project_id,task_type,status,priority,requirements,target,correlation_id,created_by,expires_at) VALUES($1,$2,$3,$4,'PENDING',$5,$6,$7,$8,$9,$10)`, v.TaskID, v.TenantID, v.ProjectID, v.TaskType, v.Priority, v.Requirements, v.Target, v.CorrelationID, v.CreatedBy, v.ExpiresAt)
	if err != nil {
		return mapPQ(err)
	}
	if err = insertAuditOutbox(ctx, tx, v.TenantID, actor, "TASK_CREATED", "task", v.TaskID, requestID, "task.created.v1", map[string]interface{}{"task_id": v.TaskID, "task_type": v.TaskType, "priority": v.Priority}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) GetTask(ctx context.Context, tenant, id string) (taskcore.Task, error) {
	var v taskcore.Task
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM tasks WHERE tenant_id=$1 AND task_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}

func (s *RegistryStore) ListTasks(ctx context.Context, tenant string, limit int, cursor, status, deviceID string) ([]taskcore.Task, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	args := []interface{}{tenant, limit}
	q := `SELECT * FROM tasks WHERE tenant_id=$1`
	add := func(clause string, value interface{}) {
		args = append(args, value)
		q += fmt.Sprintf(clause, len(args))
	}
	if cursor != "" {
		add(` AND task_id>$%d`, cursor)
	}
	if status != "" {
		add(` AND status=$%d`, status)
	}
	if deviceID != "" {
		add(` AND assigned_device_id=$%d`, deviceID)
	}
	q += ` ORDER BY task_id LIMIT $2`
	result := []taskcore.Task{}
	return result, s.DB.SelectContext(ctx, &result, q, args...)
}

func (s *RegistryStore) PendingTasks(ctx context.Context, limit int) ([]taskcore.Task, error) {
	if limit < 1 {
		limit = 50
	}
	result := []taskcore.Task{}
	err := s.DB.SelectContext(ctx, &result, `SELECT * FROM tasks WHERE status='PENDING' AND expires_at>NOW() ORDER BY CASE priority WHEN 'CRITICAL' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'NORMAL' THEN 3 ELSE 4 END,created_at LIMIT $1`, limit)
	return result, err
}

func (s *RegistryStore) EligibleDevices(ctx context.Context, tenant string, requirements taskcore.Requirements) ([]DeviceCandidate, error) {
	args := []interface{}{tenant}
	q := `SELECT d.device_id,d.device_type_id FROM devices d WHERE d.tenant_id=$1 AND d.lifecycle_status='ACTIVE'`
	if requirements.ProjectID != "" {
		args = append(args, requirements.ProjectID)
		q += fmt.Sprintf(` AND d.project_id=$%d`, len(args))
	}
	if len(requirements.AllowedDeviceTypes) > 0 {
		args = append(args, pq.Array(requirements.AllowedDeviceTypes))
		q += fmt.Sprintf(` AND d.device_type_id=ANY($%d)`, len(args))
	}
	for _, capability := range requirements.RequiredCapabilities {
		args = append(args, capability)
		q += fmt.Sprintf(` AND EXISTS(SELECT 1 FROM device_capabilities dc WHERE dc.tenant_id=d.tenant_id AND dc.device_id=d.device_id AND dc.capability_id=$%d AND dc.enabled)`, len(args))
	}
	q += ` AND NOT EXISTS(SELECT 1 FROM device_assignments a WHERE a.tenant_id=d.tenant_id AND a.device_id=d.device_id AND a.status='ACTIVE') ORDER BY d.device_id`
	result := []DeviceCandidate{}
	return result, s.DB.SelectContext(ctx, &result, q, args...)
}

func (s *RegistryStore) AssignTaskWithPlan(ctx context.Context, v taskcore.Task, deviceID, actor, requestID string, maxAttempts int, plan extension.ExecutionPlan) (command.Record, error) {
	if maxAttempts < 1 {
		maxAttempts = 5
	}
	tx, err := s.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return command.Record{}, err
	}
	defer tx.Rollback()
	var current string
	if err = tx.GetContext(ctx, &current, `SELECT status FROM tasks WHERE tenant_id=$1 AND task_id=$2 FOR UPDATE`, v.TenantID, v.TaskID); errors.Is(err, sql.ErrNoRows) {
		return command.Record{}, ErrNotFound
	}
	if err != nil {
		return command.Record{}, err
	}
	if current != string(taskcore.Pending) {
		return command.Record{}, ErrInvalidTransition
	}
	if plan.CommandType == "" || len(plan.Payload) == 0 || !json.Valid(plan.Payload) {
		return command.Record{}, ErrInvalidTransition
	}
	commandExpiry := v.ExpiresAt
	if plan.ValidUntil != nil && plan.ValidUntil.Before(commandExpiry) {
		commandExpiry = *plan.ValidUntil
	}
	if !commandExpiry.After(time.Now()) {
		return command.Record{}, ErrInvalidTransition
	}
	assignmentID := auth.NewID()
	_, err = tx.ExecContext(ctx, `INSERT INTO device_assignments(assignment_id,tenant_id,device_id,task_id,status,lease_expires_at) VALUES($1,$2,$3,$4,'ACTIVE',$5)`, assignmentID, v.TenantID, deviceID, v.TaskID, v.ExpiresAt)
	if err != nil {
		return command.Record{}, mapPQ(err)
	}
	var sequence int64
	err = tx.GetContext(ctx, &sequence, `INSERT INTO device_command_sequences(tenant_id,device_id,last_sequence) VALUES($1,$2,1) ON CONFLICT(tenant_id,device_id) DO UPDATE SET last_sequence=device_command_sequences.last_sequence+1 RETURNING last_sequence`, v.TenantID, deviceID)
	if err != nil {
		return command.Record{}, err
	}
	now := time.Now().UTC()
	record := command.Record{CommandID: auth.NewID(), TenantID: v.TenantID, DeviceID: deviceID, TaskID: v.TaskID, CommandType: plan.CommandType, Payload: plan.Payload, Status: string(command.Pending), SequenceNumber: sequence, CorrelationID: v.CorrelationID, CausationID: v.TaskID, AttemptCount: 0, MaxAttempts: maxAttempts, Version: 1, CreatedAt: now, AvailableAt: now, ExpiresAt: commandExpiry}
	_, err = tx.ExecContext(ctx, `INSERT INTO commands(command_id,tenant_id,device_id,task_id,command_type,payload,status,sequence_number,correlation_id,causation_id,max_attempts,created_at,available_at,expires_at) VALUES($1,$2,$3,$4,$5,$6,'PENDING',$7,$8,$9,$10,$11,$11,$12)`, record.CommandID, record.TenantID, record.DeviceID, record.TaskID, record.CommandType, record.Payload, record.SequenceNumber, record.CorrelationID, record.CausationID, record.MaxAttempts, now, record.ExpiresAt)
	if err != nil {
		return command.Record{}, err
	}
	result, err := tx.ExecContext(ctx, `UPDATE tasks SET status='ASSIGNED',assigned_device_id=$3,assigned_at=NOW(),updated_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status='PENDING'`, v.TenantID, v.TaskID, deviceID)
	if err != nil {
		return command.Record{}, err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return command.Record{}, ErrInvalidTransition
	}
	if err = insertAuditOutbox(ctx, tx, v.TenantID, actor, "TASK_ASSIGNED", "task", v.TaskID, requestID, "task.assigned.v1", map[string]string{"task_id": v.TaskID, "device_id": deviceID}); err != nil {
		return command.Record{}, err
	}
	if err = insertAuditOutbox(ctx, tx, v.TenantID, actor, "COMMAND_CREATED", "command", record.CommandID, requestID, "command.created.v1", record.Envelope()); err != nil {
		return command.Record{}, err
	}
	return record, tx.Commit()
}

func (s *RegistryStore) GetCommand(ctx context.Context, tenant, id string) (command.Record, error) {
	var v command.Record
	err := s.DB.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND command_id=$2`, tenant, id)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrNotFound
	}
	return v, err
}

func (s *RegistryStore) ListCommands(ctx context.Context, tenant string, limit int, cursor, status, taskID, deviceID string) ([]command.Record, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	args := []interface{}{tenant, limit}
	q := `SELECT * FROM commands WHERE tenant_id=$1`
	add := func(clause string, value interface{}) {
		args = append(args, value)
		q += fmt.Sprintf(clause, len(args))
	}
	if cursor != "" {
		add(` AND command_id>$%d`, cursor)
	}
	if status != "" {
		add(` AND status=$%d`, status)
	}
	if taskID != "" {
		add(` AND task_id=$%d`, taskID)
	}
	if deviceID != "" {
		add(` AND device_id=$%d`, deviceID)
	}
	q += ` ORDER BY command_id LIMIT $2`
	result := []command.Record{}
	return result, s.DB.SelectContext(ctx, &result, q, args...)
}

func (s *RegistryStore) PendingCommandsForDevice(ctx context.Context, tenant, deviceID string) ([]command.Record, error) {
	result := []command.Record{}
	err := s.DB.SelectContext(ctx, &result, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND status='PENDING' AND available_at<=NOW() AND expires_at>NOW() AND sequence_number=(SELECT MIN(sequence_number) FROM commands x WHERE x.tenant_id=$1 AND x.device_id=$2 AND x.status IN('PENDING','DELIVERED','ACKNOWLEDGED')) ORDER BY sequence_number`, tenant, deviceID)
	return result, err
}

func (s *RegistryStore) PrepareDelivery(ctx context.Context, tenant, deviceID, commandID, gatewayID string, epoch int64) (command.Record, error) {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return command.Record{}, err
	}
	defer tx.Rollback()
	var v command.Record
	err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 FOR UPDATE`, tenant, deviceID, commandID)
	if errors.Is(err, sql.ErrNoRows) {
		return command.Record{}, ErrNotFound
	}
	if err != nil {
		return command.Record{}, err
	}
	if v.Status != string(command.Pending) || v.AvailableAt.After(time.Now()) {
		return command.Record{}, ErrInvalidTransition
	}
	if time.Now().After(v.ExpiresAt) || v.AttemptCount >= v.MaxAttempts {
		return command.Record{}, ErrInvalidTransition
	}
	var first int64
	if err = tx.GetContext(ctx, &first, `SELECT MIN(sequence_number) FROM commands WHERE tenant_id=$1 AND device_id=$2 AND status IN('PENDING','DELIVERED','ACKNOWLEDGED')`, tenant, deviceID); err != nil || first != v.SequenceNumber {
		return command.Record{}, ErrConflict
	}
	v.AttemptCount++
	v.Status = string(command.Delivered)
	now := time.Now().UTC()
	v.SentAt = &now
	_, err = tx.ExecContext(ctx, `UPDATE commands SET status='DELIVERED',attempt_count=$4,sent_at=$5,version=version+1,last_error=NULL WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 AND status='PENDING'`, tenant, deviceID, commandID, v.AttemptCount, now)
	if err != nil {
		return command.Record{}, err
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO command_attempts(attempt_id,command_id,attempt_number,gateway_id,ownership_epoch) VALUES($1,$2,$3,$4,$5)`, auth.NewID(), commandID, v.AttemptCount, gatewayID, epoch)
	if err != nil {
		return command.Record{}, err
	}
	if err = insertAuditOutbox(ctx, tx, tenant, gatewayID, "COMMAND_DELIVERED", "command", commandID, "", "command.delivered.v1", map[string]interface{}{"command_id": commandID, "attempt": v.AttemptCount, "gateway_id": gatewayID, "ownership_epoch": epoch}); err != nil {
		return command.Record{}, err
	}
	return v, tx.Commit()
}

func (s *RegistryStore) ApplyCommandAck(ctx context.Context, principalTenant, principalDevice string, ack command.Ack) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var v command.Record
	err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 FOR UPDATE`, principalTenant, principalDevice, ack.CommandID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if ack.SequenceNumber != v.SequenceNumber {
		return ErrForbidden
	}
	if v.Status == string(command.Acknowledged) || command.IsTerminal(command.Status(v.Status)) {
		return tx.Commit()
	}
	if v.Status != string(command.Delivered) {
		return ErrInvalidTransition
	}
	now := time.Now().UTC()
	if ack.ReceivedAt.IsZero() {
		ack.ReceivedAt = now
	}
	accepted := ack.Status == "ACCEPTED" || ack.Status == "DUPLICATE"
	if accepted {
		_, err = tx.ExecContext(ctx, `UPDATE commands SET status='ACKNOWLEDGED',ack_status=$2,acknowledged_at=$3,version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, ack.Status, ack.ReceivedAt)
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE command_attempts SET completed_at=$2,result=$3 WHERE command_id=$1 AND attempt_number=$4`, v.CommandID, ack.ReceivedAt, ack.Status, v.AttemptCount)
		}
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE tasks SET status='IN_PROGRESS',started_at=COALESCE(started_at,$2),updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status='ASSIGNED'`, v.TaskID, ack.ReceivedAt)
		}
	} else {
		next := string(command.Failed)
		taskNext := string(taskcore.Failed)
		if ack.Status == "EXPIRED" {
			next, taskNext = string(command.Expired), string(taskcore.Expired)
		}
		_, err = tx.ExecContext(ctx, `UPDATE commands SET status=$2,ack_status=$3,completed_at=$4,last_error=NULLIF($5,''),version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, next, ack.Status, ack.ReceivedAt, ack.Reason)
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE tasks SET status=$2,failed_at=$3,failure_reason=NULLIF($4,''),updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status IN('ASSIGNED','IN_PROGRESS')`, v.TaskID, taskNext, ack.ReceivedAt, ack.Reason)
		}
		if err == nil {
			_, err = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, v.TaskID)
		}
	}
	if err != nil {
		return err
	}
	action := "COMMAND_ACKNOWLEDGED"
	if !accepted {
		action = "COMMAND_REJECTED"
	}
	if err = insertAuditOutbox(ctx, tx, principalTenant, principalDevice, action, "command", v.CommandID, "", "command.acknowledged.v1", map[string]interface{}{"command_id": v.CommandID, "sequence_number": v.SequenceNumber, "status": ack.Status}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) ApplyCommandResult(ctx context.Context, principalTenant, principalDevice string, result command.Result) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var v command.Record
	err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND device_id=$2 AND command_id=$3 FOR UPDATE`, principalTenant, principalDevice, result.CommandID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if result.SequenceNumber != v.SequenceNumber {
		return ErrForbidden
	}
	if command.IsTerminal(command.Status(v.Status)) {
		return tx.Commit()
	}
	if v.Status != string(command.Acknowledged) {
		return ErrInvalidTransition
	}
	if result.CompletedAt.IsZero() {
		result.CompletedAt = time.Now().UTC()
	}
	next, taskNext, action := string(command.Completed), string(taskcore.Completed), "COMMAND_COMPLETED"
	if result.Status != "SUCCEEDED" && result.Status != "COMPLETED" {
		next, taskNext, action = string(command.Failed), string(taskcore.Failed), "COMMAND_FAILED"
	}
	payload := result.Result
	if len(payload) == 0 {
		payload = []byte(`{}`)
	}
	_, err = tx.ExecContext(ctx, `UPDATE commands SET status=$2,result=$3,completed_at=$4,last_error=NULLIF($5,''),version=version+1 WHERE command_id=$1 AND status='ACKNOWLEDGED'`, v.CommandID, next, payload, result.CompletedAt, result.Reason)
	if err == nil {
		_, err = tx.ExecContext(ctx, `UPDATE tasks SET status=$2,completed_at=CASE WHEN $2='COMPLETED' THEN $3 ELSE completed_at END,failed_at=CASE WHEN $2='FAILED' THEN $3 ELSE failed_at END,failure_reason=NULLIF($4,''),updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status='IN_PROGRESS'`, v.TaskID, taskNext, result.CompletedAt, result.Reason)
	}
	if err == nil {
		_, err = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, v.TaskID)
	}
	if err != nil {
		return err
	}
	if err = insertAuditOutbox(ctx, tx, principalTenant, principalDevice, action, "command", v.CommandID, "", "command.result.v1", map[string]interface{}{"command_id": v.CommandID, "status": result.Status}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) CancelTask(ctx context.Context, tenant, taskID, actor, requestID string) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var status string
	if err = tx.GetContext(ctx, &status, `SELECT status FROM tasks WHERE tenant_id=$1 AND task_id=$2 FOR UPDATE`, tenant, taskID); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if status == string(taskcore.Cancelled) {
		return tx.Commit()
	}
	var commandStatus sql.NullString
	_ = tx.GetContext(ctx, &commandStatus, `SELECT status FROM commands WHERE tenant_id=$1 AND task_id=$2 ORDER BY sequence_number DESC LIMIT 1`, tenant, taskID)
	if commandStatus.Valid && commandStatus.String != string(command.Pending) {
		return ErrInvalidTransition
	}
	result, err := tx.ExecContext(ctx, `UPDATE tasks SET status='CANCELLED',updated_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status IN('PENDING','ASSIGNING','ASSIGNED')`, tenant, taskID)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrInvalidTransition
	}
	_, _ = tx.ExecContext(ctx, `UPDATE commands SET status='CANCELLED',completed_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status='PENDING'`, tenant, taskID)
	_, _ = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, taskID)
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "TASK_CANCELLED", "task", taskID, requestID, "task.cancelled.v1", map[string]string{"task_id": taskID}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) RetryTask(ctx context.Context, tenant, taskID, actor, requestID string, expiresAt time.Time) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE tasks SET status='PENDING',assigned_device_id=NULL,assigned_at=NULL,started_at=NULL,completed_at=NULL,failed_at=NULL,failure_reason=NULL,expires_at=$3,updated_at=NOW(),version=version+1 WHERE tenant_id=$1 AND task_id=$2 AND status IN('FAILED','EXPIRED')`, tenant, taskID, expiresAt)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n != 1 {
		return ErrInvalidTransition
	}
	_, _ = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, taskID)
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "TASK_RETRIED", "task", taskID, requestID, "task.retry.requested.v1", map[string]string{"task_id": taskID}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) RetryCommand(ctx context.Context, tenant, commandID, actor, requestID string, force bool) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var v command.Record
	if err = tx.GetContext(ctx, &v, `SELECT * FROM commands WHERE tenant_id=$1 AND command_id=$2 FOR UPDATE`, tenant, commandID); errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if v.Status != string(command.Delivered) && !(force && v.Status == string(command.Failed)) {
		return ErrInvalidTransition
	}
	if time.Now().After(v.ExpiresAt) || (!force && v.AttemptCount >= v.MaxAttempts) {
		return ErrInvalidTransition
	}
	if v.Status == string(command.Failed) {
		result, assignErr := tx.ExecContext(ctx, `UPDATE device_assignments a SET status='ACTIVE',lease_expires_at=$2,updated_at=NOW() WHERE task_id=$1 AND status IN('RELEASED','EXPIRED') AND NOT EXISTS(SELECT 1 FROM device_assignments x WHERE x.tenant_id=a.tenant_id AND x.device_id=a.device_id AND x.status='ACTIVE')`, v.TaskID, v.ExpiresAt)
		if assignErr != nil {
			return mapPQ(assignErr)
		}
		if n, _ := result.RowsAffected(); n != 1 {
			return ErrConflict
		}
		_, err = tx.ExecContext(ctx, `UPDATE tasks SET status='ASSIGNED',failure_reason=NULL,failed_at=NULL,updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status='FAILED'`, v.TaskID)
		if err != nil {
			return err
		}
	}
	resetAttempts := v.AttemptCount
	if force && v.Status == string(command.Failed) {
		resetAttempts = 0
	}
	_, err = tx.ExecContext(ctx, `UPDATE commands SET status='PENDING',attempt_count=$2,available_at=NOW(),last_error=NULL,completed_at=NULL,version=version+1 WHERE command_id=$1`, commandID, resetAttempts)
	if err != nil {
		return err
	}
	v.Status = string(command.Pending)
	if err = insertAuditOutbox(ctx, tx, tenant, actor, "COMMAND_RETRIED", "command", commandID, requestID, "command.retry.requested.v1", v.Envelope()); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *RegistryStore) ReconcileCommands(ctx context.Context, ackTimeout time.Duration) error {
	tx, err := s.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	rows, err := tx.QueryxContext(ctx, `SELECT * FROM commands WHERE status='DELIVERED' AND sent_at < NOW()-$1::interval FOR UPDATE SKIP LOCKED`, ackTimeout.String())
	if err != nil {
		return err
	}
	due := []command.Record{}
	for rows.Next() {
		var v command.Record
		if err = rows.StructScan(&v); err != nil {
			rows.Close()
			return err
		}
		due = append(due, v)
	}
	rows.Close()
	for _, v := range due {
		if time.Now().After(v.ExpiresAt) || v.AttemptCount >= v.MaxAttempts {
			next := string(command.Failed)
			taskNext := string(taskcore.Failed)
			reason := "delivery attempts exhausted"
			if time.Now().After(v.ExpiresAt) {
				next, taskNext, reason = string(command.Expired), string(taskcore.Expired), "command expired"
			}
			_, err = tx.ExecContext(ctx, `UPDATE commands SET status=$2,completed_at=NOW(),last_error=$3,version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, next, reason)
			if err == nil {
				_, err = tx.ExecContext(ctx, `UPDATE tasks SET status=$2,failed_at=NOW(),failure_reason=$3,updated_at=NOW(),version=version+1 WHERE task_id=$1 AND status IN('ASSIGNED','IN_PROGRESS')`, v.TaskID, taskNext, reason)
			}
			if err == nil {
				_, err = tx.ExecContext(ctx, `UPDATE device_assignments SET status='RELEASED',updated_at=NOW() WHERE task_id=$1 AND status='ACTIVE'`, v.TaskID)
			}
			if err != nil {
				return err
			}
			continue
		}
		backoffSeconds := int(math.Min(math.Pow(2, float64(max(v.AttemptCount-1, 0))), 30))
		_, err = tx.ExecContext(ctx, `UPDATE commands SET status='PENDING',available_at=NOW()+$2::interval,last_error='ack timeout',version=version+1 WHERE command_id=$1 AND status='DELIVERED'`, v.CommandID, (time.Duration(backoffSeconds) * time.Second).String())
		if err != nil {
			return err
		}
		v.Status = string(command.Pending)
		if err = insertAuditOutbox(ctx, tx, v.TenantID, "reconciler", "COMMAND_RETRIED", "command", v.CommandID, "", "command.retry.requested.v1", v.Envelope()); err != nil {
			return err
		}
	}
	_, err = tx.ExecContext(ctx, `WITH expired AS (UPDATE commands SET status='EXPIRED',completed_at=NOW(),last_error='command expired',version=version+1 WHERE status='PENDING' AND expires_at<=NOW() RETURNING task_id) UPDATE tasks SET status='EXPIRED',failed_at=NOW(),failure_reason='command expired',updated_at=NOW(),version=version+1 WHERE task_id IN(SELECT task_id FROM expired) AND status IN('ASSIGNED','IN_PROGRESS')`)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE device_assignments a SET status='EXPIRED',updated_at=NOW() FROM tasks t WHERE a.task_id=t.task_id AND a.status='ACTIVE' AND t.status IN('EXPIRED','FAILED','CANCELLED','COMPLETED')`)
	return tx.Commit()
}

func (s *RegistryStore) FailExpiredPendingTasks(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE tasks SET status='EXPIRED',failed_at=NOW(),failure_reason='task expired before assignment',updated_at=NOW(),version=version+1 WHERE status IN('PENDING','ASSIGNING') AND expires_at<=NOW()`)
	return err
}
