import { Battery, Braces, ChevronDown, Compass, Gauge, MapPin } from 'lucide-react';
import { exactTime, formatCoordinates, formatSpeed, relativeTime } from '../../lib/format';
import type { TwinComponent } from '../../types/domain';
import { BatteryIndicator } from '../status/Status';

function number(payload: Record<string, unknown>, key: string) {
  const value = payload[key];
  return typeof value === 'number' && Number.isFinite(value) ? value : undefined;
}

export function TwinComponentCard({ name, component }: { name: string; component: TwinComponent }) {
  const spatial = name === 'spatial/v1' || component.type === 'spatial/v1';
  const battery = name === 'battery/v1' || component.type === 'battery/v1';
  const lat = number(component.payload, 'latitude');
  const lon = number(component.payload, 'longitude');
  const speed = number(component.payload, 'speed_mps');
  const heading = number(component.payload, 'heading_degrees');
  const percent = number(component.payload, 'percent');

  return <article className="component-card">
    <div className="component-head"><span className="component-icon">{spatial ? <MapPin /> : battery ? <Battery /> : <Braces />}</span><div><strong>{spatial ? 'Spatial state' : battery ? 'Battery state' : name}</strong><span>{component.type} · schema v{component.schema_version}</span></div><time title={exactTime(component.observed_at)}>{relativeTime(component.observed_at)}</time></div>
    {spatial && <div className="component-values">
      <div><MapPin /><span>Coordinates</span><strong>{formatCoordinates(lat, lon)}</strong></div>
      <div><Gauge /><span>Speed</span><strong>{formatSpeed(speed)}</strong></div>
      <div><Compass /><span>Heading</span><strong>{heading === undefined ? 'Not reported' : `${heading.toFixed(0)}°`}</strong></div>
      <div><span className="profile-symbol">M</span><span>Profile</span><strong>{String(component.payload.mobility_profile || 'Unknown').replaceAll('_', ' ')}</strong></div>
    </div>}
    {battery && <div className="battery-component"><BatteryIndicator value={percent} /><div className="battery-track"><span style={{ width: `${Math.max(0, Math.min(100, percent || 0))}%` }} /></div></div>}
    {!spatial && !battery && <pre className="structured-data">{JSON.stringify(component.payload, null, 2)}</pre>}
    <details className="technical-details"><summary><ChevronDown /> Component metadata</summary><dl><div><dt>Boot ID</dt><dd>{component.boot_id}</dd></div><div><dt>Sequence</dt><dd>{component.sequence_number}</dd></div><div><dt>Observed</dt><dd>{exactTime(component.observed_at)}</dd></div></dl>{(spatial || battery) && <pre className="structured-data">{JSON.stringify(component.payload, null, 2)}</pre>}</details>
  </article>;
}

