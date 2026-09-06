import { ArrowUpRight, MapPin, X } from 'lucide-react';
import { Link } from 'react-router-dom';
import { deviceTypeLabel, exactTime, formatCoordinates, formatSpeed, relativeTime } from '../../lib/format';
import type { DeviceTwin, Project } from '../../types/domain';
import { BatteryIndicator, ConnectivityBadge, DeviceTypeIcon, LifecycleBadge } from '../status/Status';

export function DeviceDrawer({ twin, projects, onClose }: { twin: DeviceTwin; projects: Project[]; onClose: () => void }) {
  const state = twin.reported_state;
  const project = projects.find(item => item.project_id === twin.device.project_id);
  return <aside className="device-drawer" aria-label={`Device ${twin.device.display_name}`}>
    <button className="drawer-close" onClick={onClose} aria-label="Close device panel"><X /></button>
    <div className="drawer-device"><span className="large-device-icon"><DeviceTypeIcon type={twin.device.device_type_id} size={26} /></span><div><p>{deviceTypeLabel(twin.device.device_type_id)}</p><h2>{twin.device.display_name}</h2><code>{twin.device_id}</code></div></div>
    <div className="badge-row"><LifecycleBadge status={twin.device.lifecycle_status} /><ConnectivityBadge status={twin.connectivity.status} /></div>
    <div className="drawer-grid"><div><span>Battery</span><strong><BatteryIndicator value={state?.energy_percent} /></strong></div><div><span>Speed</span><strong>{formatSpeed(state?.velocity_mps)}</strong></div><div><span>Last seen</span><strong title={exactTime(twin.connectivity.last_seen_at)}>{relativeTime(twin.connectivity.last_seen_at)}</strong></div><div><span>Project</span><strong>{project?.name || (twin.device.project_id ? 'Unknown project' : 'Unassigned')}</strong></div></div>
    {state ? <div className="drawer-location"><MapPin /><span><small>Current coordinates</small><strong>{formatCoordinates(state.lat, state.lon)}</strong></span></div> : <div className="drawer-notice">This device has not published accepted telemetry.</div>}
    <div><p className="section-label">Capabilities</p><div className="capability-list">{twin.capabilities.length ? twin.capabilities.map(capability => <span key={capability.capability_id}>{capability.display_name}</span>) : <small>No capabilities assigned</small>}</div></div>
    <Link className="button primary drawer-action" to={`/devices/${encodeURIComponent(twin.device_id)}`}>Open device <ArrowUpRight /></Link>
  </aside>;
}

