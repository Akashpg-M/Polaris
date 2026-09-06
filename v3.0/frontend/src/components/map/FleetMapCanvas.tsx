import { useMemo } from 'react';
import type { DeviceTwin } from '../../types/domain';
import type { MobilityProfile } from '../../types/domain';
import { PolarisMap } from './PolarisMap';

function profile(type: string): MobilityProfile { return type === 'delivery_drone' ? 'AERIAL_DRONE' : type === 'ground_robot' ? 'GROUND_ROBOT' : type === 'fixed_iot_sensor' || type === 'static_camera' ? 'STATIC' : 'ROAD_VEHICLE'; }

export function FleetMapCanvas({ twins, selectedId, onSelect, compact = false }: {
  twins: DeviceTwin[];
  selectedId?: string;
  onSelect?: (twin: DeviceTwin) => void;
  compact?: boolean;
}) {
  const devices = useMemo(() => twins.flatMap(twin => twin.reported_state ? [{ id: twin.device_id, label: twin.device.display_name, position: { latitude: twin.reported_state.lat, longitude: twin.reported_state.lon }, profile: profile(twin.device.device_type_id), connectivity: twin.connectivity.status }] : []), [twins]);
  const byId = useMemo(() => new Map(twins.map(twin => [twin.device_id, twin])), [twins]);
  return <PolarisMap devices={devices} selectedId={selectedId} onDeviceSelect={id => { const twin = byId.get(id); if (twin) onSelect?.(twin); }} compact={compact} fitKey={`fleet:${devices.map(item => item.id).join(',')}`} />;
}
