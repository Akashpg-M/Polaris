import { describe, expect, it } from 'vitest';
import { filterLoadedTwins } from './filtering';

const items = [
  { device_id: 'truck-001', device: { display_name: 'Northbound Truck' } },
  { device_id: 'drone-009', device: { display_name: 'Depot Scout' } },
];

describe('filterLoadedTwins', () => {
  it('matches loaded device IDs without case sensitivity', () => {
    expect(filterLoadedTwins(items, 'TRUCK').map(item => item.device_id)).toEqual(['truck-001']);
  });

  it('matches display names and preserves all items for blank search', () => {
    expect(filterLoadedTwins(items, 'scout').map(item => item.device_id)).toEqual(['drone-009']);
    expect(filterLoadedTwins(items, ' ')).toEqual(items);
  });
});
