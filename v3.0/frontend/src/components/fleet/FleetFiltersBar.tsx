import { Search, SlidersHorizontal } from 'lucide-react';
import type { ConnectivityStatus, DeviceLifecycle, FleetFilters, Project } from '../../types/domain';
import { deviceTypeLabels } from '../../lib/format';

export function FleetFiltersBar({ filters, projects, onChange, showConnectivity = true }: {
  filters: FleetFilters;
  projects: Project[];
  onChange: (filters: FleetFilters) => void;
  showConnectivity?: boolean;
}) {
  return <div className="filters-bar">
    <span className="filters-title"><SlidersHorizontal /> Filters</span>
    <label className="filter-search"><Search /><input value={filters.search || ''} onChange={event => onChange({ ...filters, search: event.target.value })} placeholder="Loaded device ID or name" aria-label="Search loaded devices" /></label>
    <select value={filters.projectId || ''} onChange={event => onChange({ ...filters, projectId: event.target.value })} aria-label="Filter by project"><option value="">All projects</option>{projects.map(project => <option value={project.project_id} key={project.project_id}>{project.name}</option>)}</select>
    <select value={filters.deviceType || ''} onChange={event => onChange({ ...filters, deviceType: event.target.value })} aria-label="Filter by device type"><option value="">All device types</option>{Object.entries(deviceTypeLabels).map(([id, label]) => <option value={id} key={id}>{label}</option>)}</select>
    <select value={filters.lifecycle || ''} onChange={event => onChange({ ...filters, lifecycle: event.target.value as DeviceLifecycle | '' })} aria-label="Filter by lifecycle"><option value="">All lifecycles</option>{['REGISTERED','ACTIVE','SUSPENDED','DECOMMISSIONED'].map(value => <option value={value} key={value}>{value}</option>)}</select>
    {showConnectivity && <select value={filters.connectivity || ''} onChange={event => onChange({ ...filters, connectivity: event.target.value as ConnectivityStatus | '' })} aria-label="Filter by connectivity"><option value="">All connectivity</option>{['ONLINE','STALE','OFFLINE','NEVER_CONNECTED'].map(value => <option value={value} key={value}>{value.replace('_', ' ')}</option>)}</select>}
    {Object.values(filters).some(Boolean) && <button className="clear-filters" onClick={() => onChange({})}>Clear</button>}
  </div>;
}

