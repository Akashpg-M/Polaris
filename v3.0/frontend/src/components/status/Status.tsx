import { Battery, Camera, CarFront, CircleGauge, Cpu, Plane, RadioTower, SquareActivity, Truck } from 'lucide-react';
import { connectivityLabel, deviceTypeLabel, lifecycleLabel } from '../../lib/format';
import type { ConnectivityStatus, DeviceLifecycle, SocketStatus } from '../../types/domain';

export function ConnectivityBadge({ status }: { status: ConnectivityStatus }) {
  return <span className={`status-badge connectivity-${status.toLowerCase()}`}><span aria-hidden="true" className="status-shape" />{connectivityLabel(status)}</span>;
}

export function LifecycleBadge({ status }: { status: DeviceLifecycle }) {
  return <span className={`status-badge lifecycle-${status.toLowerCase()}`}>{lifecycleLabel(status)}</span>;
}

export function ProjectStatusBadge({ status }: { status: string }) {
  return <span className={`status-badge project-${status.toLowerCase()}`}>{status.charAt(0) + status.slice(1).toLowerCase()}</span>;
}

export function WebSocketStatus({ status, compact = false }: { status: SocketStatus; compact?: boolean }) {
  const label = status === 'CONNECTED' ? 'Live' : status === 'RECONNECTING' ? 'Reconnecting' : 'Disconnected';
  return <span className={`socket-status socket-${status.toLowerCase()}`} title={`Dashboard stream: ${label}`}>
    <span className="pulse-dot" aria-hidden="true" />{compact ? <span className="sr-only">{label}</span> : label}
  </span>;
}

export function BatteryIndicator({ value }: { value?: number }) {
  const valid = typeof value === 'number' && Number.isFinite(value);
  const bounded = valid ? Math.max(0, Math.min(100, value)) : 0;
  return <span className={`battery ${valid && bounded <= 20 ? 'battery-low' : ''}`} title={valid ? `${Math.round(bounded)}% battery` : 'Battery not reported'}>
    <Battery size={15} aria-hidden="true" /> {valid ? `${Math.round(bounded)}%` : '—'}
  </span>;
}

const icons = {
  connected_vehicle: CarFront,
  delivery_drone: Plane,
  ground_robot: SquareActivity,
  fixed_iot_sensor: RadioTower,
  static_camera: Camera,
  compute_node: Cpu,
};

export function DeviceTypeIcon({ type, size = 18 }: { type: string; size?: number }) {
  const Icon = icons[type as keyof typeof icons] || Truck;
  return <span className="device-type-icon" title={deviceTypeLabel(type)}><Icon size={size} aria-hidden="true" /><span className="sr-only">{deviceTypeLabel(type)}</span></span>;
}

export function MetricIcon({ kind }: { kind: 'fleet' | 'battery' }) {
  return kind === 'battery' ? <Battery aria-hidden="true" /> : <CircleGauge aria-hidden="true" />;
}

