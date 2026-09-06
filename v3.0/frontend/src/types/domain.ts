export type OperatorRole = 'PLATFORM_ADMIN' | 'TENANT_ADMIN' | 'OPERATOR' | 'VIEWER';
export type TenantStatus = 'ACTIVE' | 'SUSPENDED' | 'DEACTIVATED';
export type ProjectStatus = 'ACTIVE' | 'ARCHIVED' | string;
export type DeviceLifecycle = 'REGISTERED' | 'ACTIVE' | 'SUSPENDED' | 'DECOMMISSIONED';
export type ConnectivityStatus = 'NEVER_CONNECTED' | 'ONLINE' | 'STALE' | 'OFFLINE';
export type SocketStatus = 'CONNECTED' | 'RECONNECTING' | 'DISCONNECTED';

export interface Tenant {
  tenant_id: string;
  display_name: string;
  status: TenantStatus;
  metadata: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface Project {
  project_id: string;
  tenant_id: string;
  name: string;
  description?: string;
  status: ProjectStatus;
  metadata: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface Device {
  tenant_id: string;
  device_id: string;
  project_id?: string;
  device_type_id: string;
  display_name: string;
  lifecycle_status: DeviceLifecycle;
  firmware_version?: string;
  software_version?: string;
  model_version?: string;
  metadata: Record<string, unknown>;
  registered_at: string;
  updated_at: string;
  deactivated_at?: string;
}

export interface Capability {
  capability_id: string;
  display_name: string;
  description?: string;
  configuration: Record<string, unknown>;
  enabled: boolean;
}

export interface ReportedState {
  event_id: string;
  schema_version: number;
  id: string;
  device_id: string;
  tenant_id: string;
  device_boot_id: string;
  sequence_number: number;
  boot_started_at: number;
  type: number;
  status: number;
  lat: number;
  lon: number;
  velocity_mps: number;
  heading_deg: number;
  energy_percent: number;
  observed_at: number;
  ingested_at: number;
  timestamp: number;
}

export interface TwinComponent {
  type: string;
  schema_version: number;
  observed_at: string;
  boot_id: string;
  sequence_number: number;
  payload: Record<string, unknown>;
}

export interface DeviceTwin {
  tenant_id: string;
  device_id: string;
  device: Device;
  capabilities: Capability[];
  reported_state: ReportedState | null;
  components: Record<string, TwinComponent>;
  desired_state: null;
  connectivity: {
    status: ConnectivityStatus;
    last_seen_at: string | null;
  };
}

export type DashboardTelemetry = ReportedState;

export interface FleetFilters {
  projectId?: string;
  deviceType?: string;
  lifecycle?: DeviceLifecycle | '';
  connectivity?: ConnectivityStatus | '';
  capability?: string;
  search?: string;
}

export interface FleetActivity {
  eventId: string;
  deviceId: string;
  observedAt: string;
  message: string;
}

export interface AuthSession {
  token: string;
  tenantId: string;
  role: OperatorRole;
  tenant: Tenant;
}

export type TaskStatus = 'PENDING' | 'ASSIGNING' | 'ASSIGNED' | 'IN_PROGRESS' | 'COMPLETED' | 'FAILED' | 'CANCELLED' | 'EXPIRED';
export type TaskPriority = 'LOW' | 'NORMAL' | 'HIGH' | 'CRITICAL';
export type PlanningMode = 'DEVICE_LOCAL' | 'POLARIS_REQUIRED';
export type CommandStatus = 'PENDING' | 'DELIVERED' | 'ACKNOWLEDGED' | 'COMPLETED' | 'FAILED' | 'EXPIRED' | 'CANCELLED';
export type AckStatus = 'ACCEPTED' | 'REJECTED' | 'DUPLICATE' | 'EXPIRED' | 'UNSUPPORTED' | string;

export interface TaskRequirements {
  required_capabilities?: string[];
  minimum_battery?: number;
  allowed_device_types?: string[];
  max_distance_meters?: number;
  project_id?: string;
  planning_mode?: PlanningMode;
  custom_constraints?: unknown;
}

export interface TaskTarget {
  lat?: number;
  lon?: number;
  h3_cell?: string;
  data?: unknown;
  [key: string]: unknown;
}

export interface Task {
  task_id: string;
  tenant_id: string;
  project_id?: string;
  task_type: string;
  status: TaskStatus;
  priority: TaskPriority;
  requirements: TaskRequirements;
  target: TaskTarget;
  assigned_device_id?: string;
  correlation_id: string;
  created_by: string;
  version: number;
  created_at: string;
  updated_at: string;
  assigned_at?: string;
  started_at?: string;
  completed_at?: string;
  failed_at?: string;
  expires_at: string;
  failure_reason?: string;
}

export interface Command {
  command_id: string;
  tenant_id: string;
  device_id: string;
  task_id: string;
  command_type: string;
  payload: Record<string, unknown>;
  status: CommandStatus;
  sequence_number: number;
  correlation_id: string;
  causation_id: string;
  attempt_count: number;
  max_attempts: number;
  version: number;
  created_at: string;
  available_at: string;
  sent_at?: string;
  acknowledged_at?: string;
  completed_at?: string;
  expires_at: string;
  ack_status?: AckStatus;
  result?: unknown;
  last_error?: string;
}

export interface TaskPathTiming {
  candidate_selection_duration_us: number;
  routing_duration_us: number;
  persistence_duration_us: number;
  total_duration_us: number;
}

export interface TaskDetail { task: Task; commands: Command[] }
export interface CreateTaskResult { task: Task; command?: Command; timing?: TaskPathTiming }

export interface CreateTaskInput {
  project_id?: string;
  task_type: string;
  priority: TaskPriority;
  requirements: TaskRequirements;
  target: TaskTarget;
  expires_at: string;
}

export interface TaskFilters { status?: TaskStatus | ''; deviceId?: string }
export interface CommandFilters { status?: CommandStatus | ''; taskId?: string; deviceId?: string }

export interface OperationsActivity {
  id: string;
  kind: 'TASK' | 'COMMAND';
  entityId: string;
  relatedId?: string;
  deviceId?: string;
  state: TaskStatus | CommandStatus;
  observedAt: string;
  initial: boolean;
}

export type MobilityProfile = 'ROAD_VEHICLE' | 'GROUND_ROBOT' | 'AERIAL_DRONE' | 'STATIC';
export type RoutePolicy = 'SHORTEST' | 'FASTEST';
export type ModuleState = 'STARTING' | 'READY' | 'DEGRADED' | 'FAILED' | 'STOPPED';

export interface MobilityPosition {
  latitude: number;
  longitude: number;
  altitude_meters?: number;
}

export interface MobilityQuality {
  valid: boolean;
  confidence: number;
  anomalies?: string[];
}

export interface SpatialState {
  tenant_id: string;
  device_id: string;
  position: MobilityPosition;
  reported_position: MobilityPosition;
  heading_degrees?: number;
  speed_mps?: number;
  mobility_profile: MobilityProfile;
  h3_cell: number;
  observed_at: string;
  indexed_at: string;
  boot_id: string;
  boot_started_at: string;
  sequence_number: number;
  quality: MobilityQuality;
}

export interface SpatialCandidate { state: SpatialState; distance_meters: number }
export interface NearbySearch { latitude: number; longitude: number; radiusMeters: number; limit: number }
export interface NearbyResult { count: number; devices: SpatialCandidate[] }

export interface RouteRequest {
  mobility_profile: MobilityProfile;
  origin: MobilityPosition;
  destination: MobilityPosition;
  policy: RoutePolicy;
}

export interface RouteResult {
  route_id: string;
  road_graph_version: string;
  snapshot_version: number;
  policy: RoutePolicy;
  distance_meters: number;
  /** Go time.Duration JSON representation: nanoseconds. */
  estimated_time: number;
  waypoints: MobilityPosition[];
  edge_ids: number[];
  expanded_nodes: number;
}

export interface RoutingRuntime {
  requests: number;
  routing_busy: number;
  queue_depth: number;
  queue_capacity: number;
  active_tenants: number;
}

export interface ModuleStatus {
  state: ModuleState;
  message?: string;
  components?: Record<string, ModuleStatus>;
  details?: {
    road_graph_version?: string;
    road_nodes?: number;
    road_edges?: number;
    routing_snapshot_version?: number;
    routing_runtime?: RoutingRuntime;
    traffic_scope?: string;
    traffic_refresh_interval?: string;
    traffic_edge_states?: number;
    traffic_overlay_bytes?: number;
    [key: string]: unknown;
  };
}

export interface EngineReadiness {
  status: string;
  core?: string;
  modules?: Record<string, ModuleStatus>;
  runtime?: Record<string, number>;
  dependency?: string;
  error?: string;
}

export interface PredictedZone {
  id: string;
  lat: number;
  lon: number;
  radius_km: number;
  required_assets: number;
  target_class: number;
  tenant_id: string;
}
