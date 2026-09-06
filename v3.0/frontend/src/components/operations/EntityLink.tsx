import { ArrowUpRight } from 'lucide-react';
import { Link } from 'react-router-dom';

export function EntityLink({ label, id, to }: { label: string; id: string; to: string }) {
  return <Link className="entity-link" to={to}><span><small>{label}</small><code>{id}</code></span><ArrowUpRight /></Link>;
}
