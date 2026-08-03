import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';

export interface AuthUser {
  userId: string;
  roles: string[];
}

interface AuthState {
  /** Raw JWT (stored in sessionStorage via persist middleware). */
  token: string | null;
  user: AuthUser | null;
  roles: string[];
  isAuthenticated: boolean;
  login: (token: string, user: AuthUser) => void;
  logout: () => void;
  /** Stub: real silent refresh lands post-MVP (Keycloak). */
  refreshToken: () => Promise<string | null>;
}

/**
 * Auth store — JWT + user identity with sessionStorage persistence
 * (XSS-resistant: tokens do not survive tab closure; refresh tokens and
 * Keycloak flows are post-MVP).
 */
export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      token: null,
      user: null,
      roles: [],
      isAuthenticated: false,

      login: (token, user) => {
        console.info('[auth] login', { userId: user.userId, roles: user.roles });
        set({
          token,
          user,
          roles: user.roles,
          isAuthenticated: true,
        });
      },

      logout: () => {
        console.info('[auth] logout');
        set({
          token: null,
          user: null,
          roles: [],
          isAuthenticated: false,
        });
      },

      refreshToken: async () => {
        // Stub: return the same token (no refresh endpoint in M0.3).
        const { token } = get();
        if (token) {
          console.info('[auth] token refresh (stub)');
        }
        return token;
      },
    }),
    {
      name: 'edutrack-auth',
      storage: createJSONStorage(() => window.sessionStorage),
      partialize: (state) => ({
        token: state.token,
        user: state.user,
        roles: state.roles,
        isAuthenticated: state.isAuthenticated,
      }),
    },
  ),
);

/** useAuth hook — convenience accessor for the auth slice. */
export function useAuth() {
  const { user, roles, isAuthenticated, login, logout } = useAuthStore();
  return { user, roles, isAuthenticated, login, logout };
}
