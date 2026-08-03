import { test, expect } from "@playwright/test";

// Smoke test — verifies the frontend scaffold serves the SPA with the
// correct document title (M0.3 scaffold is in place; placeholder replaced).
// Real MVP Must-scenarios (M1–M10) are stubbed in mvp-must-scenarios.spec.ts.
test("app scaffold renders with correct title", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/VEDO EduTrack/);
});
