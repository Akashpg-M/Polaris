# Polaris v3 Backend — Feature and Logical Architecture Context

**Current implementation boundary:** Phase 0 through Phase 4.2, including deployment stabilization  
**Purpose of this document:** Give a human or AI agent enough functional and architectural context to understand the backend without reading its source code. This document explains behavior, rules, algorithms, state ownership, and component connections. It intentionally avoids source-level implementation details.

## 1. What Polaris is

Polaris is a multi-tenant control platform for connected physical and edge devices. Its strongest implemented use case is a geographically distributed mobility fleet containing connected road vehicles, delivery drones, ground robots, and fixed spatial devices. It also supports non-spatial compute devices for generic tasks.

The backend combines four related responsibilities:

1. **Authenticated telemetry ingestion:** Devices continuously report position, movement, battery, and operating state through persistent connections.
2. **Registry and digital-twin management:** Operators durably register tenants, projects, device identities, types, capabilities, credentials, and lifecycle state. Current reported state is combined with registry metadata to form a twin.
3. **Durable task and command orchestration:** Operators submit desired work. Polaris selects an eligible device, creates an immutable command, delivers it to the correct live connection, and tracks acknowledgement and completion.
4. **Mobility intelligence:** The optional Mobility module maintains a geographic index, discovers nearby devices, loads a real road graph, calculates bounded congestion-aware routes, and embeds versioned routes into commands.

Polaris uses an **at-least-once delivery model with idempotent processing**. It does not claim exactly-once processing. Replays are expected and are made harmless through stable identities, uniqueness rules, version comparisons, and guarded state transitions.

## 2. System shape and authority model

The runtime consists of two Polaris processes plus three stateful infrastructure systems:

| Component | Logical responsibility |
| --- | --- |
| Gateway | Authenticates devices and dashboards, validates telemetry, owns live device sockets, and performs the final fenced command write to a device. |
| Engine | Consumes telemetry, maintains projections, exposes control APIs, runs registry/twin/task logic, archives history, relays outbox events, dispatches commands, and hosts Mobility. |
| Kafka-compatible event fabric | Durable ordered transport for telemetry, commands, lifecycle observations, and dead-letter records. |
| PostgreSQL/PostGIS | Durable authority for tenants, devices, credentials, audit, telemetry history, tasks, assignments, commands, and transactional outbox state. |
| Redis | Atomic latest-state projection, connectivity timestamps, live connection ownership leases, gateway command notification, and transient dashboard publication. |

The most important ownership rules are:

- PostgreSQL is the authority for identity, lifecycle, capabilities, audit, tasks, assignments, command state, and historical telemetry.
- Kafka is the durable transport boundary. Ordering is guaranteed only within a partition, so device-specific streams are keyed by tenant and device identity.
- Redis is the authority for the current reported-state projection and current gateway lease, but Redis Pub/Sub is never considered durable.
- Mobility’s H3/R-tree state and traffic overlay are in-memory derivatives. They are rebuilt from authoritative/replayable state and never reserve a device or create a command by themselves.
- The Gateway is the only component allowed to perform a physical WebSocket command write.
- The core orchestration service always revalidates a module’s candidate proposal before PostgreSQL commits an assignment.

## 3. Principal end-to-end flows

### 3.1 Telemetry flow

```text
Registered device
  → authenticated Gateway WebSocket
  → binary telemetry validation
  → canonical versioned telemetry event
  → Kafka telemetry.ingress, keyed by tenant:device
  → independent consumers
      → atomic Redis latest-state/twin projection
      → compatibility in-memory state
      → Mobility H3/R-tree projection
      → PostgreSQL historical archive
      → traffic map matching and routing-cost refresh
  → Redis dashboard publication
  → tenant-filtered Gateway dashboard socket
  → browser dashboard
```

The consumers are deliberately independent. A telemetry record can be replayed after one side effect has already succeeded. Each persistent or current-state destination therefore handles duplicates without corrupting state.

### 3.2 Task and command flow

```text
Authenticated operator
  → create tenant-scoped task
  → durable PENDING task and outbox event
  → registry capability/lifecycle eligibility
  → Redis connectivity/battery/location eligibility
  → optional Mobility candidate ranking and route planning
  → authoritative eligibility recheck
  → one PostgreSQL transaction:
       reserve device
       allocate per-device command sequence
       persist immutable command
       mark task assigned
       write audit and outbox events
  → outbox relay publishes durable command to Kafka
  → dispatcher finds current Gateway ownership lease in Redis
  → gateway-specific Redis notification
  → Gateway revalidates device and fencing epoch
  → durable delivery attempt recorded
  → command written to device WebSocket
  → device ACK and result return on the authenticated socket
  → command/task state transitions and assignment release
```

If the device is offline, the durable command remains in PostgreSQL. Redis notification may be absent, but reconnection causes the Gateway to query and deliver pending commands in sequence order.

### 3.3 Registry-to-twin flow

```text
PostgreSQL registry metadata and capabilities
  + Redis reported components and connectivity
  → tenant-scoped digital-twin response
```

A registered device can have a twin even if it has never sent telemetry. In that case, metadata exists, reported state is absent, and connectivity is `NEVER_CONNECTED`.

## 4. Device telemetry contract

### 4.1 Device-owned telemetry frame

The device sends a compact binary Protobuf frame containing:

- device and tenant identifiers;
- device type and operational status;
- latitude and longitude;
- velocity and heading;
- battery percentage;
- a device boot identifier;
- a monotonically increasing sequence number within that boot;
- boot-start and observation timestamps;
- telemetry schema version.

Supported spatial device types are bike, auto, sedan, SUV, drone, robot, and static sensor. Registry device profiles constrain which telemetry types a credential may send. For example, a registered drone cannot impersonate a road vehicle, and a non-spatial compute node cannot publish the spatial telemetry schema.

### 4.2 Canonical platform envelope

The Gateway wraps accepted device data in a canonical version-1 platform event. The envelope carries:

- deterministic event identity;
- event type and schema version;
- authenticated tenant/device identity;
- boot and sequence ordering identity;
- observed and ingested timestamps;
- correlation and causation identifiers;
- producer identity and optional distributed-trace context;
- the original validated spatial payload.

The event ID is deterministically derived from tenant, device, boot, and sequence. The same logical device event therefore has the same identity when replayed.

The Kafka partition key is `tenant_id:device_id`. This keeps successive events for one device in one Kafka partition even when the device moves between H3 cells. H3 is used only after ingestion for geographic organization.

## 5. Gateway features and logic

### 5.1 Device authentication

A device connects using either:

- a bearer device credential; or
- a short-lived, single-use connection ticket obtained through an authorized operator flow.

Authentication resolves a trusted device principal containing tenant, device, credential, registered type, and project identity. A connection is accepted only when the credential is active and unexpired, the device is active, and the tenant is active.

The Gateway never trusts the tenant or device identity in telemetry as authorization. The payload identity must match the authenticated principal, after which the principal’s values are used as the platform identity.

### 5.2 Session revalidation and revocation response

The Gateway revalidates the credential, device, and tenant for every received frame. A revoked credential, suspended/decommissioned device, or suspended/deactivated tenant closes the existing session before more telemetry or command responses are accepted.

A database transport error is distinguished from a true revocation. An already authenticated session receives a bounded 30-second grace window during a transient registry outage. If the registry cannot confirm validity for the entire grace period, the session fails closed.

### 5.3 Telemetry validation boundary

Before Kafka publication, the Gateway rejects:

- malformed or non-binary telemetry frames;
- frames larger than 64 KiB;
- invalid tenant, device, or boot identifiers;
- zero or unsupported sequence numbers;
- unsupported schema versions;
- latitude or longitude outside valid Earth ranges;
- battery outside 0–100%;
- negative, non-finite, or implausibly large velocity values;
- unknown device types or a type inconsistent with the registered profile;
- missing/invalid boot and observation timestamps;
- observations older than 24 hours or over five minutes into the future;
- identity or boot-ID changes within one live connection.

Kafka publication is synchronous at this boundary. If publication fails, the Gateway closes the connection with a retryable reason so the device can reconnect and replay its unaccepted event.

### 5.4 Bidirectional frame handling

The same device WebSocket carries two frame families:

- binary Protobuf telemetry from device to platform;
- explicit JSON control frames for platform commands, device acknowledgements, and device results.

An acknowledgement can report `ACCEPTED`, `REJECTED`, `DUPLICATE`, `EXPIRED`, or `UNSUPPORTED`. A result can report completion/success or failure. Command and sequence identity are always checked against the authenticated tenant/device before durable state changes.

### 5.5 Connection ownership and fencing

Every accepted connection claims a Redis lease identified by tenant and device. The lease contains Gateway ID, connection ID, credential ID, expiry, and a monotonically increasing ownership epoch.

When a new connection supersedes an older connection:

- the epoch increases;
- the older local socket is closed;
- only the new epoch can refresh or release ownership;
- command delivery verifies the epoch before the durable delivery transition and again immediately before writing to the socket.

The double check prevents a stale Gateway/session from delivering after a reconnection race. Lease refresh tolerates an isolated Redis timeout but closes the socket if ownership cannot be confirmed before the lease window is exhausted.

### 5.6 Command delivery concurrency

Gateway command notifications are distributed across 16 bounded workers using a stable hash of tenant and device. All commands for one device go to the same worker, preserving device order, while different devices can be delivered concurrently.

Before each write the Gateway:

1. finds the active local session;
2. revalidates device authorization;
3. validates current Redis ownership;
4. asks PostgreSQL to transition the specific command to delivered and record an attempt with the ownership epoch;
5. rechecks ownership;
6. writes the canonical command envelope.

If the volatile socket write fails after the durable delivery transition, the reconciliation worker later returns the command to pending after the ACK timeout. This is an expected at-least-once scenario.

### 5.7 Dashboard delivery

Dashboards authenticate using an operator bearer credential or a short-lived single-use dashboard ticket. Each dashboard socket is bound to one tenant. The Gateway subscribes to the Redis `spatial:updates` channel and forwards only events whose tenant matches the socket’s tenant.

Dashboard Pub/Sub is a live projection, not an event archive. Clients obtain historical/durable information through APIs and PostgreSQL-backed views, not by expecting Redis Pub/Sub replay.

## 6. Kafka telemetry processing and reliability

### 6.1 Partition-aware state consumer

The state consumer uses manual offset commits and tracks batches separately for every Kafka partition.

- Messages within a partition are processed in offset order.
- A partition flushes after 1,000 messages or 150 ms, whichever occurs first.
- The 150 ms internal target leaves margin below the externally tested 250 ms partial-batch requirement.
- Only the highest contiguous successfully processed offset is committed.
- Processing stops at the first event that cannot reach a terminal result; no later offset is committed past it.
- A Kafka commit failure retains successful work for replay.
- Shutdown forces all pending partition batches to flush and waits for consumers to close.

Different partitions progress independently, preventing a failure in one partition from incorrectly advancing or globally blocking another partition.

### 6.2 Processing order

For valid telemetry, the state path applies the atomic Redis freshness decision before changing in-memory spatial views.

This ordering means:

- if Redis fails, volatile state is not advanced;
- if Redis succeeds but Kafka commit fails, replay is classified as a duplicate;
- a duplicate replay may still rebuild an empty in-memory view after an Engine restart;
- out-of-order, retired-boot, and conflicting-boot events never enter active in-memory state.

### 6.3 Failure classification and DLQ

Malformed envelopes, unsupported schemas, invalid identities, invalid coordinates, and other permanently invalid messages are copied to `telemetry.dead-letter.v1`. The dead-letter record preserves the original bytes, key, source topic, partition, offset, reason, and failure time.

Redis and PostgreSQL failures are treated as transient and retried up to five times with bounded incremental delay. When retries are exhausted, the original event is sent to the DLQ. The source offset advances only after DLQ publication succeeds. Separate consumer groups can independently emit dead-letter evidence for the same source record.

## 7. Latest-state and digital-twin logic

### 7.1 Atomic freshness classification

Redis stores current state under `polaris:twin:{tenant}:{device}`. One atomic operation compares incoming boot/sequence identity, writes the new state, updates last-seen time, updates twin components, and publishes the dashboard event.

Possible results are:

| Classification | Meaning |
| --- | --- |
| `ACCEPTED` | Newer sequence in the current boot, or the first known event. |
| `DUPLICATE` | Same boot and same sequence as current state. |
| `OUT_OF_ORDER` | Same boot but an older sequence. |
| `NEW_BOOT` | A different boot with a later boot-start time; the previous boot is retired. |
| `RETIRED_BOOT` | A delayed event from a known old boot or an older boot-start time. |
| `BOOT_CONFLICT` | Different boot IDs claim the same boot-start time. |

Sequence numbers are compared as decimal strings so the complete unsigned range is handled without Lua numeric precision loss.

Within a boot, one sequence number represents one logical observation. If a device sends different payload content with the same boot and sequence, the later frame is still a duplicate rather than a correction. A correction must use a newer sequence. Because the Gateway keeps one boot ID for a connection, a real device reboot is expected to establish a new connection.

### 7.2 Boot retirement

Each twin keeps a set of retired boot IDs. A valid newer boot can reset the device sequence, but the replaced boot is permanently retired. Delayed events from that boot cannot later regain authority, even if their sequence is high.

### 7.3 Twin components

Accepted telemetry atomically updates:

- a compatibility reported-state document used by dashboard and orchestration;
- `spatial/v1`, containing coordinates, heading, speed, and mobility profile;
- `battery/v1`, containing battery percentage;
- event, boot, sequence, last-seen, and connectivity metadata.

Each component has its own type, schema version, observation time, boot ID, sequence number, and opaque payload. This allows future modules to add independently versioned twin components without changing the core twin contract.

### 7.4 Connectivity lifecycle

Connectivity is separate from registry lifecycle:

- `NEVER_CONNECTED`: registered but no reported state exists;
- `ONLINE`: fresh accepted telemetry exists, or a valid current Gateway lease establishes live connectivity for orchestration;
- `STALE`: last telemetry exceeded the configured stale threshold;
- `OFFLINE`: last telemetry exceeded the configured offline threshold.

A Redis sorted set indexes devices by last-seen time. A periodic detector scans bounded pages of expired entries. Atomic guarded transitions ensure that a fresh event racing with a stale/offline scan cannot be overwritten by the older observation. Connectivity changes produce deterministic Kafka events.

STALE/OFFLINE devices are evicted from Mobility discovery. Fresh accepted telemetry changes the twin back to ONLINE and reintroduces a valid device to the derived index. A socket reconnection alone does not reinsert stale coordinates.

## 8. PostgreSQL telemetry archive

The archive consumer reads telemetry through its own Kafka group, independently of the latest-state consumer.

Each row stores:

- event, tenant, device, boot, and sequence identity;
- device type and operational status;
- raw coordinates and a PostGIS point;
- velocity, heading, and battery;
- observed, ingested, recorded, schema, and correlation metadata.

Two uniqueness rules make replay harmless:

- one row per event ID;
- one row per tenant/device/boot/sequence tuple.

A duplicate insert is treated as successful terminal processing. The consumer commits only after a successful/idempotent insert or successful DLQ publication. If the insert succeeds but offset commit fails, replay finds the same unique identity and creates no duplicate history.

The table has indexes for device/time history, recent density queries, and PostGIS spatial access.

## 9. Durable registry and tenancy

### 9.1 Registry entities

PostgreSQL durably stores:

- tenants;
- projects owned by tenants;
- a global catalog of device types;
- tenant-owned devices, optionally associated with projects;
- a global capability catalog;
- tenant/device capability assignments with per-device configuration;
- device credentials;
- operator API keys;
- one-time device and dashboard connection tickets;
- audit events and transactional outbox events.

Built-in device types include delivery drone, ground robot, connected vehicle, fixed IoT sensor, static camera, and non-spatial compute node.

Built-in capabilities include navigation, relocation command receipt, image capture, payload carriage, edge model execution, and temperature measurement. The command model also recognizes thermal-scan capabilities when such a capability is registered.

### 9.2 Tenant, project, and device lifecycle

Tenant states are `ACTIVE`, `SUSPENDED`, and `DEACTIVATED`.

Project state supports active operation and archival.

Device lifecycle is:

```text
REGISTERED → ACTIVE ↔ SUSPENDED
     └────────┴────────→ DECOMMISSIONED
```

Decommissioning is terminal. Devices are not physically deleted, preserving identity, audit, telemetry, task, and command history. Metadata can be updated only while the device is not decommissioned.

Suspending/decommissioning a device or disabling its tenant prevents authentication and removes active Mobility state. Reactivation allows new connections and fresh telemetry to restore runtime state.

### 9.3 Credential model

Tokens contain a public lookup prefix and at least 256 bits of random secret material. PostgreSQL stores only the prefix and a SHA-256 hash of the full bearer token. Verification uses constant-time comparison. The raw secret is returned only when issued or rotated and cannot be recovered from metadata listings.

Credential operations support:

- optional expiration;
- last-used tracking;
- revocation;
- atomic rotation, which revokes the old credential and creates the new credential in one transaction;
- immediate effect on active sessions through per-frame revalidation.

### 9.4 Operator roles and permissions

| Role | Read | Registry mutation | Create/cancel tasks | Administrative retry | Audit |
| --- | ---: | ---: | ---: | ---: | ---: |
| Platform administrator | Yes | Yes | Yes | Yes | Yes |
| Tenant administrator | Yes | Yes | Yes | Yes | Yes |
| Operator | Yes | No | Yes | No | No |
| Viewer | Yes | No | No | No | No |

Only a platform administrator can create tenants. A platform administrator explicitly selects tenant scope using the tenant header or query parameter. All other roles derive tenant scope from their stored identity.

Repository reads and mutations include tenant identity. Cross-tenant IDs are returned as not found rather than revealing resource existence.

### 9.5 Audit and transactional outbox

Security-sensitive registry and orchestration changes write three things in one PostgreSQL transaction:

1. the domain state change;
2. an immutable audit record containing actor, action, resource, request ID, and outcome;
3. a stable outbox event describing the change.

The relay claims bounded batches using skip-locked row ownership, prioritizes actionable command events, marks claimed records for retry, publishes them to Kafka, and marks the batch published afterward. If Kafka accepts an event but the published marker fails, the event is replayed. Downstream consumers must therefore remain idempotent.

Failed outbox attempts receive a bounded retry delay and are retained with error evidence; repeatedly failing rows eventually enter a failed state for operator inspection.

## 10. Digital-twin API behavior

A device twin merges:

- durable registry identity and display metadata;
- project and device type;
- firmware/software/model metadata;
- assigned capabilities and configuration;
- Redis connectivity;
- current reported state;
- independently versioned component envelopes.

Twin lists first select tenant-scoped devices from PostgreSQL and then hydrate their Redis state. They do not enumerate Redis keys as a source of ownership. Device/twin lists are cursor-based and bounded to at most 100 items per request.

## 11. Task orchestration

### 11.1 Task model

A task represents operator intent rather than a transport message. It includes:

- tenant/project identity;
- task type and priority (`LOW`, `NORMAL`, `HIGH`, `CRITICAL`);
- capability, battery, type, project, distance, and custom requirements;
- a target document, optionally containing coordinates;
- planning mode;
- correlation/creator identity;
- expiration and lifecycle timestamps;
- optimistic version and failure reason.

Task states are `PENDING`, `ASSIGNING`, `ASSIGNED`, `IN_PROGRESS`, `COMPLETED`, `FAILED`, `CANCELLED`, and `EXPIRED`. Terminal tasks cannot return to active states except through an explicit administrative retry decision.

### 11.2 Command-to-capability mapping

Polaris adds the necessary capability to task requirements based on command type:

- `NAVIGATE` and `RETURN_HOME` require navigation;
- `RELOCATE` and `STOP` require relocation-command support;
- `CAPTURE_IMAGE` requires image capture;
- `RUN_MODEL` requires edge model execution;
- scan commands require thermal-scan capability.

Unknown command types are rejected. A caller cannot bypass a required capability by omitting it from the task document.

### 11.3 Eligibility algorithm

A candidate is core-eligible only when all of these are true:

- it belongs to the same tenant;
- its registry lifecycle is ACTIVE;
- it matches any requested project and device types;
- every required capability is enabled;
- it has no active exclusive assignment;
- its current twin or valid live connection indicates ONLINE;
- reported battery meets the minimum;
- location exists when spatial constraints require it;
- direct distance is within the requested maximum.

The database first produces the lifecycle/capability/assignment candidate set. An optional domain provider may rank that set. Redis twin and connection state are fetched in a pipeline for a bounded maximum of 50 proposals. The core then rechecks authoritative database eligibility once immediately before attempting assignment.

Without domain ranking, candidates are ordered by nearest direct distance, then highest battery, then stable device ID. With Mobility ranking, routed candidates rank by ETA; non-routed candidates fall back to direct distance. Domain providers may reorder eligible devices but cannot silently remove all core-eligible fallback choices.

### 11.4 Race-safe exclusive assignment

Candidate selection is advisory until PostgreSQL commits. The final transaction:

- locks the pending task;
- creates an active device assignment;
- allocates the next monotonic sequence for that tenant/device;
- inserts the immutable pending command;
- marks the task assigned;
- writes audit and outbox events.

A partial unique rule permits only one active assignment for a tenant/device. Concurrent requests that selected the same device race at the database, and only one reservation succeeds.

### 11.5 Planning modes

For navigation and relocation tasks:

- `DEVICE_LOCAL` means Polaris may send the high-level target and let the device perform local planning.
- `POLARIS_REQUIRED` means a compatible platform module must produce a valid route plan. Missing routing, overload, timeout, no route, or unsupported profile is returned explicitly; Polaris does not fabricate a generic route.

Other task types use the generic planner, which preserves the validated target as the command payload.

### 11.6 Timing visibility

Task creation reports separate candidate-selection, routing/planning, persistence, and total durations. These timings are diagnostic observations and are not part of durable command identity.

## 12. Durable command orchestration

### 12.1 Command identity and ordering

A command contains stable command/task/tenant/device identity, command type, immutable payload, schema version, per-device sequence, correlation/causation identity, creation and expiry times, attempts, and lifecycle state.

Sequence allocation is atomic and monotonically increasing for each tenant/device. Only the lowest outstanding device sequence can be delivered. A later command cannot overtake a pending, delivered, or acknowledged earlier command.

The Kafka command key is also `tenant_id:device_id`, preserving device order through the event fabric.

### 12.2 Command lifecycle

Normal success follows:

```text
PENDING → DELIVERED → ACKNOWLEDGED → COMPLETED
```

Failure branches can lead to `FAILED`, `EXPIRED`, or `CANCELLED` depending on the current stage and observation.

- `DELIVERED` means a durable delivery attempt was recorded before the socket write.
- An accepted or duplicate ACK advances the command to acknowledged and the task to in progress.
- A rejected/unsupported ACK fails the command and task.
- An expired ACK expires them.
- A successful result completes both and releases the device assignment.
- A failed result fails both and releases the assignment.
- Duplicate ACKs/results against the same or terminal state are harmless.
- A wrong device or wrong sequence cannot advance the command.

### 12.3 Retry and reconciliation

A periodic reconciler:

- returns delivered commands with no timely ACK to pending;
- applies exponential backoff capped at 30 seconds;
- preserves command ID, sequence, and payload across delivery retries;
- terminates commands after configured attempts or expiry;
- expires unassigned tasks past their deadline;
- releases assignments belonging to terminal work;
- removes expired connection leases;
- reconsiders a bounded priority-ordered set of pending tasks.

Pending tasks are considered in priority order, then creation order. Offline work remains durable and can be assigned/delivered when a suitable device becomes available.

Administrative task retry represents a new assignment/command decision. Delivery retry represents another attempt of the same immutable command. A forced retry of a failed command can reactivate its assignment only if doing so does not conflict with another active reservation.

Cancellation is intentionally conservative. Pending work can be cancelled and its assignment released. Once a command has been physically delivered or acknowledged, the API does not pretend that local cancellation stopped the real device; such requests conflict unless a separate domain stop task is issued.

## 13. Capability-module framework

Polaris supports explicitly composed capability modules. This is a compile-time extension model, not runtime plugin discovery.

A module can provide:

- lifecycle and readiness state;
- candidate ranking;
- task planning;
- new typed twin components.

Module lifecycle states are `STARTING`, `READY`, `DEGRADED`, `FAILED`, and `STOPPED`, with component-level status and operational details.

Planner fallthrough is narrow: a specialized planner may explicitly decline an unsupported device profile, allowing another compatible planner to run. Real failures such as overload, timeout, missing graph, or no route do not silently fall through.

Mobility can be disabled while registry, telemetry, twins, generic tasks, command persistence, and device delivery remain functional. If Mobility is optional and the road graph cannot load, it reports degraded routing rather than taking down the entire core. If configured as mandatory, startup fails.

## 14. Mobility spatial subsystem

### 14.1 Profiles and stored state

Mobility recognizes four profiles:

- road vehicle;
- ground robot;
- aerial drone;
- static spatial device.

Each active state preserves both:

- **reported position**, which advances with every accepted observation; and
- **indexed position**, which changes only when spatial reindexing is useful.

It also stores H3 cell, heading, speed, observation/index time, boot/sequence source version, and quality/anomaly information.

### 14.2 Observation quality rules

Coordinates and speed must be finite and valid. Heading is normalized to the range 0–360 degrees.

Movement between observations produces an implied speed. Profile-specific plausibility limits are used:

- road vehicle: 90 m/s;
- ground robot: 20 m/s;
- aerial drone: 120 m/s;
- static device: 2 m/s.

An implausible jump is retained as a lower-confidence observation and tagged, rather than silently presented as high-confidence motion. Movement under 0.15 m/s is tagged as stationary deadband. Invalid coordinates or speed are rejected.

### 14.3 H3 sharding

The default reported position is mapped to H3 resolution 8. Its resolution-6 parent identifies an internal regional shard. Shards are also separated by tenant.

Updates for the same device use one of 64 stable device locks. When a device crosses a regional boundary, the old and new shards are locked in numeric order, old membership is removed, new membership is inserted, and the global location record is changed within the same protected operation. Deterministic lock ordering prevents deadlock, and the move cannot expose the device as active in two regions.

Default capacity is 10,000 active spatial devices per tenant.

### 14.4 Movement threshold

The R-tree is updated when at least one of these is true:

- the H3 cell changed;
- distance from indexed position reached five meters;
- the indexed position is older than 30 seconds.

Smaller movements still advance reported state and twin freshness. This reduces index churn without lying about the latest reported location.

### 14.5 Packed STR R-tree

Each regional shard uses a packed Sort-Tile-Recursive-style R-tree with fanout 16.

- Mutations update the shard’s authoritative item map and mark the hierarchy dirty.
- The next query lazily rebuilds a packed hierarchy.
- Entries are spatially ordered into bounded leaves and recursively packed into parent bounding boxes.
- A query first prunes nodes using latitude/longitude bounding rectangles.
- Antimeridian-crossing searches are split into two boxes.
- Candidates are then checked with exact Haversine distance.
- Results are deduplicated, sorted by distance and stable device ID, and bounded by the requested limit.

A linear exact index exists only as a test/benchmark oracle. It is not a competing production authority.

### 14.6 Nearby search

Nearby search maps the target to H3, expands a bounded grid disk, converts cells to relevant parent shards, and queries only those R-trees. The ring estimate grows with radius and is capped by configuration. Search radius and result limits are bounded to prevent unbounded tenant queries.

### 14.7 Eviction and restart rebuild

Mobility removes state when a device becomes stale/offline, a device is suspended/decommissioned, or a tenant is no longer active.

On Engine restart, Mobility scans Redis twin keys in bounded pages and accepts only entries whose:

- connectivity is ONLINE;
- device lifecycle is ACTIVE;
- tenant is ACTIVE;
- spatial component is valid.

Inputs are sorted by tenant and device before deterministic rebuild. R-tree internals are not persisted. This keeps Redis/PostgreSQL authoritative and makes the index recoverable.

## 15. Road routing subsystem

### 15.1 Road graph creation

The bundled Chennai OpenStreetMap PBF is converted into an immutable directed graph.

Supported road families include motorway/trunk, primary, secondary, tertiary, residential/living street, unclassified, and service roads. Only nodes referenced by accepted roads enter the graph.

For each road segment, Polaris:

- calculates geodesic distance;
- uses a valid tagged maximum speed or a road-class default;
- derives base travel time;
- handles normal bidirectional roads;
- handles forward one-way, reverse one-way, and roundabout direction;
- builds outgoing and incident edge lists;
- classifies multi-edge nodes as intersections.

The current Chennai graph evidence contains approximately 690,268 road nodes and 1,442,876 directed edges under graph version `chennai-v1`.

### 15.2 GPS-to-road snapping

Road nodes are projected onto a three-dimensional unit sphere and stored in a balanced KD-tree. Route endpoints query the nearest sphere point rather than scanning every road node. The 3D representation naturally handles longitude wraparound better than a raw two-dimensional longitude index.

### 15.3 Route policies and A*

Polaris supports:

- shortest-distance routing;
- fastest-time routing.

A* is the operational search algorithm. It uses a priority queue, accumulated path cost, predecessor edges, visited-node protection, and an admissible geodesic heuristic:

- straight-line distance for shortest routes;
- straight-line distance divided by graph maximum speed for fastest routes.

Dijkstra uses the same graph and cost snapshot without the heuristic and serves as a correctness/benchmark oracle. Tests compare A* and Dijkstra route cost.

Results contain route ID, graph version, traffic snapshot version, policy, distance, estimated duration, ordered waypoints, directed edge IDs, and expanded-node count.

### 15.4 Immutable traffic-cost snapshots

The road topology and base costs never mutate after startup. Dynamic travel times are stored in a complete versioned cost array.

Refreshing traffic builds a new array, validates its size and monotonically increasing version, copies it, and atomically swaps the current snapshot. Every route request loads one snapshot pointer and uses that version for its entire search. Concurrent refresh cannot produce a route assembled from mixed cost versions.

### 15.5 Bounded routing execution

Routing is guarded by:

- a fixed worker pool;
- a bounded global queue;
- a per-tenant concurrency semaphore;
- a request timeout;
- a maximum expanded-node count;
- caller cancellation.

Defaults are four workers, queue capacity 64, two concurrent routes per tenant, a two-second timeout, and 250,000 expansions. Global or tenant saturation returns `ROUTING_BUSY`; timeout returns a distinct timeout response. This protects telemetry and generic task processing from route overload.

Road routing currently accepts only the road-vehicle profile.

## 16. Live traffic and congestion logic

A dedicated Kafka consumer independently observes valid road-vehicle telemetry.

### 16.1 Map matching

For each road observation, Polaris:

1. finds the nearest road node using the KD-tree;
2. evaluates only edges incident to that node;
3. computes point-to-segment distance;
4. adds a heading mismatch penalty when heading exists;
5. selects the lowest-scoring directed edge;
6. rejects matches beyond the confidence threshold.

One observation updates at most one directed edge, preventing a single GPS point from congesting every outgoing road.

### 16.2 Speed smoothing and confidence

Matched speed uses an exponentially weighted moving average with new-sample weight 0.3. Each edge tracks sample count, smoothed and latest speed, last observation time, and map-match confidence.

### 16.3 Cost generation and decay

At each refresh:

- edge confidence decays exponentially with observation age;
- smoothed speed is compared with base road speed;
- the resulting multiplier can slow travel time but never make it faster than the base graph;
- no-observation edges use base cost;
- very old edge observations are removed;
- a complete new immutable routing snapshot is published.

Traffic scope is explicitly `SHARED_TRUSTED`: accepted road telemetry contributes to one shared road-cost model rather than separate tenant-specific congestion overlays.

## 17. Mobility-aware selection and route planning

### 17.1 Candidate ranking

For a task with target coordinates, Mobility:

1. performs H3/R-tree geographic narrowing;
2. keeps a bounded raw candidate set;
3. removes proposals outside the core eligible device list;
4. calculates fastest road ETA for only the top configured road candidates;
5. ranks routed candidates by ETA and remaining candidates by direct distance;
6. uses stable device identity to break ties.

Defaults permit 50 raw candidates and route at most eight. Battery ranking and all final lifecycle/capability checks remain in the core.

### 17.2 Route command planning

For a `POLARIS_REQUIRED` road `NAVIGATE` or `RELOCATE` task, the planner uses the selected device’s reported origin and target destination. It produces a versioned payload containing:

- route and route-schema identity;
- road graph version;
- routing snapshot version;
- generation and validity times;
- origin and destination;
- ordered waypoints;
- distance and duration;
- route policy.

Plan validity is limited to two minutes and never exceeds task expiry. Command expiry is the earlier of plan validity and task expiry.

The plan is immutable after command persistence. Retrying delivery reuses the same route. Replanning requires a new durable decision rather than silently mutating an existing command.

Ground robots and drones do not currently have platform global route planners. If planning is device-local, they can receive high-level targets. If a platform route is required for an unsupported profile, planning explicitly declines/fails.

## 18. Compatibility and analytical views

### 18.1 Compatibility node matching

The original `/nodes/match` endpoint remains for Phase 0 compatibility. It reads a 32-way hash-sharded in-memory map, performs a tenant/type-filtered Haversine scan, applies a radius, estimates ETA using a fixed 40 km/h speed, sorts by ETA, and caps results at 500.

This endpoint is a compatibility projection, not the advanced Mobility index. New geographic discovery should use `/spatial/devices/nearby`.

The compatibility view is runtime-only: it is not rebuilt from Redis at Engine startup and is not currently evicted by the connectivity detector. It can therefore be empty immediately after restart or retain stale entries until newer telemetry replaces them. Those limitations do not apply to Mobility nearby discovery, which has explicit rebuild and eviction rules.

### 18.2 Predicted zones

The current predicted-zone view is a lightweight historical-density heuristic:

- reads telemetry from the last hour;
- rounds latitude/longitude to two decimals, producing roughly kilometre-scale grid cells;
- counts observations per cell;
- returns the top three cells as two-kilometre zones with a nominal required asset count.

It is not machine learning, forecasting, or an active autonomous rebalancer. The current view also retains a fixed `alpha_logistics`/drone assumption and its density query is not tenant-filtered. It should therefore be treated as a legacy demonstration endpoint, not a fully tenant-authoritative production feature.

No autonomous fleet rebalancer is active in the current production wiring.

## 19. External API and socket surface

All control APIs are versioned under `/api/v1` and, except process probes and Gateway connection metrics, require operator authentication.

### 19.1 Registry and identity

- create/read/update tenants;
- create/list/read/update projects;
- register/list/read/update devices;
- activate, suspend, and decommission devices;
- list the capability catalog;
- assign/configure and remove device capabilities;
- issue, list, revoke, and rotate device credentials;
- issue device connection tickets;
- issue dashboard tickets;
- read one twin or list twins;
- inspect tenant audit events.

### 19.2 Orchestration

- create, list, and read tasks;
- cancel or administratively retry tasks;
- list and read commands;
- cancel or administratively retry commands;
- inspect orchestration counters.

Task/command list operations support bounded cursor and state/device/task filtering.

### 19.3 Spatial and routing

- compatibility nearest-node match;
- Mobility nearby-device query;
- canonical route creation;
- legacy route-calculation query retained for dashboard compatibility;
- predicted-zone view.

### 19.4 WebSockets

- `/ws/telemetry`: authenticated bidirectional device channel;
- `/ws/dashboard`: authenticated, tenant-scoped live dashboard stream.

### 19.5 Exact control API inventory

| Area | Endpoints |
| --- | --- |
| Tenants | `POST /api/v1/tenants`, `GET /api/v1/tenants/:tenant_id`, `PATCH /api/v1/tenants/:tenant_id` |
| Projects | `POST /api/v1/projects`, `GET /api/v1/projects`, `GET /api/v1/projects/:project_id`, `PATCH /api/v1/projects/:project_id` |
| Devices | `POST /api/v1/devices`, `GET /api/v1/devices`, `GET /api/v1/devices/:device_id`, `PATCH /api/v1/devices/:device_id` |
| Device lifecycle | `POST /api/v1/devices/:device_id/activate`, `/suspend`, `/decommission` |
| Capabilities | `GET /api/v1/capabilities`, `GET /api/v1/devices/:device_id/capabilities`, `PUT` or `DELETE /api/v1/devices/:device_id/capabilities/:capability_id` |
| Credentials | `POST`/`GET /api/v1/devices/:device_id/credentials`, credential revoke, and credential rotate endpoints |
| Connection tickets | `POST /api/v1/devices/:device_id/connection-ticket`, `POST /api/v1/dashboard-ticket` |
| Twins and audit | `GET /api/v1/devices/:device_id/twin`, `GET /api/v1/twins`, `GET /api/v1/audit-events` |
| Tasks | `POST`/`GET /api/v1/tasks`, `GET /api/v1/tasks/:task_id`, task cancel and administrative retry endpoints |
| Commands | `GET /api/v1/commands`, `GET /api/v1/commands/:command_id`, command cancel and administrative retry endpoints |
| Spatial | `GET /api/v1/nodes/match`, `GET /api/v1/spatial/devices/nearby` |
| Routing/analytics | `POST /api/v1/routes`, `GET /api/v1/routes/calculate`, `GET /api/v1/zones/predicted` |
| Metrics | Engine and Gateway orchestration metrics; Gateway active uplink count |

Control responses carry a stable request ID. A caller-provided `X-Request-ID` is preserved; otherwise Polaris creates one and reuses it for the response and associated audit work. Successful control responses use a data envelope, while failures expose a stable error code/message envelope.

## 20. Topics, consumer groups, and Redis channels

### 20.1 Kafka topics

| Topic | Purpose |
| --- | --- |
| `telemetry.ingress` | Canonical ordered telemetry events. |
| `telemetry.dead-letter.v1` | Original telemetry bytes plus consumer failure metadata. |
| `device.lifecycle.v1` | Registry and credential lifecycle events. |
| `device.connectivity.v1` | ONLINE/STALE/OFFLINE observations. |
| `task.lifecycle.v1` | Task lifecycle events. |
| `device.command.v1` | Durable commands and delivery retries, ordered by device. |
| `device.command.ack.v1` | Durable acknowledgement observations. |
| `device.command.result.v1` | Durable result observations. |

The deployment creates three partitions for each topic.

Principal consumer groups are:

- `polaris_engine_group` for latest and in-memory state;
- `polaris_archive_group` for PostgreSQL history;
- `polaris_traffic_group` for traffic observations;
- `polaris-command-dispatcher` for routing durable commands toward live Gateways.

### 20.2 Redis structures

| Pattern | Meaning |
| --- | --- |
| `polaris:twin:{tenant}:{device}` | Current boot/sequence, reported state, components, last seen, and connectivity. |
| twin retired-boots set | Prevents delayed previous boots from regaining authority. |
| `polaris:devices:last-seen` | Time-ordered connectivity scan index. |
| `polaris:connection:{tenant}:{device}` | Current Gateway ownership and fencing lease. |
| per-device connection epoch | Monotonically increasing fencing counter. |
| per-Gateway connection set | Operational view of owned sessions. |
| connection lease-expiry index | Cleanup and reconciliation index. |
| `polaris:gateway:{gateway}:commands` | Transient command notification channel. |
| `spatial:updates` | Transient normalized dashboard telemetry channel. |

## 21. Health, readiness, shutdown, and observability

### 21.1 Gateway probes

- Liveness indicates the process is serving.
- Readiness verifies Kafka reachability, Redis, and registry PostgreSQL.
- Readiness also reports goroutine count and database-pool state.

### 21.2 Engine probes

- Liveness detects stopped/stalled core consumer loops.
- Readiness verifies Kafka/Redis state processing, the PostgreSQL registry, and the independent archive consumer.
- It reports core status, module/component readiness, routing queue and request state, graph/snapshot identity, traffic-state counts, goroutines, and database-pool state.

### 21.3 Graceful shutdown

The Engine stops accepting HTTP traffic, cancels worker contexts, forces pending telemetry batches to flush, waits for the state, archive, command-dispatch, and traffic consumers, closes capability/routing workers, and then closes clients within a bounded shutdown deadline. The outbox relay receives cancellation and attempts a bounded final flush, although the Engine does not separately wait on an outbox completion handle. The Gateway cancels its command subscriber, shuts down HTTP serving, and terminates its live sessions with the process.

### 21.4 Metrics currently available

The backend exposes lightweight counters and readiness diagnostics for:

- active device connections;
- tasks and commands created;
- commands delivered and acknowledged;
- task/command completion and failure;
- routing requests, busy responses, queue depth/capacity, and active tenant limiters;
- process goroutines and database-pool usage.

This is operational diagnostic visibility, not a complete Prometheus/OpenTelemetry observability platform.

## 22. Deployment behavior

The Compose deployment includes:

- Gateway and Engine containers running as a non-root user;
- frontend Nginx reverse proxy and dashboard;
- Redpanda as the Kafka-compatible broker;
- Redis with append-only persistence;
- PostgreSQL/PostGIS;
- one-shot schema migration and Kafka topic initialization;
- health/readiness-gated startup dependencies;
- persistent volumes, restart policies, bounded JSON-log rotation, and an isolated network.

Only the frontend is publicly bound by default. Direct service and infrastructure ports bind to loopback for diagnostics. Browser API and WebSocket traffic uses same-origin reverse-proxy paths, so deployments do not embed browser calls to localhost.

## 23. Verification and evidence already present

The repository contains reproducible suites for:

- complete device → Gateway → Kafka → Engine → Redis/PostgreSQL → Dashboard flow;
- partial-batch timing and graceful commits;
- duplicate/out-of-order/new-boot/retired-boot telemetry;
- Redis and PostgreSQL replay behavior;
- malformed schema and DLQ processing;
- multi-partition independent progress;
- registry identity, spoof rejection, credential rotation/revocation, and tenant isolation;
- never-connected/online/stale/offline twin lifecycle;
- task assignment, capability mismatch, offline delivery, retry, expiry, and cancellation;
- duplicate and wrong-device ACK/result handling;
- Gateway ownership fencing and reconnect recovery;
- Kafka outage with transactional-outbox recovery;
- R-tree correctness against a linear oracle;
- A* correctness against Dijkstra;
- routing overload/backpressure and post-overload recovery;
- Mobility disablement and degraded operation;
- restart reconstruction of Mobility state;
- mixed 1,000-device workloads and a five-minute stability run.

The recorded scale evidence demonstrates architecture invariants on one development machine, not a universal throughput SLO or maximum production capacity.

## 24. Explicitly unimplemented or bounded capabilities

Another agent should not infer that the backend currently provides the following:

- exactly-once processing;
- Kubernetes, multi-host orchestration, service mesh, or multi-region failover;
- distributed H3 shard ownership;
- hot replacement of the road graph;
- global routing for drones or ground robots;
- low-level vehicle steering, motor control, collision avoidance, or drone attitude control;
- machine-learning demand prediction;
- an active autonomous fleet rebalancer;
- tenant-private traffic overlays—the implemented road traffic model is shared/trusted;
- a complete external identity provider or self-service operator-key administration API;
- durable dashboard replay through Redis Pub/Sub;
- full production monitoring, tracing export, alerting, backup automation, or TLS termination.

The application currently permits broad HTTP cross-origin access and permissive WebSocket origins. Authorization still applies at the application layer, but an internet-facing deployment should restrict origins and terminate TLS at a trusted ingress rather than exposing these development-friendly transport defaults unchanged.

The bundled graph is immutable for one Engine lifetime. Replacing it requires an Engine restart and a new explicit graph version.

### 24.1 Dormant experimental scaffolding

The repository still contains small research/compatibility concepts such as a static-zone strategy, an asynchronous “cellular automata” validation scaffold, and an older direct telemetry stream adapter. They are not composed into the Gateway/Engine production runtime and must not be presented as active safety simulation, autonomous orchestration, or an alternative telemetry authority. The active telemetry publisher is the canonical envelope publisher described above, and the active spatial/routing authority is the Mobility module plus the explicitly labelled compatibility match view.

## 25. Concise context for another agent

Polaris v3 is a multi-tenant, authenticated device telemetry and durable command platform. Devices are registered in PostgreSQL, receive revocable hashed credentials, and connect to the Gateway over a bidirectional WebSocket. Binary spatial telemetry is validated against authenticated identity, wrapped in a versioned deterministic event, and published to Kafka using `tenant:device` ordering. Independent manual-commit consumers atomically maintain Redis latest twins, rebuildable in-memory state, a Mobility H3/packed-R-tree index, PostgreSQL/PostGIS history, and a map-matched shared traffic overlay. Replay is normal: Redis boot/sequence classification and database uniqueness make it harmless, while invalid or retry-exhausted records reach a DLQ before offsets advance.

Operators use tenant-scoped RBAC APIs to manage tenants, projects, devices, capabilities, credentials, twins, tasks, and commands. Security-sensitive mutations write audit and outbox rows in the same transaction. Task orchestration filters by tenant, lifecycle, capability, assignment, connectivity, battery, project/type, and distance. Mobility may rank candidates by route ETA, but the core rechecks eligibility and PostgreSQL alone commits an exclusive assignment, per-device command sequence, immutable command, audit, and outbox event.

The outbox publishes commands to Kafka. A dispatcher consults Redis’s fenced Gateway lease and sends an ephemeral notification to the current Gateway. The Gateway revalidates authorization and ownership twice around the durable delivery transition, then writes to the device. ACK and result frames advance command/task state idempotently; a reconciler handles lost ACKs, exponential retry, expiry, offline recovery, and assignment release. Commands remain ordered per device and retry never changes command identity or route payload.

The optional Mobility module provides tenant-separated H3 regional sharding, exact geodesic packed-R-tree search, observation quality checks, deterministic Redis rebuild, a versioned immutable Chennai road graph, 3D KD-tree snapping, bounded A* shortest/fastest routing, Dijkstra correctness comparison, map matching, EWMA traffic, exponentially decaying immutable cost snapshots, route-aware candidate ranking, and versioned route-command plans. It can be disabled or partially degraded without replacing core registry/twin/command authority.

The intended deployment is currently a single Docker host. Delivery is at least once, road routing is road-vehicle-only, traffic is shared/trusted, and the predicted-zone endpoint is only a legacy grid-density demonstration—not machine learning or autonomous rebalancing.
