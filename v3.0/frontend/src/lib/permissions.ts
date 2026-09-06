import type { OperatorRole } from '../types/domain';

export type Permission = 'readFleet' | 'manageRegistry' | 'orchestrate' | 'adminRetry' | 'viewAudit' | 'selectTenant'
  | 'readTasks' | 'createTask' | 'cancelTask' | 'retryTask' | 'readCommands' | 'cancelCommand' | 'retryCommand';

const permissions: Record<OperatorRole, ReadonlySet<Permission>> = {
  PLATFORM_ADMIN: new Set(['readFleet', 'manageRegistry', 'orchestrate', 'adminRetry', 'viewAudit', 'selectTenant', 'readTasks', 'createTask', 'cancelTask', 'retryTask', 'readCommands', 'cancelCommand', 'retryCommand']),
  TENANT_ADMIN: new Set(['readFleet', 'manageRegistry', 'orchestrate', 'adminRetry', 'viewAudit', 'readTasks', 'createTask', 'cancelTask', 'retryTask', 'readCommands', 'cancelCommand', 'retryCommand']),
  OPERATOR: new Set(['readFleet', 'orchestrate', 'readTasks', 'createTask', 'cancelTask', 'readCommands', 'cancelCommand']),
  VIEWER: new Set(['readFleet', 'readTasks', 'readCommands']),
};

export function can(role: OperatorRole, permission: Permission) {
  return permissions[role].has(permission);
}

export function roleLabel(role: OperatorRole) {
  return ({
    PLATFORM_ADMIN: 'Platform administrator',
    TENANT_ADMIN: 'Tenant administrator',
    OPERATOR: 'Operator',
    VIEWER: 'Viewer',
  } satisfies Record<OperatorRole, string>)[role];
}
