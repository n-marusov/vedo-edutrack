---
archived: 2026-08-03
---
# Implementation Plan: M1 — Core Infrastructure (F0 + F1 + F2 + F3 + F6)

Branch: main (no branch — git.create_branches: false)
Created: 2026-08-03
Updated: 2026-08-03 (aif-improve: TDD restructuring — 6 test-first tasks added, 2 improved, 3 dependency fixes)

## Settings
- Testing: yes
- Logging: standard (INFO level — key events, service boundaries, errors)
- Docs: no (warn-only — `WARN [docs]` emitted, no mandatory checkpoint)
- Methodology: **TDD** — tests written BEFORE implementation, derived from `specs/requirements/` (FR acceptance criteria) and `specs/user-stories/` (Gherkin scenarios)

## Roadmap Linkage
Milestone: "M1: Core Infrastructure (F0 + F1 + F2 + F3 + F6)"
Rationale: M1 is the first Phase I milestone replacing scaffold stubs with real MVP infrastructure for ontology port, route engine, plan execution, resource catalog, and initial APIs — the foundation block for all subsequent milestones (M2, M3, M4).

## Commit Plan
- **Commit 1** (T1–T2): `feat(db): add Atlas migrations and repository layer for M1 schemas`
- **Commit 2** (T3–T5): `feat(ontology-port): real VEDO Hub GraphQL client with contract tests`
- **Commit 3** (T6–T9): `feat(route-planning): TDD pathfinder tests + Dijkstra engine + CLI commands`
- **Commit 4** (T10–T13): `feat(execution-gap): TDD gap/progress tests + diagnosis + coverage reports`
- **Commit 5** (T14–T15): `feat(resources): catalog domain and API endpoints`
- **Commit 6** (T16–T20): `feat(api-integration): contract tests + REST handlers + SPARQL + webhooks`
- **Commit 7** (T21–T23): `feat(frontend): ontology browser, route builder, progress and gap dashboards`
- **Commit 8** (T24–T27b): `test(quality): benchmarks, BDD, SBOM/vuln scan, integration and E2E tests`

---

## Progress

### Phase 1: Database Foundation
- [x] T1: Migrate CLI + Atlas schema DDL for 4 M1 schemas
- [x] T2: sqlc queries + pgx repository layer for all 4 schemas

### Phase 2: F0 — Ontology Port
- [x] T3: Real VEDO Hub GraphQL client replacing stub ontology
- [x] T4: Ontology sync CLI command + versioned subgraph cache
- [x] T5: Contract tests for VEDO Hub API boundary

### Phase 3: F1 — Route Planning Engine (TDD)
- [x] T6: Write pathfinder domain test suite (RED)
- [x] T7: Route domain engine — Dijkstra pathfinding
- [x] T8: Route application services + recompute triggers + plan fixation
- [x] T9: CLI: route compute (real DB mode) + plan get commands

### Phase 4: F2 — Plan Execution & Gap Coverage (TDD)
- [x] T10: Write gap & progress domain test suite (RED)
- [x] T11: Execution-progress domain — plan-vs-actual + forecast
- [x] T12: Gap-coverage domain — root-cause diagnosis + FGOS coverage
- [x] T13: Application services + CLI: gap diagnose and report

### Phase 5: F3 — Resources
- [x] T14: Resource domain — catalog, filters, module binding
- [x] T15: Resource application services + REST API handlers

### Phase 6: F6 — API & Integration Layer (TDD)
- [x] T16: Complete OpenAPI spec for M1 endpoints + regeneration
- [x] T17: Write API contract test suite (RED)
- [x] T18: Real API handlers replacing StubHandler
- [x] T19: Write SPARQL & webhook enforcement tests (RED)
- [x] T20: SPARQL read-only endpoint + webhook outbox infrastructure

### Phase 7: Frontend Features
- [x] T21: Ontology-port feature — graph browser
- [x] T22: Route-planning feature — route builder UI
- [x] T23: Execution-progress + gap-coverage + resources features

### Phase 8: Quality Gates
- [x] T24: Go performance benchmark suite
- [x] T25: SBOM generation + container vulnerability scanning
- [x] T26: BDD acceptance tests — Gherkin scenarios
- [x] T27a: Cross-module integration tests (testcontainers)
- [x] T27b: E2E API contract tests (Playwright)

---

## Tasks

### Phase 1: Database Foundation

- [x] **T1: Migrate CLI + Atlas schema DDL for 4 M1 schemas**
  Implement the `vedo-edutrack migrate` CLI subcommand (up/down/validate) wired to Atlas. Create DDL migrations for the 4 new bounded-context schemas:
  - `planmanagement` — learning plans, plan snapshots, checkpoints, FGOS constraints
  - `executionprogress` — module progress records, plan-vs-actual entries, deviation events
  - `gapcoverage` — gap diagnosis results, FGOS coverage snapshots, deficit lists
  - `resources` — resource catalog, module bindings, availability records
  Each schema follows the ADR schema-per-module convention with `version` columns for optimistic locking.
  **Files:** `backend/internal/cli/migrate.go`, `backend/migrations/000002_planmanagement_init.sql`, `backend/migrations/000003_executionprogress_init.sql`, `backend/migrations/000004_gapcoverage_init.sql`, `backend/migrations/000005_resources_init.sql`
  **Logging:**
  - `INFO [migrate] migration started` / `migration completed successfully` with affected count
  - `ERROR [migrate] migration failed: {reason}` with full error context
  - `DEBUG [migrate] applying migration <filename>` for each individual migration file

- [x] **T2: sqlc queries + pgx repository layer for all 4 schemas**
  Generate sqlc query files for CRUD operations on all new tables. Implement pgx-backed repository adapters in each bounded context's `adapters/repository/` package. Follow existing patterns from the identity-access module: `sqlc/queries/*.sql` → `sqlc/*.go` generated code → `repository.go` adapter wrapping generated queries.
  **Files:** `backend/internal/modules/planmanagement/adapters/repository/`, `backend/internal/modules/executionprogress/adapters/repository/`, `backend/internal/modules/gapcoverage/adapters/repository/`, `backend/internal/modules/resources/adapters/repository/`
  **Logging:**
  - `DEBUG [repository.<Module>] query: {operation} params: {...}` for each DB call
  - `INFO [repository.<Module>] connection pool stats: {stats}` on startup
  - `ERROR [repository.<Module>] query failed: {operation} error: {err}`

### Phase 2: F0 — Ontology Port (Real VEDO Hub Integration)

- [x] **T3: Real VEDO Hub GraphQL client replacing stub ontology**
  Replace the M0.3 stub (hardcoded 14 math topics) with a real GraphQL client in `ontologyport/adapters/hub/` that queries the VEDO Hub ontology-service GraphQL schema. Supported operations:
  - `graphNeighborhood` — traverse modules and links around a concept
  - `classDescendants` / `classTree` — resolve FGOS hierarchy, pedagogy concepts
  - `properties` / `property` — resolve module metadata, resource bindings, story links
  Implement domain types: `OntologyModule`, `OntologyLink` (5 types), `FgosBinding`, `PedagogyConcept`, `ResourceRef`, `StoryRef`.
  **Files:** `backend/internal/modules/ontologyport/domain/ontology.go`, `backend/internal/modules/ontologyport/adapters/hub/client.go`, `backend/internal/modules/ontologyport/adapters/hub/graph.go`
  **Constraints:** ≤3s timeout per NFR-FR-api.hub.read-ontology; 99.5% success rate tracked; Bearer token auth from config.
  **Logging:**
  - `INFO [ontologyport.hub] connected to Hub at {url}`
  - `DEBUG [ontologyport.hub] GraphQL query: {operation} variables: {...}`
  - `DEBUG [ontologyport.hub] GraphQL response: {moduleCount} modules, {duration}ms`
  - `ERROR [ontologyport.hub] query failed: {operation} status: {status} error: {err}`

- [x] **T4: Ontology sync CLI command + versioned subgraph cache**
  Implement `vedo-edutrack ontology sync` CLI command. Copies the relevant subgraph from VEDO Hub into an in-memory cache with ontology version tracking. Supports:
  - Full sync: fetch all modules for configured subjects (MVP: math, biology, physics, chemistry, history, literature, geography, computer science, social studies)
  - Incremental sync: diff against cached version, fetch only changed modules
  - Version stamp: `(ontologyId, version)` pair stored with cached graph
  Subgraph copy ≤5s for 10k modules per NFR-FR-api.hub.copy-subgraph.
  **Files:** `backend/internal/cli/ontology_sync.go`, `backend/internal/modules/ontologyport/application/sync_service.go`
  **Logging:**
  - `INFO [ontology.sync] sync started for ontology {id}`
  - `INFO [ontology.sync] fetched {count} modules, {linkCount} links in {duration}ms`
  - `INFO [ontology.sync] subgraph cached at version {version}`
  - `WARN [ontology.sync] incremental sync: {changedCount} modules changed`
  - `ERROR [ontology.sync] sync failed: {reason}`

- [x] **T5: Contract tests for VEDO Hub API boundary**
  Create contract tests verifying the Hub GraphQL boundary: schema compliance (all expected fields present), timeout behavior (≤3s abort), error handling (malformed queries, auth failures, ontology-not-found), pagination contracts. Tests run against the `hub-mock` container in CI.
  **Files:** `backend/tests/contract/hub_contract_test.go`
  **Logging:**
  - `INFO [contract.hub] test: {caseName} — PASS/FAIL`
  - `ERROR [contract.hub] contract violation: expected {field} in {query}`

### Phase 3: F1 — Route Planning Engine (TDD)

- [x] **T6: Write pathfinder domain test suite (TDD — RED phase)**
  Create table-driven test suite for the pathfinder BEFORE implementing the engine. Test cases derived from `specs/requirements/REQ-FR-plan.compute.shortest-path.md` acceptance criteria:
  - **Correctness:** match brute-force reference solution on ≥100 randomly generated graphs with known shortest paths
  - **Determinism:** same (position, goal, ontology, weights) → same Route ×10 consecutive calls
  - **Strict-prerequisite compliance:** 0 violations — no module appears before its strict prerequisites in any computed route
  - **Conflict resolution order:** on ≥50 conflict scenarios, verify essential+strict modules always chosen before soft, and soft before enrich
  - **Unreachable goal:** returns list of missing prerequisite modules + intermediate reachable goal (never empty response)
  - **Edge cases:** single-module graph, empty ontology, goal = position, goal with no paths
  - **Performance benchmark stub:** `testing.B` skeleton with 5k-module graph for later T24 activation
  All tests fail (RED) at commit — implementation in T7 makes them GREEN.
  **Files:** `backend/internal/modules/routeplanning/domain/pathfinder_test.go`, `backend/internal/modules/routeplanning/domain/testdata/` (random graph fixtures)
  **Logging:** N/A — pure domain tests, no logging needed

- [x] **T7: Route domain engine — Dijkstra pathfinding with link weights**
  Implement the core route computation domain — must pass all tests from T6:
  - `Route` aggregate: ordered list of modules from current position to goal
  - `Pathfinder`: Dijkstra-based shortest path with weighted edges — strict prerequisite (weight=1, unbreakable), soft prerequisite (weight=5), appliesTo (weight=20). enriches links excluded from MVP (per M1 spec).
  - Three horizons: far (full path) / mid (next N modules) / near (current module)
  - Essential core: mandatory modules derived from pedagogy concept
  - `ConceptWeightProfile`: pedagogy-concept-aware edge weight adjustments (deferred to Phase II but domain structure ready)
  - Pure domain layer — no I/O, no logging; all pure functions returning results + errors
  **Files:** `backend/internal/modules/routeplanning/domain/route.go`, `backend/internal/modules/routeplanning/domain/pathfinder.go`, `backend/internal/modules/routeplanning/domain/weights.go`, `backend/internal/modules/routeplanning/domain/horizon.go`
  **Acceptance gate:** `go test -count=1 ./internal/modules/routeplanning/domain/...` — all T6 tests GREEN

- [x] **T8: Route application services + recompute triggers + plan fixation**
  Implement application layer in `routeplanning/application/` and `planmanagement/application/`:
  - `ComputeService`: orchestrates ontology sync → pathfinding → horizon splitting. Accepts `(position, goal, pedagogyConcept, ontologyVersion) → Route`.
  - `RecomputeTrigger`: detects recompute conditions — progress change (>15% modules completed), goal change, ontology version bump. Emits `RouteRecalculationNeeded` domain event.
  - `PlanFixationService`: creates immutable `LearningPlan` snapshot from a computed route with timeline and checkpoints. Plan is INSERT-only (immutable per ADR), versioned.
  **Files:** `backend/internal/modules/routeplanning/application/compute_service.go`, `backend/internal/modules/routeplanning/application/trigger_service.go`, `backend/internal/modules/planmanagement/application/fixation_service.go`
  **Logging:**
  - `INFO [route.compute] computation started: learner={id} goal={moduleId}`
  - `INFO [route.compute] path found: {moduleCount} modules, {duration}ms`
  - `INFO [route.compute] recompute triggered: reason={progressChange|goalChange|ontologyUpdate}`
  - `INFO [plan.fixation] plan fixed: planId={id} modules={count} timeline={start}→{end}`
  - `WARN [route.compute] path not found: goal={moduleId} unreachable from position={moduleId}`
  - `ERROR [route.compute] computation failed: {reason}`

- [x] **T9: CLI: route compute (real DB mode) + plan get commands**
  Wire the route computation and plan management application services to CLI commands:
  - `vedo-edutrack route compute --learner=<id> --goal=<moduleId> [--pedagogy=<concept>] [--output json|table]` — real DB mode (postgres required), replaces M0.3 `--stub` inline logic
  - `vedo-edutrack plan get --plan=<id> [--learner=<id>] [--output json|table]` — retrieves fixed plan with timeline and checkpoints
  Both commands follow CLI-as-input-adapter pattern (ADR-DES.API.cli-interface): thin cobra commands that call application services.
  **Files:** `backend/internal/cli/route_compute.go`, `backend/internal/cli/plan_get.go`

### Phase 4: F2 — Plan Execution & Gap Coverage (TDD)

- [x] **T10: Write gap & progress domain test suite (TDD — RED phase)**
  Create test suites for execution-progress and gap-coverage domains BEFORE implementation. Test cases derived from FR acceptance criteria:
  **Gap diagnosis** (`REQ-FR-execute.gap.diagnose-root-cause`):
  - Root cause = first unmastered module when walking strict-prerequisite edges upward (≥5 distinct chains, 0 false roots)
  - Cascade ranking matches brute-force reference on ≥10 scenarios (ranking by N blocked modules, M subjects)
  - Reference scenario from vision.md: "Концентрация растворов" → root "Проценты" (70% mastery) → blocks chemistry, biology, social studies
  - Empty prerequisite chain → "no root cause found"
  - Performance: ≤2s on 1k-module graph
  **Progress tracking** (`REQ-FR-execute.progress.plan-vs-actual`):
  - 100% module deviation computed (planned vs actual dates in days)
  - `PlanDeviationDetected` emitted only when deviation exceeds ±15% of plan timeline (≥5 scenarios)
  - Module outside plan → excluded from comparison, flagged as divergence
  - No fixed plan → factual trajectory only, no deviation events
  **Binary forecast** (`REQ-FR-execute.forecast.binary-readiness`):
  - Accuracy ≥85% on simulated ≥100 completed plans
  - Insufficient data → "low accuracy" label, not error
  All tests fail (RED) at commit — implementation in T11–T12 makes them GREEN.
  **Files:** `backend/internal/modules/gapcoverage/domain/gap_test.go`, `backend/internal/modules/executionprogress/domain/progress_test.go`, `backend/internal/modules/executionprogress/domain/forecast_test.go`

- [x] **T11: Execution-progress domain — plan-vs-actual tracking + deviation + forecast**
  Implement `executionprogress/domain/` — must pass progress + forecast tests from T10:
  - Module progress tracking: mastery status per module (not_started / in_progress / mastered / skipped), timestamps, assessment results
  - Plan-vs-actual comparison: deviation in days per module, cause attribution
  - `PlanDeviationDetected` domain event: fired when deviation exceeds ±15% of plan timeline
  - Binary readiness forecast: on-track / not-on-track based on current velocity vs. remaining modules
  **Files:** `backend/internal/modules/executionprogress/domain/progress.go`, `backend/internal/modules/executionprogress/domain/deviation.go`, `backend/internal/modules/executionprogress/domain/forecast.go`
  **Logging:** (at application layer)
  - `INFO [execution.progress] module mastered: learner={id} module={id}`
  - `INFO [execution.progress] deviation detected: learner={id} actual={days} planned={days} delta={%}`
  - `WARN [execution.progress] at-risk forecast: learner={id} remaining={count} velocity={rate}`

- [x] **T12: Gap-coverage domain — root-cause diagnosis + FGOS coverage**
  Implement `gapcoverage/domain/` — must pass gap diagnosis tests from T10:
  - Root-cause gap diagnosis: walk strict-prerequisite edges upward from the lag point to find the first unmastered module. Rank all root causes by cascade impact (how many downstream modules are blocked). ≤2s for 1k modules per NFR-FR-execute.gap.diagnose-root-cause.
  - FGOS coverage live: compute coverage percentage from mastered modules × FGOS requirements mapping
  - Deficit list: uncovered FGOS requirements with blocking modules
  - Attestation readiness report: percentage + deficit list + forecast
  - Real-vs-formal knowledge mapping: mastered modules → FGOS coverage (not just plan-assigned)
  **Files:** `backend/internal/modules/gapcoverage/domain/gap.go`, `backend/internal/modules/gapcoverage/domain/coverage.go`, `backend/internal/modules/gapcoverage/domain/report.go`

- [x] **T13: Application services + CLI: gap diagnose and report commands**
  Wire execution-progress + gap-coverage application services to CLI:
  - `vedo-edutrack gap diagnose --learner=<id> [--plan=<id>]` — root-cause analysis with ranked gaps
  - `vedo-edutrack report --learner=<id> --type=coverage|attestation [--output json|table]` — FGOS coverage or attestation readiness report
  Implement application services: `ProgressService` (track, compare, forecast), `GapService` (diagnose), `CoverageService` (FGOS live, deficit, attestation).
  **Files:** `backend/internal/modules/executionprogress/application/progress_service.go`, `backend/internal/modules/gapcoverage/application/gap_service.go`, `backend/internal/modules/gapcoverage/application/coverage_service.go`, `backend/internal/cli/gap_diagnose.go`, `backend/internal/cli/report.go`
  **Logging:**
  - `INFO [gap.diagnose] diagnosis started: learner={id}`
  - `INFO [gap.diagnose] root causes found: {count} ranked by impact`
  - `INFO [coverage] FGOS coverage: {percentage}% ({covered}/{total})`
  - `INFO [report] attestation readiness: {percentage}% — {status}`

### Phase 5: F3 — Resources

- [x] **T14: Resource domain — catalog, filters, module binding**
  Implement `resources/domain/`:
  - `Resource` aggregate: type (video/text/interactive/textbook), format, source, style, difficulty level, duration, cost
  - `ResourceCatalog`: filter by format, source, difficulty; pagination; ≤200ms on 10k resources per NFR-FR-resources.catalog.filter-by-format
  - Module binding: associate resources with route modules via `appliesTo` / `enriches` ontology links
  - Availability check: resource available/unavailable with alternatives
  **Files:** `backend/internal/modules/resources/domain/resource.go`, `backend/internal/modules/resources/domain/catalog.go`
  **Acceptance gate:** Unit tests in `catalog_test.go` covering: filter precision/recall 100%, empty result handling, type ∈ {content, enabling} validation, cost ≥0 constraint, orphan-prevention (0 module bindings to non-existent modules).

- [x] **T15: Resource application services + REST API handlers**
  Implement application layer + oapi-codegen generated handlers:
  - Catalog queries with filters (format, style, difficulty, duration)
  - Module resource listing
  - Availability and cost endpoints
  Wire through the API handler replacing stub implementations.
  **Files:** `backend/internal/modules/resources/application/catalog_service.go`, `backend/internal/modules/resources/adapters/handler/`
  **Logging:**
  - `INFO [resources.catalog] query: filters={...} results={count} duration={ms}`
  - `DEBUG [resources.catalog] filter applied: {field}={value}`
  - `ERROR [resources.catalog] query failed: {reason}`

### Phase 6: F6 — API & Integration Layer (TDD for handlers + SPARQL/webhooks)

- [x] **T16: Complete OpenAPI spec for M1 endpoints + oapi-codegen regeneration**
  Extend `backend/api/openapi/v1.yaml` with all M1 REST endpoints:
  - `POST /routes/compute` — compute route (exists as stub, extend to full schema)
  - `GET /routes/{routeId}` — retrieve computed route
  - `POST /plans/fix` — fixate a route as immutable plan
  - `GET /plans/{planId}` — retrieve fixed plan with timeline
  - `GET /plans/{planId}/progress` — plan-vs-actual progress
  - `POST /progress/{learnerId}/modules/{moduleId}` — record module mastery
  - `GET /gaps/diagnose` — root-cause gap diagnosis
  - `GET /coverage/fgos` — FGOS coverage report
  - `GET /coverage/attestation` — attestation readiness report
  - `GET /resources` — resource catalog with filters
  - `GET /resources/{resourceId}` — single resource detail
  Regenerate `server.gen.go` and `types.gen.go` via `make gen`.
  **Files:** `backend/api/openapi/v1.yaml`
  **Note:** Once T16 is complete, frontend tasks (T21–T23) can begin with generated types and mock handlers — full integration requires T18.

- [x] **T17: Write API contract test suite (TDD — RED phase)**
  Create contract tests for all M1 REST endpoints BEFORE implementing handlers. Test cases derived from `specs/requirements/REQ-FR-api.rest.*.md`:
  - **Status codes:** 200 success, 400 invalid input, 401 missing token, 403 insufficient permissions, 404 not found, 429 rate limit with `Retry-After` header, 503 Hub unavailable
  - **Error format:** All errors use RFC 7807 Problem Details (`type`, `title`, `status`, `detail`, `instance`)
  - **Response schema:** Mandatory fields present for each endpoint (e.g., `RouteComputationResult` has `modules[]`, `horizons`, `ontologyVersion`)
  - **Determinism:** Same input ×3 consecutive calls → identical JSON response (byte-level)
  - **Rate limiting:** 429 + `Retry-After` header when exceeding configured limit
  Use Go `httptest.Server` with `oapi-codegen` generated `StrictServerInterface` for fast in-process tests + Playwright API tests for E2E contract validation.
  All tests fail (RED) at commit — implementation in T18 makes them GREEN.
  **Files:** `backend/internal/api/handler_contract_test.go`, `tests/e2e/api/contract/`

- [x] **T18: Real API handlers replacing StubHandler for all M1 endpoints**
  Replace the M0.3 `StubHandler` with real handlers wired to application services via `platform/wire.go` DI. Must pass all contract tests from T17.
  Each handler:
  - Validates JWT + RBAC permissions via existing auth middleware
  - Calls application service method
  - Maps domain types → OpenAPI response types
  - Returns RFC 7807 Problem Details on errors
  - p95 ≤ 200ms at 1000 concurrent per NFR-FR-integration.rest.compute-route
  **Files:** `backend/internal/api/handler.go`, `backend/internal/platform/wire.go`
  **Logging:**
  - `INFO [api] {method} {path} — {status} {duration}ms`
  - `WARN [api] {method} {path} — validation error: {field}={reason}`
  - `ERROR [api] {method} {path} — internal error: {err}` (with trace_id)
  **Acceptance gate:** `go test -count=1 ./internal/api/...` — all T17 contract tests GREEN

- [x] **T19: Write SPARQL & webhook enforcement tests (TDD — RED phase)**
  Create enforcement tests BEFORE implementing SPARQL endpoint and webhook outbox. Derived from FR specs:
  **SPARQL** (`REQ-FR-api.sparql.read-only`):
  - Mutating queries (INSERT/DELETE/CREATE/DROP) → 403
  - SELECT/ASK/DESCRIBE execute successfully; triple count identical before/after 50 queries
  - Rate limit 100 req/min → 429 + `Retry-After`
  - Timeout >30s → 504
  - Result truncation >10k rows → `truncated=true` flag
  **Webhooks** (`REQ-FR-api.webhooks.idempotency`, `REQ-FR-api.webhooks.module-mastered`, `REQ-FR-api.webhooks.plan-deviated`):
  - event_id is UUID v4, unchanged across retries
  - Duplicate event_id → 200 with `duplicate=true`, no state change
  - Exponential retry: 5 attempts at 1,2,4,8,16 min intervals
  - Exhausted retries → DLQ entry + alert
  - Concurrent delivery → unique constraint on `(source, event_id)` prevents duplicates
  - Payload schema: all mandatory fields present per FR specs
  All tests fail (RED) at commit — implementation in T20 makes them GREEN.
  **Files:** `backend/internal/modules/integrations/adapters/sparql/handler_test.go`, `backend/internal/modules/integrations/adapters/webhook/outbox_test.go`

- [x] **T20: SPARQL read-only endpoint + webhook outbox infrastructure**
  Implement integration infrastructure — must pass all tests from T19:
  - SPARQL endpoint: SELECT/ASK/DESCRIBE only, parameterized queries, rate-limited to 10 req/min per ADR-DES.API.communication-patterns
  - Webhook outbox: PostgreSQL-backed outbox table (separate migration), domain events → outbox rows → async dispatch, `(source, event_id)` unique constraint for idempotency
  - Event types: `module.mastered`, `plan.deviated`, `route.recalculated`
  **Files:** `backend/migrations/000006_integrations_outbox.sql`, `backend/internal/modules/integrations/adapters/sparql/handler.go`, `backend/internal/modules/integrations/adapters/webhook/outbox.go`, `backend/internal/modules/integrations/adapters/webhook/dispatcher.go`
  **Logging:**
  - `INFO [webhook] event queued: type={eventType} id={eventId}`
  - `INFO [webhook] dispatched: type={eventType} target={url} status={code}`
  - `WARN [webhook] retry: attempt={n} type={eventType}`
  - `ERROR [webhook] dispatch failed (exhausted): type={eventType} target={url}`
  - `INFO [sparql] query executed: {duration}ms resultSize={rows}`

### Phase 7: Frontend Features

> **Dependency note:** T21–T23 can begin after T16 (OpenAPI spec finalized — types and mock handlers available). Full backend integration requires T18 (real API handlers).

- [x] **T21: Ontology-port feature — graph browser component**
  Implement `frontend/src/features/ontology-port/` with real data from backend API:
  - `OntologyBrowser` component: renders ontology graph (React Flow) with modules as nodes and link types as color-coded edges
  - Module detail panel: metadata, FGOS bindings, linked resources
  - API hooks: `useOntology()` — fetches graph data via `/ontology/concepts` API
  - Replace M0.3 stub placeholder tests with real Vitest + RTL tests
  **Files:** `frontend/src/features/ontology-port/`, `frontend/src/features/ontology-port/__tests__/`
  **Pattern:** Follow existing shared components (Card, Badge, LoadingSpinner); use Zustand for local feature state; lazy-loaded route in `routes.tsx`

- [x] **T22: Route-planning feature — route builder UI**
  Implement `frontend/src/features/route-planning/`:
  - Goal selector: searchable module picker from ontology
  - Route builder: triggers `POST /routes/compute`, renders computed path as timeline + graph
  - Plan fixation UI: review route → confirm → creates immutable plan
  - Three-horizon toggle: far (full) / mid / near view of the route
  - Recompute indicator: shows when route needs recalculation
  **Files:** `frontend/src/features/route-planning/`, `frontend/src/features/route-planning/__tests__/`

- [x] **T23: Execution-progress + gap-coverage + resources features**
  Implement remaining feature modules:
  - `execution-progress`: plan-vs-actual comparison dashboard (progress bars, deviation flags, binary forecast badge)
  - `gap-coverage`: gap viewer (color-coded graph with root-cause markers), FGOS coverage radial chart, attestation readiness panel
  - `resources`: catalog with format/source/difficulty filters, module-linked resource list, availability badges
  Replace `RouteView`, `PlanView`, `ProgressView` pages with real feature components (remove "coming soon (M1)" stubs).
  **Files:** `frontend/src/features/execution-progress/`, `frontend/src/features/gap-coverage/`, `frontend/src/features/resources/`, `frontend/src/pages/`

### Phase 8: Quality Gates — Benchmarks, BDD, Security, Integration

- [x] **T24: Go performance benchmark suite**
  Create `testing.B` benchmarks for NFR-critical operations. Must run in CI (advisory gate in Phase I, blocking in Phase II). Thresholds from FR acceptance criteria:
  - Pathfinder on 5k-module graph: ≤1s per `REQ-FR-plan.compute.shortest-path` AC-4
  - Gap diagnosis on 1k-module graph: ≤2s per `REQ-FR-execute.gap.diagnose-root-cause` AC-4
  - Resource catalog filter on 10k resources: ≤200ms per `REQ-FR-resource.catalog.filter-by-format` AC-4
  - Subgraph copy on 10k modules: ≤5s per `REQ-FR-api.hub.copy-subgraph` AC-3
  Add `make bench` target; integrate into `deploy/ci/gates.yaml` as advisory `perf-bench` gate.
  **Files:** `backend/internal/modules/routeplanning/domain/pathfinder_bench_test.go`, `backend/internal/modules/gapcoverage/domain/gap_bench_test.go`, `backend/internal/modules/resources/domain/catalog_bench_test.go`, `backend/internal/modules/ontologyport/adapters/hub/copy_bench_test.go`

- [x] **T25: SBOM generation + container vulnerability scanning in CI**
  Integrate security scanning into CI pipeline (independent of feature completion — can run early):
  - SBOM: `syft` generates SPDX SBOM for the backend binary, attached as CI artifact
  - Vulnerability scan: `trivy` (or `grype`) scans the distroless container image; gate blocks on CRITICAL findings
  - Image attestation: cosign signs the container image
  - Add `security` group to `deploy/ci/gates.yaml` with blocking severity
  **Files:** `.github/workflows/ci.yml`, `deploy/ci/gates.yaml`
  **Logging:** CI-native (GitHub Actions step logs)

- [x] **T26: BDD acceptance tests — Gherkin scenarios for top-4 M1 user stories**
  Translate the 4 most critical M1 user story scenarios into executable Playwright E2E tests:
  1. **Shortest path computation** (`US-plan.compute.shortest-path`): родитель выбирает цель → система вычисляет кратчайший путь с strict-пререквизитами; недостижимая цель → сообщение об ошибке
  2. **Plan fixation** (`US-plan.fixation.snapshot`): наступает дата Checkpoint → система фиксирует snapshot → GUI показывает два слоя (план + актуальный маршрут)
  3. **Root-cause gap diagnosis** (`US-execute.gap.diagnose-root-cause`): ученик отстаёт по модулю M → подъём по strict-связям → найден корневой модуль + цепочка + каскадное влияние
  4. **Deviation alert** (`US-execute.alert.deviation` + `US-plan.recalculation.revise-delta`): отклонение >15% → plan.deviated webhook + уведомление + предложение пересмотра
  Each scenario 100% mapped to FR acceptance criteria. Use existing Playwright `fixtures/auth.ts` pattern.
  **Files:** `tests/e2e/gui/tests/bdd-shortest-path.spec.ts`, `tests/e2e/gui/tests/bdd-plan-fixation.spec.ts`, `tests/e2e/gui/tests/bdd-gap-diagnosis.spec.ts`, `tests/e2e/gui/tests/bdd-deviation-alert.spec.ts`

- [x] **T27a: Cross-module integration tests (testcontainers)**
  Activate the red `backend/tests/integration_test.go` scaffold with real testcontainers-go tests. Happy-path flow:
  - Start postgres + hub-mock containers → ontology sync → route compute → plan fixation → module mastery → progress tracking → gap diagnosis → FGOS coverage report
  - Verify each step produces correct domain events (in-process event bus)
  - Replace all `t.Skip("TODO...")` stubs
  **Files:** `backend/tests/integration/`

- [x] **T27b: E2E API contract tests (Playwright)**
  Activate `tests/e2e/api/` with Playwright API testing:
  - All M1 REST endpoints: compute route → fix plan → track progress → diagnose gap → coverage report → resource catalog
  - Error paths: invalid input (400), missing auth (401), forbidden (403), not found (404), rate limit (429)
  - Webhook delivery verification: simulate `module.mastered` → webhook dispatched → idempotent retry
  - SPARQL read-only guard: mutation attempt → 403
  Run against full compose stack (`make up`).
  **Files:** `tests/e2e/api/`
