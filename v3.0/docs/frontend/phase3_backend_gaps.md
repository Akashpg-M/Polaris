# Polaris Frontend Phase 3 — Confirmed Backend Gaps

## 1. Traffic edge visualization

- **Requirement:** Render meaningful current congestion on road geometry.
- **Current backend support:** Internal per-edge EWMA state and immutable route-cost snapshots; readiness exposes only snapshot version, edge-state count, and overlay byte estimate.
- **Missing capability:** Tenant-safe edge geometry/ID plus current speed, multiplier, confidence, sample count, and observation time.
- **Frontend impact:** Traffic shows truthful snapshot metadata and no heatmap.
- **Minimal proposed backend addition:** A bounded viewport endpoint joining current traffic state to simplified geometry, freshness, and confidence.
- **Priority:** Medium.

## 2. Durable route history

- **Requirement:** Show recent route activity and reopen exploratory routes.
- **Current backend support:** Route calculation returns a result object; command planning persists routes only inside command payloads.
- **Missing capability:** Tenant-scoped persisted exploratory route/history endpoint.
- **Frontend impact:** Standalone results live only while the page is mounted. Command routes reopen by command ID.
- **Minimal proposed backend addition:** Optional bounded route-result store keyed by tenant/request with retention and cursor pagination.
- **Priority:** Low.

## 3. Spatial aggregates and complete listing

- **Requirement:** Accurate tenant-wide spatial device/profile totals.
- **Current backend support:** Point/radius nearby query and bounded registry twin pages.
- **Missing capability:** Tenant aggregate counts and cursor-based spatial-state listing.
- **Frontend impact:** Overview labels its up-to-100 twin map as bounded and never claims tenant totals or a profile distribution.
- **Minimal proposed backend addition:** Tenant-scoped spatial summary plus cursor-based spatial-state list.
- **Priority:** Medium.

## 4. Runtime query limits

- **Requirement:** Validate nearby radius and result limit against actual deployment configuration.
- **Current backend support:** Server enforcement; repository defaults are 10 km and 50 results.
- **Missing capability:** Configured maximum radius and result count in Mobility readiness details.
- **Frontend impact:** The form uses repository defaults and the backend remains authoritative if deployment values differ.
- **Minimal proposed backend addition:** Return these non-sensitive limits in module details.
- **Priority:** High for configurable deployments.

## 5. Spatial index diagnostics

- **Requirement:** Explain shard/index health and reported-versus-indexed state.
- **Current backend support:** Nearby results expose H3 cell, both positions/times, quality, and source sequence. Internal R-tree statistics exist.
- **Missing capability:** Shard ID, per-index health/capacity/rebuild count, and H3 IDs encoded safely as strings.
- **Frontend impact:** Selected-result details show only exposed state. A uint64 H3 JSON number is not guaranteed exact in JavaScript.
- **Minimal proposed backend addition:** Encode H3/shard IDs as canonical strings and expose aggregate non-sensitive index status.
- **Priority:** Medium.

## 6. Routing runtime detail

- **Requirement:** Diagnose workers, active searches, timeouts, latency, and per-tenant limiter saturation.
- **Current backend support:** Process-local requests/busy, queue depth/capacity, and active tenant-limiter count.
- **Missing capability:** Worker count, active requests, timeout count, latency distribution, and safe limiter occupancy.
- **Frontend impact:** Diagnostics shows only returned counters and infers nothing from configuration.
- **Minimal proposed backend addition:** Extend readiness or add an authenticated aggregate diagnostics endpoint.
- **Priority:** Medium.

## 7. Predicted-zone tenant authority

- **Requirement:** Tenant-isolated experimental density zones.
- **Current backend support:** Last-hour telemetry is grouped on a rounded grid, but SQL is not tenant-filtered and returned tenant identity is hard-coded.
- **Missing capability:** Tenant predicate derived from authenticated scope and correct aggregate tenant identity.
- **Frontend impact:** The feature is prominently experimental and suppresses every row whose tenant ID differs from the active tenant. This cannot repair the backend query.
- **Minimal proposed backend addition:** Filter SQL by authenticated `tenant_id` and remove the hard-coded identity.
- **Priority:** Critical before production use.

## 8. Persisted route diagnostic parity

- **Requirement:** Reopen command routes with exploratory-result diagnostics.
- **Current backend support:** Command payload persists ID/schema, graph/snapshot, validity, origin/destination, waypoints, distance, duration, and policy.
- **Missing capability:** Expanded-node count and edge IDs are not persisted.
- **Frontend impact:** Mobility displays exact available geometry and labels missing diagnostics rather than reconstructing them.
- **Minimal proposed backend addition:** Persist the fields only if their operational value justifies payload growth.
- **Priority:** Low.

## 9. Explicit registry mobility profile

- **Requirement:** Select a device as a route origin and enforce its exact profile.
- **Current backend support:** Spatial results expose profile; ordinary twins expose device type but no authoritative mobility profile.
- **Missing capability:** Profile in twin/registry reads or direct tenant-scoped spatial state by device.
- **Frontend impact:** Routing uses an explicit profile control and does not infer profile from a display name or offer registry devices as origins.
- **Minimal proposed backend addition:** Add profile to the twin or expose `GET /spatial/devices/:id`.
- **Priority:** Medium.
