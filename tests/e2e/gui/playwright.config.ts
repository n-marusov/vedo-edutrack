import { defineConfig, devices } from "@playwright/test";

// Playwright config — VEDO EduTrack E2E GUI (tests/e2e/gui).
// See ADR-DES.STACK.e2e-testing-playwright-vs-cypress and
// ADR-IMPL.PROCESS.repository-structure §5 (M1–M10 Must-scenarios).

export default defineConfig({
  testDir: "./tests",
  timeout: 30_000,
  // CI: retry flaky tests twice; local runs: no retries.
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI
    ? [["html", { outputFolder: "playwright-report" }], ["github"]]
    : [["html", { outputFolder: "playwright-report" }]],
  use: {
    // Test stack by default; dev stack via E2E_BASE_URL override.
    baseURL: process.env.E2E_BASE_URL ?? "http://localhost:55173",
    screenshot: "only-on-failure",
    trace: "retain-on-failure",
    video: "retain-on-failure",
  },
  projects: [
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
    // Firefox/WebKit: smoke runs for key scenarios (cross-browser for Enterprise).
    { name: "firefox", use: { ...devices["Desktop Firefox"] } },
    { name: "webkit", use: { ...devices["Desktop Safari"] } },
  ],
  // Local runs auto-start the TEST stack (reuse if already running); disabled
  // under the gate runner (deploy/ci/e2e-run.sh manages the lifecycle).
  webServer: process.env.E2E_STACK_MANAGED
    ? undefined
    : {
        command: "cd ../.. && make test-up",
        url: "http://localhost:55173",
        reuseExistingServer: !process.env.CI,
        timeout: 180_000,
      },
});
