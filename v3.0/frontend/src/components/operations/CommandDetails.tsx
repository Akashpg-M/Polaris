import { Braces, Clock3, GitBranch, MapPinned, RefreshCw, Route } from 'lucide-react';
import type { Command } from '../../types/domain';
import { ackExplanation } from '../../features/operations/model';
import { exactTime, formatCoordinates } from '../../lib/format';
import { AckBadge } from './OperationsStatus';
import { Link } from 'react-router-dom';

function field(payload: Record<string, unknown>, name: string) { return payload[name]; }
function position(value: unknown) {
  if (!value || typeof value !== 'object') return 'Not returned';
  const point = value as Record<string, unknown>;
  return formatCoordinates(Number(point.latitude), Number(point.longitude));
}

export function CommandPayloadViewer({ command }: { command: Command }) {
  const [routeId, graph, snapshot] = [field(command.payload, 'route_id'), field(command.payload, 'road_graph_version'), field(command.payload, 'routing_snapshot_version')];
  const isRoute = Boolean(routeId);
  return <div className="payload-viewer">
    <div className="immutable-note"><GitBranch /><span><strong>Immutable Command Payload</strong><small>Delivery retries reuse this command ID, device sequence, and payload.</small></span></div>
    {isRoute && <div className="route-inspection"><div className="route-head"><Route /><div><small>Persisted route plan</small><strong>{String(routeId)}</strong></div></div><div className="route-fields">
      <div><span>Road graph</span><strong>{String(graph || 'Not returned')}</strong></div><div><span>Traffic snapshot</span><strong>{String(snapshot ?? 'Not returned')}</strong></div>
      <div><span>Route schema</span><strong>{String(field(command.payload, 'route_schema_version') ?? 'Not returned')}</strong></div><div><span>Policy</span><strong>{String(field(command.payload, 'policy') || 'Not returned')}</strong></div>
      <div><span>Origin</span><strong>{position(field(command.payload, 'origin'))}</strong></div><div><span>Destination</span><strong>{position(field(command.payload, 'destination'))}</strong></div>
      <div><span>Distance</span><strong>{Number.isFinite(field(command.payload, 'distance_meters')) ? `${(Number(field(command.payload, 'distance_meters')) / 1000).toFixed(2)} km` : 'Not returned'}</strong></div>
      <div><span>Estimated duration</span><strong>{Number.isFinite(field(command.payload, 'estimated_duration_ms')) ? `${(Number(field(command.payload, 'estimated_duration_ms')) / 1000).toFixed(1)} sec` : 'Not returned'}</strong></div>
      <div><span>Generated</span><strong>{exactTime(String(field(command.payload, 'generated_at') || ''))}</strong></div><div><span>Valid until</span><strong>{exactTime(String(field(command.payload, 'valid_until') || ''))}</strong></div>
    </div>{Array.isArray(field(command.payload, 'waypoints')) && <p className="route-waypoints"><MapPinned /> {(field(command.payload, 'waypoints') as unknown[]).length} immutable waypoints persisted</p>}<Link className="button secondary view-mobility-route" to={`/mobility/routes?command_id=${encodeURIComponent(command.command_id)}`}><MapPinned/> View on Mobility Map</Link></div>}
    {!isRoute && <div className="structured-target">{Object.entries(command.payload).slice(0, 12).map(([key, value]) => <div key={key}><span>{key}</span><strong>{typeof value === 'object' ? JSON.stringify(value) : String(value)}</strong></div>)}</div>}
    <details className="technical-details"><summary><Braces /> Raw immutable JSON</summary><pre className="structured-data">{JSON.stringify(command.payload, null, 2)}</pre></details>
  </div>;
}

export function CommandAttemptSummary({ command }: { command: Command }) {
  return <div className="attempt-summary"><div className="attempt-count"><RefreshCw /><span><small>Delivery attempts</small><strong>{command.attempt_count} / {command.max_attempts}</strong></span></div><div className="attempt-facts"><div><span>Latest sent</span><strong>{exactTime(command.sent_at)}</strong></div><div><span>Next available</span><strong>{exactTime(command.available_at)}</strong></div><div><span>Latest delivery error</span><strong>{command.last_error || 'None recorded'}</strong></div></div><p>Per-attempt gateway, ownership epoch, and failure history are persisted internally but are not exposed by the current read API.</p></div>;
}

export function AcknowledgementCard({ command }: { command: Command }) {
  return <div className="ack-card"><div><AckBadge status={command.ack_status} /><span>{exactTime(command.acknowledged_at)}</span></div><p>{ackExplanation(command.ack_status)}</p></div>;
}

export function ExecutionResultCard({ command }: { command: Command }) {
  if (!command.completed_at && command.result === undefined && !command.last_error) return <div className="compact-empty"><Clock3 /><span>No terminal execution result has been recorded.</span></div>;
  const successful = command.status === 'COMPLETED';
  return <div className={`result-card ${successful ? 'success' : 'failure'}`}><strong>{successful ? 'Execution completed' : 'Execution did not complete successfully'}</strong><span>{exactTime(command.completed_at)}</span>{command.last_error && <p>{command.last_error}</p>}{command.result !== undefined && <pre className="structured-data">{JSON.stringify(command.result, null, 2)}</pre>}</div>;
}
