import { describe, expect, it } from 'vitest';
import { formatDistance, formatRouteDuration, mobilityError, persistedCommandRoute, positionFreshness, validateNearby, validateRoute } from './model';
import { PolarisApiError } from '../../lib/api/errors';
import type { Command } from '../../types/domain';

describe('mobility model', () => {
  it('validates coordinates and repository-default nearby bounds', () => {
    expect(validateNearby(13.08, 80.27, 5000, 20)).toEqual({});
    expect(validateNearby(91, 80, 10001, 51)).toEqual(expect.objectContaining({ coordinates: expect.any(String), radius: expect.any(String), limit: expect.any(String) }));
  });
  it('converts distance and Go duration units for presentation', () => {
    expect(formatDistance(850)).toBe('850 m'); expect(formatDistance(8400)).toBe('8.40 km'); expect(formatRouteDuration(1_062_000_000_000)).toBe('17m 42s');
  });
  it('enforces the current road-vehicle-only route profile', () => {
    const base = { mobility_profile: 'ROAD_VEHICLE' as const, origin: { latitude: 13, longitude: 80 }, destination: { latitude: 13.1, longitude: 80.1 }, policy: 'FASTEST' as const };
    expect(validateRoute(base)).toEqual({}); expect(validateRoute({ ...base, mobility_profile: 'AERIAL_DRONE' }).profile).toContain('road vehicles');
  });
  it('differentiates routing overload and preserves request identity', () => {
    const result = mobilityError(new PolarisApiError({ code: 'ROUTING_BUSY', message: 'busy', status: 429, requestId: 'req-7' }));
    expect(result.title).toContain('saturated'); expect(result.message).toContain('protect'); expect(result.requestId).toBe('req-7');
  });
  it('labels stale and offline positions without claiming they are current', () => {
    expect(positionFreshness('STALE')).toBe('Position may be stale'); expect(positionFreshness('OFFLINE')).toBe('Last known position');
  });
  it('parses persisted command geometry without producing a route for incomplete payloads', () => {
    const command = { command_id:'cmd', payload:{route_id:'r1',road_graph_version:'g1',routing_snapshot_version:9,policy:'FASTEST',origin:{latitude:13,longitude:80},destination:{latitude:13.1,longitude:80.1},waypoints:[{latitude:13,longitude:80},{latitude:13.1,longitude:80.1}],distance_meters:1200,estimated_duration_ms:90000} } as unknown as Command;
    expect(persistedCommandRoute(command)).toEqual(expect.objectContaining({route_id:'r1',snapshot_version:9,estimated_time:90_000_000_000}));
    expect(persistedCommandRoute({ ...command, payload: {} })).toBeUndefined();
  });
});
