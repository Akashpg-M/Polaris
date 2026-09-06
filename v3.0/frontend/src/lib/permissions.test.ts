import { describe, expect, it } from 'vitest';
import { can } from './permissions';

describe('presentation permissions', () => {
  it('keeps viewer access read-only', () => {
    expect(can('VIEWER', 'readFleet')).toBe(true);
    expect(can('VIEWER', 'orchestrate')).toBe(false);
    expect(can('VIEWER', 'manageRegistry')).toBe(false);
  });

  it('reserves cross-tenant selection for platform administrators', () => {
    expect(can('PLATFORM_ADMIN', 'selectTenant')).toBe(true);
    expect(can('TENANT_ADMIN', 'selectTenant')).toBe(false);
  });

  it('allows operators to orchestrate but reserves retries for administrators', () => {
    expect(can('OPERATOR', 'createTask')).toBe(true);
    expect(can('OPERATOR', 'cancelCommand')).toBe(true);
    expect(can('OPERATOR', 'retryTask')).toBe(false);
    expect(can('TENANT_ADMIN', 'retryCommand')).toBe(true);
  });
});
