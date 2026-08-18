package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ConnectionOwnership struct {
	TenantID, DeviceID, GatewayID, ConnectionID, CredentialID string
	Epoch                                                     int64
	ConnectedAt, LeaseExpiresAt                               time.Time
}

type ConnectionOwnershipStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewConnectionOwnershipStore(client *redis.Client, ttl time.Duration) *ConnectionOwnershipStore {
	if ttl < 3*time.Second {
		ttl = 30 * time.Second
	}
	return &ConnectionOwnershipStore{client: client, ttl: ttl}
}

var claimOwnershipScript = redis.NewScript(`
local epoch = redis.call('INCR', KEYS[2])
redis.call('HSET', KEYS[1], 'gateway_id', ARGV[1], 'connection_id', ARGV[2], 'credential_id', ARGV[3], 'epoch', epoch, 'connected_at', ARGV[4], 'lease_expires_at', ARGV[5])
redis.call('PEXPIRE', KEYS[1], ARGV[6])
redis.call('SADD', KEYS[3], ARGV[7])
redis.call('ZADD', KEYS[4], ARGV[5], ARGV[7])
return epoch
`)

var refreshOwnershipScript = redis.NewScript(`
if redis.call('HGET', KEYS[1], 'gateway_id') ~= ARGV[1] or redis.call('HGET', KEYS[1], 'connection_id') ~= ARGV[2] or redis.call('HGET', KEYS[1], 'epoch') ~= ARGV[3] then return 0 end
redis.call('HSET', KEYS[1], 'lease_expires_at', ARGV[4])
redis.call('PEXPIRE', KEYS[1], ARGV[5])
redis.call('ZADD', KEYS[2], ARGV[4], ARGV[6])
return 1
`)

var releaseOwnershipScript = redis.NewScript(`
if redis.call('HGET', KEYS[1], 'gateway_id') ~= ARGV[1] or redis.call('HGET', KEYS[1], 'connection_id') ~= ARGV[2] or redis.call('HGET', KEYS[1], 'epoch') ~= ARGV[3] then return 0 end
redis.call('DEL', KEYS[1])
redis.call('SREM', KEYS[2], ARGV[4])
redis.call('ZREM', KEYS[3], ARGV[4])
return 1
`)

func ownershipKey(tenant, device string) string { return "polaris:connection:" + tenant + ":" + device }
func epochKey(tenant, device string) string {
	return "polaris:connection-epoch:" + tenant + ":" + device
}
func connectionMember(tenant, device string) string { return tenant + ":" + device }
func GatewayCommandChannel(gatewayID string) string {
	return "polaris:gateway:" + gatewayID + ":commands"
}

func (s *ConnectionOwnershipStore) Claim(ctx context.Context, tenant, device, gateway, connection, credential string) (ConnectionOwnership, error) {
	now := time.Now().UTC()
	expires := now.Add(s.ttl)
	member := connectionMember(tenant, device)
	epoch, err := claimOwnershipScript.Run(ctx, s.client, []string{ownershipKey(tenant, device), epochKey(tenant, device), "polaris:gateway:" + gateway + ":connections", "polaris:connections:lease-expiry"}, gateway, connection, credential, now.UnixMilli(), expires.UnixMilli(), s.ttl.Milliseconds(), member).Int64()
	return ConnectionOwnership{TenantID: tenant, DeviceID: device, GatewayID: gateway, ConnectionID: connection, CredentialID: credential, Epoch: epoch, ConnectedAt: now, LeaseExpiresAt: expires}, err
}

func (s *ConnectionOwnershipStore) Refresh(ctx context.Context, v ConnectionOwnership) (bool, error) {
	expires := time.Now().UTC().Add(s.ttl)
	result, err := refreshOwnershipScript.Run(ctx, s.client, []string{ownershipKey(v.TenantID, v.DeviceID), "polaris:connections:lease-expiry"}, v.GatewayID, v.ConnectionID, v.Epoch, expires.UnixMilli(), s.ttl.Milliseconds(), connectionMember(v.TenantID, v.DeviceID)).Int()
	return result == 1, err
}

func (s *ConnectionOwnershipStore) Release(ctx context.Context, v ConnectionOwnership) (bool, error) {
	result, err := releaseOwnershipScript.Run(ctx, s.client, []string{ownershipKey(v.TenantID, v.DeviceID), "polaris:gateway:" + v.GatewayID + ":connections", "polaris:connections:lease-expiry"}, v.GatewayID, v.ConnectionID, v.Epoch, connectionMember(v.TenantID, v.DeviceID)).Int()
	return result == 1, err
}

func (s *ConnectionOwnershipStore) Get(ctx context.Context, tenant, device string) (ConnectionOwnership, error) {
	values, err := s.client.HGetAll(ctx, ownershipKey(tenant, device)).Result()
	if err != nil {
		return ConnectionOwnership{}, err
	}
	if len(values) == 0 {
		return ConnectionOwnership{}, ErrNotFound
	}
	epoch, _ := strconv.ParseInt(values["epoch"], 10, 64)
	connected, _ := strconv.ParseInt(values["connected_at"], 10, 64)
	expires, _ := strconv.ParseInt(values["lease_expires_at"], 10, 64)
	return ConnectionOwnership{TenantID: tenant, DeviceID: device, GatewayID: values["gateway_id"], ConnectionID: values["connection_id"], CredentialID: values["credential_id"], Epoch: epoch, ConnectedAt: time.UnixMilli(connected), LeaseExpiresAt: time.UnixMilli(expires)}, nil
}

func (s *ConnectionOwnershipStore) Owns(ctx context.Context, v ConnectionOwnership) (bool, error) {
	current, err := s.Get(ctx, v.TenantID, v.DeviceID)
	if err != nil {
		return false, err
	}
	return current.GatewayID == v.GatewayID && current.ConnectionID == v.ConnectionID && current.Epoch == v.Epoch && current.LeaseExpiresAt.After(time.Now()), nil
}

func (s *ConnectionOwnershipStore) CleanExpired(ctx context.Context) error {
	now := time.Now().UnixMilli()
	members, err := s.client.ZRangeByScore(ctx, "polaris:connections:lease-expiry", &redis.ZRangeBy{Min: "-inf", Max: fmt.Sprint(now)}).Result()
	if err != nil {
		return err
	}
	for _, member := range members {
		parts := splitMember(member)
		if len(parts) != 2 {
			continue
		}
		current, getErr := s.Get(ctx, parts[0], parts[1])
		if getErr != nil || !current.LeaseExpiresAt.After(time.Now()) {
			_ = s.client.ZRem(ctx, "polaris:connections:lease-expiry", member).Err()
		}
	}
	return nil
}

func splitMember(member string) []string {
	for i := 0; i < len(member); i++ {
		if member[i] == ':' {
			return []string{member[:i], member[i+1:]}
		}
	}
	return nil
}
