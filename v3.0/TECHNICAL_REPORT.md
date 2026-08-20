# Polaris v3.0 — Full Codebase Technical Assessment

> Historical baseline (18 August 2026). The pipeline and architectural defects described here were the input to Phases 0-4.1 and are no longer the current implementation state. See `CODEBASE_CONSISTENCY_REPORT.md` and `PHASE_4_1_CONSISTENCY_EVIDENCE.md` for the post-rectification assessment and verified results.

**Assessment date:** 18 August 2026  
**Repository assessed:** `Polaris/v3.0`  
**Scope:** Go backend, React frontend, data contracts, infrastructure configuration, tests, runtime wiring, and maintainability risks.

## 1. Executive summary

Polaris is aiming to be a real-time, multi-tenant spatial command-and-control platform for connected mobility assets such as drones, robots, bikes, cars, and static sensors. Its intended flow is: accept binary telemetry over WebSockets; partition and stream it through Kafka/Redpanda; maintain a low-latency in-memory spatial index; archive history in PostgreSQL/PostGIS; calculate congestion-aware road routes; predict demand hotspots; autonomously rebalance fleet assets; send commands back through Redis; and visualize activity in a React map and diagnostics console.

The repository contains credible prototypes for most individual ideas. The Go backend compiles, the two existing tests pass, the OSM routing graph and Dijkstra implementation are substantial, and the Protobuf telemetry contract is consistently represented in the Go load tester and browser encoder. However, the current runtime is not an integrated end-to-end system. The active gateway replaced Kafka output with a mock event publisher. Consequently, telemetry accepted by the gateway does not reach the engine, database, traffic analyzer, or dashboard Redis channel. The engine independently consumes Kafka, but the active gateway never writes to Kafka. This is the central defect and makes several headline features non-functional in the current wiring.

The repository also mixes v3 and v4 designs. The UI and folder are branded v3.0, while the active gateway and swarm page announce v4.0. Actor mailboxes, hysteresis filtering, a “cellular automata” safety gate, and a spatial projector appear to be a v4 experiment, but important pieces are synthetic or unused. Metrics presented in the v4 diagnostics UI are calculated with random numbers rather than measured from the backend.

**Overall assessment:** a technically ambitious proof of concept with a number of useful isolated components, but not yet a coherent or production-ready platform. The highest priority should be to choose one architecture, restore one end-to-end telemetry path, define observable contracts between services, and add integration tests that prove ingestion → projection → query → dashboard → command delivery.

## 2. Project intent and problem statement

The code suggests Polaris is designed to solve five related problems:

1. **High-volume telemetry ingestion.** Devices establish persistent WebSocket uplinks and send compact Protobuf frames containing identity, tenant, type, status, location, motion, battery, and timestamp.
2. **Real-time spatial discovery.** Dispatchers query for nearby assets by tenant and asset type using an in-memory spatial engine.
3. **Dynamic routing.** A Chennai OpenStreetMap road graph supports GPS snapping and Dijkstra routing, with weights intended to respond to live congestion.
4. **Fleet orchestration.** Demand zones are derived statically or from recent telemetry density, then under-supplied zones trigger relocation commands.
5. **Operations visibility and testing.** A browser dashboard shows live assets, heat, predicted zones, analytics, and a browser-based swarm generator; a Go load tester provides a larger CLI workload.

The architectural vocabulary—gateway, engine, event fabric, cold storage, command router, spatial partitions, actors, handover, predictive strategy—points toward a distributed smart-city/fleet-control platform rather than a simple vehicle tracker.

## 3. Repository structure

### 3.1 Backend

- `backend/cmd/gateway`: WebSocket telemetry ingress, dashboard WebSocket, Redis command/dashboard subscribers, and connection metrics.
- `backend/cmd/engine`: OSM graph loading, Kafka consumers, PostgreSQL archiving, traffic analysis, nearest-node API, route API, prediction API, and rebalancing loop.
- `backend/cmd/loadtest`: CLI WebSocket/Protobuf drone simulator with connection and throughput counters.
- `backend/api/proto/v1`: Protobuf definition and generated Go model.
- `backend/algo_/geo`: Haversine, bounding-box, and coordinate projection helpers.
- `backend/algo_/graph`: road network, congestion weights, nearest-intersection lookup, priority queue, and Dijkstra routing.
- `backend/algo_/quadtree`: thread-safe spatial container used by the engine.
- `backend/internal/core/actor`: per-asset actor, bounded mailbox, messages, events, and registry.
- `backend/internal/core/routing` and `simulation`: hysteresis, weight proposals, safety interception, and simulated risk runline.
- `backend/internal/application`: spatial engine, Kafka consumers/archiver/analyzer, demand strategies, rebalancer, and H3 handover manager.
- `backend/internal/adapter`: HTTP/WebSocket handlers and Kafka/mock/projector repositories.
- `backend/internal/infra`: OSM, Redis, and PostgreSQL construction.
- `backend/deployments`: Redpanda, Redis, PostGIS compose configuration and schema initialization.

### 3.2 Frontend

- React 19, TypeScript, Vite, Tailwind CSS 4.
- Leaflet and `leaflet.heat` for the map.
- Routes for the spatial map, analytics page, and swarm tester.
- A hand-written Protobuf encoder for browser-generated telemetry.
- Chart.js and Lucide are installed but not meaningfully used in the assessed UI.

### 3.3 Data and configuration

- A large Chennai `.osm.pbf` file is committed under `backend/data`.
- A root `.env` exists, while no root `.gitignore`, example environment file, or repository-level setup guide is present.
- The `chennai_map` directory exists but is empty in this version; an older commented Docker configuration still refers to it for OSRM.

## 4. Intended system architecture

The intended v3 data plane appears to be:

`Device/Simulator → Gateway WebSocket → Protobuf validation → Kafka telemetry.ingress → Engine consumers → spatial index + traffic graph + PostgreSQL archive → REST queries / predicted zones`

The intended command and visualization planes appear to be:

- `Engine rebalancer → Redis telemetry:commands → Gateway → device WebSocket`
- `Engine/projector → Redis spatial:updates → Gateway dashboard registry → browser /ws/dashboard`

The active code does not currently realize those flows. The gateway instead routes telemetry into an actor registry whose publisher is `MockEventPublisher`. That publisher only logs JSON. No active adapter publishes telemetry to Kafka or state updates to Redis. The legacy implementation that used `KafkaStreamAdapter` and a device `ConnectionRegistry` remains as a large commented block in `cmd/gateway/main.go` and `handler/ingestion.go`.

## 5. Implemented feature inventory

### 5.1 Features that are substantially implemented

**Binary telemetry contract.** `spatial.proto` defines the central `SpatialObject`. Both the Go load tester and the frontend simulator generate matching field tags and types. The gateway rejects non-binary WebSocket frames and invalid Protobuf payloads.

**Actor-based ingress isolation.** The active gateway creates one actor per first observed node ID. Each actor has a bounded, non-blocking mailbox and owns basic location, battery, task, and last-ping state. Saturated mailboxes drop frames rather than blocking ingress.

**Gateway connection count.** An atomic counter exposes active telemetry sockets through `/api/v1/metrics/connections`.

**Spatial engine data structure.** The engine stores assets across 32 hash shards and maintains a spatial search structure. Updates remove old coordinates and insert new ones. Queries filter by tenant, radius, and requested type, calculate Haversine distance, estimate ETA, sort results, and cap output at 500.

**OSM road ingestion and routing.** The engine reads the Chennai PBF, filters common drivable highway classes, handles basic one-way roads, calculates segment distances, builds an adjacency list, snaps coordinates to the nearest intersection, and runs congestion-weighted Dijkstra routing.

**Streaming infrastructure adapters.** Kafka publishing, Kafka batch consumption, PostgreSQL archiving, a dead-letter topic, and a traffic analyzer exist as code. The engine starts the Kafka consumer, archiver, and traffic analyzer.

**Historical hotspot strategy.** PostgreSQL telemetry is grouped into rounded latitude/longitude cells for the last hour, and the top three clusters become predicted zones.

**Autonomous rebalancer.** Every 15 seconds the rebalancer checks zone supply, finds assets within a wider radius, and publishes relocation commands for deficits.

**Dashboard map.** The frontend constructs a Chennai-centered dark Leaflet map, manages marker reuse, derives a heat layer, removes stale nodes, reconnects its dashboard WebSocket, and polls predicted zones.

**Load generation.** Both browser and Go load generators open multiple WebSockets and send binary Protobuf telemetry at one message per second per device.

**Graceful HTTP shutdown.** Gateway and engine HTTP servers listen for termination signals and perform bounded shutdown.

### 5.2 Features present only as partial, synthetic, or disconnected implementations

**Live dashboard telemetry.** The WebSocket registry and browser consumer exist, but no active component publishes actor state events to the subscribed Redis `spatial:updates` channel. The gateway's mock publisher only logs.

**Kafka ingestion from gateway.** `KafkaStreamAdapter` exists but is not used by the active gateway. Therefore the engine's `KafkaConsumer`, archiver, and traffic analyzer receive nothing from normal gateway traffic.

**Command delivery to devices.** The active gateway no longer stores device WebSocket connections by node ID. Its Redis subscriber sends commands into actors rather than down device sockets. Moreover, the publisher and subscriber disagree on JSON shape: the engine publishes `{node_id, command:{directive,...}}`, while the gateway parses top-level `{node_id,directive}`. The directive becomes empty.

**QuadTree indexing.** `SafeQuadTree` never subdivides. Every point stays in the root slice, so searches and removals are linear scans. It is thread-safe, but it is not currently a real quadtree and should not be described as high-performance spatial indexing without benchmarks against realistic data.

**Cellular automata safety simulation.** The runline sleeps for 3 ms and returns a random risk score. It contains no cellular grid, traffic state, transition rule, or deterministic model. The safety interceptor is not wired into routing or command issuance.

**Hysteresis-based traffic projection.** The filter logic itself is implemented, but `SpatialProjector` maps every event to one hard-coded edge and is not constructed by either executable.

**H3 handover.** The manager can emit and consume handover messages, but it is not constructed. Received state is inserted at default coordinates with an empty tenant, and restored directives are only logged.

**Predictive “ML.”** The strategy is a SQL density aggregation on rounded coordinates, not machine learning. It has hard-coded tenant, asset type, coverage, and required capacity.

**Frontend diagnostics.** Mailbox saturation is inferred from active connections; hysteresis and CA latency are random values; approval rate is a threshold heuristic. These are simulated presentation values, not backend telemetry.

## 6. Critical inconsistencies and defects

### Critical — end-to-end telemetry pipeline is broken

The gateway's active entrypoint constructs `NewMockEventPublisher()` and sends actor events only to logs. The engine consumes `telemetry.ingress`, but the gateway does not publish that topic. The consequences are systemic:

- nearest-node queries remain empty for gateway-originated devices;
- PostgreSQL history is not populated;
- predicted zones receive no fresh data;
- the traffic analyzer receives no telemetry;
- dashboard updates are never emitted to Redis;
- autonomous rebalancing sees no live fleet state.

### Critical — command path is both structurally incompatible and physically disconnected

The engine's `RedisCommander` wraps the command under a `command` property. The gateway subscriber expects a top-level `directive`. Even if corrected, the subscriber only updates the actor's `currentTask`; it has no retained hardware WebSocket through which to send the directive. The older `ConnectionRegistry` supports device writes but is no longer used.

### High — actor identity can diverge from payload identity

The handler fixes `nodeID` from the first frame, but each subsequent frame is passed unchanged into that actor. A client may send a different `payload.Id` later. The event uses the actor's original ID but the payload's other identity/tenant fields are not validated. Empty first-frame IDs also create or repeatedly access an empty-ID actor. There is no authentication, ownership proof, tenant authorization, coordinate validation, size limit, rate limit, or heartbeat policy.

### High — spatial type filtering is semantically wrong

The Protobuf `NodeType` values are sequential enums (1, 2, 3, …), not bit flags. The quadtree query uses `(p.Class & reqClass) > 0`. This creates false matches: for example class 3 and requested class 2 pass the bitwise test. Exact enum equality is required unless the model is intentionally changed to powers-of-two flags.

### High — multi-tenant and query security is superficial

The match API treats any caller-supplied `tenant_id` query string as identity and returns that tenant's data. There is no authentication or authorization. Static and predictive zones hard-code `alpha_logistics`. PostgreSQL indexes are not tenant-prefixed, and node IDs are not uniquely constrained per tenant.

### High — engine correctness and lifecycle risks

- `BatchUpdate` stores caller-owned Protobuf pointers rather than copies.
- The quadtree bounds cover a fixed regional rectangle; out-of-bounds inserts silently fail while the shard map still stores them.
- The Kafka consumer only flushes exact batches of 1,000; low traffic or shutdown leaves the partial batch unapplied.
- Kafka read errors loop immediately without logging/backoff and can cause a hot failure loop.
- The engine requires Redis and panics if it is unavailable, while other dependencies degrade gracefully.
- Kafka readers/writers and database connections are not uniformly closed.
- The engine starts a predictive strategy connection and a separate archiver database connection without explicit pool/lifecycle configuration.

### High — routing and traffic weights are coarse or inefficient

- GPS snapping scans every OSM node for every request, despite comments calling this acceptable; city-scale graphs can make this expensive.
- A traffic observation updates every outgoing segment from the nearest intersection, not the actual traversed road segment.
- Congestion uses a single speed observation and fixed thresholds; there is no decay, aggregation window, direction, vehicle class, or minimum sample count.
- OSM parsing recognizes only `oneway=yes`, omits `-1`, roundabouts, access restrictions, speed limits, turn restrictions, and several road classes.
- Holding the graph read lock for the entire Dijkstra traversal blocks all congestion writes during a route computation.

### High — frontend data contract mismatch

Actor events serialize fields as `asset_id`, `lat`, `lon`, etc. The map expects `node.id`. Even if the actor event were published to Redis, markers would be keyed under `undefined`. The popup also assumes numeric optional values without a validated runtime schema.

### Medium — v3/v4 identity and architecture collision

The repository directory, map sidebar, containers, and report target say v3.0. The gateway log and swarm tester say v4.0. The proto `go_package` includes `/v3.0/`, while the Go module and generated file use `github.com/Akashpg-M/polaris/backend`. This confuses release identity and regeneration reproducibility.

### Medium — misleading metrics and terminology

The UI labels random values as live hysteresis, CA latency, and safety approval metrics. Comments describe a “zero-allocation” browser encoder although it allocates an ArrayBuffer and byte arrays. A SQL group-by is described as ML. The quadtree does not subdivide. The scale test is a functional smoke test, not a benchmark or scale test.

### Medium — error handling and observability gaps

Many meaningful errors are ignored: JSON marshal, Redis URL parsing, actor push from command subscriber, Kafka DLQ writes, congestion updates, shutdown errors, and database row iteration errors. Frontend fetch/WebSocket parse failures are usually silent. There are no health, readiness, version, Kafka lag, actor count, mailbox depth, dropped-frame, archive-failure, or command-delivery metrics.

### Medium — resource retention and concurrency concerns

Actors are never removed after sockets disconnect, so actor count and goroutines grow with every unique ID. Dashboard broadcasts hold a global read lock while performing network writes; one slow client can delay the whole broadcast and registration lifecycle. Reconnection timers in the map are not cancelled on component unmount, so a closed component can reconnect later.

### Medium — repository hygiene and documentation

- Root `.env` is present and no `.gitignore` is visible, creating a credential/configuration leakage risk.
- The only README is the unchanged Vite template; there is no architecture, setup, port, environment, schema, topic, or runbook documentation.
- Large superseded implementations remain commented inside source files.
- Typographical path/file names include `algo_`, `registory.go`, and `rounting.go`.
- Generated Protobuf is committed but no generation command/toolchain is documented.
- The PBF data file is stored directly in the repository without provenance or update instructions.

## 7. Testing and verification findings

### 7.1 Backend verification

`go test ./...` completed successfully after relocating the Go build cache into the writable workspace. Only two packages contain tests:

- `internal/core/actor`: one asynchronous actor event/recovery test.
- `internal/core/routing`: one combined smoke/performance-style test.

All other packages reported “no test files.” The current tests do not cover HTTP handlers, Protobuf rejection, tenant separation, Kafka/Redis/PostgreSQL adapters, quadtree behavior, OSM parsing, graph routing correctness, command schema, graceful shutdown, or end-to-end flow.

The actor recovery test manually writes historical values into private fields; there is no production replay implementation. Its publisher spy is not synchronized, so it would be unsafe under broader concurrency. The routing scale test sends only ten messages and uses timing assertions on a synthetic 3 ms random simulation; it is not evidence of stated large-scale capacity.

### 7.2 Frontend verification

The frontend could not be compiled or linted in the assessment environment because `npm` is not installed/available. This is an environment limitation, not a confirmed repository build failure. Static review nonetheless identified:

- hard-coded localhost URLs in Analytics and Swarm Tester while MapDashboard supports environment variables;
- the Analytics component is named `SwarmTester` internally;
- `leaflet.heat` is suppressed with `@ts-ignore` despite a local declaration file;
- broad `any` use for the heat layer;
- mojibake/corrupted emoji text visible throughout source;
- no frontend tests;
- large commented legacy component code in `SwarmTester.tsx`;
- silent catches that mask failed network calls.

### 7.3 Operational verification not performed

No full runtime smoke test was possible without starting Redis, Redpanda, PostgreSQL/PostGIS, the gateway, the engine, and frontend together. The static wiring defect means the intended flow would remain incomplete even if all dependencies started.

## 8. Security assessment

The current services should be treated as development-only.

- WebSocket origin checks allow every origin.
- CORS allows every origin.
- No HTTP or WebSocket authentication exists.
- Tenant identity is client-controlled.
- Node identity is client-controlled and not bound to a socket credential.
- No TLS termination or secure WebSocket configuration is defined.
- No payload size, rate, connection, or per-tenant quota is enforced.
- Redis, Kafka, and PostgreSQL compose services expose host ports with development credentials and no network segmentation.
- The root environment file may expose secrets if committed or shared.
- Dashboard popups interpolate device IDs into HTML, which can enable stored/streamed XSS if untrusted IDs reach the UI.
- Error/DLQ headers may contain raw database error text, potentially exposing schema details.

## 9. Performance and scalability assessment

The architecture contains useful scalability ideas—bounded actor mailboxes, sharded maps, asynchronous messaging, Kafka partitions, batching, and spatial filtering—but their current implementation does not substantiate high-scale claims.

- Per-device actors consume a goroutine and a mailbox of capacity 5,000. With thousands of devices, the reserved channel storage and permanent actor retention can be material.
- The “quadtree” is O(n) for search and removal.
- Each update performs a linear remove from the point slice, so sustained updates trend toward O(n) work per ping.
- Nearest-road-node snapping is O(V) per route and per traffic message.
- PostgreSQL inserts are one row per message with no prepared statement or bulk copy.
- Dashboard broadcast performs serial socket writes under a shared registry lock.
- Kafka batch consumption lacks time-based flushing.
- The browser stress tool creates one timer per simulated device and is unsuitable for serious load testing.

The Go CLI load generator is a better foundation, but it lacks latency histograms, server acknowledgements, reconnect behavior, payload validation, coordinated ramp-down, and automated pass/fail criteria. Its ramp calculation can divide by zero and becomes a zero-duration ticker for rates above 1,000.

## 10. Code quality and maintainability

The backend generally follows recognizable layers and uses interfaces in several useful places. Naming and comments often explain intent. However, architecture changes have been accumulated rather than completed: old implementations are commented out, new components are unconnected, versions conflict, and comments frequently claim production behavior that code does not deliver.

Recommended quality improvements include:

- remove superseded commented code after preserving it in version history;
- apply `gofmt` and consistent import grouping;
- rename misspelled files and ambiguous packages;
- replace ignored errors with deliberate handling and metrics;
- separate experimental v4 features behind a package, branch, or feature flag;
- document API/event schemas and generate both Go and TypeScript models from Protobuf;
- add structured configuration validation and an `.env.example`;
- establish lint, test, race, and build checks in CI.

## 11. Recommended target architecture

A minimal coherent v3 should use the existing event-driven design without mixing the experimental actor projection path until it is fully integrated:

1. Gateway validates authenticated device telemetry, stamps receipt time, and publishes the original `SpatialObject` to Kafka using an H3 key.
2. Engine consumes Kafka with time/size bounded batches and updates the live spatial index.
3. A separate consumer archives telemetry to Postgres/PostGIS and exposes lag/failure metrics.
4. A projection consumer publishes a normalized `DashboardNodeUpdate` JSON contract to Redis or a dedicated Kafka topic.
5. Gateway fans dashboard updates out through per-client buffered writers.
6. Rebalancer publishes a versioned `DeviceCommandEnvelope`.
7. Gateway retains authenticated device sessions and delivers commands, reporting acknowledgement/failure.
8. Actor and safety features are added only where they own an explicit production responsibility and have deterministic tests.

The alternative is to commit fully to the actor architecture. In that case, the actor publisher must be a real composite adapter that publishes telemetry/state to Kafka and dashboard projection channels, actors require lifecycle management and replay, and the engine should not maintain a separate contradictory source of truth.

## 12. Prioritized remediation roadmap

### Phase 0 — clarify and stabilize (1–2 days)

- Decide whether this repository is v3 or v4 and align all labels, module/package metadata, UI text, and containers.
- Add a root README, `.gitignore`, `.env.example`, architecture diagram, startup commands, ports, topics, and expected data flow.
- Remove secrets from tracked configuration and rotate any exposed credentials.
- Mark synthetic features and metrics clearly as simulations.

### Phase 1 — restore the vertical slice (3–7 days)

- Replace the gateway mock publisher with a real Kafka publisher, or implement a composite actor publisher.
- Define and use versioned telemetry, dashboard-update, and command envelopes.
- Restore a node-to-WebSocket device registry and command delivery acknowledgements.
- Normalize dashboard event field names (`id` versus `asset_id`).
- Add one Docker-based integration test proving: telemetry sent → Kafka consumed → match API returns node → archive row written → dashboard update received → command delivered.

### Phase 2 — correctness and security (1–2 weeks)

- Replace enum bitwise matching with exact equality or true flag definitions.
- Validate IDs, tenant, enums, coordinates, velocity, battery, frame size, and timestamp.
- Add authentication and server-derived tenant identity.
- Restrict CORS/WebSocket origins and document TLS deployment.
- Add actor eviction, socket deadlines, ping/pong, reconnect policy, and dependency health/readiness.
- Fix partial Kafka batch flushing and retry/backoff behavior.

### Phase 3 — spatial/routing quality (1–3 weeks)

- Implement actual quadtree subdivision or use an established R-tree/H3 index.
- Spatially index OSM nodes and map telemetry to road segments, not intersections.
- Add traffic aggregation, direction, decay, sample thresholds, and stale-weight reset.
- Handle OSM access, oneway variants, roundabouts, speeds, and turn restrictions as required.
- Add deterministic routing and spatial-index test fixtures.

### Phase 4 — observability and production readiness (1–2 weeks)

- Export Prometheus/OpenTelemetry metrics for throughput, latency, drops, mailbox depth, Kafka lag, database failures, active actors/sockets, route performance, and command acknowledgements.
- Replace UI-generated metrics with measured backend values.
- Add structured health endpoints, dashboards, alerts, resource limits, and capacity tests.
- Add CI for Go format/vet/test/race, Protobuf regeneration checks, TypeScript build, ESLint, frontend tests, container builds, and integration tests.

### Phase 5 — advanced features

- Implement deterministic safety simulation from actual road/traffic state or remove the cellular-automata claim.
- Implement durable actor recovery and handover ownership semantics if distributed actors remain a goal.
- Replace rounded-coordinate hotspot counts with a documented forecasting method, evaluation metrics, tenant-aware parameters, and model/version observability.

## 13. Suggested acceptance criteria for a production-oriented v3 milestone

- A documented one-command local startup works from a clean checkout.
- A generated client and server share the same Protobuf/schema source.
- A device cannot read or affect another tenant.
- A telemetry frame appears in live match results and dashboard within a defined latency SLO.
- The same frame is durably archived, or failure is retried/DLQ'd and observable.
- A relocation command reaches the intended device and produces an acknowledgement.
- Dependency loss produces bounded retries and correct readiness state, not silent loops or panics.
- Spatial type matching, movement updates, stale eviction, and routing are covered by deterministic tests.
- Metrics shown in the UI come from backend instrumentation.
- Load tests publish reproducible throughput, latency percentiles, error rate, CPU, and memory results.

## 14. Final conclusion

Polaris v3 has a strong prototype concept and several valuable building blocks: Protobuf WebSocket telemetry, a Go routing graph, event-stream adapters, PostGIS history, actor isolation, fleet rebalancing, and a useful operational UI direction. Its primary problem is not absence of code; it is architectural discontinuity. The current executables wire together two different generations of the design, leaving the main data and command paths incomplete.

The fastest path forward is to stop adding advanced features temporarily and make one narrow vertical slice real, observable, secure, and tested. Once a device update can reliably traverse the entire platform and a command can reliably return, the existing routing, prediction, simulation, and handover ideas can be evaluated and integrated on a stable foundation.

## Appendix A — verification record

- Backend command: `go test ./...`
- Result: pass; 2 packages with tests, remaining packages without tests.
- Frontend commands intended: `npm run build`, `npm run lint`
- Result: not executed because `npm` was unavailable in the assessment environment.
- Runtime integration: not executed; required infrastructure was not started, and static analysis found the active gateway-to-engine pipeline disconnected.

## Appendix B — highest-impact issue checklist

- [ ] Replace active gateway mock publisher.
- [ ] Restore Kafka telemetry publication.
- [ ] Restore live dashboard update publication.
- [ ] Define one command envelope and deliver to device sockets.
- [ ] Fix `NodeType` bitwise filtering.
- [ ] Align `asset_id`/`id` dashboard contract.
- [ ] Replace random UI metrics with real metrics.
- [ ] Implement real spatial subdivision/indexing.
- [ ] Add actor eviction and socket lifecycle controls.
- [ ] Add authentication and server-derived tenancy.
- [ ] Add end-to-end integration test and root documentation.
- [ ] Resolve v3/v4 branding and architecture split.
