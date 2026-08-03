import { test as base, expect, type Page } from "@playwright/test";

// BDD E2E helper (T26) — authenticates through the real demo login UI
// (frontend/src/pages/Login.tsx), which mints a dev JWT via /auth/token.

/**
 * BDD test — base extended with an authenticated page (learner role).
 * Logs in through the demo login form, then follows the redirect target.
 */
export const test = base.extend<{ authedPage: Page }>({
  authedPage: async ({ page }, use) => {
    await loginAs(page, "e2e-learner", "learner");
    await use(page);
  },
});

/** Log in through the demo login UI (user_id + role → JWT → /dashboard). */
export async function loginAs(page: Page, userId: string, role: string): Promise<void> {
  await page.goto("/login");
  await page.getByLabel("User ID").fill(userId);
  await page.getByLabel("Role").selectOption({ label: role });
  await page.getByRole("button", { name: "Sign in" }).click();
  // Redirect lands on the dashboard.
  await page.waitForURL("**/dashboard", { timeout: 15_000 });
}

export { expect };
