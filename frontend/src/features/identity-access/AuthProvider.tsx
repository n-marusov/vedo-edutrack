import { decodeJwt } from 'jose';
import { useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { useAuth, useAuthStore } from '../../store/authStore';

/**
 * AuthProvider — mounts once, validates any persisted token (JWT `exp`
 * claim via jose) and logs the user out when the token is expired.
 * Exposes the auth state through useAuth().
 */
export function AuthProvider({ children }: { children: ReactNode }) {
  const [ready, setReady] = useState(false);
  const logout = useAuthStore((s) => s.logout);

  useEffect(() => {
    const token = useAuthStore.getState().token;
    if (token) {
      try {
        const claims = decodeJwt(token);
        const exp = claims.exp;
        if (exp !== undefined && exp * 1000 <= Date.now()) {
          console.warn('[auth] token expired — logging out');
          logout();
        }
      } catch (err) {
        console.warn('[auth] failed to decode stored token — logging out', err);
        logout();
      }
    }
    setReady(true);
  }, [logout]);

  if (!ready) {
    return null;
  }

  return <AuthBridge>{children}</AuthBridge>;
}

/**
 * AuthBridge is a placeholder that keeps the hook boundary stable:
 * consumers call useAuth() directly from the store (single source of truth).
 * Wrapping in a component here ensures children re-render on store changes.
 */
function AuthBridge({ children }: { children: ReactNode }) {
  useAuth();
  return children;
}

export { useAuth } from '../../store/authStore';
