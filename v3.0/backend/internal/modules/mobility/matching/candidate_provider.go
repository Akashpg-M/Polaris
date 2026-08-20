package matching

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/routing"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
)

const TargetPositionKey = "mobility.target_position"

type Provider struct {
	Spatial               *spatial.Manager
	Routing               routing.RoutingEngine
	RawLimit, RoutedLimit int
	MaxRadius             float64
}

func (p *Provider) Name() string { return "mobility.candidates/v1" }
func (p *Provider) Supports(req extension.CandidateRequest) bool {
	_, latOK := numeric(req.Context["target_latitude"])
	_, lonOK := numeric(req.Context["target_longitude"])
	return latOK && lonOK && p.Spatial != nil
}
func (p *Provider) Candidates(ctx context.Context, req extension.CandidateRequest) ([]extension.Candidate, error) {
	lat, latOK := numeric(req.Context["target_latitude"])
	lon, lonOK := numeric(req.Context["target_longitude"])
	if !latOK || !lonOK {
		return nil, errors.New("mobility target is missing")
	}
	target := model.Position{Latitude: lat, Longitude: lon}
	radius := p.MaxRadius
	if v, ok := req.Context["maximum_distance_meters"].(float64); ok && v > 0 && v < radius {
		radius = v
	}
	near, err := p.Spatial.Nearby(req.TenantID, target, radius, p.RawLimit)
	if err != nil {
		return nil, err
	}
	type ranked struct {
		candidate     extension.Candidate
		eta, distance float64
		hasRoute      bool
	}
	rankedItems := make([]ranked, 0, len(near))
	eligible := make(map[string]struct{}, len(req.EligibleDeviceIDs))
	for _, id := range req.EligibleDeviceIDs {
		eligible[id] = struct{}{}
	}
	routed := 0
	for _, c := range near {
		if len(eligible) > 0 {
			if _, ok := eligible[c.State.DeviceID]; !ok {
				continue
			}
		}
		item := ranked{candidate: extension.Candidate{DeviceID: c.State.DeviceID, Attributes: map[string]any{"distance_meters": c.DistanceMeters, "spatial_observed_at": c.State.ObservedAt, "mobility_profile": c.State.MobilityProfile}}, distance: c.DistanceMeters}
		if p.Routing != nil && routed < p.RoutedLimit && c.State.MobilityProfile == model.MobilityRoadVehicle {
			routed++
			routingStarted := time.Now()
			route, e := p.Routing.Route(ctx, routing.RouteRequest{TenantID: req.TenantID, MobilityProfile: c.State.MobilityProfile, Origin: c.State.ReportedPosition, Destination: target, Policy: routing.RouteFastest})
			if req.Timing != nil {
				req.Timing.RoutingDuration += time.Since(routingStarted)
			}
			if e == nil {
				item.hasRoute = true
				item.eta = route.EstimatedTime.Seconds()
				item.candidate.Attributes["route_eta_seconds"] = item.eta
				item.candidate.Attributes["routing_snapshot_version"] = route.SnapshotVersion
			}
		}
		score := item.distance
		if item.hasRoute {
			score = item.eta
		}
		item.candidate.DomainScore = &score
		rankedItems = append(rankedItems, item)
	}
	sort.Slice(rankedItems, func(i, j int) bool {
		if rankedItems[i].hasRoute != rankedItems[j].hasRoute {
			return rankedItems[i].hasRoute
		}
		if *rankedItems[i].candidate.DomainScore != *rankedItems[j].candidate.DomainScore {
			return *rankedItems[i].candidate.DomainScore < *rankedItems[j].candidate.DomainScore
		}
		return rankedItems[i].candidate.DeviceID < rankedItems[j].candidate.DeviceID
	})
	limit := req.Limit
	if limit <= 0 || limit > len(rankedItems) {
		limit = len(rankedItems)
	}
	out := make([]extension.Candidate, limit)
	for i := 0; i < limit; i++ {
		out[i] = rankedItems[i].candidate
	}
	return out, nil
}

func numeric(v any) (float64, bool) { n, ok := v.(float64); return n, ok }
