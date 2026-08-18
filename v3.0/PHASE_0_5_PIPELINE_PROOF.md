# Polaris Phase 0.5 — Reproducible Pipeline Proof

**Run date:** 18 August 2026  
**Proof vehicle:** `SMOKE-1787037276471160700`  
**Coordinates:** `13.006700, 80.220600`

## Result

The restored pipeline was executed against the complete Docker Compose stack. Gateway, engine, frontend, Redpanda, Redis, and PostgreSQL all reached healthy state. The proof event traversed binary WebSocket ingestion, Kafka, partition-aware engine state, Redis dashboard projection, the match API, and PostgreSQL.

The browser screenshot checkpoint is not attached because the configured in-app browser service reported no available browser. The frontend container built successfully, reached healthy state, and the dashboard event itself was independently received by the smoke WebSocket client. The dashboard now exposes stable `data-testid` fields for vehicle ID and coordinates so the screenshot can be repeated when a browser session is available.

## Checkpoints

1. **Gateway health:** PASS — `GET http://localhost:6080/healthz` returned `{"status":"ok"}`.
2. **Engine health:** PASS — `GET http://localhost:6081/healthz` returned `{"status":"ok"}`.
3. **Binary Protobuf ingress:** PASS — gateway mapped `SMOKE-1787037276471160700` after the smoke client sent a binary `SpatialObject` frame.
4. **Kafka receipt and identity ordering key:** PASS — latest record metadata showed topic `telemetry.ingress`, key `alpha_logistics:SMOKE-1787037276471160700`, value size 84, partition 0, offset 4.
5. **Engine spatial state:** PASS — the engine consumer processed and committed offset 4.
6. **Match API:** PASS — `/api/v1/nodes/match` returned count 1 with the exact vehicle ID and coordinates.
7. **Redis normalized dashboard event:** PASS — the smoke dashboard WebSocket received the matching normalized event before completing.
8. **Dashboard rendering:** FUNCTIONALLY VERIFIED / SCREENSHOT BLOCKED — frontend production build and health check passed; latest-vehicle DOM fields were added. Screenshot capture was blocked by unavailable browser infrastructure.
9. **PostgreSQL archive:** PASS — `telemetry_history` contains the smoke vehicle with latitude `13.0067`, longitude `80.2206`, velocity `12.5`, battery `91`.
10. **Partial batch deadline:** PASS — flush started at 150 ms and the manual Kafka commit completed at 166 ms, below the 250 ms requirement.
11. **Manual consumer commits:** PASS — both `polaris_engine_group` and `polaris_archive_group` reported current offset equal to log end offset and total lag 0.
12. **Clean shutdown:** PASS — controlled SIGTERM logged `spatial consumer shutdown flush complete`, followed by `Engine safely terminated`; the stack was restarted and returned healthy.

## Timing

- Partial-batch queue to flush start: **150 ms**
- Partial-batch queue to Kafka commit completion: **166 ms**
- Smoke client send to dashboard receipt plus match visibility: **1,298 ms**

The end-to-end measurement includes gateway publication, broker delivery, consumer-group fetch, the 150 ms partial-batch window, Redis delivery, dashboard WebSocket receipt, and match polling. It is not yet an SLO; the broker/fetch portion should be instrumented separately before latency optimization.

## Kafka state

- Topic: `telemetry.ingress`
- Partitions: 1
- Replicas: 1
- Spatial group: `polaris_engine_group`, stable, lag 0
- Archive group: `polaris_archive_group`, stable, lag 0
- Stable raw telemetry key: `tenant_id + ":" + vehicle_id`

H3 is no longer used as the producer partition key. It remains available downstream for spatial aggregation, regional ownership, routing/congestion partitioning, and handover analysis.

## Reliability semantics now enforced

- Messages are accumulated independently per Kafka topic/partition.
- Each partition batch is sorted and processed in offset order.
- A partition flush occurs at 1,000 messages or an internal 150 ms deadline, leaving margin under the 250 ms external maximum.
- State update and Redis projection must succeed sequentially.
- Only the highest contiguous successful offset for that partition is committed.
- A failed item blocks commits beyond its offset.
- Commit failure retains the batch for at-least-once replay.
- Graceful shutdown forces pending partition batches to flush and waits for the spatial and archive consumers to close.

## Reproduction

```powershell
cd C:\Users\akash\PROJECTS\Polaris\v3.0
powershell -NoProfile -ExecutionPolicy Bypass `
  -File .\backend\deployments\smoke-test.ps1
```

The first build can take several minutes because the backend compiles the native H3 dependency. BuildKit cache mounts make subsequent Go builds materially faster.

## Follow-up risk

The architecture is intentionally at-least-once. A successful PostgreSQL insert followed by a failed Kafka commit can replay and create a duplicate because telemetry events do not yet carry a unique event ID and the database insert is not idempotent. Event identity plus `INSERT ... ON CONFLICT` is the next critical reliability task.
