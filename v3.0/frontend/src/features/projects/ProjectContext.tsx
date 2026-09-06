import { createContext, useContext, useMemo, useState, type ReactNode } from 'react';
import { useAuth } from '../auth/AuthContext';

interface ProjectValue {
  projectId: string;
  setProjectId: (projectId: string) => void;
}

const ProjectContext = createContext<ProjectValue | null>(null);

export function ProjectProvider({ children }: { children: ReactNode }) {
  const { session } = useAuth();
  const tenantId = session?.tenantId || '';
  const [projectId, setProjectIdState] = useState(() => sessionStorage.getItem(`polaris.project.${tenantId}`) || '');

  const value = useMemo(() => ({
    projectId,
    setProjectId: (next: string) => {
      setProjectIdState(next);
      if (next) sessionStorage.setItem(`polaris.project.${tenantId}`, next);
      else sessionStorage.removeItem(`polaris.project.${tenantId}`);
    },
  }), [projectId, tenantId]);

  return <ProjectContext.Provider value={value}>{children}</ProjectContext.Provider>;
}

export function useProjectContext() {
  const value = useContext(ProjectContext);
  if (!value) throw new Error('ProjectProvider is missing');
  return value;
}
