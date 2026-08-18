# Polaris v3.0

Polaris is a real-time spatial control plane for connected fleet assets. Devices send Protobuf telemetry to the gateway; Redpanda/Kafka carries the canonical event; the engine updates live spatial state, projects dashboard JSON through Redis, and a separate consumer archives every event in PostgreSQL/PostGIS.

## Architecture

```text
Simulator / device
    | binary SpatialObject v1 (device boot + sequence identity)
    v
Gateway validation + platform envelope :6080
    | key = tenant_id:device_id
    v
Redpanda telemetry.ingress (3 partitions)
                         |              |
                         |              +--> idempotent PostgreSQL/PostGIS archive
                         v
                 atomic Redis latest-state + dashboard publish
                         |
                         +--> idempotent in-memory spatial state --> REST :6081
                         +--> Gateway /ws/dashboard --> Browser
```

The delivery contract is **at least once with idempotent consumers**, not exactly once. Telemetry is keyed by `tenant_id:device_id`, preserving per-device order across H3-cell movement. H3 remains downstream-only. Each Kafka partition is processed in offset order and commits only its highest contiguous successful offset after durable projection. Batches flush at 1,000 records or 150 ms. Permanent failures are published to `telemetry.dead-letter.v1` before their source offset is committed; transient failures retry without crossing the failed offset.

The gateway wraps the device frame in a schema-versioned envelope containing event, boot, sequence, observation, ingestion, correlation, causation, producer, and trace metadata. PostgreSQL enforces unique event and device-sequence identities with `ON CONFLICT DO NOTHING`. Redis atomically classifies and applies `ACCEPTED`, `DUPLICATE`, `OUT_OF_ORDER`, `NEW_BOOT`, `RETIRED_BOOT`, and `BOOT_CONFLICT` transitions at `polaris:twin:{tenant_id}:{device_id}`. Redis is decided before volatile spatial state: if a consumer restarts after Redis succeeds but before Kafka commit, replay is a Redis duplicate that safely reconstructs the engine.

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

## Development checks

```powershell
cd backend
go test ./...
```

Environment variables used by services are `GATEWAY_PORT`, `ENGINE_PORT`, `KAFKA_BROKER_URL`, `REDIS_URL`, and `POSTGRES_URL`.
