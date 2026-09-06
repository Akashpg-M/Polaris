import { ArrowRight, Boxes, FolderKanban, Radio, ShieldCheck, WifiOff } from 'lucide-react';
import { Link } from 'react-router-dom';
import { DeviceTypeIcon, WebSocketStatus } from '../components/status/Status';
import { EmptyState, ErrorState, PageLoader } from '../components/ui/States';
import { PageHeader } from '../components/ui/PageHeader';
import { FleetMapCanvas } from '../components/map/FleetMapCanvas';
import { deviceTypeLabel, relativeTime } from '../lib/format';
import { useDashboard } from '../features/dashboard/DashboardContext';
import { useProjects, useTwins } from '../features/fleet/queries';
import { useProjectContext } from '../features/projects/ProjectContext';

function Distribution({ title, values }: { title: string; values: Array<{ label: string; value: number; tone: string }> }) {
  const total = Math.max(1, values.reduce((sum, item) => sum + item.value, 0));
  return <section className="panel distribution"><div className="panel-title"><h2>{title}</h2><span>{total === 1 && values.every(item => item.value === 0) ? 0 : total} loaded</span></div><div className="distribution-list">{values.map(item => <div key={item.label}><div><span><i className={item.tone} />{item.label}</span><strong>{item.value}</strong></div><div className="distribution-track"><span className={item.tone} style={{ width: `${item.value / total * 100}%` }} /></div></div>)}</div></section>;
}

export default function Overview() {
  const { projectId } = useProjectContext();
  const twins = useTwins({ projectId }, 100);
  const projects = useProjects();
  const dashboard = useDashboard();
  if (twins.isLoading || projects.isLoading) return <PageLoader label="Loading fleet overview" />;
  if (twins.error) return <ErrorState error={twins.error} retry={() => void twins.refetch()} />;
  const fleet = twins.data || [];
  const connectivity = ['ONLINE', 'STALE', 'OFFLINE', 'NEVER_CONNECTED'].map((status, index) => ({ label: status.replace('_', ' '), value: fleet.filter(item => item.connectivity.status === status).length, tone: ['emerald','amber','red','slate'][index] }));
  const typeCounts = Object.entries(fleet.reduce<Record<string, number>>((result, item) => ({ ...result, [item.device.device_type_id]: (result[item.device.device_type_id] || 0) + 1 }), {})).sort((a, b) => b[1] - a[1]);
  const active = fleet.filter(item => item.device.lifecycle_status === 'ACTIVE').length;
  const online = fleet.filter(item => item.connectivity.status === 'ONLINE').length;
  const stale = fleet.filter(item => item.connectivity.status === 'STALE').length;
  const offline = fleet.filter(item => item.connectivity.status === 'OFFLINE').length;

  return <div className="page-stack">
    <PageHeader eyebrow="Fleet command view" title="Overview" description="The current registered and reported state for your tenant workspace." actions={<span className="loaded-scope">Bounded view · up to 100 devices</span>} />
    {dashboard.status !== 'CONNECTED' && <div className="stream-warning"><WifiOff /><div><strong>Live updates {dashboard.status === 'RECONNECTING' ? 'are reconnecting' : 'are disconnected'}.</strong><span>Authoritative API data remains available, but current state may be stale.</span></div><WebSocketStatus status={dashboard.status} /></div>}
    <section className="metric-grid">
      <article><span className="metric-icon blue"><Boxes /></span><div><small>Registered in loaded view</small><strong>{fleet.length}</strong><p>{active} lifecycle active</p></div></article>
      <article><span className="metric-icon emerald"><Radio /></span><div><small>Online</small><strong>{online}</strong><p>Fresh reported state</p></div></article>
      <article><span className="metric-icon amber"><Radio /></span><div><small>Stale</small><strong>{stale}</strong><p>Past freshness threshold</p></div></article>
      <article><span className="metric-icon red"><Radio /></span><div><small>Offline</small><strong>{offline}</strong><p>Past offline threshold</p></div></article>
      <article><span className="metric-icon violet"><FolderKanban /></span><div><small>Projects</small><strong>{projects.data?.length || 0}</strong><p>Tenant workspaces</p></div></article>
    </section>
    {!fleet.length ? <EmptyState title="No registered devices" message="No devices are available in the selected tenant and project context." /> : <>
      <div className="overview-split"><Distribution title="Connectivity" values={connectivity} /><Distribution title="Device types" values={typeCounts.map(([type, value], index) => ({ label: deviceTypeLabel(type), value, tone: ['blue','violet','cyan','emerald','amber','slate'][index % 6] }))} /></div>
      <section className="panel map-preview"><div className="panel-title"><div><p className="eyebrow">Current reported positions</p><h2>Live fleet map</h2></div><Link to="/fleet/map">Open map <ArrowRight /></Link></div><FleetMapCanvas twins={fleet} compact /></section>
      <div className="overview-split lower">
        <section className="panel fleet-snapshot"><div className="panel-title"><h2>Fleet snapshot</h2><span>Loaded registry categories</span></div>{typeCounts.map(([type, count]) => <div className="snapshot-row" key={type}><span><DeviceTypeIcon type={type} />{deviceTypeLabel(type)}</span><strong>{count}</strong></div>)}</section>
        <section className="panel activity-panel"><div className="panel-title"><h2>Recent live activity</h2><span>This browser session</span></div>{dashboard.activities.length ? dashboard.activities.slice(0, 6).map(item => <div className="activity-row" key={item.eventId}><span className="activity-mark"><ShieldCheck /></span><div><strong>{item.deviceId}</strong><p>{item.message}</p></div><time>{relativeTime(item.observedAt)}</time></div>) : <div className="compact-empty"><Radio /><span>Waiting for accepted tenant telemetry. This feed is live and does not replay history.</span></div>}</section>
      </div>
    </>}
  </div>;
}

