# Polaris Frontend Phase 3 — Repository-grounded sequence

1. Extend typed contracts and the API client for canonical nearby search, canonical route creation, engine readiness, and experimental predicted zones.
2. Expose the existing engine `/readyz` through the same-origin Vite/Nginx frontend path so browser diagnostics remain Docker-compatible.
3. Introduce shared Leaflet primitives for device markers, coordinate selection, search radius, routes, and zones; make the existing fleet map consume the same map component.
4. Add Mobility formatting, validation, profile policy, persisted-route parsing, and error classification as independently tested logic.
5. Add explicit-query hooks: no automatic route recalculation, no continuous nearby search, modest readiness polling, and bounded predicted-zone reads.
6. Enable Mobility navigation and routes for Overview, Nearby, Routing, Traffic, Diagnostics, and Experimental Analytics.
7. Connect persisted Phase 2 command routes to the Mobility route viewer by command ID without recomputation.
8. Add responsive/accessibility styles, focused unit/component coverage, backend gap documentation, and closure evidence.

The implementation deliberately omits a traffic heatmap, route history, tenant-wide spatial totals, candidate eligibility claims, graph switching, and legacy matcher product UI because the required backend contracts do not exist.
