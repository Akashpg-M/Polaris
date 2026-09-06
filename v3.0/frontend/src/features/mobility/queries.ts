import { useAuth } from '../auth/AuthContext';
import { useQuery } from '../../lib/query/queryClient';
import type { NearbySearch } from '../../types/domain';

export function useMobilityReadiness() {
  const { session, api } = useAuth();
  return useQuery({ key: ['mobility-status', session?.tenantId || ''], query: signal => api!.mobilityReadiness(signal), enabled: Boolean(session && api), staleTime: 10_000, refetchInterval: 15_000 });
}

export function useNearby(search?: NearbySearch) {
  const { session, api } = useAuth();
  return useQuery({ key: ['mobility-nearby', session?.tenantId || '', search || null], query: signal => api!.nearby(search!, signal), enabled: Boolean(session && api && search), staleTime: 0 });
}

export function usePredictedZones() {
  const { session, api } = useAuth();
  return useQuery({ key: ['predicted-zones', session?.tenantId || ''], query: signal => api!.predictedZones(signal), enabled: Boolean(session && api), staleTime: 60_000 });
}
