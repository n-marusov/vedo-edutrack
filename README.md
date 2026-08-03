# VEDO EduTrack

[![Go](https://img.shields.io/badge/Go-1.26-blue?logo=go&logoColor=white)](backend/go.mod)
[![React](https://img.shields.io/badge/React-19-61dafb?logo=react&logoColor=white)](frontend/package.json)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.7-3178c6?logo=typescript&logoColor=white)](frontend/package.json)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169e1?logo=postgresql&logoColor=white)](deploy/docker-compose.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> The route is the unit, not the course. A learning-path engine that computes personalized educational routes over a shared knowledge ontology and drives them to a goal.

VEDO EduTrack is an educational route service built **on top of the VEDO Hub** ontology platform. It reads knowledge ontologies — school topics, cross-subject links, FGOS (Russian state standard) mappings, resources, stories — through the Hub's REST API and MCP server, and adds the educational mechanics that platforms and LMSs don't have: route planning, plan execution, root-cause gap diagnosis, live FGOS coverage, and knowledge-map visualization.

It does **not** produce or host content, and it does **not** replace an LMS. It is a different layer — the infrastructure of learning paths.

> **Status:** pre-MVP. The M0.2 engineering platform (repo layout, Docker stack, CI gates, test scaffolds, quality-gate runner) is complete. Product features (route engine, gap diagnosis, visualization) land from M0.3.

## Why an LMS is not enough

LMSs and learning platforms answer *"which course should I take?"* and measure progress *inside* a course. **EduTrack answers "how do I reach my goal?"** — it doesn't package content, it computes an individual route over a shared knowledge graph and manages its execution.

| Dimension | LMS / learning platforms | VEDO EduTrack |
|-----------|--------------------------|---------------|
| **Unit** | Course — a content package, the same for everyone | **Route** — a path through the knowledge graph to a goal |
| **Knowledge** | Course catalogs / isolated course graphs | **Open ontology** with cross-subject links |
| **Adaptivity** | Difficulty tuning inside a course | **Rebuilding the route itself** (a living route) |
| **Goal** | "Finish the course" (internal) | Role, exam, attestation (external) |
| **Tracking** | Course progress (% of lessons) | **Plan vs actual**, deviation, forecast, route rebuild |
| **Lagging behind** | "Repeat the lesson" | **Root-cause gap diagnosis** (climbing strict-prerequisite links) |
| **Reporting** | Gradebook / none | **Live FGOS & professional-standard coverage**, deficits, attestation readiness |
| **Content** | Produces and hosts | Does not produce — binds existing materials as resources |
| **Scope** | One subject / one course | **Multiple subjects** with cross-subject links |

### Where EduTrack wins

- **Families & schools** — a knowledge map with progress colors, an FGOS-checked route, and a plan-vs-actual forecast instead of hand-built Excel programs and a stress sprint before attestation.
- **EdTech platforms** — a ready-to-fork open ontology and a route engine: semantic recommendations (click-through 15–20% vs 3–5% for tag-based filtering), integration in days instead of months of manually building an isolated graph.
- **Corporate HR/L&D** — deficit analysis to a target role, onboarding grounded in real project context, compliance woven into the route (not a checkbox course), and measurable ROI.

## What it does

- **Route planning (F1)** — `Route = f(position, goal, pedagogy concept, ontology) → route`. Shortest path with strict/soft/enrich link types, three horizons, plan fixation as a snapshot with a timeline, auto-recompute on progress / goal change / ontology update.
- **Plan execution (F2)** — plan-vs-actual with deviation reasons, readiness forecast (on-track / at-risk / off-track), **root-cause gap diagnosis**, live FGOS coverage, prioritized deficits, attestation-readiness reports.
- **Knowledge map & dashboards (F4)** — 2D knowledge graph with progress/gap colors, learner/parent/HR/methodologist dashboards, group panel, route builder.
- **Resources & context (F3, F5)** — materials, stories, and project ideas bound to route modules; budget and availability planning.
- **Integrations (F6)** — REST API, read-only SPARQL, webhooks (`module.mastered`, `plan.deviated`, `route.recalculated`), LMS connectors, SSO/SAML (Keycloak), MCP server for AI agents.

## Architecture

EduTrack is a **modular monolith** (10 bounded contexts in one Go binary) with a React SPA frontend. The knowledge graph always lives in VEDO Hub — EduTrack reads it through the Hub API and never stores or edits ontologies.

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.26 · chi + oapi-codegen (OpenAPI-first) · single binary `vedo-edutrack` (cobra CLI) |
| Frontend | React 19 + TypeScript (Vite SPA) · Zustand · Tailwind CSS v4 |
| Data | PostgreSQL 16 · sqlc + Atlas migrations (drift = 0) |
| Edge & infra | Docker (distroless) · docker-compose · Traefik (TLS, rate limiting, blue-green) |
| Observability | OpenTelemetry → Prometheus / Loki / Tempo + Grafana |
| CI/CD | GitHub Actions · unified quality-gate runner (fast / delivery tiers) |

See the [C4 container diagram](specs/c4/container-overview.md) and [deployment diagrams](specs/c4/deployment-dev.md) for details.

## Getting started

### Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Docker | 24+ | Docker Desktop or daemon |
| GNU Make | 4.x | Git Bash or WSL on Windows |
| Go | 1.26 (per `backend/go.mod`) | only for `make build`/`make test` |
| Node.js | 24 (per `.nvmrc`) | `pnpm` for the frontend workspace |
| pnpm | 11 (per `package.json`) | `corepack enable` |

### Start the dev stack

```bash
git clone <repo-url> && cd vedo-edutrack
pnpm install            # workspace tooling + frontend deps
make hooks              # install lefthook git hooks (optional but recommended)
make up                 # start all 9 services (backend, frontend, postgres, OTel, traefik)
```

`make up` starts everything and waits for health checks. For hot reload instead:

```bash
make dev
```

### What you get

| Service | URL | Notes |
|---------|-----|-------|
| Frontend (Vite HMR) | http://localhost:5173 | SPA; proxies `/api` to the backend |
| Backend (air hot-reload) | http://localhost:8080 | `/healthz` (liveness), `/readyz` (readiness) |
| Traefik edge | http://localhost:8082 | dashboard (dev only); ingress on 80/443 |
| PostgreSQL | localhost:5432 | `edutrack` / `edutrack` |
| Grafana | http://localhost:3000 | Prometheus / Loki / Tempo datasources provisioned |
| Prometheus / Loki / Tempo | 9090 / 3100 / 3200 | observability backends |

Configuration lives in `deploy/.env.example` — copy it to `deploy/.env` and adjust (ports, credentials, sampling rate). Build version and URLs are injected at runtime (see [dynamic config](specs/adr/ADR-DES.INFRA.dynamic-config-injection.md)).

> [!NOTE]
> At M0.2 the backend serves only the health endpoints; domain APIs land in M0.3. Red test scaffolds intentionally fail until then — `make test` reports the pending work.

## Make targets

The `Makefile` is the single entry point. Run `make help` for the full list.

| Target | What it does |
|--------|--------------|
| `make up` / `make down` | Start / stop the full dev stack (compose, waits for health) |
| `make dev` | Dev mode — stack up with hot reload (air + Vite) |
| `make build` | Production build check (Go binary + SPA, version injected) |
| `make docker-build` | Production images — backend (distroless) + frontend (Go embed) |
| `make docker-build-all` | All production images (+ nginx frontend variant for SaaS/CDN) |
| `make test` / `make test-e2e` | All tests (Go unit, Vitest, Playwright) |
| `make lint` / `make format` | golangci-lint + Biome; gofmt + Biome format |
| `make gen` | Regenerate code (OpenAPI/sqlc — wired in M0.3) |
| `make dev-check` | Fast quality-gate tier (compile, types, lint, format, mermaid, secrets) |
| `make check` | Full delivery gate tier for task handoff |
| `make migrate` / `make migrate-down` | Atlas migrations via the CLI |
| `make hooks` | Install lefthook git hooks |
| `make gates` | Run all delivery gates (`make check` synonym) |
| `make gates-list` | List gates selected for the current phase and tier |
| `make gates-json` | Delivery gates, machine-readable JSON output |
| `make gates-<group>` | Per-group gates: `lint`, `typecheck`, `test`, `coverage`, `gen`, `db`, `validate`, `security`, `build` |
| `make ci` | Full local CI run (mirrors GitHub Actions) |
| `make clean` | Stop stack and remove build artifacts |

Run `make gates-list` to see which gates apply to the current phase.
Every target prints a colored verdict (`✅ [PASS]` / `❌ [FAIL]` / `⚠️ [WARN]` / `⏭️ [SKIP]`) — set `NO_COLOR=1` to disable colors (CI-safe); emojis are kept so verdicts stay scannable in plain logs.

### Docker images

`docker-build*` targets build **production** images only — the dev stack (`make up` / `make dev`) runs containers from base images via compose and needs no image build. Every image is tagged with the build version (`VERSION` — derived from `git describe`, override with `VERSION=x.y.z make docker-build`).

| Target | What it builds | Image tag(s) | Dockerfile | Environment |
|--------|----------------|--------------|------------|-------------|
| `make docker-build` | Default production delivery: backend + SPA embedded into the Go binary | `vedo-edutrack:<VERSION>` | `backend/Dockerfile` | Production — default; Enterprise on-prem single-binary path |
| `make docker-build-all` | All production images, incl. the nginx frontend variant | `vedo-edutrack:<VERSION>` · `vedo-edutrack-nginx:<VERSION>` | `backend/Dockerfile` · `frontend/Dockerfile` | Production — both contours (Community SaaS + Enterprise on-prem) |
| `make docker-build-backend` | Backend runtime image with embedded SPA (distroless, nonroot, healthcheck) | `vedo-edutrack:<VERSION>` | `backend/Dockerfile` | Production — backend service |
| `make docker-build-frontend-nginx` | Static SPA served by unprivileged nginx (UID 101, port 8080) | `vedo-edutrack-nginx:<VERSION>` | `frontend/Dockerfile` | Community SaaS / CDN (variant B) |

## Quality gates

All checks are declared in one manifest, `deploy/ci/gates.yaml`, and executed by `deploy/ci/run-gates.sh`:

- **Fast tier** (`make dev-check`) — the dev loop: compile, typecheck, lint, format, touched-module tests, generated-code consistency, mermaid validation, secret scan. No Postgres/Docker/E2E.
- **Delivery tier** (`make check` / `make gates`) — task handoff: everything in fast + full unit (race), integration (Postgres), Playwright E2E, distroless image build (≤ 20 MB), gosec, `pnpm audit`, Atlas migration validation, coverage.
- **Per-group** (`make gates-<group>`) — run gates for a single group, e.g.: `make gates-lint`, `make gates-security`, `make gates-test`. Run `make gates-list` for the full list.

A task is *ready for delivery* when the delivery tier completes with **zero blocking failures**.

## Built with AI Factory

This project is developed with the **AI Factory** agent toolchain: requirements, use cases, user stories, and ADRs are formalized in `specs/` and wired together by a machine-readable traceability model (`traceability.ttl` — 0 orphans, US → UC → FR → ADR → code → tests). Every artifact — specs, ADRs, C4 diagrams, this README — is generated and maintained by AI Factory skills, and the quality-gate runner is what makes an AI-assisted workflow verifiable.

## Documentation

| Document | Description |
|----------|-------------|
| [`specs/vision.md`](specs/vision.md) | Product vision and business requirements (Russian, authoritative) |
| [`AGENTS.md`](AGENTS.md) | Repository map for AI agents: structure, entry points, conventions |
| [`specs/`](specs/) | Requirements (60 FR + 54 NFR), use cases, user stories, ADRs, DDD, C4 diagrams |
| [`specs/c4/`](specs/c4/) | C4 architecture: system context, containers, components, deployment |
| [`deploy/README.md`](deploy/README.md) | Container & deployment strategy (Community SaaS vs Enterprise on-prem) |
| [`traceability.ttl`](traceability.ttl) | OWL 2 DL traceability model |

## License

[MIT](LICENSE) © 2026 n-marusov
