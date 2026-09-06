# Polaris Frontend Phase 1 — Confirmed Backend Gaps

## 1. Operator session introspection

- **Frontend requirement:** Establish authenticated user, exact role, and tenant without trusting manually entered claims.
- **Available:** Bearer authentication and authorization on every protected route.
- **Missing:** A read-only endpoint returning the resolved operator principal and permitted tenant context.
- **Why it matters:** The frontend can validate access but cannot independently discover or verify the role used for role-aware presentation.
- **Suggested minimal change:** `GET /api/v1/session` returning API-key ID, role, tenant ID, and whether explicit tenant selection is required.
- **Priority:** High before wider mutation UX.

## 2. Tenant discovery for platform administrators

- **Frontend requirement:** Select an authorized tenant context.
- **Available:** Read one known tenant and explicitly scope a platform administrator using `X-Tenant-ID`.
- **Missing:** Tenant listing/search.
- **Why it matters:** A platform administrator must type a known tenant ID; the UI cannot offer a trustworthy selector.
- **Suggested minimal change:** Paginated `GET /api/v1/tenants` restricted to platform administrators.
- **Priority:** Medium.

## 3. Fleet aggregates

- **Frontend requirement:** Accurate registered/lifecycle/connectivity/type/project totals.
- **Available:** Bounded twin/device pages.
- **Missing:** Tenant-level aggregates and total counts.
- **Why it matters:** Client aggregation describes only loaded pages and must not be presented as complete for fleets over the page limit.
- **Suggested minimal change:** One tenant-scoped `GET /api/v1/fleet/summary` query with lifecycle, connectivity, type, and project counts.
- **Priority:** Medium; high for fleets over 100 devices.

## 4. Complete cursor metadata

- **Frontend requirement:** Reliable forward/back navigation and knowledge of whether another page exists.
- **Available:** Exclusive device-ID cursor and bounded arrays.
- **Missing:** `next_cursor`, `has_more`, and total count.
- **Why it matters:** The frontend can infer a possible next cursor only when a full page is returned and must maintain its own previous-cursor stack.
- **Suggested minimal change:** Return `{ items, next_cursor, has_more }` inside the data envelope.
- **Priority:** Medium.

## 5. Twin connectivity filtering semantics

- **Frontend requirement:** Page through all ONLINE/OFFLINE/etc. twins predictably.
- **Available:** `connectivity_status` filter on `/twins`.
- **Missing:** Connectivity filtering before the registry page limit is applied.
- **Why it matters:** A filtered page can contain fewer results even when matching devices exist after that registry page.
- **Suggested minimal change:** Persist/query connectivity through an indexed projection or perform a bounded join that filters before pagination.
- **Priority:** Medium.

## 6. Telemetry history API

- **Frontend requirement:** Durable recent telemetry table and battery/speed history.
- **Available:** PostgreSQL/PostGIS telemetry archive and live/current twin state.
- **Missing:** Tenant-scoped read API for telemetry history.
- **Why it matters:** Live WebSocket observations are transient and cannot truthfully represent history.
- **Suggested minimal change:** Cursor/time-bounded `GET /api/v1/devices/:device_id/telemetry` with maximum range and limit.
- **Priority:** Medium. Phase 1 shows the latest observation and clearly marks history unavailable.

## 7. Device/project search and device-type catalog

- **Frontend requirement:** Search beyond the currently loaded page and render backend-owned type labels.
- **Available:** Exact/cursor device listing, project listing, capability catalog, and stable seeded device type IDs.
- **Missing:** Search parameters and device-type catalog endpoint.
- **Why it matters:** Local search covers only loaded items; hard-coded labels require safe fallback for future types.
- **Suggested minimal change:** Bounded `q` search for device ID/display name and `GET /api/v1/device-types`.
- **Priority:** Low to medium.

## 8. Durable recent fleet activity

- **Frontend requirement:** Activity visible after browser reload.
- **Available:** Transient tenant-filtered dashboard telemetry and administrator-only audit events.
- **Missing:** Read-oriented tenant fleet activity feed for all read roles.
- **Why it matters:** Phase 1 can only show events observed during the current browser session.
- **Suggested minimal change:** Cursor-paginated projection of non-sensitive device connectivity/telemetry lifecycle summaries.
- **Priority:** Low.

