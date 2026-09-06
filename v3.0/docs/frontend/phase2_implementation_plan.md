# Polaris Frontend Phase 2 — Repository-Grounded Sequence

1. Extend domain types and the same-origin API client with the exact task, command, create-result, timing, filter, cancellation, and retry contracts.
2. Extend the tenant query cache with data-aware polling so mounted active views refresh while terminal entity details stop issuing requests.
3. Add centralized Operations permissions, status/priority/ACK presentation, confirmed-timestamp lifecycles, requirements, targets, immutable payloads, timing diagnostics, entity links, and accessible confirmations.
4. Enable role-aware Operations navigation and routes for task list/create/detail, command list/detail, and session activity.
5. Implement the structured task wizard using backend enums and units, automatic locked capabilities, client UX validation, an explicit review, and non-optimistic submission.
6. Implement bounded cursor lists using only supported server filters. Preserve filters in URLs and label missing totals/cursor metadata.
7. Implement task/command details, conservative cancel, distinct task-orchestration and command-delivery retry language, cross-entity navigation, attempt aggregates, ACK/result inspection, and route-payload inspection where already persisted.
8. Use 10-second list/activity polling and 5-second active-detail polling because no browser lifecycle stream exists. Do not poll terminal details.
9. Add focused unit/component tests, run strict lint/type/build gates, document confirmed backend gaps, and close with reproducible evidence. Docker/Nginx same-origin deployment remains unchanged.