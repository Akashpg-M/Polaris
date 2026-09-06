# Polaris Frontend Phase 3 Closure Report

## Outcome

Polaris now presents Mobility as a coherent spatial operations capability: an operator can inspect module state, choose a geographic point, query the canonical tenant-aware spatial index, inspect reported versus indexed device positions, calculate shortest or fastest road routes, render returned geometry, inspect graph/snapshot diagnostics, and reopen an immutable route persisted in a Phase 2 command without recalculation.

The product boundary remains explicit: Mobility proposes and plans, Core revalidates eligibility, and PostgreSQL commits assignment authority.

## Usable routes

- `/mobility` — module state, bounded spatial-twin map, graph/snapshot facts, and workflow entry points.
- `/mobility/nearby` — coordinate/map selection, radius/limit validation, canonical nearby query, map/table results, spatial detail drawer, and device deep-link.
- `/mobility/routes` — road-vehicle route form, shortest/fastest policies, comparison, returned geometry and metadata, and differentiated failures.
- `/mobility/routes?command_id=...` — exact persisted command route inspection with no route request or mutation.
- `/mobility/traffic` — shared-trusted traffic model and snapshot metadata; no fabricated edge heatmap.
- `/mobility/diagnostics` — module/components, queue pressure, request/busy counters, limiter count, graph and snapshot state.
- `/mobility/experimental` — tenant-filtered presentation of the legacy predicted-zone density heuristic with prominent authority warnings.

The Phase 1 `/fleet/map` route now consumes the same `PolarisMap` primitive used by Mobility.

## Exact endpoint inventory

Phase 3 calls only:

- `GET /api/engine/spatial/devices/nearby`
- `POST /api/engine/routes`
- `GET /api/engine/zones/predicted`
- `GET /api/engine/readyz`, proxied by Vite/Nginx to Engine `/readyz`
- existing `GET /api/engine/devices/:id/twin` for a selected nearby device
- existing `GET /api/engine/commands/:id` for persisted command-route inspection

It does not call compatibility `/nodes/match` or legacy `/routes/calculate`.

## Spatial flow

```text
operator coordinate + radius + limit
  → authenticated tenant-scoped nearby endpoint
  → H3 region narrowing
  → packed R-tree radius query
  → exact Haversine verification and stable ordering
  → returned indexed/reported state
  → shared map + accessible table + selected-device twin hydration
```

Nearby is a point-in-time query. Editing its request clears the prior authoritative presentation; “Search again” explicitly repeats even an unchanged request. The UI does not silently add candidates from WebSocket traffic or call them task-eligible.

## Routing flow

```text
origin + destination + ROAD_VEHICLE + policy
  → bounded per-tenant routing admission
  → queue/worker execution
  → KD-tree road-node snapping
  → A* over immutable Chennai graph
  → one immutable distance/traffic-cost snapshot
  → route geometry + graph/snapshot/work metadata
  → shared map + textual result
```

`FASTEST` uses the route-cost snapshot; `SHORTEST` minimizes distance. Unsupported profiles are rejected in the form and remain backend-enforced. Busy, timeout, unavailable, no-route, no-road-node, outside-region, unsupported-profile, planner, and unknown failures have separate safe presentation with request ID when returned. No straight-line, third-party, or estimated fallback is drawn.

Policy comparison performs two independent requests and does not pretend both results share a snapshot. Standalone results are retained only in page state and are never polled or silently recalculated.

## Persisted command integration

Route-backed Command Detail now links to Mobility using `command_id`. Mobility reloads the tenant-scoped command, parses its stored route document, and renders only the persisted origin, destination, waypoints, distance, duration, policy, graph, and snapshot. Missing edge IDs or expansion counts are labelled unavailable. No `POST /routes` occurs in this inspection path.

## Traffic and diagnostics truthfulness

Readiness exposes module/component state, graph size/version, traffic snapshot/scope/refresh, aggregate traffic state count, queue depth/capacity, cumulative requests/busy responses, and active limiter count. These are process runtime values, not durable history.

No browser API exposes current per-edge geometry, cost, speed, confidence, samples, or observation time. Traffic therefore explains the shared-trusted model and immutable snapshots but intentionally provides no road congestion layer. Mobility degradation stays local to Mobility pages and does not turn registry, twins, telemetry, or generic orchestration into a global fatal state.

## Experimental analytics

Predicted zones are labelled a legacy density heuristic: last-hour telemetry grouped on an approximately 1.1 km rounded grid, with fixed radius/asset assumptions and no autonomous action. Because the backend query is not tenant-filtered and emits a hard-coded tenant, the browser suppresses non-matching rows and prominently reports that the endpoint is not production-safe. It is never labelled ML or AI prediction.

## Verification evidence

Executed on 6 September 2026:

| Check | Result |
| --- | --- |
| Strict ESLint (`--max-warnings 0`) | Passed |
| Vitest | 8 files passed; 37 tests passed |
| TypeScript + Vite production build | Passed; 1,793 modules transformed |
| Route-level code splitting | Passed; all six Mobility pages emitted as independent chunks |
| Compose expansion | Passed (`docker compose ... config --quiet`) |
| Frontend Docker image | Passed; `deployments-frontend` built from Node 22 and Nginx 1.27 stages |
| Nginx syntax | Passed inside the built image with Compose service aliases |
| Container health | `GET /healthz` returned 200 |
| Deep SPA fallback | `GET /mobility/routes` returned 200 and the application root document |

Tests cover coordinate/radius/limit validation, metres and Go-duration conversion, profile support, stale/offline labels, routing error mapping and request identity, persisted-route parsing, module degradation, navigation, route metadata, and immutable-route messaging. Existing Phase 1–2 tests remain green.

A live authenticated backend/browser Mobility E2E was not run because this task did not start or mutate the full deployment. The built image, proxy configuration, route fallback, compiler, tests, and static container checks are reproducible without claiming live domain data was exercised.

## Confirmed gaps

The detailed list is in `phase3_backend_gaps.md`. Highest priority items are tenant-authoritative predicted-zone SQL, dynamically exposed nearby limits, and—if a map layer is desired—a bounded traffic edge/geometry API. Other gaps include route history, tenant spatial aggregates, exact string H3 IDs, richer routing runtime counters, persisted-route diagnostic parity, and registry-readable mobility profile.

## Reusable foundation

Phase 3 leaves shared map primitives, coordinate fields, profile badges, module status, route summary/error presentation, query keys, and route parsing ready for later observability and diagnostics work. It does not begin Phase 4 registry administration, credentials, audit, simulators, or full system-health tooling.
