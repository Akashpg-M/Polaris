package spatial

import (
	"hash/fnv"
	"sort"
	"strings"
	"sync"

	"github.com/Akashpg-M/polaris/backend/algo_/geo"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/Akashpg-M/polaris/backend/internal/core/events"
	"google.golang.org/protobuf/proto"
)

// ShardCount dictates how many independent memory partitions exist.
// 32 is a standard default to minimize lock contention across highly concurrent goroutines.
const ShardCount = 32

type EngineShard struct {
	mu       sync.RWMutex
	nodes    map[string]*pb.SpatialObject
	versions map[string]stateVersion
}

type stateVersion struct {
	bootID        string
	bootStartedAt int64
	sequence      uint64
	retired       map[string]struct{}
}

type Classification string

const (
	Accepted     Classification = "ACCEPTED"
	Duplicate    Classification = "DUPLICATE"
	OutOfOrder   Classification = "OUT_OF_ORDER"
	NewBoot      Classification = "NEW_BOOT"
	RetiredBoot  Classification = "RETIRED_BOOT"
	BootConflict Classification = "BOOT_CONFLICT"
)

type Engine struct {
	shards []*EngineShard
}

// MatchResult is the DTO sent back to the dispatcher
type MatchResult struct {
	NodeID     string      `json:"node_id"`
	Type       pb.NodeType `json:"node_type"`
	Class      uint16      `json:"asset_class"`
	Lat        float64     `json:"lat"`
	Lon        float64     `json:"lon"`
	DistanceKm float64     `json:"distance_km"`
	ETASec     int         `json:"eta_seconds"`
	RouteType  string      `json:"route_type"`
}

func NewEngine() *Engine {
	shards := make([]*EngineShard, ShardCount)
	for i := 0; i < ShardCount; i++ {
		shards[i] = &EngineShard{nodes: make(map[string]*pb.SpatialObject), versions: make(map[string]stateVersion)}
	}
	return &Engine{shards: shards}
}

// getShard picks the correct memory partition using an FNV-1a hash of the NodeID
func (e *Engine) getShard(nodeID string) *EngineShard {
	h := fnv.New32a()
	h.Write([]byte(nodeID))
	return e.shards[h.Sum32()%ShardCount]
}

func (e *Engine) BatchUpdate(payloads []*pb.SpatialObject) {
	if len(payloads) == 0 {
		return
	}

	for _, p := range payloads {
		shard := e.getShard(p.Id)

		shard.mu.Lock()
		shard.nodes[p.Id] = p
		shard.mu.Unlock()
	}
}

// FindNearest is a compatibility projection for the Phase 0 endpoint. Redis is
// the freshness authority and Mobility owns indexed candidate discovery. This
// bounded linear scan deliberately avoids retaining the former non-subdividing
// "QuadTree" as a second production spatial authority.
func (e *Engine) FindNearest(tenantID string, lat, lon, radiusKm float64, reqType pb.NodeType) []MatchResult {
	var results []MatchResult
	for _, shard := range e.shards {
		shard.mu.RLock()
		for key, node := range shard.nodes {
			if node.TenantId != tenantID || node.Type != reqType {
				continue
			}
			dist := geo.Haversine(lat, lon, node.Lat, node.Lon)
			if dist <= radiusKm {
				results = append(results, MatchResult{NodeID: strings.TrimPrefix(key, tenantID+":"), Type: node.Type, Class: uint16(node.Type), Lat: node.Lat, Lon: node.Lon, DistanceKm: dist, ETASec: int((dist / 40.0) * 3600), RouteType: "COMPATIBILITY_ESTIMATE"})
			}
		}
		shard.mu.RUnlock()
	}

	sort.Slice(results, func(i, j int) bool { return results[i].ETASec < results[j].ETASec })
	if len(results) > 500 {
		return results[:500]
	}
	return results
}

func stateKey(tenantID, deviceID string) string { return tenantID + ":" + deviceID }

// ApplyEnvelope is a defensive projection guard. KafkaConsumer invokes it only
// after Redis—the canonical Phase 1 freshness authority—accepts the event.
func (e *Engine) ApplyEnvelope(envelope *events.TelemetryEnvelope) Classification {
	key := stateKey(envelope.TenantID, envelope.DeviceID)
	shard := e.getShard(key)
	shard.mu.Lock()
	current, exists := shard.versions[key]
	classification := Accepted
	apply := false
	if !exists {
		current = stateVersion{bootID: envelope.DeviceBootID, bootStartedAt: envelope.BootStartedAt, sequence: envelope.SequenceNumber, retired: make(map[string]struct{})}
		apply = true
	} else if envelope.DeviceBootID == current.bootID {
		switch {
		case envelope.SequenceNumber > current.sequence:
			current.sequence = envelope.SequenceNumber
			apply = true
		case envelope.SequenceNumber == current.sequence:
			classification = Duplicate
		default:
			classification = OutOfOrder
		}
	} else if _, retired := current.retired[envelope.DeviceBootID]; retired {
		classification = RetiredBoot
	} else if envelope.BootStartedAt > current.bootStartedAt {
		current.retired[current.bootID] = struct{}{}
		current.bootID = envelope.DeviceBootID
		current.bootStartedAt = envelope.BootStartedAt
		current.sequence = envelope.SequenceNumber
		classification = NewBoot
		apply = true
	} else if envelope.BootStartedAt == current.bootStartedAt {
		classification = BootConflict
	} else {
		classification = RetiredBoot
	}
	if apply {
		payload := proto.Clone(envelope.Payload).(*pb.SpatialObject)
		shard.nodes[key] = payload
		shard.versions[key] = current
	}
	shard.mu.Unlock()
	return classification
}
