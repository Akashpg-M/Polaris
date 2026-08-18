# Polaris v3.0

Polaris is a real-time spatial control plane for connected fleet assets. Phase 2 adds a durable, tenant-scoped device registry and makes authenticated registry identity the authority for every event. Devices send Protobuf telemetry to the gateway; Redpanda/Kafka carries the canonical event; the engine updates live spatial state, projects dashboard JSON through Redis, and a separate consumer archives every event in PostgreSQL/PostGIS.

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

Task API --> deterministic capability/twin selection --> exclusive assignment
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

Tasks express business intent, capability requirements, battery/distance constraints, priority, target, project scope, and expiry independently from device commands. Assignment is deterministic—distance, battery, then stable device ID—and an exclusive partial index prevents two active assignments from claiming the same device. The assignment, per-device command sequence, command envelope, audit event, and outbox intent are committed before any socket delivery.

`device.command.v1` is keyed by `tenant_id:device_id`. Redis stores only live gateway ownership, lease expiry, and a monotonically increasing fencing epoch; PostgreSQL and Kafka remain the durable sources. Redis Pub/Sub accelerates routing to the current gateway, while periodic database reconciliation and reconnect reconciliation guarantee store-and-forward recovery if a notification is missed.

The authenticated telemetry WebSocket is bidirectional. Telemetry remains binary Protobuf. Server commands and device `COMMAND_ACK`/`COMMAND_RESULT` messages use explicit versioned JSON envelopes. ACK means received/accepted; completion is a separate result. Retry keeps the same `command_id` and sequence. The simulator caches completed command IDs and resends its prior ACK/result without repeating execution.

Task and command APIs are tenant scoped and role protected. Operators can create/cancel tasks; tenant and platform administrators can force command/task retries; viewers are read-only. Local cancellation is allowed only before delivery. Commands that have reached a device cannot be silently cancelled as though execution stopped.

## Run locally with Docker Compose

Requirements: Docker Desktop with Compose v2.

```powershell
docker compose -f backend/deployments/docker-compose.yml up -d --build --wait
```

Services expose gateway `6080`, engine `6081`, frontend `5173`, Redpanda `9092`, Redis `6379`, and PostgreSQL `5432`. `/healthz` is liveness; `/readyz` probes dependencies required for assigned work. Compose waits on readiness plus one-shot Kafka initialization and database migration jobs. Check status with:

```powershell
docker compose -f backend/deployments/docker-compose.yml ps
```

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
