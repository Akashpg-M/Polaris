import { ArrowRight, MapPinned, Route, TrafficCone } from 'lucide-react';
import { Link } from 'react-router-dom';
import { MobilityNav } from '../components/mobility/MobilityNav';
import { ModuleStatusView } from '../components/mobility/MobilityUi';
import { FleetMapCanvas } from '../components/map/FleetMapCanvas';
import { PageHeader } from '../components/ui/PageHeader';
import { ErrorState } from '../components/ui/States';
import { useTwins } from '../features/fleet/queries';
import { useMobilityReadiness } from '../features/mobility/queries';

export default function MobilityOverview() {
  const readiness = useMobilityReadiness(); const twins = useTwins({}, 100); const module = readiness.data?.modules?.mobility; const details = module?.details;
  const spatialTwins = (twins.data || []).filter(item => item.reported_state && Number.isFinite(item.reported_state.lat) && Number.isFinite(item.reported_state.lon));
  return <div className="page-stack mobility-page"><PageHeader eyebrow="Spatial operations" title="Mobility" description="Geographic discovery and bounded traffic-aware road routing, kept separate from core assignment authority." /><MobilityNav />
    {readiness.error ? <ErrorState error={readiness.error} retry={() => void readiness.refetch()} /> : <ModuleStatusView module={module} />}
    <section className="metric-grid mobility-metrics"><article><span className="metric-icon blue"><MapPinned /></span><div><small>Spatial twins in loaded view</small><strong>{spatialTwins.length}</strong><p>Bounded to 100 registry twins; not a tenant total</p></div></article><article><span className="metric-icon cyan"><Route /></span><div><small>Road graph</small><strong>{String(details?.road_graph_version || 'Unavailable')}</strong><p>{Number(details?.road_nodes || 0).toLocaleString()} nodes · {Number(details?.road_edges || 0).toLocaleString()} edges</p></div></article><article><span className="metric-icon violet"><TrafficCone /></span><div><small>Traffic snapshot</small><strong>{String(details?.routing_snapshot_version ?? 'Unavailable')}</strong><p>{String(details?.traffic_scope || 'Scope not returned')}</p></div></article></section>
    <div className="mobility-action-grid"><Link to="/mobility/nearby"><MapPinned /><span><strong>Nearby Search</strong><small>Find indexed devices around a point using the canonical tenant-aware Mobility index.</small></span><ArrowRight /></Link><Link to="/mobility/routes"><Route /><span><strong>Route Explorer</strong><small>Calculate shortest-distance or fastest-time road routes.</small></span><ArrowRight /></Link><Link to="/mobility/traffic"><TrafficCone /><span><strong>Traffic State</strong><small>Inspect the exposed shared routing-cost snapshot metadata.</small></span><ArrowRight /></Link></div>
    <section className="panel mobility-map-panel"><div className="panel-title"><div><h2>Current reported fleet positions</h2><span>Shared Phase 1 map · bounded API hydration plus telemetry updates</span></div></div><FleetMapCanvas twins={spatialTwins} compact /></section>
    <div className="authority-rule"><strong>Mobility proposes and plans.</strong><span>Core revalidates eligibility. PostgreSQL commits assignment authority.</span></div>
  </div>;
}
