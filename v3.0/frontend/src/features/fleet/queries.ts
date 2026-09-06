import { useAuth } from '../auth/AuthContext';
import { useQuery } from '../../lib/query/queryClient';
import type { FleetFilters } from '../../types/domain';

export function useProjects() {
  const { session, api } = useAuth();
  return useQuery({
    key: ['projects', session?.tenantId || ''],
    query: signal => api!.projects(signal),
    enabled: Boolean(session && api),
    staleTime: 60_000,
  });
}

export function useCapabilities() {
  const { session, api } = useAuth();
  return useQuery({ key: ['capabilities', session?.tenantId || ''], query: signal => api!.capabilities(undefined, signal), enabled: Boolean(session && api), staleTime: 120_000 });
}

export function useProject(projectId: string) {
  const { session, api } = useAuth();
  return useQuery({
    key: ['project', session?.tenantId || '', projectId],
    query: signal => api!.project(projectId, signal),
    enabled: Boolean(session && api && projectId),
    staleTime: 60_000,
  });
}

export function useTwins(filters: FleetFilters = {}, limit = 50, cursor = '') {
  const { session, api } = useAuth();
  return useQuery({
    key: ['twins', session?.tenantId || '', filters, limit, cursor],
    query: signal => api!.twins(filters, limit, cursor, signal),
    enabled: Boolean(session && api),
    staleTime: 20_000,
  });
}

export function useTwin(deviceId: string) {
  const { session, api } = useAuth();
  return useQuery({
    key: ['twin', session?.tenantId || '', deviceId],
    query: signal => api!.twin(deviceId, signal),
    enabled: Boolean(session && api && deviceId),
    staleTime: 20_000,
  });
}
