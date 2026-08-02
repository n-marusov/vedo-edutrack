# VEDO EduTrack

> Product vision source: `vision.md` (Russian). This description summarizes the project scope for AI Factory tooling.

## Overview

VEDO EduTrack is an educational route service built **on top of the VEDO Hub** ontology platform. It reads knowledge ontologies from VEDO Hub via its REST API and MCP server, and adds unique educational mechanics:

- **Route engine**: `Route = f(learner, goal, pedagogy concept, ontology) → route`. Computes personalized learning paths by walking the knowledge graph, respecting link types, the essential core, and learner pace.
- **Gap diagnosis**: finds the root cause of a learner's lag by climbing strict-prerequisite links up the graph to the first unmastered module.
- **Learning plan**: snapshots of the route at checkpoints, "plan vs actual" tracking, deviation forecasts, notifications.
- **Educational visualization**: knowledge map with progress/ gap color-coding, dashboards for learners, parents, HR and methodologists, route builder.
- **Resource & context matching**: attaches materials, stories and project ideas to route modules, filtered by format, difficulty and learning style.
- **FGOS / professional-standard coverage**: live coverage of formal requirements, deficit lists, readiness forecasts.

## Core Features

- **F1. Knowledge ontology** (via VEDO Hub): topics, 5 types of cross-subject links (`hasStrictPrerequisite` / `hasSoftPrerequisite` / `enriches` / `appliesTo` / `isAnalogousTo`), FGOS mappings, resources, stories, project ideas, pedagogy concepts.
- **F2. Route engine**: SPARQL queries to VEDO Hub, shortest-path with strict-link constraints, auto-recompute triggers (progress, goal change, ontology update), three horizons (far / mid / near).
- **F3. Learning plan**: checkpoint snapshots, plan-vs-actual, deviation tracking, basic readiness forecast.
- **F4. Checkpoints & FGOS/professional standards**: live requirement coverage, deficit lists, attestation readiness reports.
- **F5. Resources, stories & project ideas**: catalog bound to topics, format/source filtering, 50+ stories, 30+ project ideas at launch.
- **F6. Knowledge graph & route visualization**: 2D knowledge map with progress colors, gap map view, group dashboards (parents / school director / HR), learner dashboard, route builder.
- **F7. Methodologist community** (in VEDO Hub): forks, merge requests, contributor profiles.
- **F8. Corporate module**: onboarding plans, time-to-productivity metric, gap analysis to target role.
- **F9. Integrations**: REST API for EdTech (`GET /routes/compute`, `GET /progress/{learner_id}`, `GET /fgos/coverage/{learner_id}`), read-only SPARQL endpoint, webhooks (`module.mastered`, `plan.deviated`).

## MVP Scope

- 1000+ topics for grades 5–11 (math, biology, physics, chemistry, history, literature, geography, computer science, social studies), 500+ base cross-subject links, full FGOS mapping, 3–5 pedagogy concepts.
- Route engine: shortest path with strict links, gap diagnosis, recompute on progress/goal change.
- Learning plan: checkpoint snapshots, plan-vs-actual, binary readiness forecast.
- FGOS coverage, deficits, attestation readiness report.
- Visualization: knowledge map, gap map, group panel, learner dashboard, route builder.
- Community: fork + basic merge request + contributor profile.
- Corporate: onboarding plan + time-to-productivity metric.
- Integrations: REST API, read-only SPARQL, webhooks.

**Excluded from MVP**: pedagogy-concept-aware routes, all 5 link types (MVP: strict + soft + appliesTo), ontology-update cascades, risk-level forecasts, learning-style resource matching, semantic merge comparison, DOI, out-of-the-box LMS connectors (API only), career tracks & compliance module (onboarding only), corporate dashboard.

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
