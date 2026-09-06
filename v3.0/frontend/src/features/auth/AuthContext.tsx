import { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from 'react';
import { PolarisApi } from '../../lib/api/client';
import { useQueryClient } from '../../lib/query/queryClient';
import type { AuthSession, OperatorRole } from '../../types/domain';

const storageKey = 'polaris.operator.session.v1';

function restoreSession(): AuthSession | null {
  try {
    const raw = sessionStorage.getItem(storageKey);
    if (!raw) return null;
    const value = JSON.parse(raw) as AuthSession;
    return value.token && value.tenantId && value.role && value.tenant ? value : null;
  } catch {
    return null;
  }
}

interface AuthValue {
  session: AuthSession | null;
  api: PolarisApi | null;
  signIn: (token: string, tenantId: string, role: OperatorRole) => Promise<void>;
  signOut: () => void;
  changeTenant: (tenantId: string) => Promise<void>;
}

const AuthContext = createContext<AuthValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient();
  const [session, setSession] = useState<AuthSession | null>(restoreSession);

  const signOut = useCallback(() => {
    sessionStorage.removeItem(storageKey);
    queryClient.clear();
    setSession(null);
  }, [queryClient]);

  const api = useMemo(() => session
    ? new PolarisApi({ token: session.token, tenantId: session.tenantId }, signOut)
    : null, [session, signOut]);

  const signIn = useCallback(async (token: string, tenantId: string, role: OperatorRole) => {
    const normalizedToken = token.trim();
    const normalizedTenant = tenantId.trim();
    const validationApi = new PolarisApi({ token: normalizedToken, tenantId: normalizedTenant });
    const tenant = await validationApi.tenant(normalizedTenant);
    const next = { token: normalizedToken, tenantId: normalizedTenant, role, tenant };
    queryClient.clear();
    sessionStorage.setItem(storageKey, JSON.stringify(next));
    setSession(next);
  }, [queryClient]);

  const changeTenant = useCallback(async (tenantId: string) => {
    if (!session || session.role !== 'PLATFORM_ADMIN') throw new Error('Only platform administrators can change tenant context.');
    const normalized = tenantId.trim();
    const validationApi = new PolarisApi({ token: session.token, tenantId: normalized });
    const tenant = await validationApi.tenant(normalized);
    const next = { ...session, tenantId: normalized, tenant };
    queryClient.clear();
    sessionStorage.setItem(storageKey, JSON.stringify(next));
    setSession(next);
  }, [queryClient, session]);

  return <AuthContext.Provider value={{ session, api, signIn, signOut, changeTenant }}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) throw new Error('AuthProvider is missing');
  return value;
}

