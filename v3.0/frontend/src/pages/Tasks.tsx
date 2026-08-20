import { useCallback, useEffect, useState } from 'react';

type Task = {
  task_id: string; task_type: string; status: string; priority: string;
  assigned_device_id?: string; created_at: string; failure_reason?: string;
};
type Command = {
  command_id: string; task_id: string; status: string; attempt_count: number;
  sequence_number: number; acknowledged_at?: string; completed_at?: string;
  payload?: { route_id?: string; road_graph_version?: string; routing_snapshot_version?: number; distance_meters?: number; estimated_duration_ms?: number };
};

const api = import.meta.env.VITE_ENGINE_API || 'http://localhost:6081/api/v1';

export default function Tasks() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [commands, setCommands] = useState<Command[]>([]);
  const [error, setError] = useState('');
  const [target, setTarget] = useState({ lat: '13.0067', lon: '80.2206' });

  const headers = useCallback(() => {
    const token = localStorage.getItem('polaris_operator_token') || '';
    return { Authorization: `Bearer ${token}`, 'X-Tenant-ID': 'alpha_logistics', 'Content-Type': 'application/json' };
  }, []);

  const refresh = useCallback(async () => {
    if (!localStorage.getItem('polaris_operator_token')) { setError('Set polaris_operator_token in localStorage first.'); return; }
    try {
      const [taskResponse, commandResponse] = await Promise.all([
        fetch(`${api}/tasks?limit=100`, { headers: headers() }),
        fetch(`${api}/commands?limit=100`, { headers: headers() })
      ]);
      if (!taskResponse.ok || !commandResponse.ok) throw new Error('Orchestration API rejected the request');
      setTasks((await taskResponse.json()).data || []);
      setCommands((await commandResponse.json()).data || []);
      setError('');
    } catch (e) { setError(e instanceof Error ? e.message : 'Unable to load orchestration state'); }
  }, [headers]);

  useEffect(() => { refresh(); const timer = setInterval(refresh, 2000); return () => clearInterval(timer); }, [refresh]);

  const createTask = async () => {
    const response = await fetch(`${api}/tasks`, {
      method: 'POST', headers: headers(), body: JSON.stringify({
        task_type: 'RELOCATE', priority: 'HIGH',
        requirements: { required_capabilities: ['receive_relocation_command'], minimum_battery: 30, max_distance_meters: 10000 },
        target: { lat: Number(target.lat), lon: Number(target.lon) },
        expires_at: new Date(Date.now() + 5 * 60_000).toISOString()
      })
    });
    if (!response.ok) { const body = await response.json(); setError(body.error?.message || 'Task creation failed'); return; }
    await refresh();
  };

  const cancel = async (taskId: string) => {
    const response = await fetch(`${api}/tasks/${taskId}/cancel`, { method: 'POST', headers: headers(), body: '{}' });
    if (!response.ok) setError('Only pending, not-yet-delivered work can be cancelled locally.');
    await refresh();
  };

  return <div className="h-full overflow-y-auto bg-slate-900 p-8 text-slate-100">
    <div className="mb-6 flex items-end justify-between">
      <div><h2 className="text-2xl font-bold">Durable Task Orchestration</h2><p className="mt-1 text-xs text-slate-400">Capability selection, delivery, ACK, result, retry and expiry state</p></div>
      <div className="flex gap-2">
        <input value={target.lat} onChange={e => setTarget(v => ({ ...v, lat:e.target.value }))} className="w-28 rounded bg-slate-800 p-2 text-xs" aria-label="Latitude" />
        <input value={target.lon} onChange={e => setTarget(v => ({ ...v, lon:e.target.value }))} className="w-28 rounded bg-slate-800 p-2 text-xs" aria-label="Longitude" />
        <button onClick={createTask} className="rounded bg-blue-600 px-4 py-2 text-xs font-bold hover:bg-blue-500">Create relocate task</button>
      </div>
    </div>
    {error && <div className="mb-4 rounded border border-red-800 bg-red-950/50 p-3 text-xs text-red-300">{error}</div>}
    <div className="overflow-hidden rounded-xl border border-slate-700 bg-slate-800/60">
      <table className="w-full text-left text-xs"><thead className="bg-slate-950 text-slate-400"><tr><th className="p-3">Task</th><th>State</th><th>Device</th><th>Command timeline</th><th>Created</th><th /></tr></thead>
        <tbody>{tasks.map(task => { const related = commands.filter(c => c.task_id === task.task_id); return <tr key={task.task_id} className="border-t border-slate-700 align-top">
          <td className="p-3"><div className="font-bold">{task.task_type}</div><div className="font-mono text-[10px] text-slate-500">{task.task_id}</div></td>
          <td><span className="rounded bg-slate-950 px-2 py-1">{task.status}</span>{task.failure_reason && <div className="mt-2 text-red-300">{task.failure_reason}</div>}</td>
          <td className="font-mono">{task.assigned_device_id || 'Awaiting eligible device'}</td>
          <td>{related.length ? related.map(c => <div key={c.command_id} className="mb-1"><span className="text-cyan-300">#{c.sequence_number}</span> {c.status} · attempt {c.attempt_count}{c.payload?.route_id && <div className="text-[10px] text-emerald-300">route {c.payload.route_id} · graph {c.payload.road_graph_version} · snapshot v{c.payload.routing_snapshot_version} · {(c.payload.distance_meters || 0).toFixed(0)} m</div>}</div>) : 'No command created'}</td>
          <td>{new Date(task.created_at).toLocaleTimeString()}</td>
          <td><button disabled={!['PENDING','ASSIGNED'].includes(task.status)} onClick={() => cancel(task.task_id)} className="text-red-300 disabled:text-slate-600">Cancel</button></td>
        </tr>})}</tbody>
      </table>
    </div>
  </div>;
}
