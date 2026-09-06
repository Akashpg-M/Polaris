import { useState, type FormEvent } from 'react';
import { ArrowRight, KeyRound, ShieldCheck } from 'lucide-react';
import { readableError } from '../../lib/api/errors';
import type { OperatorRole } from '../../types/domain';
import { useAuth } from './AuthContext';

export function AccessPage() {
  const { signIn } = useAuth();
  const [token, setToken] = useState('');
  const [tenantId, setTenantId] = useState('alpha_logistics');
  const [role, setRole] = useState<OperatorRole>('VIEWER');
  const [error, setError] = useState<unknown>();
  const [working, setWorking] = useState(false);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setWorking(true);
    setError(undefined);
    try { await signIn(token, tenantId, role); }
    catch (value) { setError(value); }
    finally { setWorking(false); }
  };

  const parsed = error ? readableError(error) : null;
  return <main className="access-page">
    <section className="access-story">
      <div className="brand-lockup"><span className="brand-mark">P</span><span>POLARIS</span></div>
      <p className="eyebrow">Physical operations cloud</p>
      <h1>See every device.<br />Understand its state.</h1>
      <p className="access-copy">A tenant-isolated control plane for registered fleets, current digital twins, and live spatial operations.</p>
      <div className="trust-line"><ShieldCheck /><span>Identity and permissions are always enforced by the Polaris backend.</span></div>
    </section>
    <section className="access-panel" aria-labelledby="access-title">
      <div>
        <p className="eyebrow">Operator access</p>
        <h2 id="access-title">Open your fleet workspace</h2>
        <p>Use an operator key provisioned by your Polaris administrator.</p>
      </div>
      <form onSubmit={submit} className="form-stack">
        <label>Operator token<input type="password" value={token} onChange={event => setToken(event.target.value)} autoComplete="current-password" placeholder="pol_…" required /></label>
        <label>Tenant ID<input value={tenantId} onChange={event => setTenantId(event.target.value)} autoComplete="organization" required /></label>
        <label>Role for interface presentation
          <select value={role} onChange={event => setRole(event.target.value as OperatorRole)}>
            <option value="VIEWER">Viewer</option><option value="OPERATOR">Operator</option><option value="TENANT_ADMIN">Tenant administrator</option><option value="PLATFORM_ADMIN">Platform administrator</option>
          </select>
        </label>
        <p className="field-note">The backend currently has no session-profile endpoint. This role controls presentation only; backend authorization remains authoritative.</p>
        {parsed && <div className="inline-error" role="alert"><strong>{parsed.message}</strong>{parsed.requestId && <small>Request ID: {parsed.requestId}</small>}</div>}
        <button className="button primary access-button" disabled={working || !token.trim() || !tenantId.trim()}>{working ? 'Validating access…' : <><KeyRound size={17} /> Continue to Polaris <ArrowRight size={17} /></>}</button>
      </form>
      <p className="session-note">The bearer token stays in this browser tab’s session storage and is cleared on sign out.</p>
    </section>
  </main>;
}

