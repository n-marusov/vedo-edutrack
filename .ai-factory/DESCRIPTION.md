# VEDO EduTrack

> Product vision source: `vision.md` (Russian). This description summarizes the project scope for AI Factory tooling.

## Overview

VEDO EduTrack is an educational route service built **on top of the VEDO Hub** ontology platform. It reads knowledge ontologies from VEDO Hub via its REST API and MCP server, and adds unique educational mechanics:

- **Route planning**: `Route = f(learner, goal, pedagogy concept, ontology) → route`. Computes personalized learning paths by walking the knowledge graph, respecting link types, the essential core, and learner pace. Auto-recompute on triggers. Plan fixation as a route snapshot with timeline.
- **Plan execution**: tracks progress against the fixed plan — "plan vs actual" comparison, deviation alerts, readiness forecasts. Gap diagnosis: finds the root cause of a learner's lag by climbing strict-prerequisite links up the graph to the first unmastered module. Live FGOS/professional-standard coverage, deficit lists, attestation readiness reports.
- **Educational visualization**: knowledge map with progress/gap color-coding, dashboards for learners, parents, HR and methodologists, route builder.
- **Resource & context matching**: attaches materials, stories and project ideas to route modules, filtered by format, difficulty and learning style.

## Core Features

- **F0. Knowledge ontology port** (via VEDO Hub): reading the ontology through Hub's REST API and MCP server — modules, 5 link types (`hasStrictPrerequisite` / `hasSoftPrerequisite` / `enriches` / `appliesTo` / `isAnalogousTo`), FGOS mappings, resources, stories, project ideas, pedagogy concepts. Copying the relevant subgraph for in-memory computation. Infrastructure foundation — not a user-facing domain.
- **F1. Route planning**: `Route = f(position, goal, pedagogy concept, ontology) → route`. Shortest-path with strict/soft/enrich link weights, three horizons (far / mid / near), auto-recompute triggers (progress, goal change, ontology update), cascade recompute. Plan fixation: route snapshot with timeline. Input constraints: checkpoints and FGOS/framework requirements. Pedagogy-concept-aware routing. Inter-subject schedule sync. Resource matching.
- **F2. Plan execution**: progress tracking against the fixed plan — plan-vs-actual with deviation reasons. Completion forecast (on-track / at-risk / off-track). Deviation alerts. Gap diagnosis: root-cause analysis by climbing strict-prerequisite links. Assessment items and IRT calibration. Live FGOS/professional-standard coverage, deficit lists, attestation readiness reports. Real-knowledge-to-framework mapping.
- **F3. Resources**: all resource types bound to modules — content resources (video, text, interactive, textbooks) and enabling resources (tutors, lab equipment, access passes, budget). Filters by format, style, difficulty, duration. Availability checks, cost estimation, route budget.
- **F4. Visualization**: 2D knowledge graph with progress colors, gap map view, learner dashboard, parent/HR/methodologist dashboards, route builder, group management panel.
- **F5. Real-world connection**: stories, context and project ideas served at the moment of module mastery (via `appliesTo`/`enriches` graph links). Qualities map (what traits a module develops). Motivation through practical relevance.
- **F6. Integrations**: REST API, read-only SPARQL endpoint, LMS connectors (WebTutor, iSpring, SAP SuccessFactors), webhooks (`module.mastered`, `plan.deviated`, `route.recalculated`), SSO/SAML (Keycloak), MCP server for AI agents.

## MVP Scope

- 1000+ topics for grades 5–11 (math, biology, physics, chemistry, history, literature, geography, computer science, social studies), 500+ base cross-subject links, full FGOS mapping, 3–5 pedagogy concepts.
- Route planning: shortest path with strict links, recompute on progress/goal change, plan fixation.
- Plan execution: plan-vs-actual, binary readiness forecast, gap diagnosis, FGOS coverage, deficits, attestation readiness report.
- Resources: catalog of materials bound to modules, format/source filtering, 50+ stories, 30+ project ideas at launch.
- Visualization: knowledge map, gap map, group panel, learner dashboard, route builder.
- Integrations: REST API, read-only SPARQL, webhooks.

**Excluded from MVP**: pedagogy-concept-aware routes, all 5 link types (MVP: strict + soft + appliesTo), ontology-update cascades, risk-level forecasts, learning-style resource matching, inter-subject schedule sync, out-of-the-box LMS connectors (API only).

## Tech Stack

> **Status: TBD — deferred by decision (2026-08-02).** The stack will be chosen by the user before implementation planning.

- **Programming language:** *TBD* (candidates: TypeScript / Python / Java)
- **Framework:** *TBD* (candidates: Next.js / FastAPI+NestJS / Spring Boot)
- **Database:** *TBD* (candidate: PostgreSQL for learner/plan/progress data; knowledge graph lives in VEDO Hub, queried via SPARQL API)
- **ORM:** *TBD* (candidate: Prisma / SQLAlchemy / JPA)
- **External platform:** VEDO Hub (REST API + MCP server, SPARQL/Cypher endpoint, ontology storage/versioning)

## Architecture Notes

- EduTrack is a **service layer** over VEDO Hub: it never stores or edits ontologies — it reads them through the Hub API and adds route/plan/gap/visualization mechanics.
- Two product contours share one API contract: **Community** (public ontology, families, EdTech) and **Enterprise** (private corporate ontology, onboarding, compliance).
- Route is a **function, not a document**: recomputed on any trigger (module mastered, ontology updated, goal changed).
- Route and trajectory are two always-visible layers: progress is measured against the plan, recommendations come from the current route, and the trajectory records the actual path taken.
- Essentialism: each context defines a mandatory core; order and pace vary, dependency logic is unbreakable.
- FGOS/professional standards act as a filter — coverage is always queryable.

## Non-Functional Requirements

- **Logging:** configurable via `LOG_LEVEL`.
- **Error handling:** structured error responses.
- **Security:** role-based access (learner / parent / school / methodologist / HR), private corporate ontology isolation, PII protection (152-ФЗ context for Enterprise).
- **Performance:** route recompute on graph traversal must stay interactive; SPARQL queries via VEDO Hub API.
- **Observability:** readiness forecasts and plan deviations as first-class events (`module.mastered`, `plan.deviated` webhooks).
