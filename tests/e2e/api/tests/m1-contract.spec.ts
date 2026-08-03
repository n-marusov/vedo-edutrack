import { expect, test } from "@playwright/test";
import type { APIRequestContext } from "@playwright/test";

// E2E API contract tests (T27b) — run against the compose stack (`make up`).
// Flows: health → dev token → route compute → resource catalog → gap/coverage.

// Absolute base — the Playwright request fixture resolves paths against this.
const API_BASE = process.env.E2E_API_URL ?? "http://localhost:8080/api/v1";
// Health/readiness live at the router root (public, outside /api/v1 auth group).
const ROOT_BASE = process.env.E2E_ROOT_URL ?? "http://localhost:8080";

/** Dev token for a persona (RBAC T8). */
async function devToken(request: APIRequestContext, role = "learner"): Promise<string> {
  const response = await request.post(`${API_BASE}/auth/token`, {
    data: { user_id: `e2e-api-${role}`, roles: [role] },
  });
  expect(response.status()).toBe(200);
  const body = (await response.json()) as { access_token: string };
  return body.access_token;
}

test("health and readiness endpoints respond", async ({ request }) => {
  const health = await request.get(`${ROOT_BASE}/healthz`);
  expect(health.status()).toBe(200);
  const body = (await health.json()) as { status: string };
  expect(body.status).toBe("ok");
});

test("route compute returns a strict-prerequisite chain", async ({ request }) => {
  const token = await devToken(request);
  const response = await request.post(`${API_BASE}/routes/compute`, {
    headers: { Authorization: `Bearer ${token}` },
    data: { learner_id: "e2e-api-learner", goal_topic_id: "math-5-11" },
  });
  expect(response.status()).toBe(200);
  const body = (await response.json()) as { route: { topic_id: string; order: number }[] };
  expect(body.route.length).toBeGreaterThan(0);
  expect(body.route[0].topic_id).toBe("math-5-1");
  expect(body.route.at(-1)?.topic_id).toBe("math-5-11");
});

test("resource catalog is queryable", async ({ request }) => {
  const token = await devToken(request);
  const response = await request.get(`${API_BASE}/resources?format=video`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  expect(response.status()).toBe(200);
  const body = (await response.json()) as { items: unknown[]; total: number };
  expect(body.total).toBeGreaterThan(0);
});

test("FGOS coverage returns a valid report", async ({ request }) => {
  const token = await devToken(request);
  const response = await request.get(`${API_BASE}/learners/l1/coverage/fgos`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  expect(response.status()).toBe(200);
  const body = (await response.json()) as { covered: number; total: number; percent: number };
  expect(body.total).toBeGreaterThan(0);
  expect(body.percent).toBeGreaterThanOrEqual(0);
  expect(body.percent).toBeLessThanOrEqual(100);
});

test("gap diagnosis returns root causes", async ({ request }) => {
  const token = await devToken(request);
  const response = await request.get(`${API_BASE}/learners/l1/gaps?lag_module_id=chemistry`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  expect(response.status()).toBe(200);
  const body = (await response.json()) as { status: string; root_causes: { module_id: string }[] };
  expect(body.root_causes.length).toBeGreaterThan(0);
  expect(body.root_causes[0].module_id).toBe("percent");
});

test("SPARQL read-only guard rejects mutations", async ({ request }) => {
  const token = await devToken(request);
  const response = await request.get(`${API_BASE}/sparql?query=INSERT%20DATA%20%7B%3Cs%3E%20%3Cp%3E%20%3Co%3E%7D`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  expect(response.status()).toBe(400);
});

test("missing auth token is rejected", async ({ request }) => {
  const response = await request.get(`${API_BASE}/resources`);
  expect([401, 403]).toContain(response.status());
});
