import { ArrowLeft, Battery, Box, Braces, Clock3, Cpu, Gauge, MapPin, Radio, Shield, Tag } from 'lucide-react';
import { Link, useParams, useSearchParams } from 'react-router-dom';
import { FleetMapCanvas } from '../components/map/FleetMapCanvas';
import { BatteryIndicator, ConnectivityBadge, DeviceTypeIcon, LifecycleBadge } from '../components/status/Status';
import { TwinComponentCard } from '../components/twin/TwinComponentCard';
import { EmptyState, ErrorState, PageLoader } from '../components/ui/States';
import { deviceTypeLabel, exactTime, formatCoordinates, formatSpeed, relativeTime } from '../lib/format';
import { useProjects, useTwin } from '../features/fleet/queries';

const tabs = ['overview', 'twin', 'telemetry', 'capabilities'] as const;
type Tab = typeof tabs[number];

function Metadata({ label, value, mono = false }: { label: string; value?: string | null; mono?: boolean }) {
  return <div className="metadata-row"><dt>{label}</dt><dd className={mono ? 'mono-coordinate' : ''}>{value || 'Not provided'}</dd></div>;
}

export default function DeviceDetail() {
  const { deviceId = '' } = useParams();
  const [params, setParams] = useSearchParams();
  const requested = params.get('tab') as Tab;
  const tab: Tab = tabs.includes(requested) ? requested : 'overview';
  const twin = useTwin(deviceId);
  const projects = useProjects();
  if (twin.isLoading) return <PageLoader label="Loading device twin" />;
  if (twin.error || !twin.data) return <ErrorState error={twin.error || new Error('Device was not found.')} retry={() => void twin.refetch()} />;
  const value = twin.data;
  const device = value.device;
  const state = value.reported_state;
  const project = projects.data?.find(item => item.project_id === device.project_id);
  const hasSpatial = Boolean(state && Number.isFinite(state.lat) && Number.isFinite(state.lon));
  const selectTab = (next: Tab) => setParams(next === 'overview' ? {} : { tab: next });

  return <div className="page-stack device-detail">
    <Link to="/devices" className="back-link"><ArrowLeft /> Back to devices</Link>
    <header className="device-hero"><div className="device-identity"><span className="hero-device-icon"><DeviceTypeIcon type={device.device_type_id} size={30} /></span><div><p>{deviceTypeLabel(device.device_type_id)}</p><h1>{device.display_name}</h1><code>{device.device_id}</code></div></div><div className="device-hero-status"><div className="badge-row"><LifecycleBadge status={device.lifecycle_status} /><ConnectivityBadge status={value.connectivity.status} /></div><span title={exactTime(value.connectivity.last_seen_at)}>{value.connectivity.last_seen_at ? `Last seen ${relativeTime(value.connectivity.last_seen_at)}` : 'No accepted telemetry yet'}</span></div></header>
    <nav className="detail-tabs" aria-label="Device detail sections">{tabs.map(item => <button key={item} className={tab === item ? 'active' : ''} onClick={() => selectTab(item)}>{item.charAt(0).toUpperCase() + item.slice(1)}</button>)}<span className="disabled-tab" title="Available in a later frontend phase">Tasks · Later</span></nav>

    {tab === 'overview' && <div className="detail-layout">
      <section className="detail-main page-stack">
        <section className="panel"><div className="panel-title"><div><p className="eyebrow">Latest accepted observation</p><h2>Current reported state</h2></div>{state && <span className="sequence-chip">Sequence {state.sequence_number}</span>}</div>
          {state ? <div className="current-state-grid"><div><Battery /><span>Battery</span><strong><BatteryIndicator value={state.energy_percent} /></strong></div><div><MapPin /><span>Position</span><strong>{formatCoordinates(state.lat, state.lon)}</strong></div><div><Gauge /><span>Speed</span><strong>{formatSpeed(state.velocity_mps)}</strong></div><div><Radio /><span>Heading</span><strong>{Number.isFinite(state.heading_deg) ? `${state.heading_deg.toFixed(0)}°` : 'Not reported'}</strong></div><div><Clock3 /><span>Observed</span><strong title={exactTime(state.observed_at)}>{relativeTime(state.observed_at)}</strong></div><div><Braces /><span>Schema</span><strong>Telemetry v{state.schema_version}</strong></div></div> : <div className="never-connected"><Radio /><div><strong>Never connected</strong><p>This device is registered in Polaris but has not yet published accepted telemetry.</p></div></div>}
        </section>
        <section className="panel"><div className="panel-title"><h2>Device metadata</h2><span>Durable registry</span></div><dl className="metadata-list"><Metadata label="Device ID" value={device.device_id} mono /><Metadata label="Tenant" value={device.tenant_id} mono /><Metadata label="Type" value={deviceTypeLabel(device.device_type_id)} /><Metadata label="Project" value={project?.name || (device.project_id ? device.project_id : 'Unassigned')} /><Metadata label="Firmware" value={device.firmware_version} /><Metadata label="Software" value={device.software_version} /><Metadata label="Model" value={device.model_version} /><Metadata label="Registered" value={exactTime(device.registered_at)} /><Metadata label="Last registry update" value={exactTime(device.updated_at)} /></dl>{Object.keys(device.metadata || {}).length > 0 && <details className="technical-details"><summary>Registry metadata</summary><pre className="structured-data">{JSON.stringify(device.metadata, null, 2)}</pre></details>}</section>
      </section>
      <aside className="detail-aside page-stack"><section className="panel map-card"><div className="panel-title"><h2>Location</h2></div>{hasSpatial ? <FleetMapCanvas twins={[value]} compact /> : <div className="compact-empty"><MapPin /><span>This device does not expose a current spatial component.</span></div>}</section><section className="panel"><div className="panel-title"><h2>Capabilities</h2><Link to={`?tab=capabilities`}>View all</Link></div><div className="capability-list vertical">{value.capabilities.length ? value.capabilities.slice(0, 6).map(capability => <span key={capability.capability_id}><Shield />{capability.display_name}</span>) : <small>No capabilities assigned.</small>}</div></section></aside>
    </div>}

    {tab === 'twin' && <div className="page-stack"><div className="twin-authority"><div><Box /><span><strong>Registry</strong><small>Durable identity and configuration</small></span></div><i>+</i><div><Radio /><span><strong>Reported state</strong><small>Latest operational projection</small></span></div></div>{Object.keys(value.components).length ? <div className="component-grid">{Object.entries(value.components).map(([name, component]) => <TwinComponentCard name={name} component={component} key={name} />)}</div> : <EmptyState title="No reported components" message="The durable registry exists, but this device has not published an accepted component state." />}</div>}

    {tab === 'telemetry' && <section className="panel telemetry-latest"><div className="panel-title"><div><p className="eyebrow">Current state only</p><h2>Latest telemetry</h2></div><span className="phase-label">History API required</span></div>{state ? <div className="data-table-wrap"><table className="data-table"><thead><tr><th>Observed</th><th>Coordinates</th><th>Speed</th><th>Heading</th><th>Battery</th><th>Sequence</th><th>Boot</th></tr></thead><tbody><tr><td>{exactTime(state.observed_at)}</td><td className="mono-coordinate">{formatCoordinates(state.lat, state.lon)}</td><td>{formatSpeed(state.velocity_mps)}</td><td>{state.heading_deg.toFixed(0)}°</td><td><BatteryIndicator value={state.energy_percent} /></td><td>{state.sequence_number}</td><td><code>{state.device_boot_id}</code></td></tr></tbody></table></div> : <EmptyState title="No telemetry received" message="This device has never published accepted telemetry." />}<p className="architecture-note"><Clock3 /> Polaris archives telemetry durably, but no tenant-scoped history read API is currently exposed. This view does not fabricate history from transient WebSocket events.</p></section>}

    {tab === 'capabilities' && <section className="panel"><div className="panel-title"><div><p className="eyebrow">Read-only in Phase 1</p><h2>Assigned capabilities</h2></div><span>{value.capabilities.length} assigned</span></div>{value.capabilities.length ? <div className="capability-grid">{value.capabilities.map(capability => <article key={capability.capability_id}><span className="capability-icon"><Cpu /></span><div><h3>{capability.display_name}</h3><code>{capability.capability_id}</code><p>{capability.description || 'No catalog description.'}</p></div><LifecycleBadge status={capability.enabled ? 'ACTIVE' : 'SUSPENDED'} /><details className="technical-details"><summary><Tag /> Device configuration</summary><pre className="structured-data">{JSON.stringify(capability.configuration || {}, null, 2)}</pre></details></article>)}</div> : <EmptyState title="No capabilities assigned" message="A tenant administrator can assign capabilities through the registry API in a later frontend phase." />}</section>}
  </div>;
}

