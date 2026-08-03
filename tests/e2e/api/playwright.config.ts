import { defineConfig } from "@playwright/test";

// Playwright API config — VEDO EduTrack E2E API contract tests (tests/e2e/api, T27b).
// Run against the full compose stack (`make up`): backend on :8080.
export default defineConfig({
  testDir: "./tests",
  timeout: 30_000,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI ? [["github"]] : [["list"]],
  use: {
    baseURL: process.env.E2E_API_URL ?? "http://localhost:8080/api/v1",
    extraHTTPHeaders: {
      Accept: "application/json",
    },
  },
  projects: [{ name: "api", testMatch: /.*\.spec\.ts/ }],
  // Auto-start the dev stack for local runs; disabled under the gate runner
  // (deploy/ci/e2e-run.sh manages the stack lifecycle via E2E_STACK_MANAGED).
  webServer: process.env.E2E_STACK_MANAGED
    ? undefined
    : {
        command: "cd ../.. && make up",
        url: "http://localhost:8080/healthz",
        reuseExistingServer: !process.env.CI,
        timeout: 120_000,
      },
});
