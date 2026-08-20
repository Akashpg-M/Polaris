# Phase 4.2 Final Stabilization and Closure Evidence

Measured on 20 August 2026 with Docker Desktop on Windows/amd64. These results are architecture and regression evidence from one development machine, not production SLOs or maximum-capacity claims.

## Outcome

Phase 4 is closed and frozen. Polaris preserves PostgreSQL as the durable task/command authority, Kafka as the ordered telemetry and command transport, Redis as the idempotent latest-state projection, and Mobility as an optional capability module. The closure work added visibility and validation; it did not introduce a competing execution path or a new platform feature.

## Closure changes

- The mixed harness classifies `routing_busy`, `timeout`, `cancelled`, `conflict`, `no_route`, `client_error`, `server_error`, `transport_error`, and `unexpected` separately for nearby, route, task, and command operations.
- Task responses expose lightweight candidate-selection, routing/planning, persistence, and total timings.
- Command test envelopes carry volatile relay-published and gateway-received timestamps. These timestamps are not persisted in the durable command and cannot change command identity on retry.
- Engine and gateway readiness responses expose goroutine and database-pool state; Mobility readiness also exposes routing requests, backpressure, queue depth/capacity, and active tenant limiters.
- Gateway command delivery uses 16 fixed hash-sharded workers. Hashing `tenant_id:device_id` preserves per-device order while removing the previous global serial-delivery bottleneck.
- Candidate evaluation is capped at 50 proposals, Redis twin/connection reads are pipelined, and authoritative registry eligibility is rechecked once per attempt. PostgreSQL remains the final assignment authority.
- The capability lookup has a tenant/capability/device index. The bulk-only soak seed runs `ANALYZE` before measurement so PostgreSQL does not execute a stale-cardinality nested-loop plan.
- The closure gate now checks the archive consumer group as well as the engine, traffic, and command-dispatcher groups.

## Defects found and corrected by the final regression

1. The smoke test could observe Redis/dashboard state before the independent PostgreSQL archive consumer joined and committed. The archive checkpoint now polls for up to 60 seconds and reports `polaris_archive_group` state if it expires.
2. The Phase 4 road proof leaked `SMOKE_NODE_TYPE=3` into a nested generic drone smoke test. Smoke inputs are now explicit and the Phase 4 helper restores its temporary environment before running nested regressions.
3. Calling `RequestID` more than once in one request generated different IDs for audit and API response. The ID is now cached in the Gin request context, with tests for generated and caller-supplied IDs.

## 1,000-device mixed validation

Workload: 400 road vehicles, 200 ground robots, 250 static devices, 150 non-spatial devices; 45 seconds; 120 tasks.

| Result | Measured value |
|---|---:|
| Authenticated connections | 1,000 |
| Connection errors | 0 |
| Telemetry frames | 7,882 |
| Task attempts | 120 |
| Commands delivered | 120 |
| Physical executions | 120 |
| Duplicate physical executions | 0 |
| Command identity mutations | 0 |
| Unexpected/server/transport errors | 0 / 0 / 0 |
| Expected bounded outcomes | 31 route `ROUTING_BUSY`, 4 task conflicts, 3 task `ROUTING_BUSY` |

Engine, archive, traffic, and command-dispatcher consumer lag was verified at zero. Database checks found no duplicate active assignment, cross-tenant command/task mismatch, or missing route graph/snapshot identity.

### Latency breakdown

| Path or stage | p50 | p95 | p99 |
|---|---:|---:|---:|
| Nearby query | 34.457 ms | 175.332 ms | 975.355 ms |
| Route | 47.164 ms | 210.499 ms | 444.735 ms |
| Task total | 267.879 ms | 1.849 s | 2.729 s |
| Task candidate selection | 39.558 ms | 407.253 ms | 658.568 ms |
| Task routing/planning | 0.013 ms | 417.004 ms | 660.004 ms |
| Task persistence | 96.173 ms | 670.634 ms | 1.413 s |
| Command total | 3.077 s | 5.259 s | 5.326 s |
| Persist to Kafka | 710.074 ms | 1.347 s | 1.472 s |
| Kafka to gateway | 2.082 s | 4.720 s | 4.829 s |
| Gateway to ACK | 30.293 ms | 115.766 ms | 1.816 s |

The first instrumented attempt exposed a multi-second candidate-selection path caused by repeated fleet-wide eligibility queries and stale PostgreSQL statistics after bulk seeding. Bounded evaluation, one authoritative recheck, the covering index, and seed-time `ANALYZE` removed that accidental behavior. Command timing then exposed global sequential gateway delivery; fixed hash-sharded workers reduced the five-minute run's Kafka-to-gateway p95 to 1.616 seconds. No stage waited indefinitely or grew continuously.

Raw evidence: `PHASE_4_2_SOAK_RESULT.json` and `PHASE_4_2_CONTAINER_STATS.jsonl`.

## Routing overload and recovery

The system-level overload proof used one routing worker, queue capacity 4, and 80 concurrent route requests.

- Baseline and post-overload routes returned 200.
- 75 requests received `ROUTING_BUSY`, 4 completed, and 1 reached its bounded timeout.
- Unexpected errors were zero.
- 40 telemetry frames were accepted during saturation.
- A generic `RUN_MODEL` task returned 201 and its command completed during saturation.
- Final routing queue depth was 0/4; active tenant limiter entries were zero; engine readiness remained `READY`.

Raw evidence: `PHASE_4_2_ROUTING_OVERLOAD_RESULT.json`.

## Mobility restart and rebuild

The restart proof populated active road, robot, and static devices plus inactive, non-spatial, and foreign-tenant controls. It advanced the road vehicle to sequence 2, injected stale sequence 1 at a different position, restarted Engine, and observed the readiness transition.

- Mobility was unavailable/not ready during rebuild and returned to `READY` afterward.
- The three valid active spatial devices were recovered exactly once.
- Inactive, non-spatial, and foreign-tenant devices were excluded from the tenant result.
- The foreign tenant saw only its own spatial device.
- The road vehicle retained sequence 2 and its newer position.

Raw evidence: `PHASE_4_2_MOBILITY_REBUILD_RESULT.json`.

## Five-minute stability validation

Workload: 1,000 authenticated devices, 300 seconds, 51,000 telemetry frames, and 240 tasks.

| Result | Measured value |
|---|---:|
| Connections established / errors | 1,000 / 0 |
| Commands / physical executions | 240 / 240 |
| Duplicate deliveries / identity mutations | 0 / 0 |
| Expected bounded outcomes | 22 route busy, 6 task conflicts, 5 task busy |
| Unexpected/server/transport errors | 0 / 0 / 0 |
| Final engine routing queue | 0 / 64 |
| Final engine DB pool | 10 open, 0 in use |
| Final gateway DB pool | 10 open, 0 in use |
| Final engine/gateway goroutines | 75 / 36 |

Forty-six resource samples were captured. Engine memory was 766.8 MiB initially, 712.9 MiB at midpoint, and 712.9 MiB in the last active sample. Gateway memory rose with the 1,000 live sessions and stabilized from 79.2 MiB at midpoint to 80.8 MiB in the last active sample. PostgreSQL ranged from 171.4 to 251.8 MiB, Redis from 11.3 to 20.1 MiB, and Redpanda from 486.6 to 536.7 MiB. Gateway goroutines peaked at 2,067 with active sessions and returned to 36 after disconnect. The routing queue's sampled maximum was zero and all consumer lag gates drained to zero. This short run found no obvious continuously growing resource or stranded work.

Five-minute p95 values were 125.356 ms nearby, 145.978 ms route, 1.712 seconds task total, and 3.441 seconds command delivery. The command breakdown was 1.715 seconds persist-to-Kafka, 1.616 seconds Kafka-to-gateway, and 1.192 seconds gateway-to-ACK at p95.

Raw evidence: `PHASE_4_2_STABILITY_SOAK_RESULT.json` and `PHASE_4_2_STABILITY_CONTAINER_STATS.jsonl`.

## Critical invariant and regression result

| Gate | Result |
|---|---|
| Tenant isolation and authenticated identity authority | PASS |
| Exclusive PostgreSQL device assignment | PASS |
| Gateway ownership fencing | PASS |
| Immutable command retry and route-plan identity | PASS |
| One immutable routing snapshot per request | PASS |
| A* cost equals Dijkstra oracle | PASS |
| STR R-tree equals linear oracle | PASS |
| Core capability authority after Mobility proposal | PASS |
| `DEVICE_LOCAL` / `POLARIS_REQUIRED` semantics | PASS |
| Mobility-disabled generic platform flow | PASS |
| Phase 0.5 through Phase 4 full Compose regression | PASS |
| `go test ./...` | PASS |
| `go vet ./...` | PASS |
| Focused race suite over Mobility, extension, orchestration, stream, and handlers | PASS |

The full chained Compose run completed with exit code 0. It includes Phase 0.5 telemetry, Phase 1 replay/reliability, Phase 2 identity/tenancy/twins, Phase 3 assignment/command/fencing/outage recovery, and Mobility-enabled Phase 4. The separate Mobility-disabled Compose run also completed with exit code 0. Overload and restart/rebuild proofs passed independently.

## Freeze boundary

Phase 4.2 deliberately does not add graph hot swapping, distributed H3/graph ownership, new spatial indexes, drone or ground-robot global planners, multi-region operation, Kubernetes, service mesh, machine-learning prediction, autonomous rebalancing, or a large observability stack. Delivery remains at least once with idempotent consumers; Polaris does not claim exactly-once processing.

