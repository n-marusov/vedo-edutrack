import { test, expect } from "../fixtures/bdd";

// BDD: US-plan.compute.shortest-path / UC-plan.compute.shortest-path-to-goal (M1)
// Родитель выбирает цель → система вычисляет кратчайший путь с strict-пререквизитами;
// недостижимая цель → сообщение об ошибке.
test("shortest path with strict prerequisites is computed", async ({ authedPage }) => {
  await authedPage.goto("/dashboard/route");
  await expect(authedPage.getByRole("heading", { name: "Route Builder" })).toBeVisible();

  const goalInput = authedPage.getByLabel("Goal module id");
  await goalInput.fill("math-5-11");
  await authedPage.getByRole("button", { name: "Compute route" }).click();

  // The stub backend computes the strict-prerequisite chain for math-5-11.
  const timeline = authedPage.getByTestId("route-timeline");
  await expect(timeline).toBeVisible({ timeout: 10_000 });
  await expect(timeline).toContainText("math-5-11");
});

test("unreachable goal surfaces an error", async ({ authedPage }) => {
  await authedPage.goto("/dashboard/route");
  const goalInput = authedPage.getByLabel("Goal module id");
  await goalInput.fill("unreachable-module");
  await authedPage.getByRole("button", { name: "Compute route" }).click();
  await expect(authedPage.getByRole("alert")).toBeVisible({ timeout: 10_000 });
});
