import { useAuth } from '../auth/AuthContext';
import { useQuery } from '../../lib/query/queryClient';
import type { CommandFilters, TaskDetail, TaskFilters } from '../../types/domain';
import { isCommandTerminal, isTaskTerminal } from './model';

const taskDetailPoll = (data: TaskDetail | undefined) => data && isTaskTerminal(data.task.status) ? false : 5_000;
const commandDetailPoll = (data: import('../../types/domain').Command | undefined) => data && isCommandTerminal(data.status) ? false : 5_000;

export function useTasks(filters: TaskFilters = {}, limit = 50, cursor = '', polling = true) {
  const { session, api } = useAuth();
  return useQuery({
    key: ['tasks', session?.tenantId || '', filters, limit, cursor],
    query: signal => api!.tasks(filters, limit, cursor, signal),
    enabled: Boolean(session && api), staleTime: 5_000, refetchInterval: polling ? 10_000 : false,
  });
}

export function useTask(taskId: string) {
  const { session, api } = useAuth();
  return useQuery<TaskDetail>({
    key: ['task', session?.tenantId || '', taskId],
    query: signal => api!.task(taskId, signal),
    enabled: Boolean(session && api && taskId), staleTime: 3_000,
    refetchInterval: taskDetailPoll,
  });
}

export function useCommands(filters: CommandFilters = {}, limit = 50, cursor = '', polling = true) {
  const { session, api } = useAuth();
  return useQuery({
    key: ['commands', session?.tenantId || '', filters, limit, cursor],
    query: signal => api!.commands(filters, limit, cursor, signal),
    enabled: Boolean(session && api), staleTime: 5_000, refetchInterval: polling ? 10_000 : false,
  });
}

export function useCommand(commandId: string) {
  const { session, api } = useAuth();
  return useQuery({
    key: ['command', session?.tenantId || '', commandId],
    query: signal => api!.command(commandId, signal),
    enabled: Boolean(session && api && commandId), staleTime: 3_000,
    refetchInterval: commandDetailPoll,
  });
}
