package spatial

import (
	"errors"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	h3 "github.com/uber/h3-go/v4"
)

var (
	ErrStaleVersion   = errors.New("spatial source version is not newer")
	ErrTenantCapacity = errors.New("tenant spatial capacity reached")
)

type ManagerConfig struct {
	H3Resolution, ShardResolution int
	MinMoveMeters                 float64
	MaxIndexAge                   time.Duration
	MaxH3Rings                    int
	MaxRadiusMeters               float64
	MaxDevicesPerTenant           int
}
type location struct {
	region uint64
	state  model.SpatialState
}
type MobilityShard struct {
	mu       sync.RWMutex
	RegionID uint64
	Index    SpatialIndex
	devices  map[string]model.SpatialState
}

type Manager struct {
	cfg         ManagerConfig
	shardsMu    sync.RWMutex
	shards      map[string]map[uint64]*MobilityShard
	locationsMu sync.RWMutex
	locations   map[string]location
	deviceLocks [64]sync.Mutex
}

func NewManager(c ManagerConfig) *Manager {
	return &Manager{cfg: c, shards: map[string]map[uint64]*MobilityShard{}, locations: map[string]location{}}
}
func hashDevice(v string) uint {
	var h uint = 2166136261
	for i := 0; i < len(v); i++ {
		h = (h ^ uint(v[i])) * 16777619
	}
	return h
}
func (m *Manager) cell(p model.Position) (uint64, uint64, error) {
	c, e := h3.LatLngToCell(h3.NewLatLng(p.Latitude, p.Longitude), m.cfg.H3Resolution)
	if e != nil {
		return 0, 0, e
	}
	parent, e := c.Parent(m.cfg.ShardResolution)
	return uint64(c), uint64(parent), e
}
func (m *Manager) getShard(tenant string, region uint64, create bool) *MobilityShard {
	m.shardsMu.RLock()
	s := m.shards[tenant][region]
	m.shardsMu.RUnlock()
	if s != nil || !create {
		return s
	}
	m.shardsMu.Lock()
	defer m.shardsMu.Unlock()
	if m.shards[tenant] == nil {
		m.shards[tenant] = map[uint64]*MobilityShard{}
	}
	if m.shards[tenant][region] == nil {
		m.shards[tenant][region] = &MobilityShard{RegionID: region, Index: NewRTreeSpatialIndex(), devices: map[string]model.SpatialState{}}
	}
	return m.shards[tenant][region]
}

func (m *Manager) Upsert(in model.SpatialState) error {
	if e := ValidatePosition(in.ReportedPosition); e != nil {
		return e
	}
	cell, region, e := m.cell(in.ReportedPosition)
	if e != nil {
		return e
	}
	in.H3Cell = cell
	lock := &m.deviceLocks[hashDevice(in.TenantID+"\x00"+in.DeviceID)%uint(len(m.deviceLocks))]
	lock.Lock()
	defer lock.Unlock()
	k := indexKey(in.TenantID, in.DeviceID)
	m.locationsMu.RLock()
	old, exists := m.locations[k]
	m.locationsMu.RUnlock()
	if exists && !in.NewerThan(old.state) {
		return ErrStaleVersion
	}
	in.Quality = ValidateObservation(func() *model.SpatialState {
		if exists {
			return &old.state
		}
		return nil
	}(), in)
	if in.HeadingDegrees != nil {
		normalized := math.Mod(*in.HeadingDegrees, 360)
		if normalized < 0 {
			normalized += 360
		}
		in.HeadingDegrees = &normalized
	}
	if !in.Quality.Valid {
		return errors.New("invalid mobility observation")
	}
	now := time.Now().UTC()
	if in.IndexedAt.IsZero() {
		in.IndexedAt = now
	}
	in.Position = in.ReportedPosition
	if exists && old.region == region && old.state.H3Cell == cell && DistanceMeters(old.state.Position, in.ReportedPosition) < m.cfg.MinMoveMeters && now.Sub(old.state.IndexedAt) < m.cfg.MaxIndexAge {
		in.Position = old.state.Position
		in.IndexedAt = old.state.IndexedAt
	}
	if !exists {
		m.locationsMu.RLock()
		count := 0
		for key := range m.locations {
			if len(key) > len(in.TenantID) && key[:len(in.TenantID)] == in.TenantID && key[len(in.TenantID)] == '\x00' {
				count++
			}
		}
		m.locationsMu.RUnlock()
		if count >= m.cfg.MaxDevicesPerTenant {
			return ErrTenantCapacity
		}
	}
	newShard := m.getShard(in.TenantID, region, true)
	oldShard := newShard
	if exists && old.region != region {
		oldShard = m.getShard(in.TenantID, old.region, false)
	}
	locked := []*MobilityShard{}
	if oldShard != nil && oldShard != newShard {
		if oldShard.RegionID < newShard.RegionID {
			locked = []*MobilityShard{oldShard, newShard}
		} else {
			locked = []*MobilityShard{newShard, oldShard}
		}
	} else {
		locked = []*MobilityShard{newShard}
	}
	for _, s := range locked {
		s.mu.Lock()
	}
	defer func() {
		for i := len(locked) - 1; i >= 0; i-- {
			locked[i].mu.Unlock()
		}
	}()
	if exists && oldShard != nil && oldShard != newShard {
		_ = oldShard.Index.Remove(in.TenantID, in.DeviceID)
		delete(oldShard.devices, in.DeviceID)
	}
	if e = newShard.Index.Upsert(in); e != nil {
		return e
	}
	newShard.devices[in.DeviceID] = in
	m.locationsMu.Lock()
	m.locations[k] = location{region, in}
	m.locationsMu.Unlock()
	return nil
}

func (m *Manager) Remove(tenant, device string) error {
	k := indexKey(tenant, device)
	lock := &m.deviceLocks[hashDevice(k)%uint(len(m.deviceLocks))]
	lock.Lock()
	defer lock.Unlock()
	m.locationsMu.Lock()
	loc, ok := m.locations[k]
	if ok {
		delete(m.locations, k)
	}
	m.locationsMu.Unlock()
	if !ok {
		return nil
	}
	s := m.getShard(tenant, loc.region, false)
	if s == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.devices, device)
	return s.Index.Remove(tenant, device)
}
func (m *Manager) RemoveTenant(tenant string) error {
	m.locationsMu.RLock()
	devices := []string{}
	for key := range m.locations {
		if strings.HasPrefix(key, tenant+"\x00") {
			devices = append(devices, strings.TrimPrefix(key, tenant+"\x00"))
		}
	}
	m.locationsMu.RUnlock()
	sort.Strings(devices)
	for _, device := range devices {
		if err := m.Remove(tenant, device); err != nil {
			return err
		}
	}
	return nil
}
func (m *Manager) Get(tenant, device string) (model.SpatialState, bool) {
	m.locationsMu.RLock()
	defer m.locationsMu.RUnlock()
	v, ok := m.locations[indexKey(tenant, device)]
	return v.state, ok
}

func (m *Manager) Nearby(tenant string, p model.Position, radius float64, limit int) ([]SpatialCandidate, error) {
	if e := ValidatePosition(p); e != nil {
		return nil, e
	}
	if radius <= 0 || radius > m.cfg.MaxRadiusMeters {
		return nil, errors.New("search radius outside configured bounds")
	}
	cell, _, e := m.cell(p)
	if e != nil {
		return nil, e
	}
	rings := int(math.Ceil(radius/700)) + 1
	if rings > m.cfg.MaxH3Rings {
		rings = m.cfg.MaxH3Rings
	}
	cells, e := h3.Cell(cell).GridDisk(rings)
	if e != nil {
		return nil, e
	}
	regions := map[uint64]struct{}{}
	for _, c := range cells {
		parent, pe := c.Parent(m.cfg.ShardResolution)
		if pe == nil {
			regions[uint64(parent)] = struct{}{}
		}
	}
	found := map[string]SpatialCandidate{}
	for region := range regions {
		s := m.getShard(tenant, region, false)
		if s == nil {
			continue
		}
		s.mu.RLock()
		v, qerr := s.Index.WithinRadius(p, radius)
		s.mu.RUnlock()
		if qerr != nil {
			return nil, qerr
		}
		for _, c := range v {
			if old, ok := found[c.State.DeviceID]; !ok || c.DistanceMeters < old.DistanceMeters {
				found[c.State.DeviceID] = c
			}
		}
	}
	out := make([]SpatialCandidate, 0, len(found))
	for _, v := range found {
		out = append(out, v)
	}
	sortCandidates(out)
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

func (m *Manager) EvictInactive(tenant, device, lifecycle, connectivity string) error {
	if lifecycle != "ACTIVE" || connectivity == "STALE" || connectivity == "OFFLINE" {
		return m.Remove(tenant, device)
	}
	return nil
}
func (m *Manager) Rebuild(states []model.SpatialState) error {
	sort.Slice(states, func(i, j int) bool {
		if states[i].TenantID != states[j].TenantID {
			return states[i].TenantID < states[j].TenantID
		}
		return states[i].DeviceID < states[j].DeviceID
	})
	for _, s := range states {
		if e := m.Upsert(s); e != nil && !errors.Is(e, ErrStaleVersion) {
			return e
		}
	}
	return nil
}
func (m *Manager) Snapshot() []model.SpatialState {
	m.locationsMu.RLock()
	defer m.locationsMu.RUnlock()
	out := make([]model.SpatialState, 0, len(m.locations))
	for _, v := range m.locations {
		out = append(out, v.state)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].TenantID != out[j].TenantID {
			return out[i].TenantID < out[j].TenantID
		}
		return out[i].DeviceID < out[j].DeviceID
	})
	return out
}
