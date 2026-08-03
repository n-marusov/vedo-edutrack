---
archived: 2026-08-03
---
# Implementation Plan: M0.1 — Domain Model & Architecture Baseline

Branch: none (git.create_branches=false)
Created: 2026-08-02

## Settings
- Testing: no
- Logging: standard (INFO)
- Docs: no (warn-only)

## Roadmap Linkage
Milestone: "M0.1: Domain Model & Architecture Baseline" (Phase 0 — Foundation)
Rationale: M0.1 is the second milestone of Phase 0, directly dependent on the completed M0.0. This plan delivers the domain model, architecture decisions, C4 diagrams, RBAC matrix, traceability model, and VEDO Hub / EduTrack boundary specification required before engineering platform setup (M0.2).

## Progress

### Phase 1: Domain Discovery
- [x] T1: Bounded contexts and context map
- [x] T2: Aggregates, entities, and domain events

### Phase 2: Architecture Decision Records
- [x] T3: Stack selection ADRs
- [x] T4: Database and communication patterns ADRs
- [x] T5: Repository structure and RBAC model ADRs

### Phase 3: C4 Architecture Diagrams
- [x] T6: System Context and Container diagrams
- [x] T7: Core Component diagrams

### Phase 4: Cross-Cutting Concerns
- [x] T8: Role-permission matrix
- [x] T9: Traceability model extension
- [x] T10: VEDO Hub / EduTrack responsibility boundary

### Phase 5: Validation & Finalization
- [x] T11: Cross-reference validation and M0.1 completion

## M0.1 Exit Criteria (from ROADMAP.md)

- [x] Bounded contexts, aggregates, domain events, and context map are documented
- [x] C4 System Context, Container, and core Component diagrams are available
- [x] ADRs cover stack, framework, database, communication patterns, repository structure, and RBAC model
- [x] Role-permission matrix covers all personas from `specs/vision.md`
- [x] Traceability model links US → UC → FR → COMP → TEST
- [x] VEDO Hub / EduTrack responsibility boundary is explicit: EduTrack reads ontology data and computes educational mechanics; Hub owns ontology storage, editing, versioning, forks, and social contribution flows

## Context Summary

### Source Material
- `specs/vision.md` — authoritative product vision, personas, F0–F6 functions, domain events (2.4)
- `specs/glossary.md` — domain terminology (7 sections, 50+ terms)
- `specs/requirements/` — 60 FR + 54 NFR with acceptance criteria
- `specs/user-stories/` — 47 US in Gherkin format
- `specs/use-cases/` — 42 UC with actor/channel/flow metadata
- `traceability.ttl` — existing TBox (classes, properties, integrity constraints); 0 ABox instances for COMP/TEST/ADR/C4
- `.ai-factory/DESCRIPTION.md` — architecture notes, NFR corpus summary
- `.ai-factory/ROADMAP.md` — Phase 0 structure, M0.1 exit criteria

### Existing Conventions
- **ADR**: `ADR-<LEVEL>.<AREA>.<semantic-tag>` (LEVEL: BIZ|DES|IMPL)
- **C4**: `<LEVEL>-<name>.md` (LEVEL: context|container|component), Mermaid + legend
- **UC**: `UC-<L1>.<L2>.<L3>`
- **US**: `US-<domain>.<subdomain>.<action>`
- **FR**: `REQ-FR-<domain>.<qualifier>.<action>`
- **NFR**: `REQ-NFR-<area>.<qualifier>.<attribute>`
- **Traceability**: `traceability.ttl` (OWL 2 DL Turtle), classes: Vision, Glossary, ADR, UC, FR, NFR, US, C4Diagram, Component, Test (UnitTest, IntegrationTest, E2ETest, Gate, LoadTest), Documentation; properties: derivesFrom, isSourceOf, refines, traceableTo

### What Exists vs. What Needs to Be Created

| Area | Exists | Needs Creation |
|------|--------|----------------|
| DDD models | Glossary (terms) | Bounded contexts, aggregates, entities, domain events catalog, context map |
| ADR | `specs/adr/README.md` (conventions only) | 6 ADRs: stack, framework, database, communication, repo structure, RBAC |
| C4 diagrams | `specs/c4/README.md` (conventions only) | System Context, Container, core Component diagrams |
| RBAC | Personas in vision.md, RBAC ADR (T5) | Role-permission matrix as NFR specs: role-catalog, permission-matrix, ops-admin-separation |
| Traceability | TBox in traceability.ttl | ABox: ADR instances, C4Diagram instances; COMP/TEST classes ready for future use |
| Boundary | Implicit in vision.md and glossary, comm ADR (T4) | Explicit boundary as NFR specs: ownership-boundary, ontology-read-only |
| ROADMAP | M0.1 entry exists (unchecked) | Mark as completed |

---

## Tasks

### Phase 1: Domain Discovery

- [x] **T1: Identify bounded contexts and create context map**
  Identify bounded contexts from the functional domain (F0–F6), glossary terms, and the two operational contours (Community / Enterprise). Produce a context map diagram showing relationships (partnership, customer-supplier, shared kernel, ACL) between bounded contexts.
  
  **Deliverables:**
  - `specs/ddd/context-map.md` — Mermaid context map + bounded context descriptions ✅
  
  **Boundary & domain alignment:**
  - Use `specs/vision.md` §2.2 (Level 1–2 function decomposition) and `specs/glossary.md` §4 (architecture) as primary sources ✅
  - Each bounded context must trace to at least one F0–F6 domain and at least one UC ✅ (10 контекстов, все трассируются)
  - Document the VEDO Hub ACL (Anti-Corruption Layer) as an explicit bounded-context boundary ✅ (`ontology-port`)
  
  **Logging:** Record each identified bounded context at INFO with its rationale and mapping to F0–F6. Record context-map relationship decisions at INFO with justification.

- [x] **T2: Model aggregates, entities, value objects, and domain events**
  For each bounded context from T1, model the core aggregates (transactional boundaries), their root entities, internal entities, and value objects. Produce a domain events catalog from vision.md §2.4 plus event-triggered behaviors in F1.2/F1.6/F2.1.
  
  **Deliverables:**
  - `specs/ddd/aggregates.md` — aggregates, entities, value objects per bounded context ✅
  - `specs/ddd/domain-events.md` — event catalog: name, trigger, payload sketch, consuming contexts ✅
  
  **Key domain invariants to capture:**
  - Route = f(position, goal, pedagogy concept, ontology) → route (function, not document) ✅
  - Plan = route snapshot + timeline (fixed at checkpoint; route continues recomputing independently) ✅
  - Trajectory = actual path taken (always visible alongside route and plan) ✅
  - Gap diagnosis = climb strict-prerequisite links upward to first unmastered module ✅
  - Essential core = mandatory subset defined by learning context; order/pace vary, dependency logic is unbreakable ✅
  - Two contours (Community / Enterprise) share one API contract ✅
  
  **Logging:** Log each aggregate and its root entity at INFO. Log each domain event at INFO with producer and consumer bounded contexts. Log any conflicting interpretations of glossary terms at WARN.

---

### Phase 2: Architecture Decision Records

- [x] **T3: Stack selection ADRs**
  Produce ADRs for the programming language and framework. Evaluate candidates from `DESCRIPTION.md` (TypeScript / Python / Java for language; Next.js / FastAPI+NestJS / Spring Boot for framework) against NFR constraints from `specs/requirements/REQ-NFR-*.md` (p95 ≤ 200ms, horizontal scaling, security requirements, i18n).
  
  **Deliverables:**
  - `specs/adr/ADR-DES.STACK.language-vs-vs.md` — language selection (TypeScript vs Python vs Java) ✅ (решение: **Go** + TS-фронт; кандидаты расширены Go, Python, Java — см. DESCRIPTION.md)
  - `specs/adr/ADR-DES.STACK.framework-vs-vs.md` — framework selection ✅ (решение: **chi + oapi-codegen** (бэкенд) + **React + TS** (фронт); отклонены tRPC/Connect/Vue/Next.js)
  
  **ADR format** (per `specs/adr/README.md`): Status (ПРИНЯТО), Date, Контекст, Требование-источник, Решение, Рассмотренные альтернативы, Последствия. ✅
  
  **Constraints:**
  - Must reference relevant NFRs: `REQ-NFR-api.performance.latency-p95`, `REQ-NFR-ops.performance.scalability`, `REQ-NFR-process.dev.test-coverage`, `REQ-NFR-security.compliance.owasp-application-security` ✅
  - Must account for i18n-readiness (`REQ-NFR-ops.compliance.i18n-readiness`): RU + EN, 0-code language addition ✅
  - Deterministic core + LLM module: the chosen stack must support both a deterministic routing engine and an optional LLM integration module ✅ (Go-ядро + изолированный LLM-адаптер)
  - Community contour (public ontology, families, EdTech) vs Enterprise contour (corporate, isolated, 152-ФЗ) may influence stack decisions ✅ (Go-бинарник для on-prem Enterprise; React для Community-скорости)
  
  **Logging:** Log each candidate evaluation dimension at INFO. Log the final decision at INFO with a summary of trade-offs. Log rejected alternatives with explicit rejection rationale at INFO.

- [x] **T4: Database and communication patterns ADRs**
  Produce ADRs for the database technology and inter-service communication patterns.
  
  **Deliverables:**
  - `specs/adr/ADR-DES.DATA.storage-strategy.md` — database selection (PostgreSQL candidate per `DESCRIPTION.md`), schema design approach
  - `specs/adr/ADR-DES.API.communication-patterns.md` — sync REST + async event patterns, API versioning strategy, idempotency guarantees for webhooks, SPARQL read-only endpoint
  
  **Domain-specific constraints:**
  - EduTrack stores learner/plan/progress data, NOT ontologies (knowledge graph lives in VEDO Hub)
  - Route is computed in-memory from copied subgraph — database stores only snapshots (plans) and progress records
  - Webhook idempotency per `REQ-NFR-api.availability.webhook-idempotency` and `REQ-FR-api.webhooks.idempotency`
  - SPARQL endpoint is read-only per `REQ-FR-api.sparql.read-only` and `REQ-NFR-security.compliance.owasp-application-security` (parameterization required)
  
  **Logging:** Log candidate evaluation at INFO. Log final decisions with trade-off summaries at INFO.

- [x] **T5: Repository structure and RBAC model ADRs**
  Produce ADRs for the repository layout and role-based access control model.
  
  **Deliverables:**
  - `specs/adr/ADR-IMPL.PROCESS.repository-structure.md` — monorepo layout aligned with bounded contexts from T1, folder conventions per selected stack (T3)
  - `specs/adr/ADR-DES.SECURITY.rbac-model.md` — role hierarchy, permission granularity, auth enforcement points
  
  **RBAC constraints:**
  - Personas from `specs/vision.md` §3.1: learner, parent (account owner), teacher/methodologist, school director, HR/L&D manager, EdTech platform integrator, content contributor
  - Role model from `REQ-NFR-security.compliance.role-based-access`: learner / parent / school / methodologist / HR
  - Two contours (Community / Enterprise) may have different role sets and permission scopes
  - Enterprise SSO/SAML via Keycloak per `REQ-FR-api.sso.keycloak`
  
  **Logging:** Log layout decisions at INFO with rationale from bounded contexts. Log each role and its permission scope at INFO. Log any unresolved role-model questions at WARN.

---

### Phase 3: C4 Architecture Diagrams

- [x] **T6: System Context and Container diagrams**
  Produce C4 Level 1 (System Context) and Level 2 (Container) diagrams following conventions in `specs/c4/README.md`.
  
  **Deliverables:**
  - `specs/c4/context-system.md` — System Context diagram: EduTrack system, external actors, VEDO Hub as external system, two deployment contours ✅
  - `specs/c4/container-overview.md` — Container diagram: web app, API server (modular monolith), PostgreSQL, VEDO Hub (external). Route/execution/resource engines are modules of the API server per `ADR-DES.INFRA.modular-monolith-approach`; cache is PostgreSQL + in-memory (Redis post-MVP) ✅
  
  **Diagram format** (per `specs/c4/README.md`):
  - Mermaid diagram in ` ```mermaid ` block
  - Legend: node and relationship descriptions
  - Context: scenario the diagram was built for
  - Links to functions F0–F6 from `specs/vision.md`
  
  **Key invariants to show:**
  - VEDO Hub is ALWAYS an external system/container, never internal to EduTrack
  - Two contours (Community / Enterprise) are separate instances with different actors
  - Route is computed, not stored — computation engine is shown, not a "Route DB"
  
  **Logging:** Log which actors and external systems are included at INFO. Log design decisions (e.g., number of containers) at INFO.

- [x] **T7: Core Component diagrams**
  Produce C4 Level 3 (Component) diagrams for the core bounded contexts.
  
  **Deliverables** (master-index per container, see `specs/c4/README.md`; core contexts covered as components inside the API server):
  - `specs/c4/component-api-server.md` — API server (Go / chi modular monolith): 8 infrastructure components + 10 bounded-context modules (Clean Architecture circles); covers route planning (`RouteComputationService`, `HorizonService`, `CascadeRecomputeService`), execution (`TrajectoryService`, `PlanVsActualService`, `ForecastService`), gap & coverage (`GapDiagnosisService`, `CoverageService`, `AttestationReadinessService`), and the Ontology Port ACL
  - `specs/c4/component-web-app.md` — Web app (React SPA): Domain / Application (Zustand stores) / Interface Adapters / Frameworks circles, mirroring the backend clean architecture
  
  **Component boundaries:**
  - Components align with bounded contexts from T1 ✅ (all 10 contexts as modules)
  - Each component traces to at least one UC and one FR ✅ (F0–F6 tables in both files)
  - The Ontology Port (F0) is the ACL between EduTrack and VEDO Hub — explicitly shown as a boundary ✅ (`OntologyReaderPort` / `OntologyUpdatedSubscriberPort` + hubClient adapter)
  
  **Logging:** Log each component and its responsibilities at INFO. Log component-to-bounded-context mapping at INFO.

---

### Phase 4: Cross-Cutting Concerns

- [x] **T8: Role-permission matrix**
  Produce a role-permission matrix covering all personas from `specs/vision.md` §3.1, aligned with the RBAC ADR from T5.
  
  **Deliverables:**
  - `specs/requirements/REQ-NFR-security.compliance.role-catalog.md` — role catalog: universal archetypes (contract) + personas from `specs/vision.md` §3.1 mapped to roles as seed instances
  - `specs/requirements/REQ-NFR-security.compliance.permission-matrix.md` — permission matrix: roles × functional areas × permission levels (CRUD) × scope
  - `specs/requirements/REQ-NFR-security.compliance.ops-admin-separation.md` — admin/ops separation
  
  **Matrix structure:**
  - Rows: archetypes (self, dependents-owner, staff, management, integration, admin, ops) — universal contract layer; concrete roles are seed instances
  - Columns: functional areas (route-compute, plan-view, plan-manage, progress-track, gap-diagnose, coverage-view, resource-manage, visualization, user-manage, ontology-read, webhook-configure)
  - Cells: permission level (C — create, R — read, U — update, D — delete, — — none)
  - Separate sections for Community and Enterprise contours where permission scopes differ
  - Trace each persona to at least one UC actor and one NFR role requirement
  
  **Logging:** Log each role at INFO with summary of permission count. Log any persona that maps ambiguously to the RBAC model at WARN.

- [x] **T9: Traceability model extension**
  Extend `traceability.ttl` to add ABox instances for the new M0.1 artifacts (ADRs, C4 diagrams) and ready the model for future COMP and TEST artifacts.
  
  **Deliverables:**
  - Updated `traceability.ttl` — ABox instances for all ADRs and C4 diagrams produced in T3–T7
  
  **Instance creation rules:**
  - Each ADR → `tr:ArchitectureDecision` instance with `tr:derivesFrom` links to relevant FR/NFR
  - Each C4 diagram → `tr:C4Diagram` instance with `tr:derivesFrom` links to relevant ADRs and FRs
  - Verify existing TBox has `tr:ArchitectureDecision` (line 63) and `tr:C4Diagram` (line 93) classes — both already defined
  - Verify that `tr:Component` (line 98) and `tr:Documentation` (line 157) classes are present for future M0.2/M0.3 use
  - No COMP or TEST instances yet — these are M0.2+ scope
  - Follow the existing chain convention: `Vision → UC → US → FR → ADR → COMP → TEST`
  
  **Logging:** Log each new instance at INFO with its URI and linked artifacts. Log any missing class or property at ERROR. Log the total instance count after update at INFO.

- [x] **T10: VEDO Hub / EduTrack responsibility boundary**
  Produce an explicit boundary specification documenting what EduTrack owns vs. what VEDO Hub owns, resolving ambiguities from `specs/vision.md`, `specs/glossary.md` §4, and `.ai-factory/DESCRIPTION.md`.
  
  **Deliverables:**
  - `specs/requirements/REQ-NFR-api.compliance.ownership-boundary.md` — responsibility boundary with ownership matrix
  - `specs/requirements/REQ-NFR-api.compliance.ontology-read-only.md` — ontology read-only contract (ontology port / ACL)
  
  **Boundary dimensions to cover:**
  - Data ownership: EduTrack stores learner profiles, plans, progress, gap diagnoses, coverage reports; Hub stores ontologies, modules, links, FGOS mappings, resources, stories, pedagogy concepts
  - Computation ownership: EduTrack computes routes, gaps, forecasts, coverage; Hub serves ontology queries, versioning, forking, merging
  - API ownership: EduTrack exposes route/plan/execution/coverage APIs; Hub exposes ontology CRUD, SPARQL, MCP, social contribution APIs
  - Event ownership: EduTrack emits module.mastered, plan.deviated, route.recalculated; Hub emits ontology.updated, fork.created, merge-request.opened
  - Deployment ownership: EduTrack deploys its own containers; Hub is a separate platform with its own SLA
  - The ontology port (F0) is the formal ACL — explicitly document the contract surface: which Hub endpoints EduTrack calls, data formats, error handling, version compatibility
  
  **Logging:** Log each boundary dimension at INFO with the responsible party. Log any gray-area decisions at INFO with rationale.

---

### Phase 5: Validation & Finalization

- [x] **T11: Cross-reference validation and M0.1 completion**
  Validate consistency across all M0.1 artifacts and mark the milestone as completed in ROADMAP.md.
  
  **Validation checklist:**
  - [x] Bounded contexts (T1) ↔ C4 components (T7): every component maps to exactly one bounded context
  - [x] Aggregates (T2) ↔ ADR storage strategy (T4): aggregate boundaries align with transactional boundaries
  - [x] Domain events (T2) ↔ ADR communication patterns (T4): events have defined channels (sync/async/webhook)
  - [x] Role-permission matrix (T8) ↔ RBAC ADR (T5): permissions are consistent
  - [x] RBAC ADR (T5) ↔ NFR role model: covers learner/parent/school/methodologist/HR
  - [x] C4 context (T6) ↔ Boundary spec (T10): VEDO Hub is external on all diagrams
  - [x] Traceability (T9): every ADR and C4 diagram has a tt:R instance; 0 orphan instances
  - [x] All exit criteria from ROADMAP M0.1 are satisfied
  - [x] No artifact references a non-existent file or convention
  
  **ROADMAP update:**
  - In `.ai-factory/ROADMAP.md`, change `- [ ] **M0.1: Domain Model & Architecture Baseline**` to `- [x] **M0.1: Domain Model & Architecture Baseline** ✅ completed YYYY-MM-DD`
  
  **Logging:** Log each validation check at INFO with pass/fail. Log any failures at ERROR with the specific inconsistency and recommended fix. Log the final M0.1 completion status at INFO.

---

## Commit Plan

- **Commit 1** (after tasks 1-3): `docs(m0.1): add domain model, aggregates, and stack ADRs`
- **Commit 2** (after tasks 4-6): `docs(m0.1): add database/comm ADRs, RBAC ADR, and C4 system/container diagrams`
- **Commit 3** (after tasks 7-9): `docs(m0.1): add component diagrams, role matrix, and traceability ABox`
- **Commit 4** (after tasks 10-11): `docs(m0.1): add boundary spec and finalize M0.1 in ROADMAP`

---

## Notes

- All artifacts are documentation/design — no source code is written in this milestone.
- The stack is TBD: ADRs T3–T4 resolve this.
- `specs/vision.md` is the authoritative source — every artifact must trace back to it.
- The ROADMAP network plan (M0.0 → M0.1 → M0.2 → M0.3) is strictly sequential — M0.2 depends on the ADRs and repository structure produced here.
- The glossary (`specs/glossary.md`) is the single source of truth for domain terminology — use its definitions, do not redefine terms.
