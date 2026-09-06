import type { Command, ConnectivityStatus, MobilityPosition, MobilityProfile, RoutePolicy, RouteRequest, RouteResult } from '../../types/domain';
import { readableError } from '../../lib/api/errors';

export const DEFAULT_MAX_RADIUS_METERS = 10_000;
export const DEFAULT_MAX_NEARBY_RESULTS = 50;
export const mobilityProfiles: MobilityProfile[] = ['ROAD_VEHICLE', 'GROUND_ROBOT', 'AERIAL_DRONE', 'STATIC'];
export const routePolicies: RoutePolicy[] = ['FASTEST', 'SHORTEST'];

export function validPosition(position: MobilityPosition) {
  return Number.isFinite(position.latitude) && position.latitude >= -90 && position.latitude <= 90
    && Number.isFinite(position.longitude) && position.longitude >= -180 && position.longitude <= 180;
}

export function validateNearby(latitude: number, longitude: number, radiusMeters: number, limit: number) {
  const errors: Record<string, string> = {};
  if (!validPosition({ latitude, longitude })) errors.coordinates = 'Latitude must be -90 to 90 and longitude -180 to 180.';
  if (!Number.isFinite(radiusMeters) || radiusMeters <= 0 || radiusMeters > DEFAULT_MAX_RADIUS_METERS) errors.radius = 'Radius must be greater than 0 and no more than 10 km.';
  if (!Number.isInteger(limit) || limit < 1 || limit > DEFAULT_MAX_NEARBY_RESULTS) errors.limit = 'Limit must be a whole number from 1 to 50.';
  return errors;
}

export function validateRoute(input: RouteRequest) {
  const errors: Record<string, string> = {};
  if (!validPosition(input.origin)) errors.origin = 'Enter a valid origin.';
  if (!validPosition(input.destination)) errors.destination = 'Enter a valid destination.';
  if (input.mobility_profile !== 'ROAD_VEHICLE') errors.profile = 'Platform road routing currently supports road vehicles only.';
  if (!routePolicies.includes(input.policy)) errors.policy = 'Choose a supported route policy.';
  return errors;
}

export function formatDistance(metres?: number) {
  if (!Number.isFinite(metres)) return 'Not returned';
  return metres! < 1000 ? `${Math.round(metres!)} m` : `${(metres! / 1000).toFixed(metres! < 10_000 ? 2 : 1)} km`;
}

export function formatRouteDuration(nanoseconds?: number) {
  if (!Number.isFinite(nanoseconds) || nanoseconds! < 0) return 'Not returned';
  const seconds = Math.round(nanoseconds! / 1_000_000_000);
  if (seconds < 60) return `${seconds} sec`;
  const minutes = Math.floor(seconds / 60);
  const remainder = seconds % 60;
  return `${minutes}m ${remainder}s`;
}

export function profileLabel(profile: MobilityProfile) {
  return ({ ROAD_VEHICLE: 'Road vehicle', GROUND_ROBOT: 'Ground robot', AERIAL_DRONE: 'Aerial drone', STATIC: 'Static spatial device' })[profile];
}

export function positionFreshness(status: ConnectivityStatus) {
  if (status === 'ONLINE') return 'Latest reported position';
  if (status === 'STALE') return 'Position may be stale';
  if (status === 'OFFLINE') return 'Last known position';
  return 'No connected position';
}

const routingMessages: Record<string, { title: string; message: string }> = {
  ROUTING_BUSY: { title: 'Routing capacity is saturated', message: 'Polaris is rejecting excess work to protect core services. Try again shortly.' },
  ROUTING_TIMEOUT: { title: 'Route calculation timed out', message: 'The route search exceeded the configured time limit.' },
  ROUTING_UNAVAILABLE: { title: 'Mobility routing is unavailable', message: 'Core Polaris may still be operational. Check Mobility diagnostics.' },
  NO_ROUTE: { title: 'No road route found', message: 'No valid route could be found between these locations.' },
  NO_ROAD_NODE: { title: 'No supported road location', message: 'One or both coordinates could not be associated with a road node.' },
  OUTSIDE_ROUTING_REGION: { title: 'Outside the loaded road region', message: 'The selected coordinates fall outside the currently loaded road graph.' },
  OUTSIDE_REGION: { title: 'Outside the loaded road region', message: 'The selected coordinates fall outside the currently loaded road graph.' },
  UNSUPPORTED_MOBILITY_PROFILE: { title: 'Unsupported mobility profile', message: 'Platform road routing currently supports road vehicles only.' },
  UNSUPPORTED_PROFILE: { title: 'Unsupported mobility profile', message: 'Platform road routing currently supports road vehicles only.' },
  PLANNER_UNAVAILABLE: { title: 'Planner unavailable', message: 'Task planning is currently unavailable; no fallback route was generated.' },
};

export function mobilityError(error: unknown) {
  const api = readableError(error);
  return { ...(routingMessages[api.code || ''] || { title: api.code || 'Mobility request failed', message: api.message }), requestId: api.requestId, status: api.status };
}

function record(value: unknown): Record<string, unknown> | undefined { return value && typeof value === 'object' && !Array.isArray(value) ? value as Record<string, unknown> : undefined; }
function point(value: unknown): MobilityPosition | undefined {
  const item = record(value); const latitude = Number(item?.latitude); const longitude = Number(item?.longitude);
  return validPosition({ latitude, longitude }) ? { latitude, longitude } : undefined;
}

export function persistedCommandRoute(command?: Command): RouteResult | undefined {
  if (!command) return undefined;
  const payload = command.payload;
  const origin = point(payload.origin); const destination = point(payload.destination);
  const waypoints = Array.isArray(payload.waypoints) ? payload.waypoints.map(point).filter((item): item is MobilityPosition => Boolean(item)) : [];
  if (!payload.route_id || !origin || !destination || !waypoints.length) return undefined;
  return {
    route_id: String(payload.route_id), road_graph_version: String(payload.road_graph_version || 'Not returned'),
    snapshot_version: Number(payload.routing_snapshot_version || 0), policy: payload.policy === 'SHORTEST' ? 'SHORTEST' : 'FASTEST',
    distance_meters: Number(payload.distance_meters || 0), estimated_time: Number(payload.estimated_duration_ms || 0) * 1_000_000,
    waypoints, edge_ids: [], expanded_nodes: 0,
  };
}
