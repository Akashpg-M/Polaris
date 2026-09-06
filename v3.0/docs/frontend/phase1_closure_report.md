# Polaris Frontend Phase 1 Closure Report

## Outcome

Phase 1 replaces the original simulator-oriented prototype shell with a responsive, tenant-aware operational frontend grounded in the active Polaris Gateway and Engine contracts. It presents durable registry identity and current Redis-reported state as separate authorities, hydrates current state before subscribing to transient live updates, and labels every bounded or unavailable data surface honestly.

## Implemented capability

### Foundation and access

- Responsive desktop/mobile application shell with collapsible navigation, dark/light theme, tenant identity, persistent project context, stream health, and sign-out.
- Out-of-band operator access accepts a bearer token, tenant ID, and presentation role. Access is validated by reading the selected tenant before a session is created.
- Bearer tokens are stored in tab-scoped `sessionStorage`, never durable `localStorage`, and are removed on sign-out.
- Central role/permission vocabulary mirrors backend roles. It is presentation-only; backend middleware remains authoritative.
- Platform administrators can change to a known tenant ID. Other roles remain fixed to the authenticated tenant.

### API and state model

- Typed, same-origin API client for tenant, project, device, capability, twin, dashboard-ticket, and Gateway connection contracts.
- Automatic bearer, tenant, and request-ID headers; normalized response-envelope parsing; request cancellation; structured error/request-ID propagation; automatic session removal after HTTP 401.
- Lightweight tenant-keyed query cache with stale times, in-flight request deduplication, cancellation, subscriptions, targeted refresh, and cache isolation on tenant/session change.
- Deliberate loading, empty, unauthorized/error, retry, and bounded-result states.

### Fleet, twins, and projects

- Overview shows loaded fleet size, registry-active count, connectivity distribution, type distribution, project count, spatial preview, and session-only live activity.
- Devices provides server-side project/type/lifecycle/connectivity filtering, inferred forward cursor navigation, a previous-cursor stack, and clearly scoped loaded-page search.
- Device detail separates registry metadata, reported state, component envelopes, capabilities, and connectivity. Never-connected devices have an explicit state.
- Telemetry detail shows only the latest accepted observation because the backend has no history read API; the UI does not synthesize history from WebSocket traffic.
- Projects and project detail are read-only. Counts are explicitly described as the loaded bounded view.
- Twin cards and component panels preserve unknown component payloads as inspectable structured data, allowing schema evolution without pretending to understand module-owned fields.

### Live map and dashboard stream

- Current twin state is fetched from `/twins` before the map is rendered.
- A one-time dashboard ticket opens the tenant-bound WebSocket without putting the bearer token in the URL.
- Incoming messages are identity-, numeric-, coordinate-, battery-, timestamp-, and sequence-validated before cache projection.
- Known twin list/detail entries are patched incrementally. An event for an unknown loaded device triggers authoritative twin hydration.
- Reconnect uses bounded exponential backoff with jitter, obtains a new single-use ticket, and refreshes active twin queries after connection recovery.
- The Leaflet map maintains one imperative marker layer. It applies grid clustering at the current zoom, re-clusters after zoom changes, uses device-type glyphs and connectivity color, and opens a selected-device drawer.
- The disconnected state is explicit because Redis Pub/Sub has no replay and the Gateway contract has no application heartbeat.

## Architectural adaptations from the supplied brief

- No frontend login/current-user contract was invented. Operator credentials remain externally provisioned until backend session introspection exists.
- No tenant list was fabricated. Platform admins enter a known tenant ID.
- No global totals, charts, history, or durable activity were derived from partial pages. The UI labels all client summaries as bounded/loaded views.
- No additional server-state or map-cluster runtime package was required. A small query cache and grid clusterer keep the Phase 1 dependency surface narrow while retaining cancellation, deduplication, tenant isolation, and efficient imperative map updates.
- Legacy simulator, analytics, swarm, and task pages remain unmounted for future migration and are excluded from the Phase 1 lint surface. No Phase 1 navigation exposes them.

## Verification evidence

Executed on 6 September 2026:

| Check | Result |
| --- | --- |
| TypeScript project build (`tsc -b`) | Passed |
| ESLint on active Phase 1 source (`--max-warnings 0`) | Passed |
| Vitest unit suites | 4 files passed, 14 tests passed |
| Vite production bundle | Passed; 1,768 modules transformed |
| Production dependency audit | 0 known vulnerabilities after non-breaking transitive fix |
| Docker Compose runtime check | Not run; Docker Desktop daemon was unavailable |
| Automated browser visual pass | Not completed; the configured in-app browser runtime rejected its own trusted dependency path |

Unit coverage currently protects dashboard-frame validation, tenant presentation permissions, formatting of zero/missing telemetry, unknown device-type fallback, and case-insensitive loaded-page filtering.

## Confirmed remaining backend dependencies

The detailed gap analysis is in `phase1_backend_gaps.md`. The most important future contracts are operator session introspection, platform tenant discovery, truthful fleet aggregates, complete cursor metadata, pre-pagination connectivity filtering, telemetry history, global device search, and a durable read-oriented activity feed.

## Phase 2 reuse boundary

Phase 2 can reuse the shell, API/error conventions, auth/session boundary, permission vocabulary, project context, tenant-keyed query cache, fleet primitives, status components, twin cards, and live dashboard projection. Registry mutations should be added as separate feature modules with backend-derived roles, mutation-specific cache invalidation, confirmation/audit UX, and secret one-time-display handling. The dormant prototype pages should be migrated or removed when their corresponding product phases are implemented; they must not be reconnected directly to the Phase 1 shell.
