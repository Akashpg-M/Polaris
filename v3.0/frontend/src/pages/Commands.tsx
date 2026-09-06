import { ChevronLeft, ChevronRight, Search } from 'lucide-react';
import { Link, useSearchParams } from 'react-router-dom';
import { CommandStatusBadge } from '../components/operations/OperationsStatus';
import { EmptyState, ErrorState, TableSkeleton } from '../components/ui/States';
import { PageHeader } from '../components/ui/PageHeader';
import { commandStatuses } from '../features/operations/model';
import { useCommands } from '../features/operations/queries';
import { exactTime, relativeTime } from '../lib/format';
import type { CommandFilters, CommandStatus } from '../types/domain';
import { useState } from 'react';

const pageSize = 50;

export default function Commands() {
  const [params, setParams] = useSearchParams();
  const filters: CommandFilters = { status: (params.get('status') || '') as CommandStatus | '', taskId: params.get('task') || '', deviceId: params.get('device') || '' };
  const key = JSON.stringify(filters);
  const [pages, setPages] = useState<Record<string,string[]>>({});
  const cursors = pages[key] || [''];
  const commands = useCommands(filters, pageSize, cursors.at(-1) || '');
  const update = (next: CommandFilters) => { const value = new URLSearchParams(); if (next.status) value.set('status',next.status); if (next.taskId) value.set('task',next.taskId); if (next.deviceId) value.set('device',next.deviceId); setParams(value); };
  const move = (change: (current:string[])=>string[]) => setPages(current => ({...current,[key]:change(current[key]||[''])}));
  return <div className="page-stack"><PageHeader eyebrow="Device execution" title="Commands" description="Immutable device-specific instructions created only after Polaris assigns executable work." actions={<span className="loaded-scope">At-least-once delivery</span>} />
    <div className="operation-filters"><select value={filters.status || ''} onChange={event => update({...filters,status:event.target.value as CommandStatus | ''})} aria-label="Filter commands by state"><option value="">All command states</option>{commandStatuses.map(status=><option key={status}>{status}</option>)}</select><label><Search /><input value={filters.taskId || ''} onChange={event=>update({...filters,taskId:event.target.value})} placeholder="Exact task ID" aria-label="Filter by task" /></label><label><Search /><input value={filters.deviceId || ''} onChange={event=>update({...filters,deviceId:event.target.value})} placeholder="Exact device ID" aria-label="Filter by device" /></label>{Object.values(filters).some(Boolean)&&<button onClick={()=>setParams({})}>Clear</button>}<span>Refreshes every 10 seconds while open</span></div>
    {commands.isLoading?<TableSkeleton/>:commands.error?<ErrorState error={commands.error} retry={()=>void commands.refetch()}/>:!commands.data?.length?<EmptyState title="No commands found" message="Commands appear after Polaris successfully assigns executable work to a device."/>:<div className="data-table-wrap"><table className="data-table operations-table"><thead><tr><th>Command</th><th>State</th><th>Task</th><th>Device / sequence</th><th>Attempts</th><th>Created</th><th>Expiry</th></tr></thead><tbody>{commands.data.map(command=><tr key={command.command_id}><td><Link className="operation-id" to={`/commands/${command.command_id}?return=${encodeURIComponent(params.toString())}`}><strong>{command.command_type}</strong><code>{command.command_id}</code></Link></td><td><CommandStatusBadge status={command.status}/>{command.last_error&&<small className="row-error">{command.last_error}</small>}</td><td><Link to={`/tasks/${command.task_id}`}><code>{command.task_id}</code></Link></td><td><Link to={`/devices/${command.device_id}`}><code>{command.device_id}</code></Link><small className="sequence-small">Sequence {command.sequence_number}</small></td><td>{command.attempt_count} / {command.max_attempts}</td><td title={exactTime(command.created_at)}>{relativeTime(command.created_at)}</td><td title={exactTime(command.expires_at)}>{relativeTime(command.expires_at)}</td></tr>)}</tbody></table></div>}
    <div className="pagination"><button className="button secondary" disabled={cursors.length===1} onClick={()=>move(value=>value.slice(0,-1))}><ChevronLeft/> Previous</button><span>Cursor page {cursors.length}<small>No total or next-cursor metadata is exposed</small></span><button className="button secondary" disabled={(commands.data?.length||0)<pageSize} onClick={()=>{const last=commands.data?.at(-1)?.command_id;if(last)move(value=>[...value,last]);}}>Next <ChevronRight/></button></div>
  </div>;
}
