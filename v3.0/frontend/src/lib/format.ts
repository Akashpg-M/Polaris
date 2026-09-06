import type { ConnectivityStatus, DeviceLifecycle } from '../types/domain';

export function relativeTime(value?: string | number | null, now = Date.now()) {
  if (value === undefined || value === null || value === '') return 'Never';
  const timestamp = typeof value === 'number' ? value : Date.parse(value);
  if (!Number.isFinite(timestamp)) return 'Unknown';
  const seconds = Math.max(0, Math.floor((now - timestamp) / 1000));
  if (seconds < 10) return 'Just now';
  if (seconds < 60) return `${seconds} sec ago`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes} min ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hr ago`;
  const days = Math.floor(hours / 24);
  return `${days} day${days === 1 ? '' : 's'} ago`;
}

export function exactTime(value?: string | number | null) {
  if (value === undefined || value === null || value === '') return 'No timestamp';
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? 'Invalid timestamp' : date.toLocaleString();
}

export function formatCoordinates(lat?: number, lon?: number) {
  return Number.isFinite(lat) && Number.isFinite(lon) ? `${lat!.toFixed(6)}, ${lon!.toFixed(6)}` : 'Not reported';
}

export function formatSpeed(value?: number) {
  return Number.isFinite(value) ? `${(value! * 3.6).toFixed(1)} km/h` : 'Not reported';
}

export function lifecycleLabel(status: DeviceLifecycle) {
  return status.charAt(0) + status.slice(1).toLowerCase();
}

export function connectivityLabel(status: ConnectivityStatus) {
  return status === 'NEVER_CONNECTED' ? 'Never connected' : status.charAt(0) + status.slice(1).toLowerCase();
}

export const deviceTypeLabels: Record<string, string> = {
  connected_vehicle: 'Connected vehicle',
  delivery_drone: 'Delivery drone',
  ground_robot: 'Ground robot',
  fixed_iot_sensor: 'Fixed IoT sensor',
  static_camera: 'Static camera',
  compute_node: 'Compute node',
};

export function deviceTypeLabel(value: string) {
  return deviceTypeLabels[value] || value.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
}

