# Phase 1 Reliability Evidence

Date: 2026-08-18 (Asia/Calcutta)

## Delivery contract

Polaris provides **at-least-once delivery with idempotent consumers**. It does not claim exactly-once processing. A dependency side effect followed by a Kafka commit failure is expected to replay; deterministic event identity, database uniqueness, and version-aware latest-state operations make that replay harmless.

## Implemented controls

- Canonical `polaris.telemetry.observed` JSON envelope around the Protobuf `SpatialObject`, schema version 1.
- Required device boot and sequence fields in the Protobuf wire contract; gateway-owned ingestion, producer, correlation, causation, trace, and event metadata.
- Deterministic `event_id` from tenant/device/boot/sequence and Kafka key `tenant_id:device_id`.
- Gateway rejects malformed/oversize frames, invalid identifiers, unsupported schemas, invalid coordinates/battery/velocity/type/timestamps, and connection identity changes before Kafka.
- Three-partition ingress topic and partition-local ordered batches with a 150 ms partial flush.
- Manual commits only after the highest contiguous terminal result. Commit failure retains the batch for replay.
- Atomic Redis Lua transition at `polaris:twin:{tenant}:{device}`. Only accepted/current replay events can update volatile spatial state or publish dashboard state.
- PostgreSQL unique indexes on `event_id` and `(tenant_id, device_id, device_boot_id, sequence_number)` plus `INSERT ... ON CONFLICT DO NOTHING`.
- Transient dependency errors retry the same offset up to five times. Permanent or retry-exhausted events retain their original bytes in `telemetry.dead-letter.v1`; the source offset advances only after DLQ publication.
- Separate `/healthz` liveness and `/readyz` dependency readiness endpoints.

## Verification results

| Scenario | Evidence | Result |
|---|---|---|
| Full Compose path | Final rebuilt stack: Simulator → Gateway → Kafka partition 2 → Engine/Redis → PostgreSQL → Dashboard, ID `SMOKE-1787045926510410800`, 1,431 ms | PASS |
| Health/readiness | Gateway and engine returned `{"status":"live"}` and `{"status":"ready"}` | PASS |
| Three Kafka partitions | `telemetry.ingress` partitions 0, 1, and 2 present | PASS |
| Duplicate archive replay | Live integration test executes the same insert twice; one row remains | PASS |
| Duplicate/out-of-order/new/retired/conflicting boot state | Unit state machine and live Redis Lua classification tests | PASS |
| Redis transient failure | Injected projector fails twice; no spatial mutation occurs before Redis recovery | PASS |
| Kafka commit failure | Injected commit failure leaves the successful batch intact for replay | PASS |
| Unsupported/corrupted event | Unit DLQ test plus real malformed Kafka offset 5; DLQ high-watermark advanced and both source groups reached offset 6 with zero lag | PASS |
| Graceful pending flush | Cancellation test forced pending partitions 0 and 2 to flush and commit independently | PASS |
| PostgreSQL migration | Existing Phase 0 volume migrated; both unique indexes and all envelope columns verified | PASS |
| Redis projection | Hash contains event/boot/sequence/reported-state/last-seen fields; Lua integration classifications pass | PASS |

Commands:

```powershell
cd backend
go test ./...
go test -count=1 -tags=integration -v ./internal/application/stream
./deployments/reliability-test.ps1
```

## Failure ordering

The latest-state path intentionally performs the durable atomic Redis decision before updating volatile in-memory spatial state. If Redis succeeds and Kafka commit fails, replay returns `DUPLICATE` and safely rebuilds a restarted engine. If Redis is unavailable, spatial state is not advanced. PostgreSQL remains independent: insert success followed by commit failure replays into `ON CONFLICT DO NOTHING` and is terminally successful.

DLQ messages may be emitted independently by consumer groups for the same malformed source record. Each retains source identity, failure reason, key, and original bytes; this is consumer-failure evidence, not an exactly-once DLQ claim.
