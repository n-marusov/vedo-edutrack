# Project Roadmap

> VEDO EduTrack — an educational route service on top of VEDO Hub that builds personalized, cross-subject learning paths from a shared knowledge ontology.

## Current Focus

**Now:** Phase 0 — Foundation (M0.0 complete ✅)  
**Next milestone:** M0.1: Domain Model & Architecture Baseline  
**Blocking decision:** Technology stack selection before M0.2/M0.3  
**External blocker:** Starter ontology readiness in VEDO Hub before M1/M2

## Milestone Status

| Total | Completed | Current | Next |
|-------|-----------|---------|------|
| 14 | 1 (M0.0) | M0.1: Domain Model & Architecture Baseline | M0.2: Engineering Platform |

## External Dependencies

These dependencies are outside the direct ownership of EduTrack implementation, but they block or constrain milestone completion.

| Dependency | Required by | Readiness Criteria |
|------------|-------------|--------------------|
| Starter ontology in VEDO Hub | M1, M2 | 1000+ school topics for grades 5–11, 500+ cross-subject links, FGOS mapping, 3–5 pedagogy concepts, MVP resources/stories/project ideas exposed through Hub API |
| VEDO Hub REST API and MCP server | M1, M4, M5, M6 | Stable read APIs for modules, links, FGOS bindings, resources, stories, pedagogy concepts; MCP access for AI agents; documented auth and versioning behavior |
| VEDO Hub SPARQL/Cypher endpoint | M1, M4, M6 | Read-only query access suitable for route computation and partner integrations |
| VEDO Hub fork/merge/community mechanics | M6, M8 | Public ontology fork flow, contribution review, semantic diff/merge, contributor profiles; EduTrack consumes these capabilities but does not implement ontology authoring |
| Enterprise identity provider / customer LMS access | M3, M6, M7, M9 | Testable SSO/SAML and LMS integration environments for corporate pilots |

## Milestones

### Phase 0 — Foundation (pre-MVP, 3–4 weeks)

> Engineering baseline: requirements, architecture, platform, and scaffold. The work that must happen before production feature development begins.

- [x] **M0.0: Requirements Baseline** ✅ completed 2026-08-02
  Establish a complete, traceable MVP requirements baseline from `specs/vision.md` and `.ai-factory/DESCRIPTION.md`.

  **Exit criteria:**
  - ✅ User stories are defined for planning, execution, visualization, resources, practice, API, and accessibility flows (47 `US-*`).
  - ✅ Use cases are defined for route building, plan execution, visualization, auth, and API integration (42 `UC-*`).
  - ✅ Functional and non-functional requirements are defined with stable IDs (56 `FR-*`, 13 `NFR-*`).
  - ✅ MVP acceptance criteria and MoSCoW prioritization are documented (`specs/requirements/MVP-ACCEPTANCE-CRITERIA.md`).
  - ✅ `specs/requirements/README.md` defines naming and traceability conventions.
  - ✅ Quality matrix audit: 0 critical + 0 high gaps (`.ai-factory/quality-matrix.md`).
  - ✅ Traceability: `traceability.ttl` fully populated, 0 broken orphans (US → UC → FR chains complete).

  **Dependencies:** `specs/vision.md` (authoritative), `.ai-factory/DESCRIPTION.md`.
  **Business goals:** G1, G2, G3, G4, G6.

- [ ] **M0.1: Domain Model & Architecture Baseline**
  Define the domain, system boundaries, architecture decisions, and traceability model before implementation.

  **Exit criteria:**
  - Bounded contexts, aggregates, domain events, and context map are documented.
  - C4 System Context, Container, and core Component diagrams are available.
  - ADRs cover stack, framework, database, communication patterns, repository structure, and RBAC model.
  - Role-permission matrix covers all personas from `specs/vision.md`.
  - Traceability model links US → UC → FR → COMP → TEST.
  - VEDO Hub / EduTrack responsibility boundary is explicit: EduTrack reads ontology data and computes educational mechanics; Hub owns ontology storage, editing, versioning, forks, and social contribution flows.

  **Dependencies:** M0.0.
  **Business goals:** G1, G2, G3, G4, G5, G6.

- [ ] **M0.2: Engineering Platform**
  Create the technical foundation for local development, CI, testing, and containerized execution.

  **Exit criteria:**
  - Repository layout follows the architecture baseline and bounded contexts.
  - Development environment starts with one command.
  - Build automation covers up, down, build, test, and lint workflows.
  - CI pipeline runs lint, tests, coverage, E2E checks, and build gates.
  - Unit, integration, and E2E test scaffolds exist with intentionally red or placeholder tests where implementation is pending.
  - Container strategy is documented and aligned with selected stack.

  **Dependencies:** M0.1, stack ADR, repository-structure ADR.
  **Business goals:** Engineering enabler for all goals.

- [ ] **M0.3: Runnable Product Scaffold**
  Deliver a minimal authenticated product skeleton that proves the chosen stack, local environment, API shell, ontology stub, route stub, and role-aware UI shell work together.

  **Exit criteria:**
  - Health checks verify API, database, and identity-provider reachability.
  - Authentication middleware validates JWTs and supports role-based guards.
  - Role-aware dashboard shells exist for MVP personas.
  - Ontology stub returns a small fixed module graph.
  - Route-computation stub demonstrates route calculation on the fixed graph via `POST /routes/compute`.
  - Landing page communicates the product value proposition and routes users toward registration.
  - Local environment reports healthy containers/services and scaffold-level tests pass.

  **Dependencies:** M0.2.
  **Business goals:** Engineering enabler for all goals.

**Network plan:** M0.0 → M0.1 → M0.2 → M0.3 — strictly sequential.

---

### Phase I — MVP

> Core product: route engine, visualization, family education app, corporate onboarding pilot, and integration layer.

- [ ] **M1: Core Infrastructure (F0 + F1 + F2 + F3 + F6)**
  Replace scaffold stubs with real MVP infrastructure for reading VEDO Hub ontology data, computing routes, tracking plans/progress, diagnosing gaps, and exposing initial APIs.

  **Exit criteria:**
  - EduTrack reads modules, links, FGOS bindings, resources, stories, and pedagogy concepts from VEDO Hub through approved APIs.
  - Route engine supports MVP pathfinding using strict, soft, and `appliesTo` links.
  - Route recomputes on progress and goal-change triggers.
  - Plan snapshots, plan-vs-actual tracking, binary readiness forecast, and root-cause gap diagnosis are implemented.
  - Resource catalog is available for route modules.
  - MVP APIs expose route computation, progress, and FGOS coverage consistently; route computation uses `POST /routes/compute` or the final ADR-approved equivalent.
  - Contract tests protect the VEDO Hub API boundary.

  **Dependencies:** M0.3, starter ontology readiness in VEDO Hub.
  **Business goals:** G1, G3.

- [ ] **M2: Family Education — «Дай пять» (F2 + F4 + F5)**
  Deliver the Community-facing family education product slice for route building, knowledge-map visualization, FGOS coverage, gap diagnosis, and learner/parent dashboards.

  **Exit criteria:**
  - Parent builds a route for one child in ≤5 minutes.
  - Parent can manage 2+ children from one dashboard.
  - Learner dashboard shows route, progress, next steps, and motivational cross-subject links.
  - Knowledge map shows progress, gaps, and cascade impact.
  - FGOS coverage, deficit list, and attestation-readiness report are visible in real time.
  - Launch content includes at least 50 stories and 30 project ideas.
  - Demo editor or demo flow supports product validation without requiring custom customer setup.

  **Success metrics:** route built in ≤5 minutes; Community NPS ≥ 40.
  **Dependencies:** M1.
  **Business goals:** G1, G3, G4.

- [ ] **M3: Corporate Onboarding Pilot — «Вектор Компетенций»**
  Validate the corporate onboarding application using the shared route, plan, framework-coverage, and dashboard mechanics.

  **Exit criteria:**
  - Personalized onboarding route combines professional skills, domain context, corporate values, internal services, and socialization steps.
  - Basic HR/corporate dashboard shows onboarding progress and readiness.
  - Time-to-productivity metric is captured and reportable.
  - Pilot integration path exists for HR/LMS systems through REST APIs and webhooks.
  - One pilot company can run a controlled onboarding scenario.

  **Success metrics:** target time to productivity is 2 weeks instead of 4 weeks.
  **Dependencies:** M1, M4 partial integration capabilities.
  **Business goals:** G6.

- [ ] **M4: Integration & Webhook Layer (F6)**
  Establish the MVP integration surface for partners, pilots, and AI agents without duplicating VEDO Hub ontology-management responsibilities.

  **Exit criteria:**
  - Production-ready REST API contracts are documented for MVP route, progress, coverage, and integration flows.
  - API foundation includes authentication, versioning approach, rate-limit policy, and basic SLA expectations.
  - Read-only SPARQL/query access is available through the approved Hub/EduTrack boundary.
  - Webhooks include at least `module.mastered`, `plan.deviated`, and `route.recalculated` where supported by events.
  - MCP interface for AI agents is documented and testable.
  - API examples and integration sandbox/demo data are available for early partners.

  **Dependencies:** M1.
  **Business goals:** G2, G5.

**Network plan:** M1 → M2, M3, M4. M4 partially blocks M3 HR/LMS integration. M2–M4 can proceed in parallel after M1.

---

### Phase II — Enrichment & B2B

- [ ] **M5: Route Enrichment from Extended Hub Ontology**
  Upgrade EduTrack route, resource, forecast, and motivation mechanics by consuming richer ontology data from VEDO Hub.

  **Exit criteria:**
  - Route engine uses all five cross-subject link types: `hasStrictPrerequisite`, `hasSoftPrerequisite`, `enriches`, `appliesTo`, `isAnalogousTo`.
  - Pedagogy-concept-aware routes support spiral learning, project immersion, and practice-to-theory patterns.
  - Routes recompute on ontology-update signals from VEDO Hub.
  - Resource matching supports format/style dimensions such as visual, text, and audio.
  - Forecasting upgrades from binary readiness to green/yellow/red risk levels with target ±10% accuracy.
  - Content layer scales to 200+ stories and 100+ project ideas.
  - Qualities map (`develops` tagging) supports upbringing/qualities coverage.

  **Dependencies:** M1, M2, relevant enriched ontology fields/events from VEDO Hub.
  **Business goals:** G1, G3, G4.

- [ ] **M6: EdTech Platform Integration**
  Mature the integration experience for EdTech platforms and LMS partners beyond the MVP API foundation.

  **Exit criteria:**
  - EdTech partner can onboard by forking public ontology in VEDO Hub and receiving first EduTrack API route response within ≤1 week.
  - Integration guide, sandbox data, API examples, and partner checklist are complete.
  - LMS connectors cover WebTutor and iSpring.
  - Enterprise SSO/SAML through Keycloak is supported for partner scenarios.
  - API usage, errors, and partner-facing operational signals are observable.
  - Content-agnostic boundary is preserved: partners keep ownership of content and audience; EduTrack provides route/progress/coverage mechanics.

  **Success metrics:** 20+ platforms use the public ontology/integration flow.
  **Dependencies:** M4, VEDO Hub fork/community mechanics.
  **Business goals:** G2.

---

### Phase III — Verticals & Community

- [ ] **M7: Corporate Application & Compliance**
  Expand the corporate contour from onboarding pilot to full competency, career-track, compliance, and ROI analytics use cases.

  **Exit criteria:**
  - Career tracks compute gaps from current position to target role.
  - Regulatory requirements bind to modules and real work scenarios.
  - Corporate dashboard covers onboarding funnel, time-to-productivity by department, compliance coverage, and ROI analytics.
  - System supports 50+ employees on simultaneous onboarding under control.
  - Enterprise reporting demonstrates target ROI model.

  **Success metrics:** 50+ simultaneous onboarding employees; ROI target 16:1.
  **Dependencies:** M3, M4, M6.
  **Business goals:** G4, G6.

- [ ] **M8: Community Consumption & Network Effect**
  Surface the value of the VEDO Hub contributor ecosystem through EduTrack routes, dashboards, and partner flows while leaving ontology authoring and merge mechanics in Hub.

  **Exit criteria:**
  - EduTrack consumes community-enriched ontology updates from VEDO Hub.
  - Routes and dashboards visibly benefit from newly merged topics, links, resources, and stories.
  - Contributor attribution from Hub can be shown where product-relevant.
  - Semantic diff/merge outcomes from Hub can be reflected in EduTrack integration flows.
  - Partner program for EdTech platforms is operational.

  **Success metrics:** ecosystem target of 5000+ topics, 8000+ links, and 200+ active contributors, owned by the Hub/community roadmap and consumed by EduTrack.
  **Dependencies:** M6, VEDO Hub community mechanics.
  **Business goals:** G5.

---

### Phase IV — Enterprise & Scale

- [ ] **M9: Enterprise Deployment & Compliance**
  Support private enterprise deployments and regulated corporate environments.

  **Exit criteria:**
  - On-premise/private-cloud deployment model is documented and validated with VEDO Hub Enterprise.
  - Dedicated API endpoints and SLA model are available for enterprise customers.
  - SAP SuccessFactors integration is supported.
  - Private corporate ontology isolation works through VEDO Hub Enterprise boundaries.
  - Predictive analytics cover deficit forecasting and churn-risk scenarios where data is available.

  **Dependencies:** M7, M8, enterprise deployment requirements.
  **Business goals:** G5, G6.

- [ ] **M10: Multilingual & Global Readiness**
  Prepare EduTrack for multilingual ontology consumption, localized UI, and international expansion.

  **Exit criteria:**
  - UI localization framework supports Russian and English.
  - Ontology content can be consumed in Russian and English where VEDO Hub provides multilingual data.
  - API and documentation conventions support multilingual labels and descriptions.
  - Product analytics can segment usage by locale.

  **Dependencies:** M6, M8, multilingual ontology readiness in VEDO Hub.
  **Business goals:** G5.

---

## Network Plan

```text
M0.0 → M0.1 → M0.2 → M0.3 → M1
                                  ├─→ M2 ─→ M5
                                  ├─→ M4 ─→ M6 ─→ M8 ─→ M10
                                  └─→ M3 ─→ M7 ─→ M9
```

- Phase 0 is strictly sequential because requirements, architecture, platform, and scaffold depend on one another.
- M1 is the MVP technical pivot: product, visualization, corporate pilot, and integrations all depend on real route/plan/gap mechanics.
- M4 can run in parallel with M2 and M3 after M1, but M3 needs partial M4 capabilities for HR/LMS integration.
- M5 enriches the Community product after the MVP route/visualization loop is validated.
- M6 matures partner integration after the API foundation is stable.

---

## Assumptions & Risks

| Risk / Assumption | Impact | Mitigation |
|-------------------|--------|------------|
| Starter ontology is not ready on time | Blocks M1/M2 validation | Track ontology readiness as an explicit external dependency with acceptance criteria |
| VEDO Hub API/MCP contract changes | Breaks route engine and integrations | Add contract tests, API versioning policy, and boundary ADRs in M1/M4 |
| Technology stack remains TBD too long | Blocks M0.2/M0.3 | Complete stack ADR during M0.1 before platform work starts |
| EduTrack duplicates Hub responsibilities | Product boundary drift and implementation waste | Maintain explicit Hub/EduTrack responsibility matrix in architecture docs and API contracts |
| EdTech partners hesitate to share ontology contributions | Slower network effect | Preserve private forks and content ownership; expose clear value of public contribution through Hub mechanics |
| Enterprise pilots require long integration cycles | Delays corporate validation | Keep M3 pilot scope narrow and use M4 integration sandbox/webhooks to reduce setup cost |
| Forecast accuracy cannot be proven early | Risk to G1/G4/G6 metrics | Start with binary forecast in MVP, collect data, then upgrade to risk-level forecast in M5 |

---

## Completed

| Milestone | Date |
|-----------|------|
| *(none yet)* | |

---

## Milestone ↔ Business Goals Traceability

| Milestone | Business Goals |
|-----------|----------------|
| M0.0: Requirements Baseline | G1, G2, G3, G4, G6 |
| M0.1: Domain Model & Architecture Baseline | G1, G2, G3, G4, G5, G6 |
| M0.2: Engineering Platform | Engineering enabler for all business goals |
| M0.3: Runnable Product Scaffold | Engineering enabler for all business goals |
| M1: Core Infrastructure | G1, G3 |
| M2: Family Education | G1, G3, G4 |
| M3: Corporate Onboarding Pilot | G6 |
| M4: Integration & Webhook Layer | G2, G5 |
| M5: Route Enrichment from Extended Hub Ontology | G1, G3, G4 |
| M6: EdTech Platform Integration | G2 |
| M7: Corporate Application & Compliance | G4, G6 |
| M8: Community Consumption & Network Effect | G5 |
| M9: Enterprise Deployment & Compliance | G5, G6 |
| M10: Multilingual & Global Readiness | G5 |

## Milestone ↔ Feature Areas Traceability

| Milestone | Feature Areas |
|-----------|---------------|
| M0.0 | Requirements, acceptance criteria, traceability |
| M0.1 | Architecture, domain model, RBAC, Hub/EduTrack boundary |
| M0.2 | Dev platform, CI, test scaffold, containers |
| M0.3 | App/API scaffold, auth shell, ontology stub, route stub |
| M1 | F0, F1, F2, F3, F6 |
| M2 | F2, F4, F5 |
| M3 | F1, F2, F4, F6 for corporate onboarding pilot |
| M4 | F6 |
| M5 | F0, F1, F2, F3, F5 |
| M6 | F6 |
| M7 | F1, F2, F4, F6 for corporate/compliance scenarios |
| M8 | F6 plus community-enriched ontology consumption |
| M9 | F6, enterprise deployment, compliance analytics |
| M10 | Localization, multilingual ontology consumption |
