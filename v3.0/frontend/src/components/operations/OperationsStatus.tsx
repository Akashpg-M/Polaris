import type { AckStatus, CommandStatus, TaskPriority, TaskStatus } from '../../types/domain';
import { ackExplanation, titleCase } from '../../features/operations/model';

export function TaskStatusBadge({ status }: { status: TaskStatus }) {
  return <span className={`operation-badge task-${status.toLowerCase()}`}><span aria-hidden="true" />{titleCase(status)}</span>;
}

export function CommandStatusBadge({ status }: { status: CommandStatus }) {
  return <span className={`operation-badge command-${status.toLowerCase()}`}><span aria-hidden="true" />{titleCase(status)}</span>;
}

export function TaskPriorityBadge({ priority }: { priority: TaskPriority }) {
  return <span className={`priority-badge priority-${priority.toLowerCase()}`}>{priority}</span>;
}

export function AckBadge({ status }: { status?: AckStatus }) {
  return <span className={`ack-badge ack-${(status || 'none').toLowerCase()}`} title={ackExplanation(status)}>{status || 'Not received'}</span>;
}
