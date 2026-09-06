import { ArrowUpRight, Calendar, FolderKanban } from 'lucide-react';
import { Link } from 'react-router-dom';
import { ProjectStatusBadge } from '../components/status/Status';
import { EmptyState, ErrorState, PageLoader } from '../components/ui/States';
import { PageHeader } from '../components/ui/PageHeader';
import { exactTime } from '../lib/format';
import { useProjects } from '../features/fleet/queries';

export default function Projects() {
  const projects = useProjects();
  if (projects.isLoading) return <PageLoader label="Loading projects" />;
  if (projects.error) return <ErrorState error={projects.error} retry={() => void projects.refetch()} />;
  return <div className="page-stack"><PageHeader eyebrow="Fleet organization" title="Projects" description="Read-only tenant workspaces that organize registered devices." actions={<span className="loaded-scope">Backend limit · 100 projects</span>} />{!projects.data?.length ? <EmptyState title="No projects" message="No projects exist for this tenant. Project administration is reserved for a later phase." /> : <div className="project-grid">{projects.data.map(project => <Link className="project-card" to={`/projects/${project.project_id}`} key={project.project_id}><div className="project-icon"><FolderKanban /></div><div className="project-card-title"><div><h2>{project.name}</h2><code>{project.project_id}</code></div><ArrowUpRight /></div><p>{project.description || 'No project description has been provided.'}</p><div className="project-meta"><ProjectStatusBadge status={project.status} /><span><Calendar /> Updated {exactTime(project.updated_at)}</span></div></Link>)}</div>}<p className="architecture-note"><FolderKanban /> Device and online counts are omitted because the backend does not expose project aggregates; Polaris avoids an N+1 request pattern.</p></div>;
}

