import { Battery, Braces, Clock3, Cpu, MapPin, Navigation, Ruler, Shapes } from 'lucide-react';
import type { TaskPathTiming, TaskRequirements, TaskTarget } from '../../types/domain';
import { deviceTypeLabel, formatCoordinates } from '../../lib/format';
import { microseconds, requiredCapability } from '../../features/operations/model';

export function TaskRequirementsView({ taskType, requirements }: { taskType: string; requirements: TaskRequirements }) {
  const automatic = requiredCapability(taskType);
  const operatorCaps = (requirements.required_capabilities || []).filter(value => value !== automatic);
  return <div className="requirements-view">
    <section><p className="section-label">Automatically required by Polaris</p>{automatic ? <div className="requirement-chip locked"><Navigation /><span><strong>{automatic}</strong><small>Command capability · locked</small></span></div> : <span className="muted">No automatic capability rule.</span>}</section>
    <section><p className="section-label">Operator requirements</p><div className="requirement-grid">
      {requirements.allowed_device_types?.map(type => <div className="requirement-chip" key={type}><Shapes /><span><strong>{deviceTypeLabel(type)}</strong><small>Allowed device type</small></span></div>)}
      {requirements.minimum_battery ? <div className="requirement-chip"><Battery /><span><strong>{requirements.minimum_battery}% or higher</strong><small>Minimum battery</small></span></div> : null}
      {requirements.max_distance_meters ? <div className="requirement-chip"><Ruler /><span><strong>{(requirements.max_distance_meters / 1000).toLocaleString()} km or nearer</strong><small>Maximum target distance</small></span></div> : null}
      {requirements.project_id ? <div className="requirement-chip"><Cpu /><span><strong>{requirements.project_id}</strong><small>Project constraint</small></span></div> : null}
      {operatorCaps.map(capability => <div className="requirement-chip" key={capability}><Cpu /><span><strong>{capability}</strong><small>Additional capability</small></span></div>)}
      {!requirements.allowed_device_types?.length && !requirements.minimum_battery && !requirements.max_distance_meters && !requirements.project_id && !operatorCaps.length && <span className="muted">No additional operator constraints.</span>}
    </div></section>
    {requirements.custom_constraints !== undefined && <details className="technical-details"><summary><Braces /> Custom constraints</summary><pre className="structured-data">{JSON.stringify(requirements.custom_constraints, null, 2)}</pre></details>}
  </div>;
}

export function TaskTargetViewer({ target }: { target: TaskTarget }) {
  const spatial = Number.isFinite(target.lat) && Number.isFinite(target.lon);
  const copy = () => spatial && void navigator.clipboard?.writeText(`${target.lat}, ${target.lon}`);
  return <div className="target-viewer">{spatial && <div className="coordinate-target"><span><MapPin /></span><div><small>Destination coordinates</small><strong>{formatCoordinates(target.lat, target.lon)}</strong></div><button className="button secondary" onClick={copy}>Copy</button></div>}
    {!spatial && <div className="structured-target">{Object.entries(target).map(([key, value]) => <div key={key}><span>{key}</span><strong>{typeof value === 'object' ? JSON.stringify(value) : String(value)}</strong></div>)}</div>}
    <details className="technical-details"><summary><Braces /> Raw target document</summary><pre className="structured-data">{JSON.stringify(target, null, 2)}</pre></details>
  </div>;
}

export function TaskTimingDiagnostics({ timing }: { timing?: TaskPathTiming }) {
  if (!timing) return <div className="compact-empty"><Clock3 /><span>Request-scoped orchestration timing was not returned by this read operation.</span></div>;
  const values = [['Candidate selection', timing.candidate_selection_duration_us], ['Routing / planning', timing.routing_duration_us], ['Persistence', timing.persistence_duration_us], ['Total', timing.total_duration_us]] as const;
  return <div className="timing-grid">{values.map(([label, value]) => <div key={label}><span>{label}</span><strong>{microseconds(value)}</strong></div>)}</div>;
}
