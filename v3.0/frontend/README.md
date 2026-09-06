# Polaris Control Plane Frontend

The Polaris frontend is a tenant-scoped operational control plane. Phase 1 exposes the registered fleet and digital twins. Phase 2 adds durable task and immutable command orchestration. Phase 3 adds canonical spatial discovery, bounded traffic-aware road routing, and truthful Mobility diagnostics without crossing into registry mutation, credential management, audit administration, or simulator tooling.

## Product surfaces

- Operator access using an out-of-band bearer token and known tenant ID
- Responsive application shell with persistent theme, sidebar, tenant, and project context
- Bounded fleet overview with lifecycle/connectivity/type summaries
- Cursor-based device inventory with project, type, lifecycle, connectivity, and loaded-page search filters
- Device detail with registry metadata, current state, twin components, capabilities, and an explicitly current-only telemetry view
- Read-only project list and project fleet detail
- Digital-twin cards that distinguish registry lifecycle from online/stale/offline/never-connected state
- API-hydrated Leaflet fleet map with lightweight grid clustering and a device detail drawer
- Tenant-bound dashboard WebSocket updates with single-use tickets, safe parsing, cache patching, exponential reconnect, and authoritative refresh after reconnect
- Role-aware task list, task creation wizard, task detail, cancellation, and administrative orchestration retry
- Command list/detail with per-device ordering, immutable payload and route-plan inspection, ACK/result state, conservative cancellation, and administrative delivery retry
- Confirmed-timestamp task/command timelines, request correlation metadata, and bounded current Operations activity
- Controlled REST polling for active orchestration because the telemetry dashboard socket does not carry task/command lifecycle events
- Mobility Overview, canonical Nearby Search, shortest/fastest Route Explorer, Traffic metadata, Routing Diagnostics, and explicitly experimental density analytics
- Shared Leaflet primitives across Live Map and Mobility, plus immutable command-route inspection without recalculation

The UI never treats browser role selection as an authorization boundary. The backend still authenticates and authorizes every request. The role field is temporary presentation context because the current backend does not expose session introspection.

## Run the complete application

From the repository root with Docker Desktop running:

```powershell
.\backend\deployments\start.ps1
```

Open `http://localhost:5173`. Use the `DEV_PLATFORM_ADMIN_TOKEN` generated in `backend/deployments/.env`, select `PLATFORM_ADMIN`, and enter a known tenant ID. A fresh database has no tenant fleet; the full smoke test creates an `alpha_logistics` demo tenant and device:

```powershell
.\backend\deployments\smoke-test.ps1
```

Normal deployment operations:

```powershell
.\backend\deployments\status.ps1
.\backend\deployments\logs.ps1
.\backend\deployments\verify-deployment.ps1
.\backend\deployments\stop.ps1
```

Nginx serves the single-page app and routes `/api/engine/*`, `/api/gateway/*`, and `/ws/*` to internal Compose services. No container image contains a baked-in host `localhost` API address.

## Frontend development

Requirements: Node.js 22 and npm.

```powershell
cd frontend
npm ci
npm run dev
```

The Vite development proxy expects Gateway on `127.0.0.1:6080` and Engine on `127.0.0.1:6081`.

Run the quality gates with:

```powershell
npm run lint
npm test
npm run build
```

## Contract and scope documentation

- [`../docs/frontend/frontend_contract_inventory.md`](../docs/frontend/frontend_contract_inventory.md)
- [`../docs/frontend/phase1_backend_gaps.md`](../docs/frontend/phase1_backend_gaps.md)
- [`../docs/frontend/phase1_implementation_plan.md`](../docs/frontend/phase1_implementation_plan.md)
- [`../docs/frontend/phase1_closure_report.md`](../docs/frontend/phase1_closure_report.md)
- [`../docs/frontend/phase2_implementation_plan.md`](../docs/frontend/phase2_implementation_plan.md)
- [`../docs/frontend/phase2_backend_gaps.md`](../docs/frontend/phase2_backend_gaps.md)
- [`../docs/frontend/phase2_closure_report.md`](../docs/frontend/phase2_closure_report.md)
- [`../docs/frontend/phase3_implementation_plan.md`](../docs/frontend/phase3_implementation_plan.md)
- [`../docs/frontend/phase3_backend_gaps.md`](../docs/frontend/phase3_backend_gaps.md)
- [`../docs/frontend/phase3_closure_report.md`](../docs/frontend/phase3_closure_report.md)
