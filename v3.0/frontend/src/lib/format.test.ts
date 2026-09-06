import { describe, expect, it } from 'vitest';
import { deviceTypeLabel, formatCoordinates, formatSpeed, relativeTime } from './format';

describe('fleet presentation formatters', () => {
  it('keeps zero coordinates and battery-adjacent numeric values truthful', () => {
    expect(formatCoordinates(0, 0)).toBe('0.000000, 0.000000');
    expect(formatSpeed(0)).toBe('0.0 km/h');
  });

  it('does not invent missing telemetry', () => {
    expect(formatCoordinates()).toBe('Not reported');
    expect(relativeTime(null)).toBe('Never');
  });

  it('labels known and unknown device types consistently', () => {
    expect(deviceTypeLabel('delivery_drone')).toBe('Delivery drone');
    expect(deviceTypeLabel('warehouse_cart')).toBe('Warehouse Cart');
  });
});
