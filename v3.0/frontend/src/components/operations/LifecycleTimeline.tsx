import { Check, Circle, X } from 'lucide-react';
import type { Command, Task } from '../../types/domain';
import { exactTime } from '../../lib/format';
import { CommandStatusBadge, TaskStatusBadge } from './OperationsStatus';

interface Point { label: string; time?: string; confirmed: boolean; terminal?: boolean }

function Timeline({ points }: { points: Point[] }) {
  return <ol className="lifecycle-timeline">{points.map((point, index) => <li className={point.confirmed ? 'confirmed' : 'unconfirmed'} key={`${point.label}-${index}`}>
    <span className="timeline-mark">{point.confirmed ? point.terminal ? <X /> : <Check /> : <Circle />}</span>
    <div><strong>{point.label}</strong><span>{point.time ? exactTime(point.time) : point.confirmed ? 'Confirmed; timestamp not exposed' : 'Not confirmed'}</span></div>
  </li>)}</ol>;
}

export function TaskLifecycleTimeline({ task }: { task: Task }) {
  const points: Point[] = [{ label: 'Created', time: task.created_at, confirmed: true }];
  if (task.assigned_at || ['ASSIGNED','IN_PROGRESS','COMPLETED'].includes(task.status)) points.push({ label: 'Assigned', time: task.assigned_at, confirmed: true });
  if (task.started_at || ['IN_PROGRESS','COMPLETED'].includes(task.status)) points.push({ label: 'In progress', time: task.started_at, confirmed: true });
  if (task.status === 'COMPLETED') points.push({ label: 'Completed', time: task.completed_at, confirmed: true });
  else if (['FAILED','EXPIRED','CANCELLED'].includes(task.status)) points.push({ label: task.status === 'FAILED' ? 'Failed' : task.status === 'EXPIRED' ? 'Expired' : 'Cancelled', time: task.status === 'CANCELLED' ? undefined : task.failed_at, confirmed: true, terminal: true });
  else points.push({ label: 'Completed', confirmed: false });
  return <div className="timeline-block"><TaskStatusBadge status={task.status} /><Timeline points={points} /></div>;
}

export function CommandLifecycleTimeline({ command }: { command: Command }) {
  const points: Point[] = [{ label: 'Command persisted', time: command.created_at, confirmed: true }];
  if (command.sent_at || ['DELIVERED','ACKNOWLEDGED','COMPLETED'].includes(command.status)) points.push({ label: 'Delivery attempted', time: command.sent_at, confirmed: true });
  if (command.acknowledged_at || ['ACKNOWLEDGED','COMPLETED'].includes(command.status)) points.push({ label: 'Device acknowledged', time: command.acknowledged_at, confirmed: true });
  if (command.status === 'COMPLETED') points.push({ label: 'Execution completed', time: command.completed_at, confirmed: true });
  else if (['FAILED','EXPIRED','CANCELLED'].includes(command.status)) points.push({ label: title(command.status), time: command.completed_at, confirmed: true, terminal: true });
  else points.push({ label: 'Execution completed', confirmed: false });
  return <div className="timeline-block"><CommandStatusBadge status={command.status} /><Timeline points={points} /></div>;
}

function title(value: string) { return value.charAt(0) + value.slice(1).toLowerCase(); }
