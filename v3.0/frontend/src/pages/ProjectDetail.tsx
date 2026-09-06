import { ArrowLeft, Boxes, Calendar, FolderKanban } from 'lucide-react';
import { Link, useParams } from 'react-router-dom';
import { BatteryIndicator, ConnectivityBadge, DeviceTypeIcon, LifecycleBadge, ProjectStatusBadge } from '../components/status/Status';
import { EmptyState, ErrorState, PageLoader } from '../components/ui/States';
import { deviceTypeLabel, exactTime, relativeTime } from '../lib/format';
import { useProject, useTwins } from '../features/fleet/queries';

export default function ProjectDetail() {
  const { projectId = '' } = useParams();
  const project = useProject(projectId);
  const twins = useTwins({ projectId }, 100);
  if (project.isLoading || twins.isLoading) return <PageLoader label="Loading project fleet" />;
  if (project.error || !project.data) return <ErrorState error={project.error || new Error('Project was not found.')} retry={() => void project.refetch()} />;
  const online = twins.data?.filter(item => item.connectivity.status === 'ONLINE').length || 0;
  return <div className="page-stack"><Link to="/projects" className="back-link"><ArrowLeft /> Back to projects</Link><header className="project-hero"><span className="hero-device-icon"><FolderKanban /></span><div><p>Project workspace</p><h1>{project.data.name}</h1><code>{project.data.project_id}</code></div><ProjectStatusBadge status={project.data.status} /></header><div className="project-summary"><article><Boxes /><span><small>Devices in loaded view</small><strong>{twins.data?.length || 0}</strong></span></article><article><span className="online-dot" /><span><small>Online in loaded view</small><strong>{online}</strong></span></article><article><Calendar /><span><small>Last registry update</small><strong>{exactTime(project.data.updated_at)}</strong></span></article></div>
    <section className="panel"><div className="panel-title"><div><p className="eyebrow">Project fleet</p><h2>Registered devices</h2></div><span>Up to 100</span></div>{!twins.data?.length ? <EmptyState title="No project devices" message="No devices are assigned to this project." /> : <div className="data-table-wrap"><table className="data-table"><thead><tr><th>Device</th><th>Lifecycle</th><th>Connectivity</th><th>Battery</th><th>Last seen</th></tr></thead><tbody>{twins.data.map(twin => <tr key={twin.device_id}><td><Link className="device-cell" to={`/devices/${encodeURIComponent(twin.device_id)}`}><DeviceTypeIcon type={twin.device.device_type_id} /><span><strong>{twin.device.display_name}</strong><small>{deviceTypeLabel(twin.device.device_type_id)} · {twin.device_id}</small></span></Link></td><td><LifecycleBadge status={twin.device.lifecycle_status} /></td><td><ConnectivityBadge status={twin.connectivity.status} /></td><td><BatteryIndicator value={twin.reported_state?.energy_percent} /></td><td>{relativeTime(twin.connectivity.last_seen_at)}</td></tr>)}</tbody></table></div>}</section>
    <section className="panel project-description"><h2>About this project</h2><p>{project.data.description || 'No description has been provided.'}</p>{Object.keys(project.data.metadata || {}).length > 0 && <details className="technical-details"><summary>Project metadata</summary><pre className="structured-data">{JSON.stringify(project.data.metadata, null, 2)}</pre></details>}</section>
  </div>;
}

