import { engineApi, engineReadyz, gatewayApi } from '../../config';
import type { Capability, Command, CommandFilters, CreateTaskInput, CreateTaskResult, Device, DeviceTwin, EngineReadiness, FleetFilters, NearbyResult, NearbySearch, PredictedZone, Project, RouteRequest, RouteResult, Task, TaskDetail, TaskFilters, Tenant } from '../../types/domain';
import { PolarisApiError } from './errors';

interface SuccessEnvelope<T> { data: T; request_id?: string }
interface ErrorEnvelope { error?: { code?: string; message?: string }; request_id?: string }

export interface ApiIdentity {
  token: string;
  tenantId: string;
}

function requestId() {
  return globalThis.crypto?.randomUUID?.() ?? `web-${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

function queryString(values: Record<string, string | number | undefined>) {
  const search = new URLSearchParams();
  Object.entries(values).forEach(([key, value]) => {
    if (value !== undefined && value !== '') search.set(key, String(value));
  });
  const encoded = search.toString();
  return encoded ? `?${encoded}` : '';
}

async function parseResponse<T>(response: Response): Promise<T> {
  let body: SuccessEnvelope<T> & ErrorEnvelope;
  try {
    body = await response.json();
  } catch {
    throw new PolarisApiError({ message: 'Polaris returned an unreadable response.', status: response.status });
  }
  if (!response.ok) {
    throw new PolarisApiError({
      message: body.error?.message || `Polaris rejected the request (${response.status}).`,
      code: body.error?.code,
      requestId: body.request_id,
      status: response.status,
      details: body,
    });
  }
  return body.data;
}

export class PolarisApi {
  private identity: ApiIdentity;
  private unauthorized?: () => void;

  constructor(identity: ApiIdentity, unauthorized?: () => void) {
    this.identity = identity;
    this.unauthorized = unauthorized;
  }

  private async request<T>(path: string, options: RequestInit = {}, base = engineApi): Promise<T> {
    const headers = new Headers(options.headers);
    headers.set('Authorization', `Bearer ${this.identity.token}`);
    headers.set('X-Tenant-ID', this.identity.tenantId);
    headers.set('X-Request-ID', requestId());
    if (options.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json');
    let response: Response;
    try {
      response = await fetch(`${base}${path}`, { ...options, headers });
    } catch (error) {
      if (error instanceof DOMException && error.name === 'AbortError') throw error;
      throw new PolarisApiError({ message: 'Unable to reach Polaris. Check the deployment and try again.', status: 0 });
    }
    if (response.status === 401) this.unauthorized?.();
    return parseResponse<T>(response);
  }

  tenant(tenantId = this.identity.tenantId, signal?: AbortSignal) {
    return this.request<Tenant>(`/tenants/${encodeURIComponent(tenantId)}`, { signal });
  }

  projects(signal?: AbortSignal) {
    return this.request<Project[]>('/projects', { signal });
  }

  project(projectId: string, signal?: AbortSignal) {
    return this.request<Project>(`/projects/${encodeURIComponent(projectId)}`, { signal });
  }

  devices(filters: FleetFilters = {}, limit = 50, cursor = '', signal?: AbortSignal) {
    const query = queryString({
      limit,
      cursor,
      project_id: filters.projectId,
      device_type: filters.deviceType,
      lifecycle_status: filters.lifecycle,
      capability: filters.capability,
    });
    return this.request<Device[]>(`/devices${query}`, { signal });
  }

  device(deviceId: string, signal?: AbortSignal) {
    return this.request<Device>(`/devices/${encodeURIComponent(deviceId)}`, { signal });
  }

  capabilities(deviceId?: string, signal?: AbortSignal) {
    return this.request<Capability[]>(deviceId
      ? `/devices/${encodeURIComponent(deviceId)}/capabilities`
      : '/capabilities', { signal });
  }

  twins(filters: FleetFilters = {}, limit = 50, cursor = '', signal?: AbortSignal) {
    const query = queryString({
      limit,
      cursor,
      project_id: filters.projectId,
      device_type: filters.deviceType,
      lifecycle_status: filters.lifecycle,
      connectivity_status: filters.connectivity,
      capability: filters.capability,
    });
    return this.request<DeviceTwin[]>(`/twins${query}`, { signal });
  }

  twin(deviceId: string, signal?: AbortSignal) {
    return this.request<DeviceTwin>(`/devices/${encodeURIComponent(deviceId)}/twin`, { signal });
  }

  dashboardTicket(signal?: AbortSignal) {
    return this.request<{ ticket: string; expires_in_seconds: number }>('/dashboard-ticket', {
      method: 'POST', body: '{}', signal,
    });
  }

  async activeGatewayConnections(signal?: AbortSignal) {
    const response = await fetch(`${gatewayApi}/metrics/connections`, { signal });
    if (!response.ok) throw new PolarisApiError({ message: 'Gateway connection metric is unavailable.', status: response.status });
    return response.json() as Promise<{ active_uplinks: number }>;
  }

  tasks(filters: TaskFilters = {}, limit = 50, cursor = '', signal?: AbortSignal) {
    const query = queryString({ limit, cursor, status: filters.status, device_id: filters.deviceId });
    return this.request<Task[]>(`/tasks${query}`, { signal });
  }

  task(taskId: string, signal?: AbortSignal) {
    return this.request<TaskDetail>(`/tasks/${encodeURIComponent(taskId)}`, { signal });
  }

  createTask(input: CreateTaskInput, correlationId?: string, signal?: AbortSignal) {
    const headers = correlationId ? { 'X-Correlation-ID': correlationId } : undefined;
    return this.request<CreateTaskResult>('/tasks', { method: 'POST', headers, body: JSON.stringify(input), signal });
  }

  cancelTask(taskId: string, signal?: AbortSignal) {
    return this.request<{ task_id: string; status: 'CANCELLED' }>(`/tasks/${encodeURIComponent(taskId)}/cancel`, { method: 'POST', body: '{}', signal });
  }

  retryTask(taskId: string, ttlSeconds = 300, signal?: AbortSignal) {
    return this.request<CreateTaskResult>(`/tasks/${encodeURIComponent(taskId)}/retry`, { method: 'POST', body: JSON.stringify({ ttl_seconds: ttlSeconds }), signal });
  }

  commands(filters: CommandFilters = {}, limit = 50, cursor = '', signal?: AbortSignal) {
    const query = queryString({ limit, cursor, status: filters.status, task_id: filters.taskId, device_id: filters.deviceId });
    return this.request<Command[]>(`/commands${query}`, { signal });
  }

  command(commandId: string, signal?: AbortSignal) {
    return this.request<Command>(`/commands/${encodeURIComponent(commandId)}`, { signal });
  }

  retryCommand(commandId: string, signal?: AbortSignal) {
    return this.request<{ command_id: string; status: 'PENDING' }>(`/commands/${encodeURIComponent(commandId)}/retry`, { method: 'POST', body: '{}', signal });
  }

  cancelCommand(commandId: string, signal?: AbortSignal) {
    return this.request<{ command_id: string; status: 'CANCELLED' }>(`/commands/${encodeURIComponent(commandId)}/cancel`, { method: 'POST', body: '{}', signal });
  }

  nearby(search: NearbySearch, signal?: AbortSignal) {
    const query = queryString({ lat: search.latitude, lon: search.longitude, radius_meters: search.radiusMeters, limit: search.limit });
    return this.request<NearbyResult>(`/spatial/devices/nearby${query}`, { signal });
  }

  route(input: RouteRequest, signal?: AbortSignal) {
    return this.request<RouteResult>('/routes', { method: 'POST', body: JSON.stringify(input), signal });
  }

  predictedZones(signal?: AbortSignal) {
    return this.request<PredictedZone[]>('/zones/predicted', { signal });
  }

  async mobilityReadiness(signal?: AbortSignal) {
    const headers = new Headers({ 'X-Request-ID': requestId() });
    let response: Response;
    try {
      response = await fetch(engineReadyz, { headers, signal });
    } catch (error) {
      if (error instanceof DOMException && error.name === 'AbortError') throw error;
      throw new PolarisApiError({ message: 'Engine readiness is unreachable.', status: 0 });
    }
    let body: EngineReadiness;
    try { body = await response.json() as EngineReadiness; }
    catch { throw new PolarisApiError({ message: 'Engine readiness returned an unreadable response.', status: response.status }); }
    if (!response.ok) throw new PolarisApiError({ message: body.error || `Engine is not ready (${response.status}).`, status: response.status, details: body });
    return body;
  }
}
