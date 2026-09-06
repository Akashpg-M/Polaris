import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { lazy, Suspense } from 'react';
import { AppShell } from './components/layout/AppShell';
import { PageLoader } from './components/ui/States';
import { AccessPage } from './features/auth/AccessPage';
import { AuthProvider, useAuth } from './features/auth/AuthContext';
import { DashboardProvider } from './features/dashboard/DashboardContext';
import { ProjectProvider } from './features/projects/ProjectContext';
import { QueryProvider } from './lib/query/queryClient';
const DeviceDetail = lazy(() => import('./pages/DeviceDetail'));
const Devices = lazy(() => import('./pages/Devices'));
const FleetMap = lazy(() => import('./pages/FleetMap'));
const NotFound = lazy(() => import('./pages/NotFound'));
const Overview = lazy(() => import('./pages/Overview'));
const ProjectDetail = lazy(() => import('./pages/ProjectDetail'));
const Projects = lazy(() => import('./pages/Projects'));
const Twins = lazy(() => import('./pages/Twins'));
const Tasks = lazy(() => import('./pages/Tasks'));
const CreateTask = lazy(() => import('./pages/CreateTask'));
const TaskDetail = lazy(() => import('./pages/TaskDetail'));
const Commands = lazy(() => import('./pages/Commands'));
const CommandDetail = lazy(() => import('./pages/CommandDetail'));
const OperationsActivity = lazy(() => import('./pages/OperationsActivity'));
const MobilityOverview = lazy(() => import('./pages/MobilityOverview'));
const NearbySearch = lazy(() => import('./pages/NearbySearch'));
const RouteExplorer = lazy(() => import('./pages/RouteExplorer'));
const MobilityTraffic = lazy(() => import('./pages/MobilityTraffic'));
const MobilityDiagnostics = lazy(() => import('./pages/MobilityDiagnostics'));
const MobilityExperimental = lazy(() => import('./pages/MobilityExperimental'));

function AuthenticatedApp() {
  const { session } = useAuth();
  if (!session) return <AccessPage />;

  return (
    <ProjectProvider key={session.tenantId}>
      <DashboardProvider>
        <AppShell>
          <Suspense fallback={<PageLoader label="Loading workspace" />}><Routes>
            <Route path="/" element={<Overview />} />
            <Route path="/devices" element={<Devices />} />
            <Route path="/devices/:deviceId" element={<DeviceDetail />} />
            <Route path="/fleet/map" element={<FleetMap />} />
            <Route path="/twins" element={<Twins />} />
            <Route path="/projects" element={<Projects />} />
            <Route path="/projects/:projectId" element={<ProjectDetail />} />
            <Route path="/tasks" element={<Tasks />} />
            <Route path="/tasks/new" element={<CreateTask />} />
            <Route path="/tasks/:taskId" element={<TaskDetail />} />
            <Route path="/commands" element={<Commands />} />
            <Route path="/commands/:commandId" element={<CommandDetail />} />
            <Route path="/operations/activity" element={<OperationsActivity />} />
            <Route path="/mobility" element={<MobilityOverview />} />
            <Route path="/mobility/nearby" element={<NearbySearch />} />
            <Route path="/mobility/routes" element={<RouteExplorer />} />
            <Route path="/mobility/traffic" element={<MobilityTraffic />} />
            <Route path="/mobility/diagnostics" element={<MobilityDiagnostics />} />
            <Route path="/mobility/experimental" element={<MobilityExperimental />} />
            <Route path="*" element={<NotFound />} />
          </Routes></Suspense>
        </AppShell>
      </DashboardProvider>
    </ProjectProvider>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <QueryProvider>
        <AuthProvider>
          <AuthenticatedApp />
        </AuthProvider>
      </QueryProvider>
    </BrowserRouter>
  );
}
