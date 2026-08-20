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
