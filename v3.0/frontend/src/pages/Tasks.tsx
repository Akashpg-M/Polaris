import { ChevronLeft, ChevronRight, Plus, Search } from 'lucide-react';
import { Link, useSearchParams } from 'react-router-dom';
import { TaskPriorityBadge, TaskStatusBadge } from '../components/operations/OperationsStatus';
import { EmptyState, ErrorState, TableSkeleton } from '../components/ui/States';
import { PageHeader } from '../components/ui/PageHeader';
import { useAuth } from '../features/auth/AuthContext';
import { taskStatuses } from '../features/operations/model';
import { useTasks } from '../features/operations/queries';
import { useProjects } from '../features/fleet/queries';
import { exactTime, relativeTime } from '../lib/format';
import { can } from '../lib/permissions';
import type { TaskFilters, TaskStatus } from '../types/domain';
import { useState } from 'react';

const pageSize = 50;

export default function Tasks() {
  const { session } = useAuth();
  const [params, setParams] = useSearchParams();
  const filters: TaskFilters = { status: (params.get('status') || '') as TaskStatus | '', deviceId: params.get('device') || '' };
  const key = JSON.stringify(filters);
  const [pages, setPages] = useState<Record<string, string[]>>({});
  const cursors = pages[key] || [''];
  const cursor = cursors.at(-1) || '';
  const tasks = useTasks(filters, pageSize, cursor);
  const projects = useProjects();
  const projectNames = new Map((projects.data || []).map(project => [project.project_id, project.name]));
  const update = (next: TaskFilters) => { const value = new URLSearchParams(); if (next.status) value.set('status', next.status); if (next.deviceId) value.set('device', next.deviceId); setParams(value); };
  const move = (change: (current: string[]) => string[]) => setPages(current => ({ ...current, [key]: change(current[key] || ['']) }));
  return <div className="page-stack"><PageHeader eyebrow="Durable orchestration" title="Tasks" description="Operator intent, eligibility, exclusive assignment, and execution state remain visible as distinct stages." actions={session && can(session.role, 'createTask') ? <Link className="button primary" to="/tasks/new"><Plus /> Create task</Link> : undefined} />
    <div className="operation-filters"><select value={filters.status || ''} onChange={event => update({ ...filters, status: event.target.value as TaskStatus | '' })} aria-label="Filter tasks by state"><option value="">All task states</option>{taskStatuses.map(status => <option key={status}>{status}</option>)}</select><label><Search /><input value={filters.deviceId || ''} onChange={event => update({ ...filters, deviceId: event.target.value })} placeholder="Exact assigned device ID" aria-label="Filter tasks by assigned device" /></label>{(filters.status || filters.deviceId) && <button onClick={() => setParams({})}>Clear</button>}<span>Refreshes every 10 seconds while open</span></div>
    {tasks.isLoading ? <TableSkeleton /> : tasks.error ? <ErrorState error={tasks.error} retry={() => void tasks.refetch()} /> : !tasks.data?.length ? <EmptyState title="No tasks found" message={session && can(session.role, 'createTask') ? 'Create work for an eligible Polaris device.' : 'No durable work matches the current filters.'} /> : <div className="data-table-wrap"><table className="data-table operations-table"><thead><tr><th>Task</th><th>Priority</th><th>State</th><th>Project</th><th>Assigned device</th><th>Created</th><th>Expiry</th></tr></thead><tbody>{tasks.data.map(task => <tr key={task.task_id}><td><Link className="operation-id" to={`/tasks/${task.task_id}?return=${encodeURIComponent(params.toString())}`}><strong>{task.task_type}</strong><code>{task.task_id}</code></Link></td><td><TaskPriorityBadge priority={task.priority} /></td><td><TaskStatusBadge status={task.status} />{task.failure_reason && <small className="row-error">{task.failure_reason}</small>}</td><td>{task.project_id ? projectNames.get(task.project_id) || <code>{task.project_id}</code> : <span className="muted">Any project</span>}</td><td>{task.assigned_device_id ? <Link to={`/devices/${task.assigned_device_id}`}><code>{task.assigned_device_id}</code></Link> : <span className="waiting-label">Waiting for eligibility</span>}</td><td title={exactTime(task.created_at)}>{relativeTime(task.created_at)}</td><td title={exactTime(task.expires_at)}>{relativeTime(task.expires_at)}</td></tr>)}</tbody></table></div>}
    <div className="pagination"><button className="button secondary" disabled={cursors.length === 1} onClick={() => move(value => value.slice(0,-1))}><ChevronLeft /> Previous</button><span>Cursor page {cursors.length}<small>No total or next-cursor metadata is exposed</small></span><button className="button secondary" disabled={(tasks.data?.length || 0) < pageSize} onClick={() => { const last = tasks.data?.at(-1)?.task_id; if (last) move(value => [...value,last]); }}>Next <ChevronRight /></button></div>
  </div>;
}
