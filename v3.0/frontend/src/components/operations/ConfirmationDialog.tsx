import { useEffect, useRef, type ReactNode } from 'react';
import { AlertTriangle, X } from 'lucide-react';

export function ConfirmationDialog({ open, title, children, confirmLabel, danger = false, working = false, onConfirm, onClose }: {
  open: boolean; title: string; children: ReactNode; confirmLabel: string; danger?: boolean; working?: boolean; onConfirm: () => void; onClose: () => void;
}) {
  const confirm = useRef<HTMLButtonElement | null>(null);
  const panel = useRef<HTMLElement | null>(null);
  useEffect(() => {
    if (!open) return;
    confirm.current?.focus();
    const keyboard = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && !working) onClose();
      if (event.key !== 'Tab' || !panel.current) return;
      const controls = [...panel.current.querySelectorAll<HTMLElement>('button:not(:disabled),a[href],input:not(:disabled),select:not(:disabled),textarea:not(:disabled)')];
      if (!controls.length) return;
      const first = controls[0], last = controls.at(-1)!;
      if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus(); }
      else if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus(); }
    };
    document.addEventListener('keydown', keyboard);
    return () => document.removeEventListener('keydown', keyboard);
  }, [open, working, onClose]);
  if (!open) return null;
  return <div className="dialog-backdrop" role="presentation" onMouseDown={event => { if (event.target === event.currentTarget && !working) onClose(); }}>
    <section ref={panel} className="confirmation-dialog" role="alertdialog" aria-modal="true" aria-labelledby="confirmation-title">
      <button className="dialog-close" onClick={onClose} disabled={working} aria-label="Close confirmation"><X /></button>
      <span className={`dialog-symbol ${danger ? 'danger' : ''}`}><AlertTriangle /></span><h2 id="confirmation-title">{title}</h2><div className="dialog-copy">{children}</div>
      <div className="dialog-actions"><button className="button secondary" onClick={onClose} disabled={working}>Go back</button><button ref={confirm} className={`button ${danger ? 'danger' : 'primary'}`} onClick={onConfirm} disabled={working}>{working ? 'Applying…' : confirmLabel}</button></div>
    </section>
  </div>;
}
