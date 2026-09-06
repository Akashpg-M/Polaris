import { describe, expect, it } from 'vitest';
import { parseDashboardTelemetry } from './events';

const valid = {
  event_id: 'evt-1',
  device_id: 'vehicle-7',
  tenant_id: 'tenant-a',
  device_boot_id: 'boot-2',
  sequence_number: 42,
  lat: 13.041,
  lon: 80.233,
  energy_percent: 74,
  observed_at: 1_750_000_000_000,
};

describe('parseDashboardTelemetry', () => {
  it('accepts a well-formed tenant dashboard event', () => {
    expect(parseDashboardTelemetry(valid)?.device_id).toBe('vehicle-7');
  });

  it.each([
    [{ ...valid, tenant_id: '' }],
    [{ ...valid, lat: 91 }],
    [{ ...valid, lon: -181 }],
    [{ ...valid, energy_percent: 101 }],
    [{ ...valid, sequence_number: -1 }],
    [{ ...valid, observed_at: Number.NaN }],
  ])('rejects an invalid or incomplete event', value => {
    expect(parseDashboardTelemetry(value)).toBeNull();
  });
});
