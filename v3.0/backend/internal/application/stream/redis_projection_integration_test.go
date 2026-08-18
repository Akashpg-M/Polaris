//go:build integration

package stream

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/redis/go-redis/v9"
)

func TestRedisProjectionAtomicClassifications(t *testing.T) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(opts)
	defer client.Close()
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis integration dependency unavailable: %v", err)
	}
	device := fmt.Sprintf("integration-%d", time.Now().UnixNano())
	key := "polaris:twin:tenant-1:" + device
	defer client.Del(ctx, key, key+":retired_boots")
	p := NewRedisProjector(client)
	started := time.Now().Add(-time.Minute).UnixMilli()
	first := streamEnvelope(device, 1)
	first.BootStartedAt = started
	first.Payload.BootStartedAt = started
	if got, err := p.Apply(ctx, first); err != nil || got != spatial.Accepted {
		t.Fatalf("first: %s %v", got, err)
	}
	if got, _ := p.Apply(ctx, first); got != spatial.Duplicate {
		t.Fatalf("duplicate: %s", got)
	}
	older := streamEnvelope(device, 1)
	older.SequenceNumber = 0
	older.Payload.SequenceNumber = 0
	if got, _ := p.Apply(ctx, older); got != spatial.OutOfOrder {
		t.Fatalf("out of order: %s", got)
	}
	newBoot := streamEnvelope(device, 1)
	newBoot.DeviceBootID = "boot-2"
	newBoot.Payload.DeviceBootId = "boot-2"
	newBoot.BootStartedAt = started + 1000
	newBoot.Payload.BootStartedAt = started + 1000
	if got, _ := p.Apply(ctx, newBoot); got != spatial.NewBoot {
		t.Fatalf("new boot: %s", got)
	}
	if got, _ := p.Apply(ctx, first); got != spatial.RetiredBoot {
		t.Fatalf("retired boot: %s", got)
	}
	conflict := streamEnvelope(device, 1)
	conflict.DeviceBootID = "boot-conflict"
	conflict.Payload.DeviceBootId = "boot-conflict"
	conflict.BootStartedAt = newBoot.BootStartedAt
	conflict.Payload.BootStartedAt = newBoot.BootStartedAt
	if got, _ := p.Apply(ctx, conflict); got != spatial.BootConflict {
		t.Fatalf("boot conflict: %s", got)
	}
	large := streamEnvelope(device, 9_007_199_254_740_993)
	large.DeviceBootID = "boot-2"
	large.Payload.DeviceBootId = "boot-2"
	large.BootStartedAt = newBoot.BootStartedAt
	large.Payload.BootStartedAt = newBoot.BootStartedAt
	if got, _ := p.Apply(ctx, large); got != spatial.Accepted {
		t.Fatalf("large sequence: %s", got)
	}
	largeOlder := streamEnvelope(device, 9_007_199_254_740_992)
	largeOlder.DeviceBootID = "boot-2"
	largeOlder.Payload.DeviceBootId = "boot-2"
	largeOlder.BootStartedAt = newBoot.BootStartedAt
	largeOlder.Payload.BootStartedAt = newBoot.BootStartedAt
	if got, _ := p.Apply(ctx, largeOlder); got != spatial.OutOfOrder {
		t.Fatalf("large out of order: %s", got)
	}
	values, err := client.HMGet(ctx, key, "device_boot_id", "sequence_number", "event_id", "reported_state", "last_seen_at").Result()
	if err != nil {
		t.Fatal(err)
	}
	if values[0] != "boot-2" || values[1] != "9007199254740993" || values[2] == nil || values[3] == nil || values[4] == nil {
		t.Fatalf("incomplete atomic twin hash: %#v", values)
	}
}
