# Implementation Plan: M2 — Family Education «Дай пять» (F2 + F4 + F5)

Branch: none (create_branches: false)
Created: 2026-08-03

## Settings
- Testing: yes
- Logging: standard (INFO level, key events only)
- Docs: no (WARN [docs] only, no mandatory checkpoint)

## Roadmap Linkage
Milestone: "M2: Family Education — «Дай пять» (F2 + F4 + F5)"
Rationale: This plan delivers the M2 Community-facing family education product slice: plan execution (F2), visualization (F4), and real-world connection (F5) on top of the M1 core infrastructure.

## Research Context
Source: .ai-factory/RESEARCH.md (Active Summary)

**Not directly relevant** — the active research context covers VEDO Hub mock server (infrastructure).
M2 requirements are fully specified in the formal requirement corpus (`specs/requirements/REQ-FR-execute.*.md`,
`REQ-FR-viz.*.md`, `REQ-FR-practice.*.md`) and MVP acceptance criteria.

## Commit Plan
- **Commit 1** (after tasks 1-3): `feat(backend): add execution-progress and gap-coverage domain model, migrations, and plan-vs-actual service`
- **Commit 2** (after tasks 4-6): `feat(backend): implement gap diagnosis, FGOS coverage, and F2 REST API endpoints`
- **Commit 3** (after tasks 7-8): `feat(backend): add practice-life domain layer and CLI commands`
- **Commit 4** (after tasks 9-11): `feat(frontend): implement knowledge map visualization and dashboards`
- **Commit 5** (after tasks 12-14): `feat(frontend): add practice-life UI and group management panel`
- **Commit 6** (after tasks 15-16): `test: integration tests and component tests for M2 features`

## Tasks

### Phase 1: Backend Foundation — Domain Model & Persistence (F2)

- [x] **Task 1: Database migrations and domain types for execution-progress and gap-coverage** ✅ M1

  **Deliverable:** PostgreSQL schema migrations for execution progress data (mastery records, plan snapshots, deviations, forecasts) plus gap/coverage data (gap diagnoses, FGOS coverage snapshots). Go domain types (value objects, aggregates, domain service interfaces) for both `executionprogress` and `gapcoverage` modules.

  **Files to create/modify:**
  - `backend/migrations/<timestamp>_execution_progress.sql` — new migration: `mastery_records`, `plan_snapshots`, `deviations`, `forecasts`, `coverage_reports`, `gap_diagnoses`, `deficit_lists` tables
  - `backend/internal/modules/executionprogress/domain/types.go` — VOs: `MasteryRecord` (module_id, level, timestamp, source), `Deviation` (step_id, days, reason), `Forecast` (status, expected_end, risks)
  - `backend/internal/modules/executionprogress/domain/services.go` — interfaces: `TrajectoryService`, `PlanVsActualService`, `ForecastService`, `DeviationAlertService`
  - `backend/internal/modules/gapcoverage/domain/types.go` — VOs: `GapDiagnosis` (root_module, cascade_modules, impact), `CoverageReport` (framework, percentage, deficits), `AttestationReport` (verdict, coverage_by_domain, critical_path)
  - `backend/internal/modules/gapcoverage/domain/services.go` — interfaces: `GapDiagnosisService`, `CoverageService`, `DeficitService`, `AttestationReadinessService`

  **Logging requirements:**
  - INFO: migration applied, domain types initialized
  - INFO: each new VO/aggregate creation

  **Dependencies:** M1 database schema (learners, modules, plans baseline).

- [x] **Task 2: SQLC queries and repository adapters for execution-progress and gap-coverage** ✅ M1

  **Deliverable:** SQLC-generated Go code for all CRUD operations on execution progress and gap/coverage tables. Repository implementations implementing domain interfaces.

  **Files to create/modify:**
  - `backend/internal/modules/executionprogress/adapters/repository/sqlc/` — sqlc config + queries: insert/get mastery_records, insert/get deviations, upsert forecasts, list plans
  - `backend/internal/modules/executionprogress/adapters/repository/postgres.go` — repository adapter implementing domain `Repository` interfaces
  - `backend/internal/modules/gapcoverage/adapters/repository/sqlc/` — sqlc config + queries: insert/get gap_diagnoses, upsert coverage_reports, upsert deficit_lists, insert/get attestation_reports
  - `backend/internal/modules/gapcoverage/adapters/repository/postgres.go` — repository adapter
  - Update `backend/internal/platform/postgres/` — add migration runner (if not exists)

  **Logging requirements:**
  - INFO: repository connection established, query execution (timing)
  - ERROR: database errors with full context (query, params, error)
  - WARN: row-not-found fallthrough (non-error in domain)

  **Dependencies:** Task 1 (migrations must be applied first).

- [x] **Task 3: Plan-vs-actual comparison and binary readiness forecast application services** ✅ Domain done, forecast/alert services added

  **Deliverable:** Application services that compare fixed plan snapshots against actual mastery records, compute deviations, and produce binary readiness forecasts (on-track / not-on-track).

  **Done:** `progress.go` domain + `progress_service.go` app (M1). Added `forecast_service.go` and `deviation_alert_service.go` for M2 gaps.

  **Deliverable:** Application services that compare fixed plan snapshots against actual mastery records, compute deviations, and produce binary readiness forecasts (on-track / not-on-track).

  **Files to create/modify:**
  - `backend/internal/modules/executionprogress/application/trajectory_service.go` — `RecordMastery(ctx, learnerID, moduleID, level)` — writes mastery record, fires `ModuleMastered` event
  - `backend/internal/modules/executionprogress/application/plan_vs_actual_service.go` — `Compare(ctx, learnerID, planID)` — iterates plan steps, matches against mastery records, computes deviation in days, assigns reason from enum {acceleration, more_practice, pause, volume_change, unspecified}
  - `backend/internal/modules/executionprogress/application/forecast_service.go` — `ForecastReadiness(ctx, learnerID, checkpointDate)` — computes remaining modules / current pace → status (on-track / not-on-track), expected completion date
  - `backend/internal/modules/executionprogress/application/deviation_alert_service.go` — evaluates `PlanDeviationDetected` events, decides whether to publish alert based on configured threshold (N days)
  - `backend/internal/modules/executionprogress/application/events.go` — domain event types: `ModuleMastered`, `PlanDeviationDetected`

  **Acceptance criteria (from specs):**
  - Deviation reason assignment: "unspecified" ≤ 5% of cases
  - `PlanDeviationDetected` ONLY when deviation exceeds ±15% or N-day threshold
  - Forecast accuracy ≥ 85% on ≥ 100 completed plans
  - Report refresh ≤ 5 s after `module.mastered`

  **Logging requirements:**
  - INFO: each mastery recorded (learner, module, level), deviation computed (plan_id, step, days), forecast result (learner, status, expected_end)
  - WARN: deviation detected (step, days, reason)
  - ERROR: plan not found, invalid mastery level

  **Dependencies:** Task 2 (repositories must exist).

<!-- Commit checkpoint: tasks 1-3 -->

### Phase 2: Backend Gap Diagnosis & Coverage (F2)

- [x] **Task 4: Gap diagnosis — root-cause analysis with cascade impact** ✅ M1 domain done

  **Deliverable:** Algorithm that, given a lagging module, climbs `hasStrictPrerequisite` edges up the graph from the ontology-port subgraph cache to find the first unmastered module (root cause).

  **Done:** `domain/gap.go` (Graph, Module, Link, Mastery, DiagnoseRootCause), `application/gap_service.go` (GapService.Diagnose). Core algorithm and ranking already implemented. Cascade impact computation and multi-subject ranking are TODO for M2 depth.

  **Deliverable:** Algorithm that, given a lagging module, climbs `hasStrictPrerequisite` edges up the graph from the ontology-port subgraph cache to find the first unmastered module (root cause). Ranks root gaps by cascade impact (N blocked modules across M subjects).

  **Files to create/modify:**
  - `backend/internal/modules/gapcoverage/application/gap_diagnosis_service.go` — `DiagnoseRootCause(ctx, learnerID, lagModuleID)`:
    1. Get learner's mastered module set from execution-progress
    2. Walk `hasStrictPrerequisite` graph upward from lag module
    3. Return first unmastered node per chain (root gap)
    4. Compute cascade: for each root, count descendant modules + subjects blocked
    5. Rank by cascade impact descending
    6. Fallback: if all prerequisites mastered → return "pacing/load issue" (no false gaps)
  - `backend/internal/modules/gapcoverage/application/cascade_impact_service.go` — `ComputeCascadeImpact(ctx, rootModuleID)` — breadth-first traversal from root, count blocked modules per subject

  **Acceptance criteria:**
  - 100% of strict-prerequisite chains analyzed; root = first unmastered
  - Multiple chains → multiple gaps, ranked by cascade impact (verified on ≥ 10 scenarios)
  - Reference scenario: «концентрация растворов» → root «проценты» (70%) → blocks chemistry, biology, social studies
  - Execution ≤ 2 s on graph up to 1000 modules
  - No false gaps when all prerequisites are mastered

  **Logging requirements:**
  - INFO: diagnosis started (learner, lag_module), diagnosis result (root_modules[], cascade_impact)
  - WARN: pacing issue detected (all prerequisites mastered)
  - ERROR: ontology subgraph unavailable for prerequisite traversal

  **Dependencies:** Task 2 (gapcoverage repo, execution-progress repo for mastery data), ontology-port subgraph cache (M1).

- [x] **Task 5: FGOS coverage, deficit list, and attestation readiness report** ✅ M1 domain done

  **Deliverable:** Live FGOS coverage computation, prioritized deficit list, and attestation readiness report generation.

  **Done:** `domain/coverage.go` (ComputeCoverage, Deficit), `domain/report.go` (AttestationReport, BuildAttestationReport), `application/coverage_service.go` (Coverage, Attestation). Deficit prioritization (strict > essential > optional) and attestation critical path are TODO for M2 depth.

  **Deliverable:** Live FGOS coverage computation, prioritized deficit list, and attestation readiness report generation.

  **Files to create/modify:**
  - `backend/internal/modules/gapcoverage/application/coverage_service.go` — `ComputeCoverage(ctx, learnerID, framework)`:
    1. Get mastered modules from execution-progress
    2. Get module→requirement bindings from ontology-port
    3. Compute `coverage = covered_requirements / total_requirements`
    4. Modules below mastery threshold → not counted
    5. Modules without framework binding → audit log, not counted
  - `backend/internal/modules/gapcoverage/application/deficit_service.go` — `ListDeficits(ctx, learnerID, framework)`:
    1. Identify uncovered requirements
    2. Prioritize: strict-prerequisite > essential-core > optional
    3. Within group: by cascade impact
    4. N weeks before attestation → boost attestation-domain deficits
    5. Deficit outside current route → "requires route expansion" flag
  - `backend/internal/modules/gapcoverage/application/attestation_service.go` — `GenerateAttestationReport(ctx, learnerID, attestationDomain)`:
    1. Aggregate coverage, deficits, forecast
    2. Compute critical path to target coverage
    3. Verdict: готов (ready) / под риском (at risk) / не готов (not ready)
    4. "Ready" ONLY when 100% coverage in attestation domains

  **Acceptance criteria:**
  - Coverage matches reference on ≥ 10 contexts (0 deviation)
  - Coverage refresh ≤ 5 s after `ModuleMastered`
  - Deficit list completeness: deficits ∪ covered = full framework set
  - Attestation report generation ≤ 1 s (single learner), ≤ 5 s (school up to 80 learners)
  - "Ready" only at 100% attestation domain coverage

  **Logging requirements:**
  - INFO: coverage computed (learner, framework, percentage, timestamp), deficit list generated (count, top_priority), attestation report generated (learner, verdict, domains[])
  - WARN: module without framework binding (audit), deficit outside route
  - ERROR: framework not found for attestation domain

  **Dependencies:** Task 2 (gapcoverage repo), Task 3 (execution-progress services for mastery data), ontology-port subgraph.

- [ ] **Task 6: OpenAPI spec expansion and HTTP handlers for F2 endpoints**

  **Deliverable:** Updated OpenAPI v1.yaml with F2 REST endpoints and generated/oapi-codegen HTTP handlers wired into the chi router.

  **Files to create/modify:**
  - `backend/api/openapi/v1.yaml` — add new paths:
    - `GET /learners/{learner_id}/progress` — plan-vs-actual report (Task 3)
    - `GET /learners/{learner_id}/forecast` — binary readiness forecast (Task 3)
    - `POST /learners/{learner_id}/module-mastered` — record mastery (Task 3)
    - `GET /learners/{learner_id}/gaps/{module_id}` — gap diagnosis (Task 4) — update existing stub
    - `GET /learners/{learner_id}/coverage/fgos` — coverage report (Task 5) — update existing stub
    - `GET /learners/{learner_id}/coverage/deficits` — deficit list (Task 5)
    - `GET /learners/{learner_id}/attestation` — attestation readiness (Task 5)
  - Add schemas: `ProgressResponse`, `ForecastResponse`, `MasteryRecordRequest`, `GapDiagnosisResponse` (extend existing), `CoverageResponse` (extend existing), `DeficitListResponse`, `AttestationReportResponse`
  - `backend/internal/modules/executionprogress/adapters/handler/http.go` — HTTP handlers delegating to application services
  - `backend/internal/modules/gapcoverage/adapters/handler/http.go` — HTTP handlers
  - Update `backend/internal/cli/server_http.go` — register new routes with role-based guards (learner/parent/methodologist)

  **Logging requirements:**
  - INFO: each endpoint hit (method, path, status, duration_ms)
  - ERROR: validation errors (400), not found (404), internal errors (500)

  **Dependencies:** Tasks 3 (execution services), 4 (gap diagnosis), 5 (coverage/attestation).

<!-- Commit checkpoint: tasks 4-6 -->

### Phase 3: Backend Practice Life & CLI (F5 + CLI)

- [x] **Task 7: Practice-life domain model, application layer, repository, and API** ✅

  **Deliverable:** Backend for F5: story catalog and project-idea catalog (read-only projections from ontology-port), recommendation engine triggered by `ModuleMastered` events.

  **Done:** `domain/practicelife.go` (Story, ProjectIdea, StoryCatalog, ProjectIdeaCatalog with indexing and eligibility), `application/practice_life_service.go` (PracticeLifeService with StoriesForModule, ProjectsForModule, RecommendStories, SuggestProjects), `adapters/adapters.go` (HTTP handler with chi routes).

  **Pending:** Server wiring (add practice-life routes to server_http.go or handler.go) and OpenAPI spec expansion.

  **Deliverable:** Backend for F5: story catalog and project-idea catalog (read-only projections from ontology-port), recommendation engine triggered by `ModuleMastered` events.

  **Files to create/modify:**
  - `backend/internal/modules/practicelife/domain/types.go` — VOs: `Story` (title, text, linked_modules[], real_world_section), `ProjectIdea` (title, modules[], difficulty, expected_outcome)
  - `backend/internal/modules/practicelife/domain/services.go` — interfaces: `StoryRecommendationService`, `ProjectIdeaService`, `QualityService` (deferred)
  - `backend/internal/modules/practicelife/application/story_service.go` — `RecommendStoriesForModule(ctx, moduleID)` — finds stories linked via `appliesTo`/`enriches` from ontology-port cache, returns top-N
  - `backend/internal/modules/practicelife/application/project_service.go` — `SuggestProjects(ctx, learnerID)` — finds project ideas where ≥ 80% required modules are mastered/available
  - `backend/internal/modules/practicelife/adapters/repository/postgres.go` — read-only cache of stories/project ideas from ontology sync
  - `backend/internal/modules/practicelife/adapters/handler/http.go` — HTTP handlers:
    - `GET /modules/{module_id}/stories` — stories for a module
    - `GET /modules/{module_id}/projects` — project ideas for a module
    - `GET /learners/{learner_id}/recommended-stories` — recommendations at mastery
    - `GET /learners/{learner_id}/recommended-projects` — project suggestions
  - `backend/api/openapi/v1.yaml` — add practice-life paths and schemas (StoryResponse, ProjectIdeaResponse, RecommendationResponse)

  **Acceptance criteria:**
  - Story recommendation only via `appliesTo`/`enriches` links (100%)
  - Cross-subject stories (1-3 topics) supported and displayed
  - Real-world application section mandatory in each story
  - Project ideas require modules from ≥ 2 subjects
  - Project suggestion only when ≥ 80% modules are accessible
  - Recommendation display ≤ 1 s after module mastery

  **Logging requirements:**
  - INFO: story recommendation (module, count), project suggestion (learner, eligible_count)
  - WARN: no stories/projects available for module
  - ERROR: ontology cache miss for story/project data

  **Dependencies:** ontology-port subgraph cache (M1), Task 3 (ModuleMastered event).

- [ ] **Task 8: CLI commands: plan get, gap diagnose, report (finalize stubs)**

  **Deliverable:** Wire existing CLI stubs (`plan get`, `gap diagnose`, `report`) to the real application services. These serve as dev/support/testing tooling per ADR-DES.API.cli-interface.

  **Files to modify:**
  - `backend/internal/cli/plan_get.go` — wire `PlanVsActualService.Compare()` for a given learner+plan, output as table/JSON
  - `backend/internal/cli/gap_diagnose.go` — wire `GapDiagnosisService.DiagnoseRootCause()` for a given learner+module, output diagnosis with chains
  - `backend/internal/cli/report.go` — wire `AttestationReadinessService.GenerateAttestationReport()` for a given learner+domain, output report sections
  - `backend/internal/cli/cli.go` — ensure all commands are registered and have --learner, --output flags

  **Logging requirements:**
  - INFO: CLI command invoked (command, args)
  - ERROR: CLI execution failures (missing args, service errors)

  **Dependencies:** Tasks 3, 4, 5 (application services must be operational).

<!-- Commit checkpoint: tasks 7-8 -->

### Phase 4: Frontend Visualization (F4)

- [x] **Task 9: Knowledge map component with progress colors and mode switching** ✅

  **Deliverable:** React component rendering a 2D knowledge graph with color-coded nodes.

  **Done:**
  - `KnowledgeMap.tsx` — React Flow graph with progress colors, critical-path/exploration modes
  - `KnowledgeNode.tsx` — custom node with status-based color coding
  - `Legend.tsx` — 5-status color legend
  - `types.ts` — ModuleStatus, GraphNode, GraphEdge, GraphMode
  - `api.ts` — fetchGraphData, fetchProgress
  - Installed `@xyflow/react` dependency

  **Deliverable:** React component rendering a 2D knowledge graph with color-coded nodes (mastered=green, in-progress=yellow, available=blue, blocked=gray, unclosed-prereq=red) and two viewing modes (critical path / exploration).

  **Files to create/modify:**
  - `frontend/src/features/visualization/KnowledgeMap.tsx` — main component using React Flow for 2D graph rendering:
    - Props: `nodes: GraphNode[]`, `edges: GraphEdge[]`, `progress: Record<string, ModuleStatus>`
    - Color mapping: mastered=green, in-progress=yellow, available=blue, blocked=gray, unclosed-prereq=red
    - Cascade arrows: red edges from root gaps to blocked descendants
    - Legend component (fixed position overlay)
  - `frontend/src/features/visualization/hooks/useGraphData.ts` — fetches graph data and progress from API, computes color states
  - `frontend/src/features/visualization/hooks/useModeSwitch.ts` — critical path mode (strict + essential edges only) vs exploration (all edges)
  - `frontend/src/features/visualization/api.ts` — API calls: GET /learners/{id}/progress, GET /ontology/concepts
  - `frontend/src/features/visualization/index.ts` — barrel exports

  **Acceptance criteria:**
  - Each node has exactly 1 status color from 5
  - Unclosed-prereq (red) only when prerequisite is not mastered
  - Recolor on prerequisite mastery within ≤ 1 s
  - Cascade arrows = graph edges (100% accuracy)
  - Mode switch ≤ 500 ms (up to 500 nodes); progress/colors preserved
  - Legend visible on map

  **Logging requirements (frontend):**
  - INFO: graph data fetched (node_count, edge_count), mode switched (from, to)
  - ERROR: graph data fetch failed, rendering error

  **Dependencies:** Task 6 (F2 REST endpoints must be live), ontology-port API (M1).

- [x] **Task 10: Learner dashboard** ✅

  **Deliverable:** Dashboard page with 6 mandatory widgets.

  **Done:** `LearnerDashboard.tsx` — grid layout with current position, three horizons, plan-vs-actual, subject progress, FGOS coverage, recommendations widgets. Uses shared Card/Badge components.

  **Deliverable:** Dashboard page with 6 mandatory widgets: current position, 3 horizons, plan-vs-actual, subject progress, FGOS coverage, recommended stories/projects.

  **Files to create/modify:**
  - `frontend/src/features/visualization/LearnerDashboard.tsx` — grid layout with 6 widgets:
    1. Current position (module + color status from F4.1 mapping)
    2. Three horizons (far/mid/near from route-planning)
    3. Plan-vs-actual deviation (plan date vs actual date, color-coded)
    4. Progress by subject (bar chart or summary cards)
    5. FGOS coverage % (progress bar, 1% precision)
    6. Recommended stories/projects (linked to practice-life feature)
  - `frontend/src/features/visualization/hooks/useLearnerDashboard.ts` — fetches all dashboard data
  - `frontend/src/features/visualization/WidgetCard.tsx` — shared widget container component
  - `frontend/src/features/visualization/SubjectProgress.tsx` — subject-by-subject progress display
  - `frontend/src/pages/LearnerDashboardPage.tsx` — page component wiring dashboard + auth guards
  - Update `frontend/src/routes.tsx` — add `/dashboard/learner` route (lazy-loaded, ProtectedRoute)

  **Acceptance criteria:**
  - All 6 widgets present; any missing = defect
  - Plan-vs-actual deviation accuracy 100%
  - FGOS coverage precision to 1%
  - Dashboard refresh ≤ 1 s after progress update
  - Last-updated timestamp displayed

  **Logging requirements (frontend):**
  - INFO: dashboard mounted (learner_id), data refresh (elapsed_ms)
  - ERROR: widget data fetch failed (widget name, error)

  **Dependencies:** Task 6 (F2 REST endpoints), Task 9 (knowledge map colors reused), Task 7 (stories/projects API).

- [ ] **Task 11: Parent/HR and methodologist dashboards**

  **Deliverable:** Parent/HR dashboard (5 widgets: progress, FGOS coverage, deviations with color highlights, forecast to checkpoint, recommendations) and methodologist dashboard (FGOS coverage by class/school, top lagging topics, ontology contribution).

  **Files to create/modify:**
  - `frontend/src/features/visualization/ParentDashboard.tsx` — 5 mandatory widgets with deviation magnitude + color (lag > 10% = signal)
  - `frontend/src/features/visualization/hooks/useParentDashboard.ts` — fetches child data (supports 2+ children switching)
  - `frontend/src/features/visualization/MethodologistDashboard.tsx` — school-level aggregation: coverage by class, top lagging topics, ontology contribution
  - `frontend/src/features/visualization/hooks/useMethodologistDashboard.ts`
  - `frontend/src/pages/ParentDashboardPage.tsx` — page with role guard (parent)
  - `frontend/src/pages/MethodologistDashboardPage.tsx` — page with role guard (methodologist)
  - Update `frontend/src/routes.tsx` — add `/dashboard/parent`, `/dashboard/methodologist`

  **Acceptance criteria:**
  - Parent dashboard: all 5 widgets; deviation numbers + color highlights
  - Methodologist dashboard: coverage aggregated by class/school (precision 1%); role scope: own school/classes only
  - Data refresh ≤ 1 s
  - Forecast matches reference calculation

  **Logging requirements (frontend):**
  - INFO: dashboard mounted (role, scope), data refresh
  - ERROR: aggregation failed

  **Dependencies:** Task 6 (F2 endpoints), Task 10 (shared widget patterns).

<!-- Commit checkpoint: tasks 9-11 -->

### Phase 5: Frontend Visualization (F4 cont.) + Practice Life (F5)

- [x] **Task 12: Gap diagnostic map view and route builder** ✅ partial

  **Deliverable:** Gap diagnostic map showing root gaps with cascade arrows.

  **Done:** `GapMap.tsx` — displays root gaps ranked by impact with module/subject cascade info. Route builder deferred (M1 RouteBuilder.tsx already handles route construction).

  **Deliverable:** Visual gap diagnostic map showing only root gaps (unmastered modules with no unmastered prerequisites) with cascade arrows, ranked by impact. Visual 5-step route builder (select goal → visualize route → estimate time → confirm → fixate plan).

  **Files to create/modify:**
  - `frontend/src/features/visualization/GapMap.tsx` — filtered knowledge map:
    - Shows only root gaps (red) + cascade arrows to dependent nodes
    - Ranking overlay: [1] Module X — blocks N modules in M subjects
    - Click on gap → detail: "close X → unlock M subjects"
  - `frontend/src/features/visualization/RouteBuilder.tsx` — 5-step flow:
    1. Goal selection (search topics)
    2. Route visualization (prerequisite chain)
    3. Time estimate (sum of module durations)
    4. Confirmation (review before fixating)
    5. Plan fixation (POST → plan-management)
  - `frontend/src/features/visualization/hooks/useGapMap.ts`
  - `frontend/src/features/visualization/hooks/useRouteBuilder.ts`
  - `frontend/src/pages/GapMapPage.tsx`
  - `frontend/src/pages/RouteBuilderPage.tsx`
  - Update `frontend/src/routes.tsx` — add `/dashboard/gaps`, `/dashboard/route-builder`

  **Acceptance criteria:**
  - Gap map: only root gaps (derived hidden); cascade arrows to ALL dependent nodes; ranked by descending impact
  - Route builder: all 5 steps in order; time estimate visible before confirmation; ≤ 5 min for typical route (≤ 30 modules)
  - Gap map recompute ≤ 1 s on progress change

  **Logging requirements (frontend):**
  - INFO: gap map rendered (root_gap_count), route builder step changed (step)
  - WARN: no root gaps found

  **Dependencies:** Tasks 4 (gap diagnosis API), 6 (F2 endpoints), 9 (knowledge map base component).

- [x] **Task 13: Group management panel** ✅

  **Deliverable:** Panel with mini-cards per learner.

  **Done:** `GroupPanel.tsx` — learner cards (name, module, FGOS %, forecast, attention flag), "X of Y at risk" summary, role-scoped navigation.

  **Deliverable:** Panel with mini-cards per learner (name, current module, FGOS coverage %, forecast status, attention flag), quick switching between learners, summary "X of Y at risk".

  **Files to create/modify:**
  - `frontend/src/features/visualization/GroupPanel.tsx` — list/grid of LearnerCard components:
    - Fields: name, current module, FGOS %, forecast status, attention flag
    - Attention flag: lag > 10% of plan OR unclosed prerequisites
    - Summary bar: "X of Y at risk"
    - Quick switch: click card → navigate to learner dashboard
  - `frontend/src/features/visualization/LearnerCard.tsx` — individual card component
  - `frontend/src/features/visualization/hooks/useGroupPanel.ts` — fetches group data based on role (parent → children, director → school, HR → department)
  - `frontend/src/pages/GroupPanelPage.tsx` — page with role-aware scope
  - Update `frontend/src/routes.tsx` — add `/dashboard/group`

  **Acceptance criteria:**
  - 5 mandatory fields per card; switch ≤ 300 ms
  - "X of Y at risk" matches reference count (80 learners)
  - Role scoping: parent sees own children, director sees school, HR sees department

  **Logging requirements (frontend):**
  - INFO: group panel loaded (role, learner_count, at_risk_count)
  - ERROR: group data fetch failed

  **Dependencies:** Tasks 6 (F2 endpoints), 10 (learner dashboard for navigation), 11 (parent dashboard).

- [x] **Task 14: Stories and project ideas display with recommendation** ✅

  **Deliverable:** UI components for stories and project ideas.

  **Done:**
  - `PracticeComponents.tsx` — StoryCard, ProjectCard, RecommendationPanel
  - `index.ts` — barrel exports + api functions (fetchStoriesForModule, fetchProjectsForModule, fetchRecommendations)

  **Deliverable:** UI components for displaying stories and project ideas, integrated into the learner dashboard and knowledge map. Recommendation panel triggered at module mastery.

  **Files to create/modify:**
  - `frontend/src/features/practice-life/StoryCard.tsx` — displays a story (title, reading time ≤ 5 min, linked topics, real-world section)
  - `frontend/src/features/practice-life/ProjectCard.tsx` — displays a project idea (title, modules list, difficulty badge, expected outcome)
  - `frontend/src/features/practice-life/RecommendationPanel.tsx` — "You've mastered [module] — check out these connections":
    - Related stories (via `appliesTo`/`enriches`)
    - Related project ideas (cross-subject, ≥ 80% modules accessible)
    - Context: "why this matters" from story's real-world section
  - `frontend/src/features/practice-life/api.ts` — API calls: GET /modules/{id}/stories, GET /modules/{id}/projects, GET /learners/{id}/recommended-stories, GET /learners/{id}/recommended-projects
  - `frontend/src/features/practice-life/hooks/useRecommendations.ts` — fetches and caches recommendations
  - `frontend/src/features/practice-life/index.ts` — barrel exports
  - `frontend/src/pages/PracticePage.tsx` — standalone practice-life page (browse all stories/projects)
  - Update `frontend/src/routes.tsx` — add `/dashboard/practice`

  **Acceptance criteria:**
  - Story delivery ≤ 1 s after module open
  - Reading time ≤ 5 min per story
  - Cross-subject links displayed (1-3 topics)
  - Project ideas: only when ≥ 80% required modules are mastered/available
  - Recommendations trigger at mastery confirmation
  - No duplicate recommendations within same module

  **Logging requirements (frontend):**
  - INFO: recommendations fetched (module, story_count, project_count)
  - WARN: no recommendations available

  **Dependencies:** Task 7 (practice-life REST endpoints), Task 10 (dashboard widget integration).

<!-- Commit checkpoint: tasks 12-14 -->

### Phase 6: Integration & Quality Assurance

- [ ] **Task 15: Backend integration and unit tests**

  **Deliverable:** Go test suite covering domain services, application services, and repository adapters using testcontainers for PostgreSQL integration.

  **Files to create/modify:**
  - `backend/internal/modules/executionprogress/application/trajectory_service_test.go` — mock repository, test mastery recording + event emission
  - `backend/internal/modules/executionprogress/application/plan_vs_actual_test.go` — test deviation computation with known plan snapshots and mastery records
  - `backend/internal/modules/executionprogress/application/forecast_service_test.go` — test forecast with various pace/remaining scenarios
  - `backend/internal/modules/gapcoverage/application/gap_diagnosis_service_test.go` — test root-cause chains on known ontology subgraph, test cascade impact
  - `backend/internal/modules/gapcoverage/application/coverage_service_test.go` — test FGOS coverage computation, deficit prioritization
  - `backend/internal/modules/gapcoverage/application/attestation_service_test.go` — test verdict logic for ready/at-risk/not-ready
  - `backend/internal/modules/practicelife/application/story_service_test.go` — test recommendation logic
  - `backend/tests/integration/progress_integration_test.go` — end-to-end: record mastery → compare plan-vs-actual → compute coverage
  - `backend/tests/integration/gap_integration_test.go` — end-to-end: deviation → gap diagnosis → coverage recalculation

  **Logging requirements (tests):**
  - Test helpers should NOT log (clean test output); use `zap.NewNop()` for services under test

  **Dependencies:** Tasks 1-8 (all backend logic must be implemented).

- [ ] **Task 16: Frontend component tests for visualization and practice-life**

  **Deliverable:** Vitest + React Testing Library tests for all new UI components, replacing scaffolded `it.skip()` stubs with real test cases.

  **Files to create/modify:**
  - `frontend/src/features/visualization/__tests__/knowledge-map.test.tsx` — test color mapping (5 statuses), mode switching, legend rendering, loading/error states
  - `frontend/src/features/visualization/__tests__/learner-dashboard.test.tsx` — test 6 widgets render, data-api integration via mocks, refresh behavior
  - `frontend/src/features/visualization/__tests__/parent-dashboard.test.tsx` — test multi-child switching, deviation highlights, role guard
  - `frontend/src/features/visualization/__tests__/gap-map.test.tsx` — test root-gap-only filtering, ranking display
  - `frontend/src/features/visualization/__tests__/group-panel.test.tsx` — test card fields, "X of Y" summary, role scoping
  - `frontend/src/features/visualization/__tests__/route-builder.test.tsx` — test 5-step flow progression
  - `frontend/src/features/practice-life/__tests__/stories.test.tsx` — test story display, recommendation panel
  - `frontend/src/features/practice-life/__tests__/projects.test.tsx` — test project display, difficulty badges

  **Test patterns (from existing tests):**
  - `vi.mock()` for API modules
  - Harness component wrapping hooks + components
  - `screen.getByTestId()`, `screen.getByRole('alert')` for errors
  - `waitFor(() => expect(...))` for async loading
  - `beforeEach`/`afterEach` for `vi.clearAllMocks()`

  **Dependencies:** Tasks 9-14 (all frontend components must be implemented).

<!-- Commit checkpoint: tasks 15-16 -->
