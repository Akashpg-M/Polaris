# Polaris Frontend Phase 2 Closure Report

## Outcome

Polaris now exposes durable orchestration as a first-class operator workflow. Authorized users can create capability-aware tasks, inspect assignment and planning decisions at the level the backend exposes, follow immutable commands through delivery/ACK/result state, navigate between task, command, and device identities, conservatively cancel pre-delivery work, and invoke clearly differentiated administrative retries.

The frontend preserves the backend model: a task is operator intent; a command is the immutable device-specific instruction produced only after eligibility, exclusive reservation, and planning succeed.

## Pages

- `/tasks` — bounded cursor task list with state and exact assigned-device filters.
- `/tasks/new` — structured five-step Task, Target, Requirements, Planning, and Review workflow plus creation result.
- `/tasks/:taskId` — lifecycle, metadata, requirements, target, assignment, related commands, and actions.
- `/commands` — bounded cursor command list with state, exact task, and exact device filters.
- `/commands/:commandId` — sequence, immutable payload, persisted route metadata, lifecycle, attempts, ACK, result, and actions.
- `/operations/activity` — bounded current durable state refreshed while mounted and explicitly not presented as historical audit.

Operations navigation is enabled in the Phase 1 shell. Viewers see Tasks, Commands, and Activity. Operator and administrator presentation roles also see Create Task. Only tenant/platform administrator presentation roles see task or command retry controls.

## API integration

Phase 2 uses:

- `POST /api/engine/tasks`
- `GET /api/engine/tasks`
- `GET /api/engine/tasks/:task_id`
- `POST /api/engine/tasks/:task_id/cancel`
- `POST /api/engine/tasks/:task_id/retry`
- `GET /api/engine/commands`
- `GET /api/engine/commands/:command_id`
- `POST /api/engine/commands/:command_id/retry`
- `POST /api/engine/commands/:command_id/cancel`

Every call retains Phase 1 bearer, tenant, request-ID, structured-error, cancellation, and automatic 401 session handling. Critical mutations are never optimistic and are never automatically retried.

## Task creation

The wizard transforms presentation fields into the exact backend DTO:

```text
structured form
  → client UX validation
  → project/task/priority/expiry
  → target document
  → requirements in backend units
  → explicit review
  → POST /tasks
  → durable Task plus optional assigned Command and request timing
```

Supported command types use the backend capability mapping. The automatic capability is displayed as locked while operator-added capabilities remain distinct. Battery is validated as 0–100 percent; distance is entered in kilometres and sent as metres; spatial coordinates use backend bounds. `DEVICE_LOCAL` and `POLARIS_REQUIRED` are sent only for NAVIGATE/RELOCATE. The UI never claims candidate availability before the backend evaluates and reserves it.

A successful `PENDING` result is shown as durable work waiting for eligibility, not as a failure. Request-scoped candidate, routing, persistence, and total timing is shown only when the mutation response returns it.

## Lifecycle and delivery semantics

Task and command timelines use only current durable state and explicit timestamps. Missing intermediate timestamps are labelled unconfirmed; no timeline is filled using assumptions.

Command Detail makes the per-device sequence prominent and explains head-of-line ordering. `DELIVERED` is a delivery attempt, `ACKNOWLEDGED` is durable ACK advancement, and `COMPLETED` requires a successful result. Aggregate attempts, latest send, next availability, ACK classification, result payload, and latest error are displayed exactly as returned.

Route-backed command payloads show route/schema/graph/snapshot/policy, origin/destination, generation/validity, distance, duration, and waypoint count when present. Payloads are read-only and explicitly immutable.

Task retry is labelled a new orchestration decision that may create a new assignment and command. Command retry is labelled another delivery decision for the same immutable command, sequence, and payload. Cancellation controls are disabled outside known pre-delivery states, while backend 409 remains authoritative under races.

## Real-time and refresh strategy

The dashboard WebSocket carries telemetry only. Phase 2 therefore uses:

- 10-second polling only while task/command list or Activity pages are mounted;
- 5-second polling for active Task/Command Detail;
- no further entity-detail network polling after a terminal state is observed;
- targeted task/command query refresh after confirmed mutations;
- no global one-second polling and no full fleet-cache flush.

Activity combines the latest bounded task/command read models and labels itself `Current durable snapshot`. It does not pretend to be audit or lifecycle history.

## Permissions

- Viewer: read Tasks, Commands, Activity, payloads, and lifecycle.
- Operator: viewer access plus task creation and legal pre-delivery cancellation.
- Tenant/platform administrator: operator access plus administrative task and command retry.
- Platform tenant switching continues through the Phase 1 tenant scope.

Because the backend still lacks session introspection, the selected frontend role remains presentation context. Backend middleware is the security boundary and may reject controls if the selected presentation role does not match the credential.

## Testing evidence

Executed on 6 September 2026:

| Check | Result |
| --- | --- |
| TypeScript (`tsc -b`) | Passed |
| ESLint (`--max-warnings 0`) | Passed |
| Vitest | 6 files passed; 26 tests passed |
| Vite production bundle | Passed; 1,782 modules transformed |
| Route-level code splitting | Passed; Operations pages emitted as independent chunks |
| Docker Compose configuration | Passed (`docker compose ... config --quiet`) |
| Frontend container build | Blocked before source build; Docker Hub OAuth connections for the Node/Nginx base images were forcibly closed |
| Docker Compose runtime | Not run; a fresh frontend image could not be pulled/built because of the registry network failure |
| Browser E2E | Not run in this environment |

Phase 2 tests cover task/command status labels, priority presentation, backend capability inference, DTO/unit conversion, form validation, transition-action guards, duplicate ACK explanation, confirmed-only timelines, requirements separation, immutable route payload presentation, and role permissions. The existing backend `phase3-command-test.ps1` remains the repository’s real service/simulator orchestration proof; it was not rerun because the required base-image registry access failed.

## Backend gaps and known limitations

The complete confirmed list is in `phase2_backend_gaps.md`. Most important are:

- no browser task/command lifecycle stream or durable Operations activity API;
- no command-attempt list API;
- no candidate preview or structured rejection diagnostics;
- narrow task/command filters and incomplete cursor metadata;
- orchestration timing is not durable/readable from Task Detail;
- command schema and transport observations are absent from the command read model;
- no mutation idempotency key for ambiguous network timeouts.

The UI does not provide standalone routing, nearby search, traffic, registry mutations, credentials, audit, failure centre, health centre, or simulator tooling.

## Phase 3 readiness

The next Mobility-focused frontend phase can reuse route payload inspection, coordinate inputs, task planning modes, entity links, lifecycle/status primitives, permissions, confirmations, tenant query keys, active/terminal polling, and the responsive Operations layouts. A future route explorer should remain separate from inspection of an already-persisted command plan.
