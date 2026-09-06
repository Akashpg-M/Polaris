import type { DashboardTelemetry } from '../../types/domain';

function finite(value: unknown): value is number {
  return typeof value === 'number' && Number.isFinite(value);
}

export function parseDashboardTelemetry(raw: unknown): DashboardTelemetry | null {
  if (!raw || typeof raw !== 'object') return null;
  const value = raw as Partial<DashboardTelemetry>;
  if (!value.event_id || !value.device_id || !value.tenant_id || !value.device_boot_id) return null;
  if (!finite(value.sequence_number) || !finite(value.lat) || !finite(value.lon) || !finite(value.energy_percent) || !finite(value.observed_at)) return null;
  if (value.lat < -90 || value.lat > 90 || value.lon < -180 || value.lon > 180) return null;
  if (value.energy_percent < 0 || value.energy_percent > 100 || value.sequence_number < 0 || value.observed_at <= 0) return null;
  return value as DashboardTelemetry;
}
