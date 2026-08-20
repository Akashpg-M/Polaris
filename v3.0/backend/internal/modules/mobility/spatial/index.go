package spatial

import (
	"sort"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type SpatialCandidate struct {
	State          model.SpatialState `json:"state"`
	DistanceMeters float64            `json:"distance_meters"`
}

type SpatialIndex interface {
	Upsert(model.SpatialState) error
	Remove(tenantID, deviceID string) error
	Nearest(target model.Position, limit int, radiusMeters float64) ([]SpatialCandidate, error)
	WithinRadius(target model.Position, radiusMeters float64) ([]SpatialCandidate, error)
}

type LinearSpatialIndex struct{ states map[string]model.SpatialState }

func NewLinearSpatialIndex() *LinearSpatialIndex {
	return &LinearSpatialIndex{states: map[string]model.SpatialState{}}
}
func indexKey(t, d string) string { return t + "\x00" + d }
func (i *LinearSpatialIndex) Upsert(s model.SpatialState) error {
	if err := ValidatePosition(s.Position); err != nil {
		return err
	}
	i.states[indexKey(s.TenantID, s.DeviceID)] = s
	return nil
}
func (i *LinearSpatialIndex) Remove(t, d string) error { delete(i.states, indexKey(t, d)); return nil }
func (i *LinearSpatialIndex) WithinRadius(p model.Position, radius float64) ([]SpatialCandidate, error) {
	if err := ValidatePosition(p); err != nil {
		return nil, err
	}
	out := []SpatialCandidate{}
	for _, s := range i.states {
		d := DistanceMeters(p, s.Position)
		if radius <= 0 || d <= radius {
			out = append(out, SpatialCandidate{s, d})
		}
	}
	sortCandidates(out)
	return out, nil
}
func (i *LinearSpatialIndex) Nearest(p model.Position, limit int, radius float64) ([]SpatialCandidate, error) {
	out, e := i.WithinRadius(p, radius)
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, e
}
func sortCandidates(v []SpatialCandidate) {
	sort.Slice(v, func(a, b int) bool {
		if v[a].DistanceMeters != v[b].DistanceMeters {
			return v[a].DistanceMeters < v[b].DistanceMeters
		}
		return v[a].State.DeviceID < v[b].State.DeviceID
	})
}
