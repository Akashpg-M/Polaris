# Phase 4 baseline

Captured 2026-08-19 before Phase 4 changes.

- Baseline commit: `59bd32c84ca505cdbd61ae671e6152c665828887`
- Database migration source: `backend/deployments/init.sql`, through the Phase 3 task/assignment/command schema.
- Kafka topics: `telemetry.ingress`, `telemetry.dead-letter.v1`, `device.lifecycle.v1`, `device.connectivity.v1`, `task.lifecycle.v1`, `device.command.v1`, `device.command.ack.v1`, and `device.command.result.v1` (three partitions each except the initialization marker).
- `go test ./...`: PASS.
- `go vet ./...`: PASS.
- Frontend build at baseline capture: host Node/npm unavailable. The same source was subsequently verified through the frontend Docker production stage (`tsc -b` + Vite): PASS.

## Frozen Phase 3 behavior

The task API persists intent first, deterministically ranks eligible devices, and commits the exclusive assignment, per-device command sequence, immutable command payload, audit record, and outbox event in PostgreSQL before delivery. ACK and result are separate. Retry reuses the command ID, sequence, and payload. Gateway ownership is fenced and commands are reconciled after missed notifications or restarts.

## Known pre-Phase-4 spatial and routing limits

Candidate distance ranking is embedded in the orchestration service. The live spatial engine uses a legacy quadtree. The road graph is mutable, performs a linear nearest-node scan, and exposes Dijkstra under a route handler; dynamic edge mutation can affect concurrent searches. Twin state is a single `reported_state` JSON document rather than versioned components.

## Non-regression contract

Phase 4 must preserve telemetry ordering and at-least-once/idempotent processing, authenticated tenant/device identity, tenant isolation, registry and twin behavior, durable tasks, exclusive assignments, per-device command order, transactional outbox delivery, gateway fencing, ACK/result semantics, bounded retry and expiry, and restart reconciliation. Domain modules may propose candidates and execution plans; PostgreSQL remains assignment authority and Phase 3 remains the sole command delivery path.
