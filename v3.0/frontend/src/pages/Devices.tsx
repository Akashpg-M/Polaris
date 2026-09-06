import { ChevronLeft, ChevronRight, ExternalLink } from 'lucide-react';
import { Link, useSearchParams } from 'react-router-dom';
import { FleetFiltersBar } from '../components/fleet/FleetFiltersBar';
import { BatteryIndicator, ConnectivityBadge, DeviceTypeIcon, LifecycleBadge } from '../components/status/Status';
import { EmptyState, ErrorState, TableSkeleton } from '../components/ui/States';
import { PageHeader } from '../components/ui/PageHeader';
import { deviceTypeLabel, exactTime, formatCoordinates, relativeTime } from '../lib/format';
import type { ConnectivityStatus, DeviceLifecycle, FleetFilters } from '../types/domain';
import { useProjects, useTwins } from '../features/fleet/queries';
import { useProjectContext } from '../features/projects/ProjectContext';
import { filterLoadedTwins } from '../features/fleet/filtering';
import { useMemo, useState } from 'react';

const pageSize = 50;

export default function Devices() {
  const [params, setParams] = useSearchParams();
  const { projectId, setProjectId } = useProjectContext();
  const filters = useMemo<FleetFilters>(() => ({
    projectId: params.get('project') || projectId,
    deviceType: params.get('type') || '',
    lifecycle: (params.get('lifecycle') || '') as DeviceLifecycle | '',
    connectivity: (params.get('connectivity') || '') as ConnectivityStatus | '',
    search: params.get('q') || '',
  }), [params, projectId]);
  const serverFilters = { ...filters, search: undefined };
  const filterKey = JSON.stringify([filters.projectId, filters.deviceType, filters.lifecycle, filters.connectivity]);
  const [cursorPages, setCursorPages] = useState<Record<string, string[]>>({});
  const cursors = cursorPages[filterKey] || [''];
  const setCursors = (update: (current: string[]) => string[]) => setCursorPages(current => ({ ...current, [filterKey]: update(current[filterKey] || ['']) }));
  const cursor = cursors[cursors.length - 1];
  const twins = useTwins(serverFilters, pageSize, cursor);
  const projects = useProjects();
  const updateFilters = (next: FleetFilters) => {
    const values = new URLSearchParams();
    if (next.projectId) values.set('project', next.projectId);
    if (next.deviceType) values.set('type', next.deviceType);
    if (next.lifecycle) values.set('lifecycle', next.lifecycle);
    if (next.connectivity) values.set('connectivity', next.connectivity);
    if (next.search) values.set('q', next.search);
    setProjectId(next.projectId || '');
    setCursorPages(current => ({ ...current, [JSON.stringify([next.projectId, next.deviceType, next.lifecycle, next.connectivity])]: [''] }));
    setParams(values);
  };
  const items = filterLoadedTwins(twins.data || [], filters.search);
  const projectNames = new Map((projects.data || []).map(project => [project.project_id, project.name]));

  return <div className="page-stack">
    <PageHeader eyebrow="Fleet registry" title="Devices" description="Lifecycle, connectivity, ownership, and latest reported state remain deliberately distinct." actions={<span className="loaded-scope">Cursor page {cursors.length} · {twins.data?.length || 0} loaded</span>} />
    <FleetFiltersBar filters={filters} projects={projects.data || []} onChange={updateFilters} />
    {twins.isLoading ? <TableSkeleton /> : twins.error ? <ErrorState error={twins.error} retry={() => void twins.refetch()} /> : !items.length ? <EmptyState title="No devices match this view" message="Adjust the filters or select another project. Search only covers the loaded cursor page." /> : <div className="data-table-wrap"><table className="data-table"><thead><tr><th>Device</th><th>Project</th><th>Registry lifecycle</th><th>Connectivity</th><th>Battery</th><th>Last seen</th><th>Current position</th><th><span className="sr-only">Open</span></th></tr></thead><tbody>{items.map(twin => <tr key={twin.device_id}>
      <td><Link className="device-cell" to={`/devices/${encodeURIComponent(twin.device_id)}`}><DeviceTypeIcon type={twin.device.device_type_id} /><span><strong>{twin.device.display_name}</strong><small>{deviceTypeLabel(twin.device.device_type_id)} · {twin.device_id}</small></span></Link></td>
      <td>{twin.device.project_id ? projectNames.get(twin.device.project_id) || <code>{twin.device.project_id.slice(0, 8)}…</code> : <span className="muted">Unassigned</span>}</td>
      <td><LifecycleBadge status={twin.device.lifecycle_status} /></td><td><ConnectivityBadge status={twin.connectivity.status} /></td>
      <td><BatteryIndicator value={twin.reported_state?.energy_percent} /></td><td title={exactTime(twin.connectivity.last_seen_at)}>{relativeTime(twin.connectivity.last_seen_at)}</td>
      <td className="mono-coordinate">{formatCoordinates(twin.reported_state?.lat, twin.reported_state?.lon)}</td><td><Link className="row-link" to={`/devices/${encodeURIComponent(twin.device_id)}`} aria-label={`Open ${twin.device.display_name}`}><ExternalLink /></Link></td>
    </tr>)}</tbody></table></div>}
    <div className="pagination"><button className="button secondary" disabled={cursors.length === 1} onClick={() => setCursors(value => value.slice(0, -1))}><ChevronLeft /> Previous</button><span>Cursor page {cursors.length}<small>No total count is exposed by the backend</small></span><button className="button secondary" disabled={(twins.data?.length || 0) < pageSize} onClick={() => { const last = twins.data?.at(-1)?.device_id; if (last) setCursors(value => [...value, last]); }}>Next <ChevronRight /></button></div>
  </div>;
}
