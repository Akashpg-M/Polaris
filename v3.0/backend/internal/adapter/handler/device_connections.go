package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/Akashpg-M/polaris/backend/internal/core/command"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type deviceSession struct {
	principal auth.DevicePrincipal
	ownership repository.ConnectionOwnership
	conn      *websocket.Conn
	writeMu   sync.Mutex
	cancel    context.CancelFunc
}

func (s *deviceSession) writeJSON(value interface{}) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.conn.WriteJSON(value)
}

func (s *deviceSession) close(code int, reason string) {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_ = s.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(code, reason), time.Now().Add(time.Second))
	_ = s.conn.Close()
}

type DeviceConnectionManager struct {
	mu        sync.RWMutex
	sessions  map[string]*deviceSession
	gatewayID string
	lease     time.Duration
	owners    *repository.ConnectionOwnershipStore
	store     *repository.RegistryStore
	redis     *redis.Client
	metrics   *orchestration.Metrics
}

func NewDeviceConnectionManager(gatewayID string, lease time.Duration, redisClient *redis.Client, store *repository.RegistryStore, metrics *orchestration.Metrics) *DeviceConnectionManager {
	if gatewayID == "" {
		gatewayID = "gateway-1"
	}
	return &DeviceConnectionManager{sessions: map[string]*deviceSession{}, gatewayID: gatewayID, lease: lease, owners: repository.NewConnectionOwnershipStore(redisClient, lease), store: store, redis: redisClient, metrics: metrics}
}

func deviceKey(tenant, device string) string { return tenant + ":" + device }

func (m *DeviceConnectionManager) Register(ctx context.Context, conn *websocket.Conn, principal auth.DevicePrincipal) (*deviceSession, error) {
	connectionID := auth.NewID()
	ownership, err := m.owners.Claim(ctx, principal.TenantID, principal.DeviceID, m.gatewayID, connectionID, principal.CredentialID)
	if err != nil {
		return nil, err
	}
	sessionCtx, cancel := context.WithCancel(context.Background())
	session := &deviceSession{principal: principal, ownership: ownership, conn: conn, cancel: cancel}
	key := deviceKey(principal.TenantID, principal.DeviceID)
	m.mu.Lock()
	previous := m.sessions[key]
	m.sessions[key] = session
	m.mu.Unlock()
	if previous != nil {
		previous.cancel()
		previous.close(websocket.CloseServiceRestart, "connection superseded by a newer ownership epoch")
	}
	m.metrics.ActiveConnections.Add(1)
	_ = m.store.Audit(ctx, principal.TenantID, m.gatewayID, "DEVICE_OWNERSHIP_CLAIMED", "device", principal.DeviceID, "", "SUCCESS")
	go m.heartbeat(sessionCtx, session)
	go m.ReconcileDevice(sessionCtx, principal.TenantID, principal.DeviceID)
	return session, nil
}

func (m *DeviceConnectionManager) Unregister(ctx context.Context, session *deviceSession) {
	if session == nil {
		return
	}
	session.cancel()
	key := deviceKey(session.principal.TenantID, session.principal.DeviceID)
	m.mu.Lock()
	if m.sessions[key] == session {
		delete(m.sessions, key)
		m.metrics.ActiveConnections.Add(-1)
	}
	m.mu.Unlock()
	released, _ := m.owners.Release(ctx, session.ownership)
	if released {
		_ = m.store.Audit(ctx, session.principal.TenantID, m.gatewayID, "DEVICE_OWNERSHIP_LOST", "device", session.principal.DeviceID, "", "SUCCESS")
	}
}

func (m *DeviceConnectionManager) heartbeat(ctx context.Context, session *deviceSession) {
	interval := m.lease / 3
	if interval < time.Second {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ok, err := m.owners.Refresh(ctx, session.ownership)
			if err != nil || !ok {
				session.close(websocket.ClosePolicyViolation, "gateway ownership lease lost")
				return
			}
		}
	}
}

func (m *DeviceConnectionManager) session(tenant, device string) *deviceSession {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[deviceKey(tenant, device)]
}

func (m *DeviceConnectionManager) Deliver(ctx context.Context, envelope command.Envelope) error {
	session := m.session(envelope.TenantID, envelope.DeviceID)
	if session == nil {
		return repository.ErrNotFound
	}
	if err := m.store.RevalidateDevice(ctx, session.principal); err != nil {
		session.close(websocket.ClosePolicyViolation, "credential, tenant, or device is inactive")
		return repository.ErrForbidden
	}
	ok, err := m.owners.Owns(ctx, session.ownership)
	if err != nil || !ok {
		return repository.ErrForbidden
	}
	record, err := m.store.PrepareDelivery(ctx, envelope.TenantID, envelope.DeviceID, envelope.CommandID, m.gatewayID, session.ownership.Epoch)
	if err != nil {
		return err
	}
	// Fence again immediately before the volatile write. A reconnect may have
	// advanced the epoch while the PostgreSQL delivery transition was running.
	ok, err = m.owners.Owns(ctx, session.ownership)
	if err != nil || !ok {
		return repository.ErrForbidden
	}
	if err = session.writeJSON(record.Envelope()); err != nil {
		return err
	}
	m.metrics.CommandsDelivered.Add(1)
	return nil
}

func (m *DeviceConnectionManager) ReconcileDevice(ctx context.Context, tenant, device string) {
	commands, err := m.store.PendingCommandsForDevice(ctx, tenant, device)
	if err != nil {
		slog.Warn("pending command reconciliation failed", "tenant_id", tenant, "device_id", device, "error", err)
		return
	}
	for _, record := range commands {
		if err = m.Deliver(ctx, record.Envelope()); err != nil {
			if err != repository.ErrInvalidTransition && err != repository.ErrConflict {
				slog.Warn("reconciled command delivery deferred", "command_id", record.CommandID, "error", err)
			}
			return
		}
	}
}

func (m *DeviceConnectionManager) StartSubscriber(ctx context.Context) {
	pubsub := m.redis.Subscribe(ctx, repository.GatewayCommandChannel(m.gatewayID))
	defer pubsub.Close()
	for message := range pubsub.Channel() {
		var envelope command.Envelope
		if json.Unmarshal([]byte(message.Payload), &envelope) != nil || envelope.FrameType != "COMMAND" {
			continue
		}
		if err := m.Deliver(ctx, envelope); err != nil && err != repository.ErrNotFound && err != repository.ErrInvalidTransition && err != repository.ErrConflict {
			slog.Warn("live command notification could not be delivered", "command_id", envelope.CommandID, "error", err)
		}
	}
}

func (m *DeviceConnectionManager) StartReconciler(ctx context.Context, interval time.Duration) {
	if interval < 100*time.Millisecond {
		interval = time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.mu.RLock()
			devices := make([][2]string, 0, len(m.sessions))
			for _, session := range m.sessions {
				devices = append(devices, [2]string{session.principal.TenantID, session.principal.DeviceID})
			}
			m.mu.RUnlock()
			for _, device := range devices {
				go m.ReconcileDevice(ctx, device[0], device[1])
			}
		}
	}
}

func (m *DeviceConnectionManager) OwnershipStore() *repository.ConnectionOwnershipStore {
	return m.owners
}
