import { Activity, FlaskConical, Gauge, MapPinned, Route, TrafficCone } from 'lucide-react';
import { NavLink } from 'react-router-dom';

const links = [
  ['/mobility', 'Overview', Activity, true], ['/mobility/nearby', 'Nearby', MapPinned, false], ['/mobility/routes', 'Routing', Route, false],
  ['/mobility/traffic', 'Traffic', TrafficCone, false], ['/mobility/diagnostics', 'Diagnostics', Gauge, false], ['/mobility/experimental', 'Experimental', FlaskConical, false],
] as const;

export function MobilityNav() { return <nav className="mobility-tabs" aria-label="Mobility sections">{links.map(([to, label, Icon, end]) => <NavLink key={to} to={to} end={end} className={({ isActive }) => isActive ? 'active' : ''}><Icon />{label}</NavLink>)}</nav>; }
