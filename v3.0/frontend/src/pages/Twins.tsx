import { ArrowUpRight, Boxes, Cpu, MapPin } from 'lucide-react';
import { Link, useSearchParams } from 'react-router-dom';
import { FleetFiltersBar } from '../components/fleet/FleetFiltersBar';
import { BatteryIndicator, ConnectivityBadge, DeviceTypeIcon, LifecycleBadge } from '../components/status/Status';
import { EmptyState, ErrorState, PageLoader } from '../components/ui/States';
import { PageHeader } from '../components/ui/PageHeader';
import { deviceTypeLabel, formatCoordinates, relativeTime } from '../lib/format';
import type { ConnectivityStatus, FleetFilters } from '../types/domain';
import { useProjects, useTwins } from '../features/fleet/queries';
import { useProjectContext } from '../features/projects/ProjectContext';
import { filterLoadedTwins } from '../features/fleet/filtering';

export default function Twins() {
  const [params, setParams] = useSearchParams();
  const { projectId, setProjectId } = useProjectContext();
  const filters: FleetFilters = { connectivity: (params.get('connectivity') || '') as ConnectivityStatus | '', projectId: params.get('project') || projectId, deviceType: params.get('type') || '', search: params.get('q') || '' };
  const twins = useTwins({ ...filters, search: undefined }, 100);
  const projects = useProjects();
  const update = (next: FleetFilters) => { const value = new URLSearchParams(); if (next.connectivity) value.set('connectivity', next.connectivity); if (next.projectId) value.set('project', next.projectId); if (next.deviceType) value.set('type', next.deviceType); if (next.search) value.set('q', next.search); setProjectId(next.projectId || ''); setParams(value); };
  if (twins.isLoading) return <PageLoader label="Loading digital twins" />;
  if (twins.error) return <ErrorState error={twins.error} retry={() => void twins.refetch()} />;
  const items = filterLoadedTwins(twins.data || [], filters.search);
  return <div className="page-stack"><PageHeader eyebrow="Operational state" title="Digital twins" description="Durable registry identity combined with the latest Redis-reported components." actions={<span className="loaded-scope">Up to 100 twins</span>} /><FleetFiltersBar filters={filters} projects={projects.data || []} onChange={update} showConnectivity />
    {!items.length ? <EmptyState title="No twins in this view" message="Registered devices can exist without telemetry; try the Never connected filter." /> : <div className="twin-grid">{items.map(twin => <Link className="twin-card" to={`/devices/${encodeURIComponent(twin.device_id)}?tab=twin`} key={twin.device_id}>
      <div className="twin-card-head"><span className="large-device-icon"><DeviceTypeIcon type={twin.device.device_type_id} /></span><div><strong>{twin.device.display_name}</strong><small>{deviceTypeLabel(twin.device.device_type_id)}</small></div><ArrowUpRight /></div>
      <div className="badge-row"><LifecycleBadge status={twin.device.lifecycle_status} /><ConnectivityBadge status={twin.connectivity.status} /></div>
      <div className="twin-card-data"><div><BatteryIndicator value={twin.reported_state?.energy_percent} /><span>Battery</span></div><div><MapPin /><strong>{formatCoordinates(twin.reported_state?.lat, twin.reported_state?.lon)}</strong><span>Position</span></div><div><Cpu /><strong>{Object.keys(twin.components).length}</strong><span>Components</span></div></div>
      <div className="twin-card-footer"><code>{twin.device_id}</code><span>{relativeTime(twin.connectivity.last_seen_at)}</span></div>
    </Link>)}</div>}
    <p className="architecture-note"><Boxes /> Registry metadata is durable in PostgreSQL. Reported components are the latest Redis projection and can be reconstructed from accepted telemetry.</p>
  </div>;
}
