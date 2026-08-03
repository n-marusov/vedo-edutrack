import { test as base } from "@playwright/test";

// Auth fixture skeleton — M0.3+ Keycloak/JWT integration.
// See ADR-DES.STACK.e2e-testing-playwright-vs-cypress (auth): JWT login via
// UI or token injection through storageState (no repeated logins per scenario).

export interface AuthContext {
  /** Injected JWT access token (null until M0.3 wires the identity provider). */
  token: string | null;
  /** MVP persona (RBAC T8): learner, parent, teacher, admin. */
  role: "learner" | "parent" | "teacher" | "admin";
}

export const test = base.extend<{ auth: AuthContext }>({
  auth: async ({}, use) => {
    // TODO(M0.3): inject JWT via storageState or perform a UI login
    // (Keycloak realm from deploy/keycloak).
    await use({ token: null, role: "learner" });
  },
});

export { expect } from "@playwright/test";
