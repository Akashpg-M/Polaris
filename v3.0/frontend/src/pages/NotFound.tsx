import { ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function NotFound() {
  return <div className="page-state"><strong>Page not found</strong><span>This surface is not part of the current Polaris frontend phase.</span><Link className="button primary" to="/"><ArrowLeft /> Return to overview</Link></div>;
}
