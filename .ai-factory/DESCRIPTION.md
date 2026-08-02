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

> **Status: выбран (ПРИНЯТО, 2026-08-02)** — зафиксирован в ADR T3–T5: `ADR-DES.STACK.language-vs-vs`, `ADR-DES.STACK.framework-vs-vs`, `ADR-DES.DATA.storage-strategy`, `ADR-DES.API.communication-patterns`, `ADR-IMPL.PROCESS.repository-structure`, `ADR-DES.SECURITY.rbac-model`, `ADR-IMPL.PROCESS.development-tooling`.

- **Programming language:** Go (бэкенд) + TypeScript (фронтенд) — `ADR-DES.STACK.language-vs-vs`
- **Framework:** chi + oapi-codegen (бэкенд, OpenAPI-first); React + TS (фронт, SPA) — `ADR-DES.STACK.framework-vs-vs`
- **Database:** PostgreSQL (learner/plan/progress data; knowledge graph lives in VEDO Hub, queried via SPARQL API)
- **Data access / migrations:** sqlc + Atlas (drift-детекция)
- **DI / logging / i18n:** wire · zap + OTel · go-i18n
- **Auth:** JWT RS256/JWKS (jwx); Keycloak — пост-MVP (Enterprise SSO)
- **Testing:** Go testing + testify · testcontainers-go · gremlins (spike) · Playwright · React Testing Library
- **Observability:** OTel (Go+Web) → Collector → Prometheus/Loki/Tempo + Grafana
- **Infra:** Docker (distroless, Go embed) · docker-compose + Traefik (blue-green) · K8s — пост-MVP
- **CI/CD:** GitHub Actions (lint → test → mutation → coverage → security → build → deploy → smoke)
- **Dev tools:** Biome (фронт, pre-commit) · gofmt + golangci-lint (бэкенд) · pre-commit framework — `ADR-IMPL.PROCESS.development-tooling`
- **External platform:** VEDO Hub (REST API + MCP server, SPARQL/Cypher endpoint, ontology storage/versioning)

## Architecture Notes

- EduTrack is a **service layer** over VEDO Hub: it never stores or edits ontologies — it reads them through the Hub API and adds route/plan/gap/visualization mechanics.
- Two product contours share one API contract: **Community** (public ontology, families, EdTech) and **Enterprise** (private corporate ontology, onboarding, compliance).
- Route is a **function, not a document**: recomputed on any trigger (module mastered, ontology updated, goal changed).
- Route and trajectory are two always-visible layers: progress is measured against the plan, recommendations come from the current route, and the trajectory records the actual path taken.
- Essentialism: each context defines a mandatory core; order and pace vary, dependency logic is unbreakable.
- FGOS/professional standards act as a filter — coverage is always queryable.

## Non-Functional Requirements

Formal NFR corpus: **49 NFR** in `specs/requirements/REQ-NFR-*.md`, each with measurable acceptance criteria, traced in `traceability.ttl`. Coverage areas:

- **Observability & ops:** structured logging (`LOG_LEVEL`, JSON, `request_id`/`trace_id`), golden-signals dashboards, alerting with P1–P4 escalation (noise ≤ 20%), distributed tracing, incident communication (status page), product metrics (NPS, forecast accuracy ±10%, FGOS coverage freshness).
- **Release & change:** deployment verification (drift = 0), canary releases with kill switch ≤ 5 min, change management (0 manual prod changes, auto-rollback), CI/CD resilience (MTR ≤ 2 h), maintenance windows (attestation-period protection), reversible DB migrations (auto-rollback ≤ 15 min).
- **Security:** role-based access (learner / parent / school / methodologist / HR), OWASP Top 10 (rate limiting, input validation, SPARQL parameterization, JWT RS256/JWKS), PII protection under 152-ФЗ (encryption at-rest/in-transit, incident notification ≤ 24 h), supply-chain security (pinning, SBOM, secrets scan), environment isolation (0 prod PII in dev), data residency in RU (242-ФЗ), archive security, production access (2-person rule, JIT).
- **Data & lifecycle:** RPO ≤ 1 h / RTO ≤ 4 h, backup validity (restore-tested), retention policy per category, PII export ≤ 30 min / deletion ≤ 30 days, processor registry, decommissioning safety, concurrent-write consistency (optimistic locking, 0 lost updates).
- **Performance & scalability:** API p95 ≤ 200 ms at 1000 concurrent, SPARQL route recompute < 1 s, horizontal scaling (10× load, autoscaling ≤ 5 min), multi-AZ (≥ 2 zones).
- **UX & ergonomics:** WCAG 2.1 AA (axe-core gate: 0 critical), supported browsers/OS (Chrome/Firefox/Safari/Edge, Windows 10+/macOS 12+/iOS 15+/Android 11+), code complexity gates (CC ≤ 10), admin console ergonomics, user competency & training, i18n-readiness (RU+EN, ICU, 0-code language addition; RTL deferred).
- **Support & docs:** tiered support SLA (Community ≤ 48 h / Pro ≤ 4 h / Enterprise ≤ 1 h), user documentation (100% scenarios), developer documentation (bus-factor ≥ 2).

Total requirements corpus: **109** (60 FR + 49 NFR), plus 42 UC and 47 US, all traced in `traceability.ttl` (0 orphans).
