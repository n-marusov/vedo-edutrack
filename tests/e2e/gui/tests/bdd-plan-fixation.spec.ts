import { test, expect } from "../fixtures/bdd";

// BDD: US-plan.fixation.snapshot / UC-plan.fixation.snapshot-plan (M7)
// Наступает дата Checkpoint → система фиксирует snapshot → GUI показывает
// план и актуальный маршрут как два слоя.
test("plan fixation shows the immutable plan layer", async ({ authedPage }) => {
  await authedPage.goto("/dashboard/plan");
  await expect(authedPage.getByRole("heading", { name: "My Plan" })).toBeVisible({ timeout: 10_000 });

  // The route layer shows the plan-fixation entry point (F1 plan fixation).
  await authedPage.goto("/dashboard/route");
  await expect(authedPage.getByRole("heading", { name: "Route Builder" })).toBeVisible({ timeout: 10_000 });
  await expect(authedPage.getByText("Plan fixation", { exact: true })).toBeVisible({ timeout: 10_000 });
});
