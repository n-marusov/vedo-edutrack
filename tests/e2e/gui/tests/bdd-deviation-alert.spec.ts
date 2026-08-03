import { test, expect } from "../fixtures/bdd";

// BDD: US-execute.alert.deviation + US-plan.recalculation.revise-delta (M6/M9)
// Отклонение >15% → план отклоняется → дашборд показывает статус риска
// (предложение пересмотра маршрута).
test("deviation alert surfaces an at-risk forecast", async ({ authedPage }) => {
  await authedPage.goto("/dashboard/progress");
  await expect(authedPage.getByRole("heading", { name: "Progress & coverage" })).toBeVisible();

  // Plan-vs-actual shows deviation flags and a readiness forecast badge.
  const list = authedPage.getByTestId("progress-list");
  await expect(list).toBeVisible({ timeout: 10_000 });
  await expect(list).toContainText("late");
});
