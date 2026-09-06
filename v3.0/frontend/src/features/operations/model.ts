import type { AckStatus, CommandStatus, CreateTaskInput, PlanningMode, TaskPriority, TaskStatus } from '../../types/domain';

export const taskStatuses: TaskStatus[] = ['PENDING', 'ASSIGNING', 'ASSIGNED', 'IN_PROGRESS', 'COMPLETED', 'FAILED', 'CANCELLED', 'EXPIRED'];
export const commandStatuses: CommandStatus[] = ['PENDING', 'DELIVERED', 'ACKNOWLEDGED', 'COMPLETED', 'FAILED', 'EXPIRED', 'CANCELLED'];
export const taskPriorities: TaskPriority[] = ['LOW', 'NORMAL', 'HIGH', 'CRITICAL'];
export const taskTypes = ['NAVIGATE', 'RELOCATE', 'RETURN_HOME', 'STOP', 'CAPTURE_IMAGE', 'RUN_MODEL', 'THERMAL_SCAN', 'START_SCAN'] as const;
export const deviceTypes = ['connected_vehicle', 'delivery_drone', 'ground_robot', 'fixed_iot_sensor', 'static_camera', 'compute_node'];

const capabilityByType: Record<string, string> = {
  NAVIGATE: 'navigate', RETURN_HOME: 'navigate',
  RELOCATE: 'receive_relocation_command', STOP: 'receive_relocation_command',
  CAPTURE_IMAGE: 'capture_image', RUN_MODEL: 'run_model',
  THERMAL_SCAN: 'thermal_scan', START_SCAN: 'thermal_scan',
};

export function requiredCapability(taskType: string) { return capabilityByType[taskType] || ''; }
export function isTaskTerminal(status: TaskStatus) { return ['COMPLETED', 'FAILED', 'CANCELLED', 'EXPIRED'].includes(status); }
export function isCommandTerminal(status: CommandStatus) { return ['COMPLETED', 'FAILED', 'EXPIRED', 'CANCELLED'].includes(status); }
export function canCancelTaskState(status: TaskStatus, commandStatus?: CommandStatus) {
  return ['PENDING', 'ASSIGNING', 'ASSIGNED'].includes(status) && (!commandStatus || commandStatus === 'PENDING');
}
export function canRetryTaskState(status: TaskStatus) { return status === 'FAILED' || status === 'EXPIRED'; }
export function canRetryCommandState(status: CommandStatus) { return status === 'DELIVERED' || status === 'FAILED'; }
export function canCancelCommandState(status: CommandStatus) { return status === 'PENDING'; }

export function titleCase(value: string) {
  return value.toLowerCase().split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
}

export function priorityDescription(priority: TaskPriority) {
  return ({ LOW: 'Routine non-urgent work.', NORMAL: 'Default operational work.', HIGH: 'Prioritized ahead of normal tasks.', CRITICAL: 'Highest orchestration priority.' })[priority];
}

export function ackExplanation(status?: AckStatus) {
  if (!status) return 'No device acknowledgement has been recorded.';
  return ({
    ACCEPTED: 'The device accepted this command.',
    DUPLICATE: 'The device reported this command as already processed. Polaris treats this as an idempotent acknowledgement.',
    REJECTED: 'The device rejected this command.',
    EXPIRED: 'The device reported this command as expired.',
    UNSUPPORTED: 'The device does not support this command.',
  } as Record<string, string>)[status] || `The device returned acknowledgement status ${status}.`;
}

export interface TaskDraft {
  projectId: string;
  taskType: string;
  priority: TaskPriority;
  expiresAt: string;
  latitude: string;
  longitude: string;
  targetJSON: string;
  minimumBattery: string;
  maximumDistanceKm: string;
  deviceTypes: string[];
  optionalCapabilities: string[];
  planningMode: PlanningMode;
  customConstraintsJSON: string;
  correlationId: string;
}

export interface ValidationResult { errors: Record<string, string>; input?: CreateTaskInput }

function parseObject(value: string, field: string, errors: Record<string, string>) {
  try {
    const parsed = JSON.parse(value || '{}');
    if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') errors[field] = 'Enter a JSON object.';
    return parsed as Record<string, unknown>;
  } catch {
    errors[field] = 'Enter valid JSON.';
    return {};
  }
}

export function validateTaskDraft(draft: TaskDraft): ValidationResult {
  const errors: Record<string, string> = {};
  if (!taskTypes.includes(draft.taskType as typeof taskTypes[number])) errors.taskType = 'Select a supported task type.';
  if (!taskPriorities.includes(draft.priority)) errors.priority = 'Select a valid priority.';
  const expiry = new Date(draft.expiresAt);
  if (Number.isNaN(expiry.getTime()) || expiry.getTime() <= Date.now()) errors.expiresAt = 'Expiry must be in the future.';
  const spatial = draft.taskType === 'NAVIGATE' || draft.taskType === 'RELOCATE';
  let target: Record<string, unknown>;
  if (spatial) {
    const lat = Number(draft.latitude), lon = Number(draft.longitude);
    if (!draft.latitude.trim() || !Number.isFinite(lat) || lat < -90 || lat > 90) errors.latitude = 'Latitude must be between -90 and 90.';
    if (!draft.longitude.trim() || !Number.isFinite(lon) || lon < -180 || lon > 180) errors.longitude = 'Longitude must be between -180 and 180.';
    target = { lat, lon };
  } else target = parseObject(draft.targetJSON, 'targetJSON', errors);
  const battery = draft.minimumBattery === '' ? 0 : Number(draft.minimumBattery);
  if (!Number.isInteger(battery) || battery < 0 || battery > 100) errors.minimumBattery = 'Battery must be a whole percentage from 0 to 100.';
  const distanceKm = draft.maximumDistanceKm === '' ? 0 : Number(draft.maximumDistanceKm);
  if (!Number.isFinite(distanceKm) || distanceKm < 0) errors.maximumDistanceKm = 'Distance cannot be negative.';
  const custom = draft.customConstraintsJSON.trim() ? parseObject(draft.customConstraintsJSON, 'customConstraintsJSON', errors) : undefined;
  if (Object.keys(errors).length) return { errors };
  const automatic = requiredCapability(draft.taskType);
  const required = [...new Set([automatic, ...draft.optionalCapabilities].filter(Boolean))];
  const requirements = {
    ...(required.length ? { required_capabilities: required } : {}),
    ...(battery ? { minimum_battery: battery } : {}),
    ...(draft.deviceTypes.length ? { allowed_device_types: draft.deviceTypes } : {}),
    ...(distanceKm ? { max_distance_meters: distanceKm * 1000 } : {}),
    ...(draft.projectId ? { project_id: draft.projectId } : {}),
    ...(spatial ? { planning_mode: draft.planningMode } : {}),
    ...(custom ? { custom_constraints: custom } : {}),
  };
  return { errors, input: { project_id: draft.projectId || undefined, task_type: draft.taskType, priority: draft.priority, requirements, target, expires_at: expiry.toISOString() } };
}

export function microseconds(value?: number) {
  if (!Number.isFinite(value)) return 'Not returned';
  return value! >= 1000 ? `${(value! / 1000).toFixed(2)} ms` : `${value} µs`;
}
