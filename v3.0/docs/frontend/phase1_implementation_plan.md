# Polaris Frontend Phase 1 — Grounded Implementation Plan

1. Replace the prototype shell with a responsive control-plane shell, protected routes, theme/sidebar state, tenant/project context, and centralized permission helpers.
2. Add a typed same-origin API client, normalized API errors, request IDs, cancellation support, and a tenant-aware query cache.
3. Add out-of-band operator access setup because the backend has no login/current-session endpoint. Validate access against the selected tenant and clearly label role as local presentation context.
4. Implement Overview, Devices, Device Detail, Digital Twins, Projects, Project Detail, and Live Map using the actual bounded twin/project contracts.
5. Hydrate current state from `/twins`, then connect through a one-time dashboard ticket. Patch tenant query caches per device, reconnect with exponential backoff, and authoritatively refresh after reconnect.
6. Retain Leaflet, add marker clustering only if a lightweight compatible dependency is justified, and otherwise use one imperative marker layer so telemetry does not rerender the React tree.
7. Add deliberate loading, empty, error, unauthorized, and forbidden states plus lifecycle/connectivity/battery/type primitives.
8. Add focused unit/component tests for parsing, permissions, formatting, generic components, filtering, and socket messages; then run lint, tests, and production container build.
9. Close with exact integration evidence, confirmed backend gaps, and Phase 2 reuse boundaries. No task, command, route, registry mutation, audit, or simulator UI will be exposed in this phase.
