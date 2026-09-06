import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import { ModuleStatusView, ProfileBadge, RouteErrorCard, RouteSummary } from './MobilityUi';
import { MobilityNav } from './MobilityNav';
import { PolarisApiError } from '../../lib/api/errors';
import type { RouteResult } from '../../types/domain';

const route: RouteResult = { route_id:'route-1',road_graph_version:'chennai-v1',snapshot_version:842,policy:'FASTEST',distance_meters:8400,estimated_time:1_062_000_000_000,waypoints:[{latitude:13,longitude:80},{latitude:13.1,longitude:80.1}],edge_ids:[1,2],expanded_nodes:13402 };

describe('mobility presentation', () => {
  it('renders navigable Mobility sections', () => { const html=renderToStaticMarkup(<MemoryRouter><MobilityNav/></MemoryRouter>); expect(html).toContain('Nearby'); expect(html).toContain('Diagnostics'); expect(html).toContain('Experimental'); });
  it('shows component-level degraded state without a global outage claim', () => { const html=renderToStaticMarkup(<ModuleStatusView module={{state:'DEGRADED',message:'road graph unavailable',components:{spatial:{state:'READY'},routing:{state:'FAILED'}}}}/>); expect(html).toContain('DEGRADED'); expect(html).toContain('spatial'); expect(html).toContain('FAILED'); });
  it('shows route policy, units, graph and immutable snapshot metadata', () => { const html=renderToStaticMarkup(<RouteSummary route={route}/>); expect(html).toContain('Fastest time'); expect(html).toContain('8.40 km'); expect(html).toContain('17m 42s'); expect(html).toContain('chennai-v1'); expect(html).toContain('842'); });
  it('labels persisted routes as immutable inspection', () => { const html=renderToStaticMarkup(<RouteSummary route={route} persisted/>); expect(html).toContain('Persisted Command Route'); expect(html).toContain('does not recalculate'); });
  it('presents routing error identity and profile text', () => { const html=renderToStaticMarkup(<><RouteErrorCard error={new PolarisApiError({code:'OUTSIDE_ROUTING_REGION',message:'outside',status:422,requestId:'req-1'})}/><ProfileBadge profile="ROAD_VEHICLE"/></>); expect(html).toContain('Outside the loaded road region'); expect(html).toContain('req-1'); expect(html).toContain('Road vehicle'); });
});
