# Phase 4.1 consistency and closure evidence

Captured on 2026-08-19. This document closes the rectification plan against the current working tree. Results are measurements from this machine, not general capacity promises.

## Implemented corrections

- Removed the legacy autonomous rebalancer, direct handover/command abstractions, actor registry, old QuadTree, duplicate spatial projector, old routing handler, obsolete graph implementations, and coarse traffic analyzer from the production tree.
- Added AST architecture tests that prevent direct socket command writes outside the gateway delivery boundary, Redis command publication outside the dispatcher, NodeType bitmask use, old authoritative routing imports, and Mobility-to-delivery coupling.
- Made NAVIGATE planning explicit: `DEVICE_LOCAL` may use the generic destination command; `POLARIS_REQUIRED` must obtain a compatible domain plan and fails explicitly otherwise.
- Persisted route schema, road graph version, traffic snapshot version, generation time, validity, and immutable route payload. Delivery retry does not replan or mutate a command.
- Kept Redis as latest reported-state authority and Mobility H3/R-tree state as a replayable derivative. Lifecycle and capability decisions remain core/PostgreSQL responsibilities.
- Made shared road-traffic policy explicit as `SHARED_TRUSTED`; map matching updates one bounded directed edge and refreshes only the immutable cost overlay.
- Bounded PostgreSQL pools, removed per-device gateway reconciliation storms, batched outbox publication/marking, and made credential revalidation distinguish invalid credentials from transient database failure.
- Added deterministic mixed-system soak, real Chennai route comparison, STR mutation/query contention benchmarks, Mobility-enabled proof, and Mobility-disabled isolation proof.
- Corrected the Phase 4 nearby integration fixture: prior devices at identical coordinates could legitimately fill its limited result page. The smoke client now accepts coordinates and the Phase 4 probe uses an isolated location.

## 1,000-device mixed Compose proof

Command:

```powershell
./backend/deployments/phase4-closure-soak.ps1 -Devices 1000 -DurationSeconds 45 -Tasks 120
```

Distribution: 400 road vehicles, 200 ground robots, 250 static devices, and 150 non-spatial devices. All 1,000 authenticated connections were established with zero connection errors. During the measured interval the clients sent 7,691 telemetry events, attempted 120 tasks, received 120 commands, and performed 120 physical executions.

Safety invariants passed:

- engine, traffic, and dispatcher Kafka lag returned to zero;
- PostgreSQL persisted at least all 120 attempted tasks despite 11 client responses crossing cancellation/deadline boundaries;
- no duplicate active assignment, duplicate physical execution, command identity mutation, or cross-tenant result was observed;
- every routed plan carried graph and snapshot identity;
- routing saturation remained bounded and surfaced as request failure rather than unbounded queue growth;
- generic non-Mobility work remained functional.

| Measurement | p50 | p95 | p99 | max |
|---|---:|---:|---:|---:|
| Command delivery | 2.082 s | 6.159 s | 6.624 s | 7.255 s |
| Nearby query | 61.092 ms | 359.999 ms | 1.319 s | 2.831 s |
| Route | 68.898 ms | 428.235 ms | 961.143 ms | 2.538 s |
| Task creation/assignment | 488.870 ms | 3.644 s | 4.316 s | 4.930 s |

Connection establishment took 40.426 seconds at the controlled 25 connections/second rate. The harness recorded 54 request errors, including intentional `ROUTING_BUSY` backpressure and requests crossing the workload deadline. These are not hidden: a future longer capacity run should split this counter by HTTP/error classification.

Fourteen container-stat samples were captured. Engine memory ended lower than its first measured sample (841.9 MiB to 655.4 MiB). Gateway memory rose while establishing 1,000 live sessions and ended at 78.3 MiB; Redis ended at 16.35 MiB, Redpanda at 444.9 MiB, and PostgreSQL at 235.7 MiB. This finite run found no persistent growth, but it is not a multi-hour leak proof.

Raw artifacts: `PHASE_4_1_SOAK_RESULT.json` and `PHASE_4_1_CONTAINER_STATS.jsonl`.

## Spatial contention decision

The packed STR R-tree was tested with 80/20, 50/50, and 20/80 write/read mixes:

| States | 80% writes | 50% writes | 20% writes |
|---:|---:|---:|---:|
| 1,000 | 40.1 us/op | 55.4 us/op | 84.6 us/op |
| 5,000 | 291.6 us/op | 439.3 us/op | 477.5 us/op |
| 10,000 | 345.5 us/op | 614.4 us/op | 935.5 us/op |

No correctness or race failure was found, so the simpler lazy packed rebuild remains. A background copy-on-write rebuild is deferred until a longer production-shaped run demonstrates tail-latency pressure.

## Real Chennai routing decision

The graph contains 690,268 road nodes and 1,442,876 directed edges. A* and Dijkstra returned equal route costs. A* expanded substantially fewer nodes and ran faster for short, medium, long, and dense-city samples; the edge-of-graph sample was slightly slower despite fewer expansions.

| Class | A* latency / expanded | Dijkstra latency / expanded |
|---|---:|---:|
| Short | 4.037 ms / 2,776 | 6.695 ms / 6,727 |
| Medium | 73.709 ms / 48,440 | 149.545 ms / 117,481 |
| Long | 281.665 ms / 147,753 | 539.173 ms / 363,945 |
| Dense | 65.402 ms / 43,774 | 159.285 ms / 123,019 |
| Edge | 1,174.738 ms / 540,987 | 1,030.949 ms / 671,852 |

A* remains operational and Dijkstra remains its oracle; no universal A* speed claim is made.

## Regression result

- `go test ./...`: PASS
- `go vet ./...`: PASS
- focused race suite over Mobility, extension, orchestration, and stream processing: PASS
- Phase 0.5 smoke: PASS
- Phase 1 reliability: PASS
- Phase 2 identity/twin/security: PASS
- Phase 3 command orchestration: PASS
- Phase 4 Mobility enabled: PASS
- Phase 4 Mobility disabled: PASS
- Phase 4.1 1,000-device mixed Compose soak: PASS

## Remaining bounded limitations

- The road graph is immutable in-process and graph replacement requires a clean engine restart.
- Global road routing covers road vehicles; aerial and indoor/ground global planners remain future modules.
- The 45-second workload demonstrates bounded behavior and architectural invariants, not a production SLO or multi-hour capacity/leak claim.
- The soak's aggregate request-error counter should be separated into backpressure, timeout, cancellation, and unexpected errors before formal performance gating.
- Delivery remains at least once with idempotent consumers. Polaris does not claim exactly-once processing.

## Outcome

Phase 4.1 removes competing legacy authorities, makes planning responsibility and data ownership explicit, proves the mixed pipeline at 1,000 authenticated devices, and freezes Mobility behind the generic durable task/command core. The remaining items are documented operational limits rather than contradictory implementation paths.


