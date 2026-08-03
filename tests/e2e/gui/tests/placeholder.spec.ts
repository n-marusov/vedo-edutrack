import { test, expect } from "@playwright/test";

test("placeholder — fails until M0.3 scaffold is ready", async ({ page }) => {
  test.fail(true, "TODO: implement E2E tests for MVP Must-scenarios (M0.3+)");
  await page.goto("/");
  await expect(page).toHaveTitle(/VEDO EduTrack/);
});
