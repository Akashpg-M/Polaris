package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"github.com/redis/go-redis/v9"
)

var latestStateScript = redis.NewScript(`
local current_boot = redis.call('HGET', KEYS[1], 'device_boot_id')
local incoming_boot = ARGV[1]
local incoming_started = tonumber(ARGV[2])
local incoming_seq = ARGV[3]
local function compare_decimal(a, b)
  a = string.gsub(a, '^0+', ''); if a == '' then a = '0' end
  b = string.gsub(b, '^0+', ''); if b == '' then b = '0' end
  if string.len(a) < string.len(b) then return -1 end
  if string.len(a) > string.len(b) then return 1 end
  if a < b then return -1 end
  if a > b then return 1 end
  return 0
end
local classification = 'ACCEPTED'
if current_boot then
  local current_started = tonumber(redis.call('HGET', KEYS[1], 'boot_started_at'))
  local current_seq = redis.call('HGET', KEYS[1], 'sequence_number')
  if current_boot == incoming_boot then
    local sequence_order = compare_decimal(incoming_seq, current_seq)
    if sequence_order == 0 then return 'DUPLICATE' end
    if sequence_order < 0 then return 'OUT_OF_ORDER' end
  elseif redis.call('SISMEMBER', KEYS[2], incoming_boot) == 1 then
    return 'RETIRED_BOOT'
  elseif incoming_started > current_started then
    redis.call('SADD', KEYS[2], current_boot)
    classification = 'NEW_BOOT'
  elseif incoming_started == current_started then
    return 'BOOT_CONFLICT'
  else
    return 'RETIRED_BOOT'
  end
end
redis.call('HSET', KEYS[1],
  'device_boot_id', incoming_boot,
  'boot_started_at', ARGV[2],
  'sequence_number', ARGV[3],
  'event_id', ARGV[4],
  'reported_state', ARGV[5],
  'last_seen_at', ARGV[6],
  'connectivity_status', 'ONLINE')
redis.call('ZADD', KEYS[3], ARGV[6], ARGV[8])
redis.call('PUBLISH', ARGV[7], ARGV[5])
return classification
`)

type RedisProjector struct{ client *redis.Client }

func NewRedisProjector(client *redis.Client) *RedisProjector { return &RedisProjector{client: client} }

func dashboardEnvelopeJSON(e *events.TelemetryEnvelope) ([]byte, error) {
	p := e.Payload
	return json.Marshal(map[string]interface{}{
		"event_id": e.EventID, "schema_version": e.SchemaVersion,
		"id": p.Id, "device_id": e.DeviceID, "tenant_id": e.TenantID,
		"device_boot_id": e.DeviceBootID, "sequence_number": e.SequenceNumber,
		"boot_started_at": e.BootStartedAt,
		"type":            p.Type, "status": p.Status, "lat": p.Lat, "lon": p.Lon,
		"velocity_mps": p.VelocityMps, "heading_deg": p.HeadingDeg,
		"energy_percent": p.EnergyPercent, "observed_at": e.ObservedAt,
		"ingested_at": e.IngestedAt, "timestamp": e.ObservedAt,
	})
}

func (p *RedisProjector) Apply(ctx context.Context, e *events.TelemetryEnvelope) (spatial.Classification, error) {
	data, err := dashboardEnvelopeJSON(e)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("polaris:twin:%s:%s", e.TenantID, e.DeviceID)
	result, err := latestStateScript.Run(ctx, p.client, []string{key, key + ":retired_boots", "polaris:devices:last-seen"},
		e.DeviceBootID, e.BootStartedAt, e.SequenceNumber, e.EventID, string(data), time.Now().UTC().UnixMilli(), DashboardUpdatesChannel, e.TenantID+":"+e.DeviceID).Text()
	return spatial.Classification(result), err
}

func (p *RedisProjector) Ready(ctx context.Context) error { return p.client.Ping(ctx).Err() }
