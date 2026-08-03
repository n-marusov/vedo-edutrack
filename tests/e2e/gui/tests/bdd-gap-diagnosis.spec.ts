import { test, expect } from "../fixtures/bdd";

// BDD: US-execute.gap.diagnose-root-cause / UC-execute.gap.diagnose-root-cause (M2)
// Ученик отстаёт по модулю → подъём по strict-связям → найден корневой модуль
// + цепочка + каскадное влияние.
test("root-cause gap diagnosis renders ranked root causes", async ({ authedPage }) => {
  await authedPage.goto("/dashboard/progress");
  await expect(authedPage.getByRole("heading", { name: "Progress & coverage" })).toBeVisible();

  // Gap diagnosis panel shows the fixture root cause (percent) with impact.
  const diagnosis = authedPage.getByTestId("gap-diagnosis");
  await expect(diagnosis).toBeVisible({ timeout: 10_000 });
  await expect(diagnosis).toContainText("percent");
});
