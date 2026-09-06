import { renderToStaticMarkup } from 'react-dom/server';
import { describe, expect, it } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import type { Command, Task } from '../../types/domain';
import { CommandPayloadViewer } from './CommandDetails';
import { CommandLifecycleTimeline, TaskLifecycleTimeline } from './LifecycleTimeline';
import { CommandStatusBadge, TaskPriorityBadge, TaskStatusBadge } from './OperationsStatus';
import { TaskRequirementsView } from './TaskDetails';

const task: Task={task_id:'task-1',tenant_id:'tenant',task_type:'NAVIGATE',status:'ASSIGNED',priority:'HIGH',requirements:{required_capabilities:['navigate'],minimum_battery:50},target:{lat:13,lon:80},assigned_device_id:'vehicle-1',correlation_id:'corr',created_by:'actor',version:2,created_at:'2026-09-06T10:00:00Z',updated_at:'2026-09-06T10:00:01Z',assigned_at:'2026-09-06T10:00:01Z',expires_at:'2026-09-06T10:05:00Z'};
const command:Command={command_id:'command-1',tenant_id:'tenant',device_id:'vehicle-1',task_id:'task-1',command_type:'NAVIGATE',payload:{route_id:'route-1',road_graph_version:'graph-v1',routing_snapshot_version:4,origin:{latitude:13,longitude:80},destination:{latitude:13.1,longitude:80.1},waypoints:[]},status:'ACKNOWLEDGED',sequence_number:42,correlation_id:'corr',causation_id:'task-1',attempt_count:2,max_attempts:5,version:3,created_at:'2026-09-06T10:00:01Z',available_at:'2026-09-06T10:00:01Z',sent_at:'2026-09-06T10:00:02Z',acknowledged_at:'2026-09-06T10:00:03Z',expires_at:'2026-09-06T10:05:00Z',ack_status:'DUPLICATE'};

describe('operations components',()=>{
  it('renders text labels in addition to status color',()=>{const html=renderToStaticMarkup(<><TaskStatusBadge status="IN_PROGRESS"/><CommandStatusBadge status="DELIVERED"/><TaskPriorityBadge priority="CRITICAL"/></>);expect(html).toContain('In Progress');expect(html).toContain('Delivered');expect(html).toContain('CRITICAL');});
  it('shows only confirmed task timestamps and leaves completion unconfirmed',()=>{const html=renderToStaticMarkup(<TaskLifecycleTimeline task={task}/>);expect(html).toContain('Assigned');expect(html).toContain('Not confirmed');});
  it('distinguishes delivery, acknowledgement and incomplete execution',()=>{const html=renderToStaticMarkup(<CommandLifecycleTimeline command={command}/>);expect(html).toContain('Delivery attempted');expect(html).toContain('Device acknowledged');expect(html).toContain('Execution completed');expect(html).toContain('Not confirmed');});
  it('locks the automatically required capability',()=>{const html=renderToStaticMarkup(<TaskRequirementsView taskType="NAVIGATE" requirements={task.requirements}/>);expect(html).toContain('Automatically required by Polaris');expect(html).toContain('locked');});
  it('renders persisted route payload metadata without edit controls',()=>{const html=renderToStaticMarkup(<MemoryRouter><CommandPayloadViewer command={command}/></MemoryRouter>);expect(html).toContain('Immutable Command Payload');expect(html).toContain('route-1');expect(html).toContain('graph-v1');expect(html).not.toContain('textarea');});
});
