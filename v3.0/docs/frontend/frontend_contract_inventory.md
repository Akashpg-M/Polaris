# Polaris Frontend Phase 1 Contract Inventory

This inventory was derived from the active Gateway and Engine routes and repository behavior. It supersedes assumptions in the frontend phase brief wherever they differ.

## Transport and response conventions

- Container browser traffic is same-origin. Engine REST calls use `/api/engine/*`, Gateway REST calls use `/api/gateway/*`, and Gateway sockets use `/ws/*` through Nginx.
- Engine control APIs accept `Authorization: Bearer <operator-token>`.
- A platform administrator must also select a tenant using `X-Tenant-ID` (or the `tenant_id` query fallback). Tenant-scoped principals derive their tenant in the backend; sending the matching header is harmless, but the frontend must never use it to imply authorization.
- Success responses use `{ "data": <payload>, "request_id": "..." }`.
- Error responses use `{ "error": { "code": "...", "message": "..." }, "request_id": "..." }`.
- The backend preserves caller-provided `X-Request-ID`; otherwise it generates one.
- Common errors are `UNAUTHENTICATED` (401), `FORBIDDEN` (403), `TENANT_REQUIRED` or `INVALID_REQUEST` (400), `NOT_FOUND` (404), `CONFLICT` or `INVALID_LIFECYCLE_TRANSITION` (409), and `INTERNAL_ERROR` (500).
- Read access is available to `PLATFORM_ADMIN`, `TENANT_ADMIN`, `OPERATOR`, and `VIEWER`. The backend—not the frontend—is the security boundary.

## Phase 1 endpoints

| Method | Path | Authentication / role | Tenant scoping | Query / request | Response data | Pagination | Frontend use |
| --- | --- | --- | --- | --- | --- | --- | --- |
| GET | `/api/engine/tenants/:tenant_id` | Bearer; any read role | Platform admin may read explicit tenant; other roles may read only their own | Path tenant ID | Tenant: `tenant_id`, `display_name`, `status`, `metadata`, `created_at`, `updated_at` | None | Validate access session; header tenant context |
| GET | `/api/engine/projects` | Bearer; any read role | Current tenant | None | Array of Project: `project_id`, `tenant_id`, `name`, optional `description`, `status`, `metadata`, timestamps | Fixed maximum 100; no cursor metadata | Project filter; Projects page |
| GET | `/api/engine/projects/:project_id` | Bearer; any read role | Current tenant and project ID | Path project ID | One Project | None | Project detail |
| GET | `/api/engine/devices` | Bearer; any read role | Current tenant | `limit` (1–100; invalid becomes 50), `cursor` (exclusive device ID), `project_id`, `device_type`, `lifecycle_status`, `capability` | Array of Device registry records | Forward cursor is inferred from last returned `device_id`; response has no next-cursor/total metadata | Registry-oriented fleet page and project device list |
| GET | `/api/engine/devices/:device_id` | Bearer; any read role | Current tenant and device ID | Path device ID | One Device registry record | None | Device metadata fallback/detail |
| GET | `/api/engine/devices/:device_id/capabilities` | Bearer; any read role | Current tenant and device ID | Path device ID | Array of assigned Capability: `capability_id`, `display_name`, optional `description`, `configuration`, `enabled` | None | Device capabilities |
| GET | `/api/engine/capabilities` | Bearer; any read role | Global catalog behind authenticated API | None | Array of capability catalog records | None | Labels and capability vocabulary |
| GET | `/api/engine/twins` | Bearer; any read role | Current tenant | Same registry filters as devices plus `connectivity_status` | Array of hydrated twins | Same forward device cursor is accepted but not returned; connectivity is filtered after bounded registry selection | Overview aggregation, Devices, Twins, initial map hydration |
| GET | `/api/engine/devices/:device_id/twin` | Bearer; any read role | Current tenant and device ID | Path device ID | DeviceTwin described below | None | Device detail tabs and map drawer hydration |
| POST | `/api/engine/dashboard-ticket` | Bearer; any read role | Current tenant | Empty JSON body accepted | `{ ticket, expires_in_seconds }` | None | Obtain single-use browser-safe WebSocket ticket |
| GET | `/api/gateway/metrics/connections` | No operator middleware on Gateway | Not tenant-scoped | None | `{ active_uplinks }` | None | Optional platform connection diagnostic; not a tenant fleet metric |
| GET | `/api/engine/healthz`, `/api/engine/readyz` | None | None | None | Process/dependency status | None | Not used as full observability; optional startup diagnostics |

## Device twin contract

`GET /devices/:device_id/twin` and each `/twins` element return:

```text
tenant_id, device_id
device: durable registry Device
capabilities: assigned Capability[]
reported_state: null or latest dashboard-shaped telemetry
components: map<string, ComponentEnvelope>
desired_state: currently null
connectivity: { status, last_seen_at }
```

Device fields are `tenant_id`, `device_id`, optional `project_id`, `device_type_id`, `display_name`, `lifecycle_status`, optional firmware/software/model versions, `metadata`, `registered_at`, `updated_at`, and optional `deactivated_at`.

Connectivity values are `NEVER_CONNECTED`, `ONLINE`, `STALE`, and `OFFLINE`. Registry lifecycle is independently `REGISTERED`, `ACTIVE`, `SUSPENDED`, or `DECOMMISSIONED`.

Known component envelopes contain `type`, `schema_version`, `observed_at`, `boot_id`, `sequence_number`, and an opaque `payload`. Current known keys are `spatial/v1` and `battery/v1`; unknown keys are valid future components.

Current reported telemetry contains `event_id`, `schema_version`, `id`, `device_id`, `tenant_id`, `device_boot_id`, `sequence_number`, `boot_started_at`, numeric `type` and `status`, `lat`, `lon`, `velocity_mps`, `heading_deg`, `energy_percent`, `observed_at`, `ingested_at`, and legacy `timestamp`.

## Dashboard WebSocket contract

| Field | Contract |
| --- | --- |
| Path | `/ws/dashboard?ticket=<single-use-ticket>` |
| Authentication | Ticket obtained from authenticated `POST /api/engine/dashboard-ticket` |
| Tenant | Bound into the ticket; Gateway forwards only matching `tenant_id` events |
| Direction | Server-to-browser telemetry stream; browser messages are not part of the product contract |
| Replay | None. Redis Pub/Sub is transient. API hydration is required before connect and after reconnect. |
| Heartbeat | No application heartbeat frame is defined |
| Message | JSON normalized reported telemetry with the fields listed above |
| Failure handling | Malformed messages must be ignored; obtain a new ticket for each reconnect because tickets are single-use |

## Device type values visible in registry

The seeded backend catalog uses `delivery_drone`, `ground_robot`, `connected_vehicle`, `fixed_iot_sensor`, `static_camera`, and `compute_node`. There is no read endpoint for the device-type catalog, so the frontend maintains presentation metadata for these known values and safely displays unknown future IDs.

## Permission model required by Phase 1

| Permission | Platform admin | Tenant admin | Operator | Viewer |
| --- | ---: | ---: | ---: | ---: |
| Read fleet/projects/twins | Yes | Yes | Yes | Yes |
| Manage registry | Yes | Yes | No | No |
| Create/cancel tasks | Yes | Yes | Yes | No |
| Administrative retry | Yes | Yes | No | No |
| Read audit | Yes | Yes | No | No |

Phase 1 is read-only. Permission utilities model future navigation, but no mutation controls are exposed.

## Confirmed contract discrepancies from the proposed plan

- There is no login/session/current-user endpoint. An operator token is provisioned out of band, and the backend does not return the principal role or tenant after authentication.
- There is no tenant-list endpoint. Platform administrators must already know/select a tenant ID.
- Device and twin list responses do not return total counts, `next_cursor`, or `has_more`.
- `/twins` performs server-side hydration but applies connectivity filtering after selecting the bounded device page.
- There is no telemetry-history read endpoint despite durable telemetry storage.
- There is no fleet aggregate, global search, or project-count endpoint.
- Gateway `active_uplinks` is global to that Gateway process and cannot be labelled as tenant online-device count.

## Phase 2 task and command contracts

All Operations endpoints use the Phase 1 bearer and tenant rules and the standard `{data, request_id}` / `{error, request_id}` envelopes.

| Method | Path | Role | Request / filters | Response data | Semantics |
| --- | --- | --- | --- | --- | --- |
| POST | `/api/engine/tasks` | Platform admin, tenant admin, operator | `project_id?`, `task_type`, `priority?`, `requirements`, arbitrary valid JSON `target`, `expires_at?`; optional `X-Correlation-ID` | `{task, command?, timing?}` | Persists intent, then immediately attempts assignment. No eligible device is an accepted `PENDING` result, not an API error. |
| GET | `/api/engine/tasks` | All read roles | `limit` 1–100, exclusive `cursor` task ID, `status`, `device_id` | `Task[]` ordered by task ID | No total, next cursor, task type, priority, project, or date filter. |
| GET | `/api/engine/tasks/:task_id` | All read roles | Path ID | `{task, commands}`; commands bounded to 100 | Tenant-scoped task and its current commands. |
| POST | `/api/engine/tasks/:task_id/cancel` | Platform admin, tenant admin, operator | No body required | `{task_id,status:"CANCELLED"}` | Legal only while task is pending/assigning/assigned and any command is still `PENDING`. It does not claim a delivered physical command stopped. |
| POST | `/api/engine/tasks/:task_id/retry` | Platform admin, tenant admin | `{ttl_seconds?}` | `{task,command?,timing?}` | Legal for `FAILED` or `EXPIRED`; clears assignment/execution state, extends expiry, makes a new assignment decision and may create a new command. |
| GET | `/api/engine/commands` | All read roles | `limit` 1–100, exclusive `cursor` command ID, `status`, `task_id`, `device_id` | `Command[]` ordered by command ID | No total, next cursor, type, project, or date filter. |
| GET | `/api/engine/commands/:command_id` | All read roles | Path ID | One `Command` | Includes aggregate attempts and ACK/result fields but no per-attempt records. |
| POST | `/api/engine/commands/:command_id/retry` | Platform admin, tenant admin | No body required | `{command_id,status:"PENDING"}` | Reuses immutable command ID, device sequence and payload. Supports `DELIVERED`, or administratively reactivates a `FAILED` command/assignment when conflict-free and unexpired. |
| POST | `/api/engine/commands/:command_id/cancel` | Platform admin, tenant admin, operator | No body required | `{command_id,status:"CANCELLED"}` | Delegates to conservative task cancellation and is therefore legal only before delivery. |
| GET | `/api/engine/metrics/orchestration` | All read roles | None | Process-local cumulative counters | Diagnostic, not durable history and not used as entity state. |

### Task document

A task contains `task_id`, `tenant_id`, optional `project_id`, `task_type`, `status`, `priority`, JSON `requirements`, JSON `target`, optional `assigned_device_id`, `correlation_id`, `created_by`, `version`, `created_at`, `updated_at`, optional `assigned_at`, `started_at`, `completed_at`, `failed_at`, `expires_at`, and optional `failure_reason`.

Task states are `PENDING`, `ASSIGNING`, `ASSIGNED`, `IN_PROGRESS`, `COMPLETED`, `FAILED`, `CANCELLED`, and `EXPIRED`. Priorities are `LOW`, `NORMAL`, `HIGH`, and `CRITICAL`.

Requirements support `required_capabilities`, `minimum_battery` (0–100), `allowed_device_types`, `max_distance_meters`, `project_id`, `planning_mode`, and arbitrary `custom_constraints`. Planning mode is valid only for `NAVIGATE` and `RELOCATE`: `DEVICE_LOCAL` sends high-level intent and `POLARIS_REQUIRED` requires a domain planner. The backend automatically adds mandatory capability rules: `NAVIGATE`/`RETURN_HOME` → `navigate`; `RELOCATE`/`STOP` → `receive_relocation_command`; `CAPTURE_IMAGE` → `capture_image`; `RUN_MODEL` → `run_model`; `THERMAL_SCAN`/`START_SCAN` → `thermal_scan`.

Eligibility is authoritative in Core: matching tenant, active lifecycle, optional project/type/capabilities, no active exclusive assignment, online connection/current twin, minimum battery, and optional target distance. Ranking is domain score/proposal order where supplied, then distance, higher battery, and stable device ID. The frontend has no candidate-preview or rejection-diagnostic endpoint.

Creation timing fields, when returned, are microseconds: `candidate_selection_duration_us`, `routing_duration_us`, `persistence_duration_us`, and `total_duration_us`.

### Command document

A command contains `command_id`, `tenant_id`, `device_id`, `task_id`, `command_type`, immutable JSON `payload`, `status`, per-device `sequence_number`, `correlation_id`, `causation_id`, `attempt_count`, `max_attempts`, `version`, `created_at`, `available_at`, optional `sent_at`, `acknowledged_at`, `completed_at`, `expires_at`, optional `ack_status`, arbitrary JSON `result`, and optional `last_error`.

Command states are `PENDING`, `DELIVERED`, `ACKNOWLEDGED`, `COMPLETED`, `FAILED`, `EXPIRED`, and `CANCELLED`. Only the lowest outstanding per-device sequence is deliverable. `DELIVERED` records a durable delivery attempt, `ACKNOWLEDGED` requires an `ACCEPTED` or idempotent `DUPLICATE` ACK, and `COMPLETED` requires a successful result. Delivery reconciliation can return `DELIVERED` to `PENDING` with backoff while retaining command identity.

The command read model does not expose `schema_version` or delivery-observation timing from the transport envelope. It exposes only aggregate attempt count, latest sent time, next available time, and latest error. ACK reason is not stored separately; terminal ACK failure detail is represented through `ack_status`, `last_error`, and command/task failure fields. Result timestamp is `completed_at`.

### Browser update availability

The dashboard WebSocket subscribes only to normalized telemetry from Redis Pub/Sub. Task/command lifecycle outbox events go to Kafka topics but have no browser subscription or durable read-oriented activity endpoint. Operations pages must therefore hydrate through REST and poll only while mounted: active lists/details may refresh modestly; terminal details stop network polling. Activity can only be labelled session/current-state activity, not audit history.

Known orchestration errors include `INVALID_TASK`, `NO_ELIGIBLE_DEVICE`, `ROUTING_BUSY`, `ROUTING_TIMEOUT`, `ROUTING_UNAVAILABLE`, `NO_ROUTE`, `NO_ROAD_NODE`, `UNSUPPORTED_PROFILE`, `OUTSIDE_REGION`, `PLANNER_UNAVAILABLE`, `INVALID_STATE_TRANSITION`, `NOT_FOUND`, `FORBIDDEN`, and `ORCHESTRATION_ERROR`.

## Phase 3 Mobility contracts

All canonical Mobility calls are read-authorized and tenant-scoped by backend middleware. The browser must keep the three spatial mechanisms distinct: canonical Mobility nearby search, road routing, and the legacy in-memory compatibility matcher.

| Method | Browser path | Backend path | Request | Response | Product status |
| --- | --- | --- | --- | --- | --- |
| GET | `/api/engine/spatial/devices/nearby` | `/api/v1/spatial/devices/nearby` | Required `lat`, `lon`; `radius_meters` defaults to 5,000; `limit` defaults to 20 and is bounded by `MOBILITY_MAX_RAW_CANDIDATES` | `{count,devices}` where each item is `{state,distance_meters}` | Canonical nearby device discovery |
| POST | `/api/engine/routes` | `/api/v1/routes` | `{mobility_profile,origin,destination,policy}`; tenant ID in a payload is overwritten by authenticated scope | RouteResult | Canonical exploratory road routing |
| GET | `/api/engine/routes/calculate` | `/api/v1/routes/calculate` | `src_lat`, `src_lon`, `tgt_lat`, `tgt_lon`; always road vehicle + fastest | Legacy success document without the standard envelope | Compatibility only; not used by Phase 3 |
| GET | `/api/engine/nodes/match` | `/api/v1/nodes/match` | Legacy coordinate/radius matching parameters | Runtime in-memory Haversine result with compatibility ETA assumptions | Compatibility only; not the Mobility index |
| GET | `/api/engine/zones/predicted` | `/api/v1/zones/predicted` | None | Zone array: `id`, `lat`, `lon`, `radius_km`, `required_assets`, numeric `target_class`, `tenant_id` | Experimental legacy density heuristic |
| GET | `/api/engine/readyz` | `/readyz` | None | Engine/core readiness, module states/components/details and runtime | Read-only Mobility status and diagnostics; frontend proxy addition required |

### Canonical nearby response

`SpatialState` exposes tenant/device identity, `position` (the indexed position), `reported_position` (latest accepted report), optional heading/speed, exact `mobility_profile`, H3 cell, observation/index timestamps, boot/version identity, and quality `{valid,confidence,anomalies}`. Results are sorted by exact Haversine distance and then stable device ID. The endpoint does not return registry display name, project, battery, lifecycle, or connectivity; those require the existing tenant-scoped twin endpoint for a selected device. It supports no device/profile/project/capability filters.

Coordinates are decimal degrees. Radius is metres. Backend validation allows latitude `[-90,90]`, longitude `[-180,180]`, positive radius up to the deployment-configured maximum (10 km by default), and a limit up to the deployment-configured raw candidate maximum (50 by default). These maxima are not exposed dynamically, so the frontend uses the repository defaults and reports that configuration gap.

### Route request and response

Profiles are `ROAD_VEHICLE`, `GROUND_ROBOT`, `AERIAL_DRONE`, and `STATIC`; current global routing accepts only `ROAD_VEHICLE`. Policies are `SHORTEST` and `FASTEST`. A successful RouteResult contains `route_id`, `road_graph_version`, `snapshot_version`, `policy`, `distance_meters`, Go `time.Duration` JSON in nanoseconds as `estimated_time`, ordered `waypoints` using `{latitude,longitude,altitude_meters?}`, `edge_ids`, and `expanded_nodes`.

The road router uses KD-tree endpoint snapping, bounded workers/per-tenant concurrency, A*, and one immutable cost snapshot throughout a request. `SHORTEST` minimizes distance; `FASTEST` uses the snapshot's current travel costs. Current routing errors are `ROUTING_BUSY` (429), `ROUTING_TIMEOUT` (504), `ROUTING_UNAVAILABLE` (503), `NO_ROUTE`, `NO_ROAD_NODE`, `OUTSIDE_ROUTING_REGION`, and `UNSUPPORTED_MOBILITY_PROFILE` (422). Unknown failures become `ROUTING_ERROR` (500). Standard failures include a request ID.

### Readiness, routing runtime, and traffic metadata

`/readyz` exposes a `mobility` module when enabled. Its state is `STARTING`, `READY`, `DEGRADED`, `FAILED`, or `STOPPED`; components currently contain `spatial` and `routing`. Details, when a graph is loaded, include `road_graph_version`, road node/edge counts, `routing_snapshot_version`, `traffic_scope`, refresh interval, traffic edge-state count, overlay byte estimate, and `routing_runtime` with cumulative request/busy counts, queue depth/capacity, and active tenant limiter count. Worker count, timeout count, active request count, individual limiter occupancy, route latency, and route history are not exposed.

Traffic scope is currently only `SHARED_TRUSTED`. The browser receives snapshot metadata and aggregate internal state counts, but no individual edge speed, confidence, sample count, observation time, cost, or geometry. Therefore Phase 3 must not draw a congestion layer.

### Predicted zones and authority limitation

Predicted zones are a last-hour PostgreSQL telemetry density view grouped by coordinates rounded to two decimals, limited to three cells, with fixed 2 km radius and nominal required asset count. Despite legacy comments/log text, it is not ML, forecasting, or autonomous action. The query is not tenant-filtered and emitted zone tenant identity is currently hard-coded. The frontend filters returned rows to the active tenant, labels the view experimental, and never treats it as tenant-safe authoritative analytics.
