import { useState, useRef, useEffect } from 'react';
import { serializeProtobufTelemetry } from '../types/polaris';
import type { LogEntry } from '../types/polaris';
import { engineApi, gatewayWs } from '../config';

interface DiagnosticMetrics {
  activeConnections: number;
  spatialMatcherLatencyMs: number;
  routerExecutionLatencyMs: number;
}

export default function SwarmTester() {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [metrics, setMetrics] = useState<DiagnosticMetrics>({
    activeConnections: 0,
    spatialMatcherLatencyMs: 0,
    routerExecutionLatencyMs: 0
  });

  const activeTimers = useRef<ReturnType<typeof setInterval>[]>([]);
  const openSockets = useRef<WebSocket[]>([]);

  const addLog = (msg: string, type: LogEntry['type'] = 'info') => {
    const time = new Date().toISOString().split('T')[1].slice(0, 12);
    setLogs((prev) => [{ time, msg, type }, ...prev].slice(0, 50)); 
  };

  // Run real-time backend benchmark assertions alongside the streaming payload
  const runBackendBenchmarks = async () => {
    try {
	  const operatorToken = localStorage.getItem('polaris_operator_token');
	  if (!operatorToken) return;
	  const authHeaders = { Authorization: `Bearer ${operatorToken}`, 'X-Tenant-ID': 'alpha_logistics' };
      // 1. Sample the Mobility H3/R-tree query latency indicator.
      const t0 = performance.now();
      await fetch(`${engineApi}/nodes/match?lat=13.0067&lon=80.2206&radius_km=5.0`, { headers: authHeaders });
      const spatialDiff = performance.now() - t0;

      // 2. Sample the Mobility routing latency indicator.
      const t1 = performance.now();
      const routeRes = await fetch(`${engineApi}/routes/calculate?src_lat=13.0067&src_lon=80.2206&tgt_lat=13.0012&tgt_lon=80.2565`, { headers: authHeaders });
      const routeJson = await routeRes.json();
      const routeDiff = performance.now() - t1;

      setMetrics(prev => ({
        ...prev,
        spatialMatcherLatencyMs: parseFloat(spatialDiff.toFixed(1)),
        routerExecutionLatencyMs: routeJson.error ? 0 : parseFloat(routeDiff.toFixed(1))
      }));
    } catch (err) {
      // Silent error boundaries during warm boot sequences
    }
  };

  const bootDrone = async (nodeId: string) => {
	const operatorToken=localStorage.getItem('polaris_operator_token');if(!operatorToken){addLog('Missing polaris_operator_token','danger');return;}
	const headers={Authorization:`Bearer ${operatorToken}`,'X-Tenant-ID':'alpha_logistics','Content-Type':'application/json'};
	const created=await fetch(`${engineApi}/devices`,{method:'POST',headers,body:JSON.stringify({device_id:nodeId,device_type_id:'delivery_drone',display_name:nodeId})});if(!created.ok&&created.status!==409)return;
	await fetch(`${engineApi}/devices/${nodeId}/activate`,{method:'POST',headers});await fetch(`${engineApi}/devices/${nodeId}/credentials`,{method:'POST',headers,body:'{}'});
	const ticketResult=await fetch(`${engineApi}/devices/${nodeId}/connection-ticket`,{method:'POST',headers,body:'{}'});if(!ticketResult.ok)return;const ticket=(await ticketResult.json()).data.ticket;
    const ws = new WebSocket(`${gatewayWs}/ws/telemetry?ticket=${encodeURIComponent(ticket)}`);
    ws.binaryType = "arraybuffer";
    
    let lat = 13.0067 + (Math.random() * 0.02 - 0.01);
    let lon = 80.2206 + (Math.random() * 0.02 - 0.01);
    let velocityMps = 12.0 + Math.random() * 8;
    let headingDeg = Math.floor(Math.random() * 360);
	const bootStartedAt=Date.now();const deviceBootId=`analytics-${nodeId}-${bootStartedAt}`;let sequenceNumber=0;
	let lastCommandSequence=0;const completedCommands=new Map<string,{ack:string;result:string}>();
	ws.onmessage=(event)=>{if(typeof event.data!=='string')return;try{const command=JSON.parse(event.data);if(command.frame_type!=='COMMAND')return;const previous=completedCommands.get(command.command_id);if(previous){ws.send(previous.ack);ws.send(previous.result);return;}let status='ACCEPTED',reason='';if(Date.now()>Date.parse(command.expires_at)){status='EXPIRED';reason='expired';}else if(command.sequence_number<=lastCommandSequence){status='REJECTED';reason='out-of-order sequence';}const ack=JSON.stringify({frame_type:'COMMAND_ACK',command_id:command.command_id,sequence_number:command.sequence_number,status,received_at:new Date().toISOString(),reason});ws.send(ack);if(status!=='ACCEPTED')return;lastCommandSequence=command.sequence_number;if((command.command_type==='RELOCATE'||command.command_type==='NAVIGATE')&&command.payload){if(typeof command.payload.lat==='number')lat=command.payload.lat;if(typeof command.payload.lon==='number')lon=command.payload.lon;}setTimeout(()=>{if(ws.readyState!==WebSocket.OPEN)return;const result=JSON.stringify({frame_type:'COMMAND_RESULT',command_id:command.command_id,sequence_number:command.sequence_number,status:'SUCCEEDED',completed_at:new Date().toISOString(),result:{execution_count:1}});completedCommands.set(command.command_id,{ack,result});ws.send(result);addLog(`${nodeId} completed ${command.command_type}`,'success');},250);}catch{addLog(`${nodeId} received invalid command frame`,'danger');}};

    ws.onopen = () => {
      addLog(`Uplink Channel [${nodeId}] Active`, 'success');
      openSockets.current.push(ws);
      setMetrics(prev => ({ ...prev, activeConnections: openSockets.current.length }));
      
      const timer = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
		  sequenceNumber+=1;
          const radians = (headingDeg * Math.PI) / 180;
          lat += (velocityMps * Math.cos(radians)) / 111000;
          lon += (velocityMps * Math.sin(radians)) / (111000 * Math.cos(lat * Math.PI / 180));

          // Compile raw hardware bytes matching proto wire specs
		  const observedAt=Date.now();const binaryBuffer = serializeProtobufTelemetry({
            id: nodeId,
            tenantId: "alpha_logistics",
            type: 5,   
            status: 3, 
            lat,
            lon,
            velocityMps,
            headingDeg,
            energyPercent: Math.floor(85 + Math.random() * 15),
			timestamp: observedAt, deviceBootId, sequenceNumber, bootStartedAt, observedAt, schemaVersion:1
          });

          ws.send(new Uint8Array(binaryBuffer));
        }
      }, 1000);
      
      activeTimers.current.push(timer);
    };

    ws.onclose = () => {
      // Safely filter connection out of operational slice arrays
      openSockets.current = openSockets.current.filter(s => s !== ws);
      setMetrics(prev => ({ ...prev, activeConnections: openSockets.current.length }));
    };
  };

  const launchSwarm = (count: number) => {
    addLog(`Deploying ${count} performance testing threads onto cluster...`, 'info');
    for (let i = 1; i <= count; i++) {
      setTimeout(() => bootDrone(`DRONE-${Math.floor(1000 + Math.random() * 9000)}`), i * 15);
    }
  };

  const purgeSwarm = () => {
    addLog("Purging all active simulated hardware signals...", 'danger');
    activeTimers.current.forEach(clearInterval);
    openSockets.current.forEach(ws => ws.close());
    activeTimers.current = [];
    openSockets.current = [];
    setMetrics({ activeConnections: 0, spatialMatcherLatencyMs: 0, routerExecutionLatencyMs: 0 });
  };

  useEffect(() => {
    const metricTicker = setInterval(runBackendBenchmarks, 2000);
    return () => {
      clearInterval(metricTicker);
      purgeSwarm();
    };
  }, []);

  return (
    <div className="p-8 h-full flex flex-col gap-6 overflow-y-auto bg-slate-900 text-slate-100">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">🚁 System Stress & Performance Ingestion Console</h2>
        <p className="text-xs text-slate-400 mt-1">Directly compiles and injects pure binary payloads into the cluster mesh</p>
      </div>

      {/* Real-time Diagnostics HUD */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-slate-800 p-5 rounded-xl border border-slate-700/60 shadow-md">
          <div className="text-slate-400 text-xs uppercase tracking-widest font-semibold">Locally Verified Connections</div>
          <div className="text-4xl font-extrabold text-blue-500 mt-2">{metrics.activeConnections}</div>
        </div>
        <div className="bg-slate-800 p-5 rounded-xl border border-slate-700/60 shadow-md">
          <div className="text-slate-400 text-xs uppercase tracking-widest font-semibold">Mobility Spatial Latency</div>
          <div className="text-4xl font-extrabold text-cyan-400 mt-2">{metrics.spatialMatcherLatencyMs} <span className="text-sm font-medium text-slate-500">ms</span></div>
        </div>
        <div className="bg-slate-800 p-5 rounded-xl border border-slate-700/60 shadow-md">
          <div className="text-slate-400 text-xs uppercase tracking-widest font-semibold">A* Router Compute</div>
          <div className="text-4xl font-extrabold text-teal-400 mt-2">
            {metrics.routerExecutionLatencyMs > 0 ? `${metrics.routerExecutionLatencyMs} ms` : <span className="text-sm text-red-400 font-semibold">Network Offline</span>}
          </div>
        </div>
      </div>

      {/* Control Actions */}
      <div className="flex gap-4 bg-slate-950 p-4 rounded-xl border border-slate-800">
        <button onClick={() => launchSwarm(5)} className="bg-blue-600 hover:bg-blue-500 px-5 py-2.5 rounded-lg font-bold transition text-sm">Deploy 5</button>
        <button onClick={() => launchSwarm(50)} className="bg-purple-600 hover:bg-purple-500 px-5 py-2.5 rounded-lg font-bold transition text-sm">Deploy 50</button>
        <button onClick={() => launchSwarm(200)} className="bg-red-600 hover:bg-red-500 px-5 py-2.5 rounded-lg font-bold transition text-sm">Stress Load (200)</button>
        <button onClick={purgeSwarm} className="border border-slate-700 hover:bg-slate-800 px-5 py-2.5 rounded-lg font-bold transition text-sm ml-auto">Teardown Cluster Stream</button>
      </div>

      {/* Hardware Logging Feed */}
      <div className="flex-1 min-h-[250px] bg-slate-950 border border-slate-800 rounded-xl p-4 overflow-y-auto font-mono text-xs shadow-inner">
        {logs.length === 0 && <div className="text-slate-600 italic">No low-level telemetry streams established. Select deployment parameter above.</div>}
        {logs.map((l, i) => (
          <div key={i} className="py-0.5 border-b border-slate-900/40 last:border-0">
            <span className="text-slate-500 mr-2">[{l.time}]</span>
            <span className={
              l.type === 'warning' ? 'text-amber-400' : 
              l.type === 'success' ? 'text-emerald-400' : 
              l.type === 'danger' ? 'text-red-400' : 'text-blue-400'
            }>{l.msg}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
