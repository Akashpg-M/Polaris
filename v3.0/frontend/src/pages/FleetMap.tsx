import { Layers3, LocateFixed, WifiOff } from 'lucide-react';
import { useMemo, useState } from 'react';
import { FleetFiltersBar } from '../components/fleet/FleetFiltersBar';
import { DeviceDrawer } from '../components/fleet/DeviceDrawer';
import { FleetMapCanvas } from '../components/map/FleetMapCanvas';
import { WebSocketStatus } from '../components/status/Status';
import { ErrorState, PageLoader } from '../components/ui/States';
import { useDashboard } from '../features/dashboard/DashboardContext';
import { useProjects, useTwins } from '../features/fleet/queries';
import { useProjectContext } from '../features/projects/ProjectContext';
import { filterLoadedTwins } from '../features/fleet/filtering';
import type { DeviceTwin, FleetFilters } from '../types/domain';

export default function FleetMap() {
  const { projectId, setProjectId } = useProjectContext();
  const [filters, setFilters] = useState<FleetFilters>({});
  const [selected, setSelected] = useState<DeviceTwin | null>(null);
  const projects = useProjects();
  const effectiveFilters = { ...filters, projectId };
  const twins = useTwins({ ...effectiveFilters, search: undefined }, 100);
  const dashboard = useDashboard();
  const spatial = useMemo(() => filterLoadedTwins(twins.data || [], filters.search).filter(twin => twin.reported_state && Number.isFinite(twin.reported_state.lat) && Number.isFinite(twin.reported_state.lon)), [twins.data, filters.search]);
  if (twins.isLoading) return <PageLoader label="Hydrating the live fleet map" />;
  if (twins.error) return <ErrorState error={twins.error} retry={() => void twins.refetch()} />;
  return <div className="map-page">
    <div className="map-toolbar"><div><p className="eyebrow">Current operational geography</p><h1>Live fleet map</h1><p>API-hydrated first, then incrementally updated from the tenant dashboard stream.</p></div><div className="map-toolbar-status"><WebSocketStatus status={dashboard.status} /><span><strong>{spatial.length}</strong> spatial devices<small>Up to 100 hydrated twins</small></span></div></div>
    <div className="map-filter-wrap"><FleetFiltersBar filters={effectiveFilters} projects={projects.data || []} onChange={next => { setProjectId(next.projectId || ''); setFilters({ ...next, projectId: undefined }); setSelected(null); }} /></div>
    {dashboard.status !== 'CONNECTED' && <div className="map-disconnected"><WifiOff /> Live updates are unavailable. Displayed API state may be stale.</div>}
    <div className="map-stage"><FleetMapCanvas twins={spatial} selectedId={selected?.device_id} onSelect={setSelected} />
      <div className="map-legend"><strong><Layers3 /> Marker key</strong><span><i className="shape road">◆</i> Vehicle</span><span><i className="shape drone">▲</i> Drone</span><span><i className="shape robot">■</i> Robot</span><span><i className="shape static">●</i> Static</span><small><LocateFixed /> Nearby markers combine as you zoom out.</small></div>
      {!spatial.length && <div className="map-empty"><LocateFixed /><strong>No spatial devices in this view</strong><span>Never-connected and non-spatial devices remain available in Fleet.</span></div>}
      {selected && <DeviceDrawer twin={selected} projects={projects.data || []} onClose={() => setSelected(null)} />}
    </div>
  </div>;
}
