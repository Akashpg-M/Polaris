# Polaris v3.0

Polaris is a tenant-isolated edge-device control plane. Its core owns identity, registry, digital twins, durable tasks, exclusive assignment, command sequencing, outbox delivery, gateway fencing, retry, and reconciliation. Compile-time capability modules add domain behavior without bypassing those guarantees. Mobility is the first module: it projects spatial telemetry into H3-sharded R-trees, provides versioned road routing, proposes nearby candidates, and plans NAVIGATE/RELOCATE commands.

## Architecture

```text
Operator API key --> tenant / project / device registry --> hashed device credential
                                                          |
Simulator / device -- bearer credential or one-time ticket-+
    | authenticated binary SpatialObject v1
    v
Gateway authentication + validation + server-derived identity :6080
    | key = tenant_id:device_id
    v
Redpanda telemetry.ingress (3 partitions)
                         |              |
                         |              +--> idempotent PostgreSQL/PostGIS archive
                         v
                 atomic Redis latest-state + dashboard publish
                         |
                         +--> idempotent in-memory spatial state --> REST :6081
                         +--> tenant-filtered Gateway /ws/dashboard --> Browser

PostgreSQL registry mutation --> audit row + outbox row (one transaction)
                                      |
                                      +--> embedded relay --> device.lifecycle.v1

PostgreSQL registry metadata + capabilities
                         + Redis reported state / connectivity
                                      |
                                      +--> tenant-isolated digital-twin API

Task API --> deterministic capability/twin selection --> explicit planning mode --> exclusive assignment
    |                         |
    |                         +--> PostgreSQL command + audit + outbox (atomic)
    |                                              |
    |                                              v
    |                                  device.command.v1 (Kafka)
    |                                              |
    v                                              v
Task reconciliation <-- ACK/result <-- authenticated bidirectional WebSocket
                                              ^
Redis leased gateway ownership + fencing epoch-+
```

The delivery contract is **at least once with idempotent consumers**, not exactly once. Telemetry is keyed by `tenant_id:device_id`, preserving per-device order across H3-cell movement. H3 remains downstream-only. Each Kafka partition is processed in offset order and commits only its highest contiguous successful offset after durable projection. Batches flush at 1,000 records or 150 ms. Permanent failures are published to `telemetry.dead-letter.v1` before their source offset is committed; transient failures retry without crossing the failed offset.

The gateway wraps the device frame in a schema-versioned envelope containing event, boot, sequence, observation, ingestion, correlation, causation, producer, and trace metadata. PostgreSQL enforces unique event and device-sequence identities with `ON CONFLICT DO NOTHING`. Redis atomically classifies and applies `ACCEPTED`, `DUPLICATE`, `OUT_OF_ORDER`, `NEW_BOOT`, `RETIRED_BOOT`, and `BOOT_CONFLICT` transitions at `polaris:twin:{tenant_id}:{device_id}`. Redis is decided before volatile spatial state: if a consumer restarts after Redis succeeds but before Kafka commit, replay is a Redis duplicate that safely reconstructs the engine.

## Phase 2 identity and registry

Registry, twin, spatial, route, and dashboard-ticket APIs require a bearer operator key. The development-only platform key is supplied at startup through `DEV_PLATFORM_ADMIN_TOKEN`; it is hashed before storage and must never be committed. Platform administrators select tenant scope with `X-Tenant-ID`; tenant-scoped roles always derive it from the authenticated key.

Device credentials contain 256 bits of cryptographic secret material and are stored only as hashes. A telemetry WebSocket is upgraded only after an active credential resolves to an active tenant and device. The authenticated principal must match the redundant payload identity, and the Kafka envelope always uses the principal identity. Browser simulators exchange operator authentication for a short-lived, hashed, one-use connection ticket. Credential, device, and tenant status are revalidated during an active telemetry session, so revocation takes effect on its next frame.

Registry mutations atomically write their audit and outbox records. The embedded relay publishes lifecycle events with at-least-once semantics. Digital twins compose PostgreSQL metadata and capabilities with the Phase 1 Redis reported state. The Redis last-seen index supports `NEVER_CONNECTED`, `ONLINE`, `STALE`, and `OFFLINE` states without scanning keys.

## Phase 3 task and command orchestration

Tasks express business intent, capability requirements, battery/distance constraints, priority, target, project scope, planning responsibility, and expiry independently from device commands. Assignment is deterministic—domain score/distance, battery, then stable device ID—and an exclusive partial index prevents two active assignments from claiming the same device. The assignment, per-device command sequence, command envelope, audit event, and outbox intent are committed before any socket delivery.

`device.command.v1` is keyed by `tenant_id:device_id`. Redis stores only live gateway ownership, lease expiry, and a monotonically increasing fencing epoch; PostgreSQL and Kafka remain the durable sources. Redis Pub/Sub accelerates routing to the current gateway, while periodic database reconciliation and reconnect reconciliation guarantee store-and-forward recovery if a notification is missed.

The authenticated telemetry WebSocket is bidirectional. Telemetry remains binary Protobuf. Server commands and device `COMMAND_ACK`/`COMMAND_RESULT` messages use explicit versioned JSON envelopes. ACK means received/accepted; completion is a separate result. Retry keeps the same `command_id` and sequence. The simulator caches completed command IDs and resends its prior ACK/result without repeating execution.

Task and command APIs are tenant scoped and role protected. Operators can create/cancel tasks; tenant and platform administrators can force command/task retries; viewers are read-only. Local cancellation is allowed only before delivery. Commands that have reached a device cannot be silently cancelled as though execution stopped.

## Phase 4 capability modules and Mobility

The engine constructs an explicit extension registry at startup. Modules have `Start`, `Ready`, and `Close` lifecycle methods; candidate providers only propose ordered device IDs; task planners only produce execution intent. There is no dynamic Go plugin loading, hidden registration, or direct module-to-device delivery.

```text
POLARIS CORE
  identity / tenancy / registry / versioned twins
  durable tasks / authoritative eligibility / PostgreSQL reservation
  immutable commands / outbox / Kafka / gateway fencing
                         |
                extension contracts
                         |
CAPABILITY MODULES       +-- default generic planner
                         +-- Mobility
                              H3 tenant/region shards
                              moving-device R-tree
                              immutable OSM graph + 3D KD-tree
                              A* + Dijkstra oracle
                              atomic traffic cost snapshots
                              candidate provider + task planner
                         |
FUTURE INTELLIGENCE / ADVISOR LAYER
```

Twin state remains generic. Redis atomically stores independently versioned `spatial/v1` and `battery/v1` envelopes while retaining the compatibility reported-state view. Modules decode only components they own; non-spatial compute nodes and cameras continue using the same registry, task, and command core.

Raw telemetry remains keyed by `tenant_id:device_id`; H3 is downstream-only. Mobility uses tenant plus coarse H3 parents as concurrency boundaries and an exact R-tree with Haversine post-filtering. Reported position updates on every accepted observation, while indexed position moves on the configured distance threshold, H3 transition, or maximum age. Stale, offline, suspended, decommissioned, and disabled-tenant devices are evicted. Restart rebuild reads current components from Redis and revalidates lifecycle in PostgreSQL.

Road topology is immutable after OSM load. A 3D KD-tree snaps to static road nodes, A* is operational, and Dijkstra remains a correctness oracle. Every request reads one immutable cost snapshot identified by both `road_graph_version` and `routing_snapshot_version`. Bounded global/per-tenant execution returns explicit `ROUTING_BUSY`, `ROUTING_TIMEOUT`, `ROUTING_UNAVAILABLE`, `NO_ROUTE`, and unsupported-profile failures. Tenant limiter entries are reference-counted and removed after requests.

Mobility may be disabled with `POLARIS_MODULE_MOBILITY_ENABLED=false`. A missing road graph degrades routing without disabling generic telemetry, twins, CAPTURE_IMAGE, RUN_MODEL, ACK, or result processing unless Mobility is mandatory. Road routing is currently implemented only for `ROAD_VEHICLE`; drones and ground robots still participate in spatial discovery. Local collision avoidance, steering, motor control, and flight attitude remain device responsibilities.

Useful APIs are `GET /api/v1/spatial/devices/nearby`, `POST /api/v1/routes`, and the authoritative `POST /api/v1/tasks`. `DEVICE_LOCAL` navigation sends a high-level destination through the generic planner. `POLARIS_REQUIRED` requires a compatible domain planner and returns `PLANNER_UNAVAILABLE` instead of silently falling back. Routed plans record route schema, road graph, traffic snapshot, generation/validity times, origin, destination, waypoints, distance, and duration. An expired plan is never silently substituted inside an existing command; task retry creates a new plan, while delivery retry reuses the immutable command.

Traffic is intentionally `SHARED_TRUSTED`: authenticated road telemetry from platform ingress contributes to one shared Chennai road-cost overlay. The overlay is an eight-byte-per-edge immutable cost slice refreshed at a bounded interval; topology is never copied during refresh.

### State and delivery ownership

| Concern | Authority | Derived/cache layer |
|---|---|---|
| Raw telemetry order and replay | Kafka `telemetry.ingress` | PostgreSQL archive |
| Latest accepted reported state and freshness | Redis twin Lua transition | compatibility engine map and Mobility H3/R-tree |
| Registry, capabilities, tasks, assignments and commands | PostgreSQL | API response models |
| Road topology | immutable versioned Mobility graph | 3D KD-tree node index |
| Dynamic traffic | immutable routing cost snapshot | map-matched EWMA observations |
| Device command delivery | PostgreSQL command + outbox + Kafka | Redis notification and fenced gateway socket |

Redis Pub/Sub is never a durable command or lifecycle source. The autonomous rebalancer, direct Redis command publisher, actor ingress state, old QuadTree, and old routing handler have been retired from production code.

## Run locally with Docker Compose

Requirements: Docker Desktop with Compose v2.

```powershell
.\backend\deployments\start.ps1
```

The script generates an ignored deployment `.env` with secure random local credentials, validates Compose, builds the images, and waits for readiness. The public dashboard is exposed on `5173`; direct gateway `6080`, engine `6081`, Redpanda `9092`, Redis `6379`, and PostgreSQL `5432` ports bind to loopback by default. Browser API and WebSocket traffic is routed through the frontend so deployments do not contain baked-in `localhost` endpoints.

`/healthz` is liveness; `/readyz` probes dependencies required for assigned work. Compose waits on readiness plus one-shot Kafka initialization and database migration jobs. Check status with:

```powershell
.\backend\deployments\status.ps1
```

Use `logs.ps1`, `stop.ps1`, and `verify-deployment.ps1` for normal operation and deployment validation. The complete configuration, secret, start/stop, verification, update, and backup procedures are in [DEPLOYMENT.md](DEPLOYMENT.md).

## End-to-end smoke test

The smoke test starts the stack, opens a dashboard socket, sends one canonical Protobuf frame through the telemetry socket, waits for the Redis-backed dashboard event, queries engine state, and verifies the PostgreSQL row:

```powershell
./backend/deployments/smoke-test.ps1
```

Success prints the complete verified path and the unique smoke node ID.

The captured Phase 0.5 evidence and checkpoint results are recorded in `PHASE_0_5_PIPELINE_PROOF.md`.

Run the Phase 1 unit, live Redis/PostgreSQL integration, and Compose checks with:

```powershell
./backend/deployments/reliability-test.ps1
```

The tested failure semantics and current evidence are recorded in `PHASE_1_RELIABILITY_EVIDENCE.md`.

Run the authenticated Phase 2 registry, credential, isolation, twin, connectivity, outbox, audit, and regression proof with:

```powershell
./backend/deployments/phase2-identity-test.ps1
```

The exact scenarios, live results, operational assumptions, and API inventory are recorded in `PHASE_2_IDENTITY_AND_TWIN_EVIDENCE.md`.

Run the complete Phase 3 task/command, ACK/result, retry, fencing, outage, replay, RBAC, and recovery proof with:

```powershell
./backend/deployments/phase3-command-test.ps1
```

The architecture contract and captured evidence are documented in `PHASE_3_COMMAND_ORCHESTRATION_EVIDENCE.md`.

## Development checks

```powershell
cd backend
go test ./...
```

Environment variables used by services include `GATEWAY_PORT`, `ENGINE_PORT`, `KAFKA_BROKER_URL`, `REDIS_URL`, `POSTGRES_URL`, `DEV_PLATFORM_ADMIN_TOKEN`, `DEVICE_STALE_AFTER`, `DEVICE_OFFLINE_AFTER`, `OFFLINE_SCAN_INTERVAL`, `CONNECTION_TICKET_TTL`, `OUTBOX_BATCH_SIZE`, `OUTBOX_POLL_INTERVAL`, `GATEWAY_ID`, `CONNECTION_LEASE_TTL`, `COMMAND_ACK_TIMEOUT`, `COMMAND_RECONCILE_INTERVAL`, and `COMMAND_MAX_ATTEMPTS`.

Mobility bounds use the `POLARIS_MODULE_MOBILITY_*` and `MOBILITY_*` variables in `backend/deployments/docker-compose.yml`. Invalid H3, candidate, capacity, queue, timeout, expansion, or traffic-age limits fail startup.

Run the Phase 4 architecture proof and measured benchmark with:

```powershell
./backend/deployments/phase4-mobility-test.ps1
./backend/deployments/phase4-mobility-test.ps1 -FullRegression
cd backend
go run ./cmd/mobility-benchmark -devices 1000 -queries 5000 -moving-percent 40 -search-radius 5000 -candidate-limit 10
go run ./cmd/routing-benchmark -graph data/chennai-metro.osm.pbf -version chennai-v1
go test ./internal/modules/mobility/... -bench . -benchmem
./deployments/phase4-closure-soak.ps1 -Devices 1000 -DurationSeconds 45 -Tasks 120
```

Run the Phase 4.2 closure proofs with:

```powershell
./backend/deployments/phase4-routing-overload-test.ps1
./backend/deployments/phase4-mobility-rebuild-test.ps1
./backend/deployments/phase4-closure-soak.ps1 -Devices 1000 -DurationSeconds 45 -Tasks 120
./backend/deployments/phase4-closure-soak.ps1 -Devices 1000 -DurationSeconds 300 -Tasks 240 -EvidenceName PHASE_4_2_STABILITY
./backend/deployments/phase4-mobility-test.ps1 -FullRegression -SkipLocalChecks
./backend/deployments/phase4-module-isolation-test.ps1
```

Phase 4 implementation evidence is recorded in `PHASE_4_CAPABILITY_MOBILITY_EVIDENCE.md`. Rectification decisions and the original mixed measurements are in `PHASE_4_1_CONSISTENCY_EVIDENCE.md`. Final error classification, timing, overload recovery, restart/rebuild, five-minute stability, and regression results are in `PHASE_4_2_FINAL_STABILIZATION_EVIDENCE.md`. The project-wide authority and inconsistency audit is in `CODEBASE_CONSISTENCY_REPORT.md`.

The recorded soak is an architecture and safety-invariant proof on one development machine. It is not a production SLO, maximum-capacity claim, or long-duration leak certification.
