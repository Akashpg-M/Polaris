import { createContext, useContext, useEffect, useMemo, useRef, useState, type ReactNode } from 'react';
import { gatewayWs } from '../../config';
import { useQueryClient, type QueryKey } from '../../lib/query/queryClient';
import type { DashboardTelemetry, DeviceTwin, FleetActivity, FleetFilters, SocketStatus, TwinComponent } from '../../types/domain';
import { useAuth } from '../auth/AuthContext';
import { parseDashboardTelemetry } from './events';

interface DashboardValue {
  status: SocketStatus;
  lastEventAt: string | null;
  activities: FleetActivity[];
}

const DashboardContext = createContext<DashboardValue | null>(null);

function observedISO(event: DashboardTelemetry) { return new Date(event.observed_at).toISOString(); }

function component(type: string, event: DashboardTelemetry, payload: Record<string, unknown>): TwinComponent {
  return {
    type,
    schema_version: 1,
    observed_at: observedISO(event),
    boot_id: event.device_boot_id,
    sequence_number: event.sequence_number,
    payload,
  };
}

function patchTwin(twin: DeviceTwin, event: DashboardTelemetry): DeviceTwin {
  if (twin.device_id !== event.device_id || twin.tenant_id !== event.tenant_id) return twin;
  const current = twin.reported_state;
  if (current && current.device_boot_id === event.device_boot_id && current.sequence_number > event.sequence_number) return twin;
  const mobility = event.type === 5 ? 'AERIAL_DRONE' : event.type === 6 ? 'GROUND_ROBOT' : event.type === 7 ? 'STATIC' : 'ROAD_VEHICLE';
  return {
    ...twin,
    reported_state: event,
    connectivity: { status: 'ONLINE', last_seen_at: new Date().toISOString() },
    components: {
      ...twin.components,
      'spatial/v1': component('spatial/v1', event, {
        latitude: event.lat,
        longitude: event.lon,
        heading_degrees: event.heading_deg,
        speed_mps: event.velocity_mps,
        mobility_profile: mobility,
      }),
      'battery/v1': component('battery/v1', event, { percent: event.energy_percent }),
    },
  };
}

function matches(twin: DeviceTwin, filters: FleetFilters) {
  const search = filters.search?.trim().toLowerCase();
  return (!filters.projectId || twin.device.project_id === filters.projectId)
    && (!filters.deviceType || twin.device.device_type_id === filters.deviceType)
    && (!filters.lifecycle || twin.device.lifecycle_status === filters.lifecycle)
    && (!filters.connectivity || twin.connectivity.status === filters.connectivity)
    && (!filters.capability || twin.capabilities.some(item => item.enabled && item.capability_id === filters.capability))
    && (!search || twin.device_id.toLowerCase().includes(search) || twin.device.display_name.toLowerCase().includes(search));
}

function twinListKey(key: QueryKey, tenantId: string) {
  return key[0] === 'twins' && key[1] === tenantId;
}

export function DashboardProvider({ children }: { children: ReactNode }) {
  const { session, api } = useAuth();
  const queryClient = useQueryClient();
  const [status, setStatus] = useState<SocketStatus>('DISCONNECTED');
  const [lastEventAt, setLastEventAt] = useState<string | null>(null);
  const [activities, setActivities] = useState<FleetActivity[]>([]);
  const seen = useRef(new Set<string>());

  useEffect(() => {
    if (!session || !api) {
      setStatus('DISCONNECTED');
      return;
    }
    let disposed = false;
    let socket: WebSocket | null = null;
    let retryTimer: ReturnType<typeof setTimeout> | undefined;
    let attempt = 0;
    const tenantId = session.tenantId;
    const seenEvents = seen.current;

    const apply = async (event: DashboardTelemetry) => {
      if (event.tenant_id !== tenantId) return;
      let found = false;
      queryClient.updateWhere<DeviceTwin[]>(
        key => twinListKey(key, tenantId),
        (current, key) => {
          const filters = (key[2] || {}) as FleetFilters;
          const limit = Number(key[3]) || 50;
          const next = current.map(twin => {
            if (twin.device_id !== event.device_id) return twin;
            found = true;
            return patchTwin(twin, event);
          }).filter(twin => matches(twin, filters));
          return next.slice(0, limit || 50).sort((a, b) => a.device_id.localeCompare(b.device_id));
        },
      );
      queryClient.set<DeviceTwin>(['twin', tenantId, event.device_id], current => current ? patchTwin(current, event) : current);

      if (!found) {
        try {
          const hydrated = await api.twin(event.device_id);
          if (disposed) return;
          const patched = patchTwin(hydrated, event);
          queryClient.set(['twin', tenantId, event.device_id], patched);
          queryClient.updateWhere<DeviceTwin[]>(
            key => twinListKey(key, tenantId),
            (current, key) => {
              const filters = (key[2] || {}) as FleetFilters;
              const limit = Number(key[3]) || 50;
              const cursor = String(key[4] || '');
              if (cursor || !matches(patched, filters) || current.some(twin => twin.device_id === patched.device_id)) return current;
              return [...current, patched].sort((a, b) => a.device_id.localeCompare(b.device_id)).slice(0, limit);
            },
          );
        } catch {
          // A live event may race a lifecycle change; authoritative refresh repairs it.
        }
      }

      setLastEventAt(new Date().toISOString());
      if (!seenEvents.has(event.event_id)) {
        seenEvents.add(event.event_id);
        setActivities(current => [{
          eventId: event.event_id,
          deviceId: event.device_id,
          observedAt: observedISO(event),
          message: `Telemetry accepted · ${event.energy_percent}% battery`,
        }, ...current].slice(0, 12));
      }
    };

    const scheduleReconnect = () => {
      if (disposed) return;
      setStatus('RECONNECTING');
      const delay = Math.min(30_000, 1_000 * 2 ** Math.min(attempt++, 5));
      retryTimer = setTimeout(connect, delay + Math.floor(Math.random() * 300));
    };

    const connect = async () => {
      if (disposed) return;
      setStatus(attempt ? 'RECONNECTING' : 'DISCONNECTED');
      try {
        const { ticket } = await api.dashboardTicket();
        if (disposed) return;
        socket = new WebSocket(`${gatewayWs}/ws/dashboard?ticket=${encodeURIComponent(ticket)}`);
        socket.onopen = () => {
          attempt = 0;
          setStatus('CONNECTED');
          queryClient.refresh(['twins', tenantId]);
          queryClient.refresh(['twin', tenantId]);
        };
        socket.onmessage = message => {
          try {
            const event = parseDashboardTelemetry(JSON.parse(String(message.data)));
            if (event) void apply(event);
          } catch {
            // Ignore malformed transient frames without breaking the live connection.
          }
        };
        socket.onerror = () => socket?.close();
        socket.onclose = scheduleReconnect;
      } catch {
        scheduleReconnect();
      }
    };

    void connect();
    return () => {
      disposed = true;
      if (retryTimer) clearTimeout(retryTimer);
      if (socket) {
        socket.onclose = null;
        socket.close();
      }
      setStatus('DISCONNECTED');
      seenEvents.clear();
      setActivities([]);
    };
  }, [api, queryClient, session]);

  const value = useMemo(() => ({ status, lastEventAt, activities }), [status, lastEventAt, activities]);
  return <DashboardContext.Provider value={value}>{children}</DashboardContext.Provider>;
}

export function useDashboard() {
  const value = useContext(DashboardContext);
  if (!value) throw new Error('DashboardProvider is missing');
  return value;
}
