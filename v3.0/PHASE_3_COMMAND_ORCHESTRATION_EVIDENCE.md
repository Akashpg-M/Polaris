# Phase 3 — Durable Task and Command Orchestration Evidence

Verified on 18 August 2026 against the existing Phase 1/2 PostgreSQL and Kafka volumes using Docker Desktop and the repository Compose stack.

## Outcome

Polaris now turns tenant-scoped task intent into durable, capability-validated, exclusively assigned device commands. Commands are committed before delivery, ordered per device, published through the transactional outbox, routed through leased and fenced gateway ownership, acknowledged and completed over authenticated bidirectional WebSockets, and recovered through bounded retry, expiry, reconnect, and database reconciliation. Delivery remains at least once; stable command identity and simulator deduplication make replay safe without claiming exactly-once execution.

## Reproduce

```powershell
cd C:\Users\akash\PROJECTS\Polaris\v3.0
powershell -NoProfile -ExecutionPolicy Bypass `
  -File .\backend\deployments\phase3-command-test.ps1
```

Final proof output:

```text
PASS: Simulator -> Gateway -> Kafka -> Engine -> Redis -> PostgreSQL -> Dashboard
      (SMOKE-1787070493948, 1404 ms)
PASS: complete
PASS: duplicate
PASS: offline
PASS: fencing
PASS: wrong-ack
PASS: capability-mismatch
PASS: receive-no-ack
PASS: resume
PASS: Phase 3 durable task assignment, fenced delivery, idempotent ACK/result,
      retry, expiry, recovery, RBAC and tenancy flow
```

## Durable data model

Phase 3 adds the following idempotent migration objects:

| Table | Responsibility |
|---|---|
| `tasks` | Business intent, requirements, target, lifecycle, actor, timestamps, expiry, version |
| `device_assignments` | Explicit exclusive device reservation and task lease |
| `device_command_sequences` | Atomic monotonically increasing sequence allocation per tenant/device |
| `commands` | Durable canonical command, attempts, ACK/result state, expiry, version |
| `command_attempts` | Delivery attempt, gateway, fencing epoch, timestamps, outcome |

`uq_device_active_assignment` is a partial unique index over `(tenant_id, device_id)` for active reservations. Commands are unique by both server-generated `command_id` and `(tenant_id, device_id, sequence_number)`.

Mutable task and command rows carry versions, while transitions use status predicates and row locks. Terminal states cannot return to active states through ordinary lifecycle operations. Administrative task retry creates a new command decision; delivery retry preserves the existing command ID and sequence.

## Assignment policy

An eligible device must satisfy all of the following:

```text
same tenant
AND lifecycle ACTIVE
AND Redis connectivity ONLINE
AND every required capability enabled
AND battery >= minimum
AND allowed device type/project constraints
AND distance <= maximum when configured
AND no active exclusive assignment
```

Candidates are ranked by shortest spatial distance, then highest battery, then stable device ID. PostgreSQL uniqueness and transaction locking make the final reservation race-safe even when concurrent selectors saw the same candidate.

Command types have an explicit capability mapping, including `RELOCATE -> receive_relocation_command`, `NAVIGATE -> navigate`, and `CAPTURE_IMAGE -> capture_image`. Unknown command types are rejected rather than bypassing capabilities.

## Persist-before-delivery transaction

Successful assignment performs the following durable work before any WebSocket write:

```text
BEGIN
  lock task
  reserve device
  allocate per-device command sequence
  insert command PENDING
  transition task to ASSIGNED
  insert audit events
  insert task/command outbox events
COMMIT
```

The outbox relay prioritizes actionable command rows, publishes actual commands to `device.command.v1`, and retains its existing at-least-once mark-after-publish semantics. Task, ACK, and result observations use separate lifecycle topics.

## Gateway ownership and fencing

An authenticated device connection atomically claims:

```text
polaris:connection:{tenant}:{device}
```

The record contains gateway ID, connection ID, credential ID, connected time, lease expiry, and a monotonically increasing epoch. The gateway refreshes the lease, releases only its own epoch, and closes an older local socket when a newer connection supersedes it.

Delivery verifies the epoch both before its durable delivery transition and immediately before the volatile socket write. This second check prevents a reconnect that races with PostgreSQL from allowing the stale owner to send. Credential, tenant, and device status are also revalidated before every delivery.

Kafka is the durable notification source and Redis Pub/Sub routes to the current gateway. Redis Pub/Sub is not treated as durable: connected-device polling and connect-time database reconciliation independently query outstanding commands in sequence order.

## Bidirectional protocol

Telemetry remains binary Protobuf for Phase 0–2 compatibility. Control frames have an explicit canonical envelope:

- Server to device: `COMMAND`, schema version, command/task/tenant/device IDs, sequence, timestamps, correlation/causation IDs, typed payload.
- Device to server: `COMMAND_ACK` with `ACCEPTED`, `REJECTED`, `DUPLICATE`, `EXPIRED`, or `UNSUPPORTED`.
- Device to server: `COMMAND_RESULT` with separate success/failure outcome and safe result data.

ACK/result identity is derived from the authenticated connection. A device cannot acknowledge another device's command. Duplicate ACKs and results resolve by command ID and are harmless. Stale ACKs never update a newer command.

The browser simulator checks expiry and sequence, caches completed command IDs, executes a command only once, and resends the prior ACK/result when the same ID is delivered again.

## Retry, expiry, cancellation, and reconciliation

- ACK timeout returns `DELIVERED` to `PENDING` with bounded exponential backoff.
- Retries preserve command ID and sequence while increasing `attempt_count`.
- `max_attempts` and `expires_at` are externalized; the first reached makes the operation terminal.
- Offline commands stay durable and are selected when the device reconnects.
- Pending commands are selected strictly by lowest outstanding device sequence.
- Command completion/failure reconciles its task and releases the assignment.
- Pending tasks are periodically reconsidered when an eligible device appears.
- Pending tasks/commands expire without late delivery.
- Pending work can be cancelled locally. Delivered or acknowledged work returns conflict rather than pretending the physical action stopped.
- Tenant/platform administrators can explicitly retry failed tasks or force a failed/delivered command retry.

## API and RBAC surface

| Operation | Platform admin | Tenant admin | Operator | Viewer |
|---|---:|---:|---:|---:|
| Create/cancel task | Yes | Yes | Yes | No |
| Read tasks/commands | Yes | Yes | Yes | Yes |
| Retry task/command | Yes | Yes | No | No |

Routes:

```text
POST /api/v1/tasks
GET  /api/v1/tasks
GET  /api/v1/tasks/:task_id
POST /api/v1/tasks/:task_id/cancel
POST /api/v1/tasks/:task_id/retry
GET  /api/v1/commands
GET  /api/v1/commands/:command_id
POST /api/v1/commands/:command_id/retry
POST /api/v1/commands/:command_id/cancel
GET  /api/v1/metrics/orchestration
```

All repository queries include tenant scope. A tenant-A task queried with a tenant-B key returned 404. A viewer task mutation returned 403.

## Failure matrix exercised

| Scenario | Verified behavior |
|---|---|
| Basic flow | Task assigned, command persisted/published/delivered, ACK, result, task completed |
| Capability mismatch | Online device without `capture_image` was not assigned; task expired |
| Offline device | Task remained pending; reconnect telemetry caused assignment and delivery |
| Lost ACK | Same command ID/sequence retried; simulated execution count remained one |
| Gateway crash | Unacknowledged command survived restart and completed after reconnect |
| Ownership race | Second connection received a higher epoch; stale socket closed |
| Wrong-device ACK | Connection closed and target command remained unchanged |
| Duplicate ACK/result | Idempotent current/terminal-state handling |
| Expiry | Expired task/command was terminal and unavailable to reconciliation |
| Cancellation | Pending task cancelled and reservation released |
| Kafka unavailable | Task/command transaction committed; outbox remained recoverable and later published |
| Outbox/Kafka replay | Replayed command event produced no duplicate command row |
| Cross tenant / viewer | Resource hidden across tenant; viewer mutation forbidden |

## Captured operational state

After proof and regression runs:

```text
gateway /healthz  live
gateway /readyz   ready
engine  /healthz  live
engine  /readyz   ready

tasks:    COMPLETED=16, EXPIRED=6, CANCELLED=2
commands: COMPLETED=16, EXPIRED=3, PENDING=0
device assignments=19
delivery attempts=21
command audit events=77
task audit events=45
published command outbox events=77

polaris-command-dispatcher TOTAL-LAG=0
polaris_engine_group       TOTAL-LAG=0
polaris_archive_group      TOTAL-LAG=0
```

All command/lifecycle topics have three partitions:

```text
device.command.v1
device.command.ack.v1
device.command.result.v1
task.lifecycle.v1
```

## Regression evidence

- `go test ./...`: pass.
- `go vet ./...`: pass.
- Frontend TypeScript/Vite production build: pass.
- Phase 1 replay, idempotency, DLQ, Redis-failure, manual-commit, multi-partition shutdown suite: pass.
- Phase 2 credential, spoofing, revocation, suspension, tenant isolation, twin, audit and outbox proof: pass.
- Final authenticated Phase 0.5 path after Phase 3: pass (`SMOKE-1787070776889`, 1283 ms in the Phase 2 regression run).
- Idempotent Phase 3 migration reapplied twice without changing prior telemetry or registry history.

## Deliberate scope choices

Phase 3 stays inside the existing engine and gateway executables. It does not introduce a workflow/BPMN engine, a new microservice, Redis command durability, an external identity provider, or exactly-once claims. Initial tasks produce one deterministic command; the task/command separation and sequence model allow short workflows in a later phase without redesigning delivery.

The WebSocket uses binary Protobuf for telemetry and explicit JSON control envelopes for commands, ACKs, and results. This avoids breaking Phase 0–2 clients while still making frame type, schema version, identity, sequence, and expiry unambiguous.

