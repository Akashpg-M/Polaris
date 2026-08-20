package twin

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

var transitionScript = redis.NewScript(`
local seen=redis.call('HGET',KEYS[1],'last_seen_at')
if not seen or seen~=ARGV[1] then return 0 end
local current=redis.call('HGET',KEYS[1],'connectivity_status')
if current==ARGV[2] then return 0 end
if ARGV[2]=='STALE' and current~='ONLINE' then return 0 end
if ARGV[2]=='OFFLINE' and current~='ONLINE' and current~='STALE' then return 0 end
redis.call('HSET',KEYS[1],'connectivity_status',ARGV[2])
return 1`)

type Detector struct {
	redis                    *redis.Client
	writer                   *kafka.Writer
	stale, offline, interval time.Duration
	onTransition             func(tenant, device, status string)
}

func (d *Detector) SetTransitionHandler(fn func(tenant, device, status string)) { d.onTransition = fn }

func NewDetector(client *redis.Client, broker string, stale, offline, interval time.Duration) *Detector {
	return &Detector{redis: client, writer: &kafka.Writer{Addr: kafka.TCP(broker), Topic: "device.connectivity.v1", Balancer: &kafka.Hash{}}, stale: stale, offline: offline, interval: interval}
}
func (d *Detector) Start(ctx context.Context) {
	defer d.writer.Close()
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			d.scan(ctx, "OFFLINE", d.offline)
			d.scan(ctx, "STALE", d.stale)
		case <-ctx.Done():
			return
		}
	}
}
func (d *Detector) scan(ctx context.Context, next string, age time.Duration) {
	cutoff := time.Now().Add(-age).UnixMilli()
	entries, err := d.redis.ZRangeByScoreWithScores(ctx, "polaris:devices:last-seen", &redis.ZRangeBy{Min: "-inf", Max: strconv.FormatInt(cutoff, 10), Offset: 0, Count: 100}).Result()
	if err != nil {
		return
	}
	for _, z := range entries {
		member, ok := z.Member.(string)
		if !ok {
			continue
		}
		parts := strings.SplitN(member, ":", 2)
		if len(parts) != 2 {
			continue
		}
		last := strconv.FormatInt(int64(z.Score), 10)
		key := "polaris:twin:" + member
		changed, err := transitionScript.Run(ctx, d.redis, []string{key}, last, next).Int()
		if err != nil || changed != 1 {
			continue
		}
		if d.onTransition != nil {
			d.onTransition(parts[0], parts[1], next)
		}
		sum := sha256.Sum256([]byte(member + ":" + next + ":" + last))
		eventID := hex.EncodeToString(sum[:])
		payload, _ := json.Marshal(map[string]interface{}{"event_id": eventID, "event_type": "device.connectivity.changed.v1", "schema_version": 1, "tenant_id": parts[0], "device_id": parts[1], "connectivity_status": next, "last_seen_at": last})
		_ = d.writer.WriteMessages(ctx, kafka.Message{Key: []byte(member), Value: payload})
	}
}
