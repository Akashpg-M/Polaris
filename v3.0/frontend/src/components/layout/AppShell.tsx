import { useEffect, useState, type ReactNode } from 'react';
import { Activity, Boxes, ChevronLeft, ChevronRight, Command, FolderKanban, History, LayoutDashboard, ListTodo, Map, MapPinned, Menu, Moon, Plus, Route, Search, ServerCog, Settings2, Sun, TrafficCone, Waypoints, X } from 'lucide-react';
import { NavLink, useNavigate } from 'react-router-dom';
import { can, roleLabel } from '../../lib/permissions';
import { useAuth } from '../../features/auth/AuthContext';
import { useDashboard } from '../../features/dashboard/DashboardContext';
import { useProjectContext } from '../../features/projects/ProjectContext';
import { useProjects } from '../../features/fleet/queries';
import { WebSocketStatus } from '../status/Status';

const primary = [
  { to: '/', label: 'Overview', icon: LayoutDashboard },
  { to: '/devices', label: 'Devices', icon: Boxes },
  { to: '/fleet/map', label: 'Live map', icon: Map },
  { to: '/twins', label: 'Digital twins', icon: Waypoints },
  { to: '/projects', label: 'Projects', icon: FolderKanban },
];
const future = [
  { label: 'Registry', icon: Settings2 },
  { label: 'Observability', icon: ServerCog },
];

export function AppShell({ children }: { children: ReactNode }) {
  const { session, signOut, changeTenant } = useAuth();
  const { status } = useDashboard();
  const { projectId, setProjectId } = useProjectContext();
  const projects = useProjects();
  const navigate = useNavigate();
  const [collapsed, setCollapsed] = useState(() => localStorage.getItem('polaris.sidebar.collapsed') === 'true');
  const [mobileOpen, setMobileOpen] = useState(false);
  const [theme, setTheme] = useState<'dark' | 'light'>(() => (localStorage.getItem('polaris.theme') as 'dark' | 'light') || 'dark');
  const [tenantDraft, setTenantDraft] = useState(session?.tenantId || '');

  useEffect(() => { document.documentElement.dataset.theme = theme; localStorage.setItem('polaris.theme', theme); }, [theme]);
  useEffect(() => { localStorage.setItem('polaris.sidebar.collapsed', String(collapsed)); }, [collapsed]);

  const changeProject = (value: string) => setProjectId(value);
  const applyTenant = async () => { if (tenantDraft && tenantDraft !== session?.tenantId) { await changeTenant(tenantDraft); setProjectId(''); navigate('/'); } };

  return <div className={`app-shell ${collapsed ? 'sidebar-collapsed' : ''}`}>
    <button className="mobile-menu button ghost" onClick={() => setMobileOpen(true)} aria-label="Open navigation"><Menu /></button>
    {mobileOpen && <button className="sidebar-scrim" onClick={() => setMobileOpen(false)} aria-label="Close navigation overlay" />}
    <aside className={`sidebar ${mobileOpen ? 'mobile-open' : ''}`}>
      <div className="sidebar-brand"><span className="brand-mark">P</span><span className="brand-words"><strong>POLARIS</strong><small>CONTROL PLANE</small></span><button className="mobile-close" onClick={() => setMobileOpen(false)} aria-label="Close navigation"><X /></button></div>
      <nav aria-label="Primary navigation">
        <p className="nav-group-label">Workspace</p>
        {primary.map(item => <NavLink key={item.to} to={item.to} end={item.to === '/'} title={collapsed ? item.label : undefined} onClick={() => setMobileOpen(false)} className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}><item.icon /><span>{item.label}</span></NavLink>)}
        <p className="nav-group-label future-label">Operations</p>
        <NavLink to="/tasks" title={collapsed ? 'Tasks' : undefined} onClick={() => setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><ListTodo/><span>Tasks</span></NavLink>
        {session && can(session.role,'createTask') && <NavLink to="/tasks/new" title={collapsed?'Create task':undefined} onClick={()=>setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><Plus/><span>Create task</span></NavLink>}
        <NavLink to="/commands" title={collapsed ? 'Commands' : undefined} onClick={() => setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><Command/><span>Commands</span></NavLink>
        <NavLink to="/operations/activity" title={collapsed ? 'Activity' : undefined} onClick={() => setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><History/><span>Activity</span></NavLink>
        <p className="nav-group-label future-label">Mobility</p>
        <NavLink to="/mobility" end title={collapsed?'Mobility overview':undefined} onClick={()=>setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><Activity/><span>Overview</span></NavLink>
        <NavLink to="/mobility/nearby" title={collapsed?'Nearby search':undefined} onClick={()=>setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><MapPinned/><span>Nearby search</span></NavLink>
        <NavLink to="/mobility/routes" title={collapsed?'Routing':undefined} onClick={()=>setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><Route/><span>Routing</span></NavLink>
        <NavLink to="/mobility/traffic" title={collapsed?'Traffic':undefined} onClick={()=>setMobileOpen(false)} className={({isActive})=>`nav-item ${isActive?'active':''}`}><TrafficCone/><span>Traffic</span></NavLink>
        <p className="nav-group-label future-label">Later phases</p>
        {future.map(item => <span key={item.label} className="nav-item disabled" title={`${item.label} · Later phase`}><item.icon /><span>{item.label}</span><small>Later</small></span>)}
      </nav>
      <div className="sidebar-bottom">
        <div className="control-health"><span className="control-icon"><Activity /></span><span><strong>Core services</strong><small>Deployment connected</small></span></div>
        <button className="collapse-button" onClick={() => setCollapsed(value => !value)} aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}>{collapsed ? <ChevronRight /> : <><ChevronLeft /><span>Collapse</span></>}</button>
      </div>
    </aside>
    <div className="workspace">
      <header className="topbar">
        <div className="tenant-context">
          <small>Tenant workspace</small>
          {session?.role === 'PLATFORM_ADMIN' ? <div className="tenant-switch"><input value={tenantDraft} onChange={event => setTenantDraft(event.target.value)} aria-label="Tenant ID" /><button onClick={() => void applyTenant()}>Apply</button></div> : <strong>{session?.tenant.display_name}</strong>}
        </div>
        <div className="topbar-center">
          <label className="project-select"><FolderKanban /><select value={projectId} onChange={event => changeProject(event.target.value)} aria-label="Project filter"><option value="">All projects</option>{projects.data?.map(project => <option key={project.project_id} value={project.project_id}>{project.name}</option>)}</select></label>
          <div className="search-placeholder" title="Search is limited to loaded page in Fleet views"><Search /><span>Search in Fleet</span><kbd>/</kbd></div>
        </div>
        <div className="topbar-actions">
          <WebSocketStatus status={status} />
          <button className="icon-button" onClick={() => setTheme(value => value === 'dark' ? 'light' : 'dark')} aria-label={`Use ${theme === 'dark' ? 'light' : 'dark'} theme`}>{theme === 'dark' ? <Sun /> : <Moon />}</button>
          <div className="identity"><span className="avatar">{session?.role.charAt(0)}</span><span><strong>{roleLabel(session!.role)}</strong><small>{session?.tenantId}</small></span></div>
          <button className="signout" onClick={signOut}>Sign out</button>
        </div>
      </header>
      <main className="page-container">{children}</main>
    </div>
  </div>;
}
