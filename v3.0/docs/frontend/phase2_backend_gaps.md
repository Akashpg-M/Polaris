# Polaris Frontend Phase 2 — Confirmed Backend Gaps

## 1. Browser task and command lifecycle stream

- **Requirement:** Near-live task and command state updates and a unified Operations activity feed.
- **Existing backend support:** Durable lifecycle mutations, transactional outbox events, and Kafka topics. The dashboard WebSocket carries tenant-filtered telemetry only.
- **Missing:** Browser-consumable task/command lifecycle stream or durable read-oriented activity endpoint.
- **UX impact:** Mounted Operations views use controlled REST polling. Activity shows current bounded state, not event history.
- **Minimal change:** Add a tenant-bound, resumable Operations stream backed by a durable projection and event cursor; do not reuse transient telemetry Pub/Sub as history.
- **Priority:** High.

## 2. Command delivery-attempt read model

- **Requirement:** Explain each attempt, Gateway, fencing epoch, completion/failure, and retry schedule.
- **Existing backend support:** `command_attempts` is persisted; command records expose aggregate `attempt_count`, `max_attempts`, latest `sent_at`, `available_at`, and `last_error`.
- **Missing:** Tenant-scoped attempt-list endpoint and attempt result fields in the command response.
- **UX impact:** The frontend truthfully shows aggregate attempts only.
- **Minimal change:** `GET /api/v1/commands/:id/attempts`, tenant-scoped and ordered by attempt number.
- **Priority:** High for delivery diagnosis.

## 3. Candidate preview and rejection diagnostics

- **Requirement:** Preview likely candidates or explain why none qualify.
- **Existing backend support:** Authoritative eligibility, bounded domain ranking, current-state checks, and deterministic reservation during assignment.
- **Missing:** Side-effect-free candidate preview and structured per-rule elimination counts.
- **UX impact:** Review explains submitted constraints but cannot claim candidates or identify the rejecting requirement.
- **Minimal change:** Optional bounded preview returning aggregate exclusion reasons without reserving a device or leaking cross-tenant identity.
- **Priority:** Medium.

## 4. Task list filters

- **Requirement:** Filter by type, priority, project, and created range.
- **Existing backend support:** Server filters only for `status` and exact assigned `device_id`.
- **Missing:** Task type, priority, project, and date filters.
- **UX impact:** Phase 2 exposes only filters that remain truthful across the dataset.
- **Minimal change:** Add indexed query parameters with bounded date ranges.
- **Priority:** Medium.

## 5. Command list filters

- **Requirement:** Filter by command type, project, and created range.
- **Existing backend support:** Server filters for `status`, exact `task_id`, and exact `device_id`.
- **Missing:** Command type, project, and date filters.
- **UX impact:** Project would require an N+1 task join in the browser, so it is omitted.
- **Minimal change:** Add indexed type/date filters and a tenant/project join in the repository query.
- **Priority:** Medium.

## 6. Cursor metadata and totals

- **Requirement:** Know whether another task/command page exists and display reliable totals.
- **Existing backend support:** Exclusive entity-ID cursor with bounded arrays.
- **Missing:** `next_cursor`, `has_more`, and total.
- **UX impact:** The frontend infers a possible next page only from a full page and maintains a local previous-cursor stack.
- **Minimal change:** Return `{items,next_cursor,has_more}`.
- **Priority:** Medium.

## 7. Durable orchestration timing

- **Requirement:** Inspect candidate, routing, persistence, and total timing later from Task Detail.
- **Existing backend support:** Creation returns request-scoped microsecond timing. Task retry does not currently return equivalent timed assignment diagnostics.
- **Missing:** Persisted or separately queryable orchestration timing.
- **UX impact:** Timing is shown on the creation result only and Task Detail labels it unavailable.
- **Minimal change:** Persist one diagnostic observation per orchestration decision or expose it through a bounded diagnostic endpoint.
- **Priority:** Low to medium.

## 8. Command envelope and transport observations in read API

- **Requirement:** Show schema version and relay/Gateway delivery timing where proven.
- **Existing backend support:** The transport envelope has schema version and volatile delivery observations.
- **Missing:** These fields in `GET /commands/:id`; Gateway receive time is not a durable command record field.
- **UX impact:** The UI labels schema version unavailable and treats transport stages as architecture background, not trace evidence.
- **Minimal change:** Persist explicitly scoped delivery observations or expose a trace endpoint; keep them separate from immutable command identity.
- **Priority:** Medium for operational diagnostics.

## 9. Detailed ACK/result observation

- **Requirement:** Inspect ACK reason, original result status, and repeated idempotent observations.
- **Existing backend support:** Command exposes `ack_status`, `acknowledged_at`, final command status, result JSON, `completed_at`, and `last_error`.
- **Missing:** Dedicated ACK reason/original result status and observation history.
- **UX impact:** The frontend explains the actual ACK classification but cannot reconstruct more detail than the durable command record.
- **Minimal change:** Add normalized latest ACK/result metadata or a bounded observation endpoint if operationally necessary.
- **Priority:** Medium.

## 10. Mutation idempotency key

- **Requirement:** Safely recover when the network times out after task creation or an administrative mutation.
- **Existing backend support:** Transactional persistence and request/correlation IDs for diagnosis.
- **Missing:** Client-supplied idempotency key with replayed response semantics for task creation and control mutations.
- **UX impact:** The frontend never automatically retries critical mutations; an operator must refresh/search before deciding whether to submit again.
- **Minimal change:** Accept `Idempotency-Key`, persist its request/result identity per tenant and operation, and return the original committed response on replay.
- **Priority:** High for production mutation safety.
