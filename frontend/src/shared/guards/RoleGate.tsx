import type { ReactNode } from 'react';
import { useAuthStore } from '../../store/authStore';

interface RoleGateProps {
  /** Required role (or roles); children render only when the user has one. */
  requiredRole: string | string[];
  children: ReactNode;
  /** Optional 403 fallback; default renders a message. */
  fallback?: ReactNode;
}

/** RoleGate — renders children only when the user has the required role. */
export function RoleGate({ requiredRole, children, fallback }: RoleGateProps) {
  const roles = useAuthStore((s) => s.roles);
  const required = Array.isArray(requiredRole) ? requiredRole : [requiredRole];
  const allowed = required.some((r) => roles.includes(r));

  if (!allowed) {
    if (fallback) return <>{fallback}</>;
    console.warn('[guard] access denied', { required, roles });
    return (
      <div className="flex min-h-[40vh] items-center justify-center">
        <div className="text-center">
          <p className="text-2xl font-semibold text-gray-900">403</p>
          <p className="mt-1 text-sm text-gray-600">You don't have permission to view this page.</p>
        </div>
      </div>
    );
  }
  return <>{children}</>;
}
