# Phase 2 — Identity, Registry, Tenancy, and Digital-Twin Evidence

Verified on 18 August 2026 from `C:\Users\akash\PROJECTS\Polaris\v3.0` with Docker Desktop and the repository Compose stack.

## Result

The reproducible Phase 2 proof passed:

```text
PASS: Simulator -> Gateway -> Kafka -> Engine -> Redis -> PostgreSQL -> Dashboard
      (SMOKE-1787060757926, 3656 ms)
PASS: basic
PASS: rejected
PASS: send
PASS: revoke-session
PASS: rejected
PASS: send
PASS: ticket
PASS: rejected
PASS: rejected
PASS: rejected
PASS: Phase 2 authenticated registry, credential lifecycle, tenant isolation,
      outbox, audit and digital-twin flow
```

Run it again with:

```powershell
cd C:\Users\akash\PROJECTS\Polaris\v3.0
powershell -NoProfile -ExecutionPolicy Bypass `
  -File .\backend\deployments\phase2-identity-test.ps1
```

The script creates a fresh high-entropy development platform key for each run, rebuilds and waits for Compose, provisions registry resources, exercises identity failures and credential lifecycle, checks persistence and Kafka lag, runs the Go tests, and builds the frontend. It does not print issued secrets.

## Implemented architecture

### Durable control-plane model

PostgreSQL owns tenants, projects, device types, devices, capabilities, device-capability configuration, device and operator credential hashes, one-time ticket hashes, audit events, and transactional outbox events. Device deletion is intentionally absent: lifecycle transitions retain historical identity.

The implemented device lifecycle is:

```text
REGISTERED -> ACTIVE -> SUSPENDED -> ACTIVE
     |           |          |
     +-----------+----------+----> DECOMMISSIONED (terminal)
```

Tenant status is `ACTIVE`, `SUSPENDED`, or `DEACTIVATED`. Gateway authentication requires both tenant and device to be active.

### Authentication and identity binding

- Device and operator tokens use `pol_<kind>_<public-prefix>.<secret>` and at least 256 random secret bits.
- PostgreSQL stores the lookup prefix and SHA-256 hash, never the bearer secret. Verification uses constant-time comparison.
- Raw credentials are returned only by issue/rotate operations. Metadata listings cannot recover them.
- Non-browser clients authenticate the WebSocket handshake with `Authorization: Bearer`.
- Browser simulators and dashboards use separately scoped, short-lived, single-use, hash-stored tickets.
- The gateway resolves a `DevicePrincipal` before WebSocket upgrade, rejects payload/principal disagreement, and constructs Kafka identity and partition key from the principal.
- Active telemetry connections revalidate credential, device, and tenant state on every frame. Revocation or suspension therefore closes the session before another event can be accepted.

This deliberately avoids an authentication cache in Phase 2. Per-frame database validation is simpler and gives deterministic revocation; a bounded cache can be introduced later with explicit invalidation if profiling demonstrates a need.

### Authorization and tenant isolation

Operator roles are `PLATFORM_ADMIN`, `TENANT_ADMIN`, `OPERATOR`, and `VIEWER`. Middleware enforces the permission matrix before handlers run. Platform administrators must explicitly select a tenant using `X-Tenant-ID`; all other principals derive scope from their stored tenant. Repository reads and writes take tenant ID and include it in SQL predicates. Cross-tenant resource access returns 404 and is safely audited.

The same protection covers registry, twin, match, routing, predictive-zone, and dashboard access. Dashboard connections resolve an operator ticket to one tenant, and Redis dashboard messages are delivered only to sockets for that tenant.

### Transactional lifecycle events and audit

Security-sensitive mutations execute their domain change, audit insert, and outbox insert in one PostgreSQL transaction. The embedded relay claims bounded batches using `FOR UPDATE SKIP LOCKED`, publishes `device.lifecycle.v1`, then marks rows published. Publish-success/mark-failure can replay, so the contract remains at least once and events carry stable IDs.

Audited operations include tenant and project creation/status, device registration/lifecycle, capability changes, credential issue/rotation/revocation, and cross-tenant denial. Authorization headers, raw secrets, and telemetry payloads are excluded.

### Digital twins and connectivity

Twin responses combine tenant-scoped PostgreSQL device metadata and capabilities with `polaris:twin:{tenant_id}:{device_id}` in Redis. A registered device without Redis state returns `reported_state: null` and `NEVER_CONNECTED`.

The Phase 1 atomic Redis projection also sets connectivity to `ONLINE` and updates the `polaris:devices:last-seen` sorted-set index. A configurable detector performs guarded `ONLINE -> STALE -> OFFLINE` transitions and publishes deterministic connectivity events. Fresh accepted telemetry returns the twin to `ONLINE`.

## API inventory

All routes below are under `/api/v1` and require an operator bearer key:

| Area | Routes |
|---|---|
| Tenants | `POST /tenants`, `GET /tenants/:tenant_id`, `PATCH /tenants/:tenant_id` |
| Projects | `POST /projects`, `GET /projects` |
| Devices | `POST /devices`, `GET /devices`, `GET /devices/:device_id`, lifecycle actions |
| Capabilities | catalog list plus get/put/delete per device |
| Credentials | issue, metadata list, revoke, atomic rotate |
| Tickets | device connection ticket and tenant-scoped dashboard ticket |
| Twins | `GET /devices/:device_id/twin`, paginated `GET /twins` |
| Security | `GET /audit-events` (administrator roles only) |
| Spatial | match, route, and predictive-zone reads |

Device and twin lists use a maximum page size of 100 and query PostgreSQL first; they never use Redis `KEYS`.

## Verification matrix

| Requirement | Evidence |
|---|---|
| Durable tenant/project/device registry | Resources survive service recreation; PostgreSQL counts verified |
| Types and capabilities | Seed catalog, FK-backed device type, assigned capability query |
| Secret handling | `plaintext_token_hashes=0`; issue/list/rotate behavior exercised |
| Unknown device blocked | Invalid token handshake rejected before upgrade |
| Suspended/decommissioned blocked | Device and tenant state scenarios rejected |
| Server-derived identity | Spoofed tenant/device frame closes socket; no spoof state created |
| Rotation and revocation | Old token rejected, new token accepted, active session rejected after revoke |
| Operator RBAC and isolation | Protected APIs reject unauthenticated access; tenant-B twin is hidden from tenant A |
| Coherent digital twin | Registry metadata/capabilities and Redis reported state verified together |
| Connectivity lifecycle | `NEVER_CONNECTED`, `ONLINE`, `STALE`, `OFFLINE`, and return online exercised |
| Recoverable registry events | Outbox rows reach `PUBLISHED`; migrations are repeatable |
| Audit | Security mutations and denial events persisted |
| Phase 0.5/1 preservation | Authenticated full-path smoke plus Go and reliability suites pass |

## Live operational evidence

Immediately after the final proof:

```text
gateway /healthz  {"status":"live"}
gateway /readyz   {"status":"ready"}
engine  /healthz  {"status":"live"}
engine  /readyz   {"status":"ready"}

tenants=2
projects=6
devices=6
credentials=14
audit=60
outbox_published=60
plaintext_token_hashes=0

polaris_engine_group  TOTAL-LAG 0 (all 3 telemetry partitions)
polaris_archive_group TOTAL-LAG 0 (all 3 telemetry partitions)
```

Kafka topics present with three partitions are `telemetry.ingress`, `telemetry.dead-letter.v1`, `device.lifecycle.v1`, and `device.connectivity.v1`. Reapplying `init.sql` twice completed successfully and preserved existing telemetry.

## Regression evidence

- `go test ./...`: pass, including token/RBAC and lifecycle unit tests.
- Phase 1 live Redis/PostgreSQL integration and multi-partition reliability suite: pass.
- Authenticated Phase 0.5 full-path smoke: pass at 3656 ms end-to-end for the final Compose run.
- Frontend production build in Compose: pass.
- Gateway and engine Compose health/readiness: pass.

The Phase 2 script uses five/eight-second connectivity thresholds to keep verification fast; Compose defaults remain configurable. Dedicated Prometheus metric export and external identity-provider integration are intentionally deferred: neither is needed to establish the Phase 2 security boundary, and omitting them keeps this phase inside the existing gateway/engine architecture.

## Phase 2 outcome

Polaris now manages authenticated, tenant-owned devices through a durable registry and capability model. The gateway derives identity from revocable credentials rather than trusting telemetry payloads, registry changes publish reliable lifecycle events through a transactional outbox, and PostgreSQL metadata combines with Redis-reported state to expose tenant-isolated digital twins with online, stale and offline lifecycle tracking.
