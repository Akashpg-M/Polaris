# Polaris v3.0 final codebase consistency report

Captured on 2026-08-19 after Phases 0 through 4.1.

## Project objective and present capability

Polaris is a multi-tenant edge-device control plane. It authenticates registered devices, validates and orders telemetry, maintains durable history and current digital twins, selects eligible devices for tasks, persists immutable commands, and delivers those commands through fenced live gateway sessions. Capability modules may add domain intelligence without acquiring authority over tenancy, reservation, persistence, or device delivery.

The implemented Mobility module adds tenant-sharded spatial projection, nearby discovery, immutable Chennai road topology, snapshot-consistent routing, shared trusted traffic observation, candidate ranking, and durable NAVIGATE/RELOCATE planning. Generic devices and tasks continue to operate when Mobility is disabled or routing is degraded.

## Authoritative ownership model

| Concern | Authority | Rebuildable/ephemeral representation |
|---|---|---|
| Telemetry order and replay | Kafka, keyed by tenant and device | consumer batches |
| Telemetry history | PostgreSQL idempotent archive | none |
| Latest accepted reported state | atomic Redis twin transition | engine map and Mobility index |
| Tenants, projects, devices, credentials, roles, capabilities | PostgreSQL registry | API models/caches |
| Tasks, assignment, command sequence and immutable command | PostgreSQL transaction | dispatch/reconciliation views |
| Command/event publication intent | PostgreSQL transactional outbox | Kafka and Redis wake-up notification |
| Live connection ownership | Redis lease plus fencing epoch | gateway socket |
| Road topology | immutable versioned graph | KD-tree node index |
| Dynamic road cost | immutable versioned cost snapshot | EWMA observation accumulator |

## Inconsistencies found and resolved

### Competing architectures

- Removed the autonomous rebalancer and old handover path so domain logic cannot bypass durable tasks and commands.
- Removed obsolete actor ingress state, old QuadTree, duplicate spatial projection, duplicate route handler/graph stack, direct commander abstraction, and coarse traffic analyzer.
- Added architecture tests to prevent these dependency and delivery bypass patterns from returning.

### Planning ambiguity

- Added explicit `DEVICE_LOCAL` and `POLARIS_REQUIRED` semantics. Device-local navigation may carry a destination only. Polaris-required navigation fails when no compatible planner exists instead of silently fabricating a generic route.
- Durable route plans now identify route schema, road graph, cost snapshot, validity window, and complete immutable execution payload.

### Selection and tenancy

- Fixed candidate ranking that could drop core-eligible unranked devices.
- Added eligible-device IDs to provider requests and filters before expensive route scoring; core still revalidates every proposal and PostgreSQL remains reservation authority.
- Pipelined twin/connection reads per task and removed stale tenant limiter entries.
- Confirmed tenant context flows from operator authentication to core requests and Mobility shard keys; untrusted payload identity cannot select another tenant/device.

### Reliability and dependency behavior

- Bounded all PostgreSQL pools to prevent connection saturation.
- Batched outbox Kafka writes and PostgreSQL publication markers while retaining at-least-once replay.
- Removed the per-device/per-second gateway reconciliation storm; durable engine and reconnect reconciliation remain.
- Gateway lease heartbeat tolerates short Redis interruptions but closes on fencing mismatch or TTL-window failure.
- Credential revalidation distinguishes actual revocation from temporary database failure, using a bounded grace interval for the latter.
- Non-spatial authenticated devices use their active fenced lease for online eligibility rather than requiring a spatial telemetry component.

### Verification defects

- Normalized local test endpoints to `127.0.0.1` to avoid Windows localhost/IPv6 nondeterminism.
- Corrected Phase 3 wrong-ACK expectations: an ACK mismatch can legitimately leave a command pending for timeout redelivery, but can never acknowledge/complete it.
- Extended the restart fixture expiry to survive Docker Desktop restart latency.
- Corrected the Phase 4 nearby fixture so accumulated equal-distance devices cannot hide the asserted device behind the endpoint limit.
- The soak harness now establishes the whole fleet before timing work, validates semantic counters rather than unreliable child-process exit metadata, and verifies database invariants independently.

## Architectural audit findings

- One command authority remains: PostgreSQL command/outbox -> Kafka -> dispatcher -> fenced gateway. Redis Pub/Sub is only an acceleration signal.
- One route authority remains: both route APIs, candidate scoring, and task planning delegate to the Mobility routing engine.
- One connectivity vocabulary remains: `NEVER_CONNECTED`, `ONLINE`, `STALE`, and `OFFLINE`; Mobility only includes or evicts derived index membership.
- One capability authority remains: registry capabilities. Mobility consumes capability requests but does not maintain competing booleans.
- Node types are enums; repository-wide static checks reject bitmask operations over them.
- Redis decides boot/sequence freshness. Mobility's comparison is a defensive monotonic projection guard and cannot make a rejected event authoritative.
- Traffic sharing is explicitly `SHARED_TRUSTED`; topology is global, while selection, quotas, devices, twins, tasks, and commands remain tenant-scoped.
- Snapshot refresh copies the edge-cost overlay, not the 690,268-node/1,442,876-edge topology.

## Verification summary

All Go unit/integration tests, vet, the focused race suite, Phase 0.5, Phase 1, Phase 2, Phase 3, both Phase 4 module modes, and the Phase 4.1 mixed soak pass. The soak established 1,000 authenticated devices and completed 120 persisted task/command executions with no duplicate physical execution, active double assignment, tenant leakage, or identity mutation. See `PHASE_4_1_CONSISTENCY_EVIDENCE.md` for measurements.

## Open limitations, not hidden defects

1. Road-graph replacement requires restart; hot swapping and distributed graph ownership are not implemented.
2. Road routing is available only for road vehicles. Other profiles use explicit device-local navigation or fail a Polaris-required request.
3. A 45-second local soak is architecture evidence, not a production SLO, maximum capacity, or long-duration leak certification.
4. Soak request failures are currently aggregated; performance gating should classify expected backpressure separately from unexpected failure.
5. Derived in-memory spatial state is rebuilt from Redis after restart, so brief post-restart spatial warm-up is expected.
6. At-least-once delivery can replay work after a commit failure; correctness depends intentionally on idempotent consumers and immutable identities.

## Final assessment

No known contradictory production path remains from the audited legacy spatial, routing, actor, rebalancing, or direct-command implementations. The current codebase has a coherent authority boundary: Kafka orders durable streams, PostgreSQL owns durable control state, Redis owns current reported state and leased connectivity, the gateway owns only fenced transport, and Mobility owns only derived domain computation. The remaining limitations are bounded and documented, and the regression suite is reproducible from Compose.

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


# Phase 4 capability and Mobility evidence

Captured on 2026-08-19 against baseline commit `59bd32c84ca505cdbd61ae671e6152c665828887`.

## 1. Objective and boundary

Phase 4 adds compile-time capability modules without changing Phase 3's execution authority. Modules propose candidates and plans; the core revalidates tenant, lifecycle, connectivity, capabilities, battery, project/type/distance constraints, and active assignment; PostgreSQL alone commits the reservation, command sequence, immutable payload, audit, and outbox event.

The core extension package contains only lifecycle, candidate, planner, registry, and twin-envelope contracts. It contains no H3, R-tree, OSM, KD-tree, A*, Dijkstra, map-matching, or traffic implementation. Composition is explicit in `cmd/engine`; there are no Go plugins, global registrars, or `init` discovery.

## 2. Extension contracts

- `Module`: named `Start`, `Ready`, and `Close` lifecycle with STARTING/READY/DEGRADED/FAILED/STOPPED states and component details.
- `CandidateProvider`: capability/type/project/tenant request, bounded ordered proposals, optional domain score and attributes.
- `TaskPlanner`: task plus typed twin, returning a versioned command payload, generation time, validity, and opaque metadata.
- `Registry`: explicit ordered construction, lifecycle fan-in, status exposure, planner fallthrough only when a planner reports that the selected device profile is unsupported.

`ROUTING_BUSY`, timeout, missing graph, no-route, and other real failures never fall through to a fabricated generic route.

## 3. Twin component model

Redis's existing freshness Lua script now writes `component:spatial/v1` and `component:battery/v1` atomically with the accepted boot/sequence and legacy reported-state view. Each component carries type, schema version, observed time, boot ID, sequence, and opaque JSON payload. The twin API returns a typed component map while keeping Phase 2 response compatibility.

## 4. Mobility spatial architecture

Mobility defines ROAD_VEHICLE, GROUND_ROBOT, AERIAL_DRONE, and STATIC profiles. Spatial states preserve reported and indexed positions separately, H3 cell, observation/index times, source boot/sequence, and explicit quality/anomaly output.

H3 resolution 8 is used downstream only. Tenant plus resolution-6 parent cells form internal shards. Same-device updates are striped; cross-shard transitions lock old/new shards in numeric order, remove old membership, insert new membership, and update the location map as one critical section. A device cannot be returned from two active regions.

The production exact index is a packed STR R-tree with geodesic post-filtering. Mutations mark a shard hierarchy dirty and the next read publishes a rebuilt hierarchy. A linear index is retained only as a correctness/benchmark oracle. Movement updates the index when the H3 cell changes, distance threshold is crossed, or indexed age expires; the twin's reported position still advances on every accepted observation.

Latitude/longitude validation, longitude normalization, Haversine distance, antimeridian bounding-box splitting, high-latitude cases, same-point behavior, heading normalization, invalid speed, implied-speed jumps, and stationary deadband are implemented and tested.

## 5. Version, lifecycle, and recovery

Projection updates require a newer sequence within one boot or a newer boot start. Phase 1 remains authoritative for retired-boot decisions. Duplicate/stale concurrent workers cannot regress Mobility state. STALE, OFFLINE, SUSPENDED, DECOMMISSIONED, suspended/deactivated tenant changes remove active state. Reconnection alone does not reinsert a device; a fresh accepted telemetry event does.

Restart rebuild scans Redis component state in bounded pages, revalidates active device and tenant lifecycle in PostgreSQL, sorts deterministic tenant/device input, and rebuilds derived H3/R-tree state. R-tree internals are not persisted.

## 6. Road topology and node index

The Mobility OSM loader keeps only nodes referenced by supported drivable ways, handles forward/reverse/roundabout one-way semantics, validates endpoints/distance/travel time, assigns road classes and speeds, and freezes topology after the builder completes. The corrected road-node set contains 690,268 nodes and 1,442,876 directed edges for `chennai-v1`; the earlier accidental 2.70M all-OSM-node index was detected and corrected during live proof.

`NodeType` is an ordinary enum, not a sequential value used as a bit mask. Multi-edge nodes are classified as intersections. Static nearest-road-node lookup uses a 3D unit-sphere KD-tree, which handles longitude wrapping without a linear graph scan.

## 7. Routing and cost snapshots

A* is the operational shortest/fastest search. Dijkstra uses the same immutable graph/cost contract as the correctness oracle. Distance uses geodesic admissible heuristics; fastest uses distance divided by graph maximum speed. Results validate edge direction and include route ID, policy, distance, duration, waypoints, edge IDs, expanded nodes, and one snapshot version.

Base topology never mutates. `RoutingCostSnapshot` edge costs are copied, validated, versioned, and atomically swapped. A request loads its pointer once. Partial or non-monotonic snapshots are rejected. Routing uses a bounded queue, fixed worker count, context cancellation, expansion limit, timeout, and a separate per-tenant semaphore. Saturation returns `ROUTING_BUSY`.

## 8. Map matching and traffic

The `polaris_traffic_group` consumer reads ordered telemetry independently. Road observations are matched against bounded incident segments using distance-to-segment plus heading compatibility; low-confidence matches are ignored. Per-edge EWMA speed, last speed/time, sample count, and confidence are maintained. Snapshot cost multipliers decay exponentially toward base cost, including a periodic refresh when no new telemetry arrives. Malformed events reach the Phase 1 DLQ before commit.

## 9. Candidate and planner flow

The Mobility provider executes H3 neighborhood narrowing, R-tree radius filtering, a bounded raw set, and route ETA for only the configured top K road candidates. Ordering is deterministic: routed ETA where available, direct distance fallback, battery during core revalidation, then stable device ID. Providers cannot reserve devices.

NAVIGATE and RELOCATE road plans contain route ID, route schema version, routing snapshot version, generated/valid-until timestamps, origin, destination, waypoints, distance, duration, and policy. Command expiry is the earlier of task expiry and plan validity. Retry reuses the stored command ID, sequence, and payload. Replanning is a new durable decision; an existing command is never silently mutated.

Unsupported drone/robot road planning explicitly declines the Mobility planner and preserves Phase 3's generic high-level command path. Low-level steering, motor control, collision avoidance, and flight attitude remain local to the edge device.

## 10. APIs, UI, configuration, and readiness

- `GET /api/v1/spatial/devices/nearby`
- `POST /api/v1/routes`
- compatible `GET /api/v1/routes/calculate`
- authoritative `POST /api/v1/tasks`

All are tenant scoped and operator protected. The dashboard distinguishes road, robot, drone, and static markers; the task table exposes route and snapshot metadata. Compose supplies explicit enablement, H3, radius/fanout, capacity, worker/queue, timeout/expansion, per-tenant, traffic-age, and graph-version/path bounds. Invalid configuration fails startup.

Live `/readyz` evidence:

```text
core: READY
mobility: READY
spatial: READY
routing: READY
road_graph_version: chennai-v1
road_nodes: 690268
road_edges: 1442876
routing_snapshot_version: 3
```

A missing graph produces Mobility DEGRADED with spatial READY and routing FAILED unless Mobility is mandatory. With Mobility disabled, the live proof passed telemetry, twin projection, CAPTURE_IMAGE planning, and durable generic command persistence.

## 11. Correctness and concurrency verification

Automated tests cover:

- randomized R-tree vs linear nearest-10 correctness;
- insert, upsert, move, remove, same region, antimeridian, high latitude, and same point;
- tenant isolation, monotonic versions, movement thresholds, eviction, and deterministic replay/rebuild;
- A* cost equals Dijkstra on deterministic graphs, one-way/no-route, nearest KD node, and NodeType regression;
- immutable snapshot replacement and concurrent readers using one version;
- explicit queue saturation, timeout, unsupported profile, invalid snapshot, and optional/mandatory graph failure;
- planner schema/snapshot/validity metadata and generic behavior with Mobility disabled;
- race detector over Mobility, extension, and orchestration packages.

Phase 3's live suite separately proves exclusive-assignment races, stable command replay, ACK/result idempotency, fencing, offline recovery, Kafka outage/outbox recovery, RBAC, and cross-tenant isolation. Its Kafka-restart evidence window is 60 seconds to account for broker readiness; the durable retry behavior itself is unchanged.

## 12. Measured benchmarks

Host: Windows/amd64, AMD Ryzen 5 5625U. These are measurements, not promises. The Windows clock rounded sub-microsecond samples to zero, so p50 `0 us` should be read as below clock resolution.

Heterogeneous in-memory distribution: 40% road, 20% robot, 25% static spatial, 15% non-spatial; 40% of mobile states moved; 1,000 radius/nearest-10 queries at 5 km.

| Total devices | Indexed | Index | p95 | p99 | Queries/s |
|---:|---:|---|---:|---:|---:|
| 1,000 | 850 | R-tree | 540 us | 578 us | 9,647 |
| 1,000 | 850 | linear | 1,003 us | 1,402 us | 7,089 |
| 5,000 | 4,250 | R-tree | 1,371 us | 1,580 us | 2,554 |
| 5,000 | 4,250 | linear | 1,754 us | 2,126 us | 1,208 |

On a 30x30 deterministic grid, A* measured 531,920 ns/op and 217,133 B/op; Dijkstra measured 459,066 ns/op and 253,562 B/op. On the real 690,268-node Chennai graph, A* was faster and expanded fewer nodes in short, medium, long, and dense-city samples; an edge-of-graph sample remained slightly slower. Dijkstra remains the cost oracle and no universal A* speedup is claimed. Full measurements are in `PHASE_4_1_CONSISTENCY_EVIDENCE.md`.

## 13. Live architecture proof

`phase4-mobility-test.ps1 -SkipLocalChecks` passed:

```text
Simulator -> Gateway -> Kafka -> Engine -> Redis -> PostgreSQL -> Dashboard
authenticated road telemetry -> H3/R-tree nearby query
POST /routes -> A* route + snapshot
POST /tasks NAVIGATE -> PostgreSQL immutable route command
```

Observed end-to-end telemetry/dashboard samples were 1,191 ms and 1,182 ms. PostgreSQL contained command `0361bb0f-58ca-486a-81d6-7a6164cb5e1f` for road vehicle `P4-ROAD-1787112089406`, route `4b005a45-455a-4dc4-8852-fa3198131428`, snapshot 2. Final Kafka evidence showed `polaris_engine_group` and `polaris_traffic_group` stable with `TOTAL-LAG 0` across all three telemetry partitions.

## 14. Regression results

- `go test ./...`: PASS
- `go vet ./...`: PASS
- focused `go test -race`: PASS
- frontend production Docker build: PASS (`tsc -b` and Vite)
- Phase 1 reliability Compose suite: PASS
- Phase 2 identity/security/twin Compose suite: PASS
- Phase 3 command orchestration Compose suite: PASS
- Phase 4 Mobility-enabled Compose proof: PASS
- Phase 4 Mobility-disabled Compose proof: PASS
- Phase 4.1 1,000-device mixed Compose soak: PASS

## 15. Acceptance matrix

| Area | Result | Evidence |
|---|---|---|
| Compile-time explicit extensions; Mobility disablement; generic tasks | PASS | registry and live isolation proof |
| Versioned optional twin components; non-spatial type | PASS | Redis Lua/API and compute-node fixture type |
| H3 sharding, R-tree, version guard, transition safety, eviction, rebuild, tenancy | PASS | implementation plus randomized/race tests |
| Immutable OSM topology, NodeType fix, KD-tree, A*, Dijkstra | PASS | graph implementation, unit and live route proof |
| Map matching, decayed traffic, immutable snapshots | PASS | traffic consumer/state/snapshot tests; group lag zero |
| Bounded workers, timeout/cancellation, per-tenant isolation | PASS | implementation and saturation tests |
| Candidate provider, top-K routing, core revalidation, PostgreSQL authority | PASS | service transaction boundary and Phase 3 races |
| NAVIGATE/RELOCATE plans and stable replay/staleness | PASS | live NAVIGATE and planner/Phase 3 replay tests |
| Module degradation and missing/corrupt graph semantics | PASS | optional/mandatory module tests |
| 1,000/5,000 heterogeneous spatial benchmark | PASS | measured harness results above |
| Simultaneous 1,000-device full Compose mixed telemetry/task/ACK workload | PASS | 1,000 connections, 7,691 telemetry frames, 120 durable executions, zero lag and safety invariants; Phase 4.1 evidence |
| Phase 1-3 regressions and build gates | PASS | results above |

## 16. Known limitations

- Global road routing is implemented for ROAD_VEHICLE only. Drones and ground robots use generic high-level command planning and spatial discovery; profile-specific global planners are future modules.
- The R-tree is a packed STR hierarchy rebuilt lazily after shard mutations. Mixed 80/20, 50/50, and 20/80 mutation/query benchmarks through 10,000 states passed; the design remains until longer production-shaped evidence justifies added copy-on-write complexity.
- The OSM graph is loaded in-process at startup and is not hot-swapped. Distributed shard ownership, multi-region routing, prediction, and autonomous rebalancing are outside Phase 4.
- A 45-second local 1,000-device run proves architectural invariants, not a production SLO, maximum capacity, or multi-hour memory-leak guarantee.
- The graph and snapshot update procedure is versioned but still requires a clean engine restart for graph replacement.

## Outcome

Polaris now manages domain behavior through explicit capability contracts while preserving its durable orchestration authority. Mobility provides tenant-isolated H3/R-tree spatial state, immutable road topology, bounded snapshot-consistent routing, map-matched decaying traffic costs, candidate proposals, and versioned route planning. Registry/twin/task/command behavior remains functional when Mobility is disabled or routing is degraded. The Phase 4.1 mixed soak closes the missing system proof; production scale claims remain deliberately limited to the recorded measurements.

## Phase 4.2 closure update

Phase 4.2 closes the remaining operational evidence gaps without changing these authorities. The load harness now distinguishes expected bounded outcomes from unexpected failures, exposes task and command stage timings, and records runtime/pool/queue state. A 1,000-device mixed run, a five-minute 51,000-frame stability run, bounded routing overload/recovery, and Mobility restart/rebuild all passed. The full Phase 0.5–4 Compose chain, Mobility-disabled mode, unit/vet checks, and focused race tests also passed.

The final regressions found and corrected three test/operability inconsistencies: a too-short independent archive checkpoint that omitted archive-group diagnostics, leaked device-profile environment between nested smoke tests, and request IDs that changed between audit and response. Detailed results and raw artifact names are in `PHASE_4_2_FINAL_STABILIZATION_EVIDENCE.md`. Phase 4 is now frozen at this boundary; its remaining limitations are explicit future scope rather than competing implementations.
