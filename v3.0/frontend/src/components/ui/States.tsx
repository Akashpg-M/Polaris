import { AlertTriangle, Inbox, LoaderCircle, RefreshCw } from 'lucide-react';
import { readableError } from '../../lib/api/errors';

export function PageLoader({ label = 'Loading fleet data' }: { label?: string }) {
  return <div className="page-state" role="status"><LoaderCircle className="spin" /><strong>{label}</strong><span>Reading the current tenant view from Polaris.</span></div>;
}

export function EmptyState({ title, message }: { title: string; message: string }) {
  return <div className="page-state"><Inbox /><strong>{title}</strong><span>{message}</span></div>;
}

export function ErrorState({ error, retry }: { error: unknown; retry?: () => void }) {
  const value = readableError(error);
  return <div className="page-state error-state" role="alert">
    <AlertTriangle />
    <strong>Unable to load this view</strong>
    <span>{value.message}</span>
    {value.requestId && <small>Request ID: <code>{value.requestId}</code></small>}
    {retry && <button className="button secondary" onClick={retry}><RefreshCw size={15} /> Retry</button>}
  </div>;
}

export function TableSkeleton({ rows = 6 }: { rows?: number }) {
  return <div className="skeleton-table" role="status" aria-label="Loading table">
    {Array.from({ length: rows }).map((_, index) => <div className="skeleton-row" key={index}><span /><span /><span /><span /></div>)}
  </div>;
}

