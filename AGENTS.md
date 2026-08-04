# AGENTS.md

> This file is a structural map for AI agents working in this repository. It is maintained by `/aif`; update it when the project structure changes significantly.

## Project Overview

VEDO EduTrack is an educational trajectory service built on top of the VEDO Hub ontology platform. It reads knowledge ontologies via the Hub's REST API and MCP server, then adds educational mechanics: trajectory generation, learning plans, gap diagnosis, FGOS coverage, and knowledge-graph visualization. Product vision lives in `vision.md` (Russian); the machine-readable summary is in `.ai-factory/DESCRIPTION.md`.

## Tech Stack

> Stack selected 2026-08-02 (ADRs T3–T5); see `specs/adr/` and `.ai-factory/DESCRIPTION.md` for details.

- **Backend:** Go + chi + oapi-codegen (OpenAPI-first), modular monolith, Clean Architecture; единый бинарник `vedo-edutrack` с cobra-подкомандами (server / mcp / migrate / seed / ontology sync / route compute / plan get / gap diagnose / report) — `ADR-DES.API.cli-interface`
- **Frontend:** React + TypeScript SPA (Vite), Zustand, Tailwind CSS v4, React Flow + Cytoscape.js
- **Database:** PostgreSQL (learner/plan/progress data; knowledge graph lives in VEDO Hub)
- **ORM/data access:** sqlc + pgx, Atlas migrations
- **Auth:** JWT RS256/JWKS (jwx); Keycloak post-MVP (Enterprise SSO)
- **Observability:** OTel (Go+Web) → Collector → Prometheus/Loki/Tempo + Grafana
- **CI/CD:** GitHub Actions; Docker (distroless) + docker-compose + Traefik; K8s post-MVP
- **Testing:** Playwright E2E (`tests/e2e/gui` + `tests/e2e/api`), Vitest + RTL (фронт), Go test + testify + testcontainers (бэкенд)
- **External platform:** VEDO Hub (REST API + MCP server, SPARQL/Cypher endpoint)

## Project Structure

```
vedo-edutrack/
├── vision.md              # Product vision & business requirements (Russian, authoritative)
├── LICENSE                # Project license
├── .ai-factory.json       # AI Factory v2.15.0 manifest: installed skills, MCP flags
├── package.json           # pnpm workspace root: вспомогательный тулинг (валидация Mermaid, scripts/)
├── pnpm-workspace.yaml    # pnpm workspace: корень + apps/* (React SPA) и tools/*
├── .nvmrc                 # Закреплённая версия Node (24)
├── traceability.ttl       # OWL 2 DL traceability model (US → UC → FR → ADR/C4 → COMP → TEST)
├── Makefile               # Единая точка входа: up/down/dev/build/test/lint/migrate/dev-check/check/ci
├── lefthook.yml           # Git-хуки (Lefthook): biome (frontend) + gofmt + golangci-lint (backend)
├── .github/workflows/     # GitHub Actions: ci.yml — тонкие обёртки над deploy/ci/run-gates.sh
├── backend/               # Go-модуль (modular monolith, единый бинарник vedo-edutrack)
│   ├── Dockerfile         # distroless multi-stage (build: golang:1.26-alpine → runtime: distroless nonroot)
│   ├── .golangci.yml      # Линт-конфиг: depguard (границы модулей), gocyclo ≤ 10, revive и др.
│   ├── .air.toml          # Go hot-reload для dev-контейнера (make dev)
│   ├── cmd/vedo-edutrack/ # Тонкий entry: cli.Execute() (composition root, wire)
│   ├── internal/
│   │   ├── cli/           # Cobra-дерево: server, mcp, migrate, seed, ontology sync, route compute,
│   │   │                  #   plan get, gap diagnose, report (входной адаптер, ADR-DES.API.cli-interface)
│   │   ├── modules/       # 10 bounded contexts (routeplanning, executionprogress, gapcoverage, …)
│   │   │                  #   + red-тест-скаффолды T13 (domain/application, handler/repository — skip)
│   │   ├── pkg/           # Общие утилиты (инфраструктурно-нейтральные)
│   │   └── platform/      # Общие адаптеры: postgres, telemetry, config, logger, auth, wire.go, ratelimit, circuitbreaker, eventbus
│   ├── migrations/        # Atlas: <timestamp>_<schema>_<desc>.sql
│   ├── api/openapi/v1.yaml# OpenAPI-спека (источник истины)
│   ├── tests/             # Интеграционные (Go, testcontainers, кросс-модульные) — скаффолд T13
│   └── go.mod
├── frontend/              # React SPA (Vite, pnpm-воркспейс)
│   ├── Dockerfile          # nginx вариант (SaaS/CDN)
│   ├── biome.json         # Линт/формат: single quotes, width 100, React-hooks правила
│   ├── vitest.config.ts   # Vitest + RTL (jsdom, coverage v8, setup src/test/setup.ts)
│   └── src/               # App, design/ (токены), features/ (10 модулей + __tests__), shared/, store/, styles/
├── deploy/                # Инфраструктура и CI/CD
│   ├── ci/                # gates.yaml (единый манифест гейтов) + run-gates.sh (fast/delivery) + хелперы
│   ├── docker-compose.yml # Dev-стек: 9 сервисов (backend+air, frontend+Vite, postgres, OTel, traefik)
│   ├── README.md          # Container strategy (M0.2 exit-criteria)
│   ├── keycloak/, observability/, postgres/, seeds/, traefik/
├── scripts/               # Вспомогательные скрипты (validate-mermaid.mjs и др.)
├── tests/                 # Системные тесты: e2e/gui (Playwright, скаффолд T15) + e2e/api + integration
├── specs/                 # Формализованные спецификации
│   ├── vision.md          # Продуктовое видение (рус., авторитетный источник)
│   ├── glossary.md        # Доменный глоссарий (единственный источник терминов)
│   ├── requirements/      # Требования: 60 FR + 54 NFR + MVP-ACCEPTANCE-CRITERIA.md (в т.ч. RBAC T8 и граница Hub T10 как NFR-спеки)
│   ├── use-cases/         # 42 use case (UC-<L1>.<L2>.<L3>)
│   ├── user-stories/      # 47 user stories (US-*, Gherkin)
│   ├── adr/               # Architecture Decision Records (ADR-DES/IMPL.*)
│   ├── ddd/               # DDD: context-map.md, aggregates.md, domain-events.md
│   └── c4/                # C4: context-system, container-overview, component-*
├── .agents/
│   └── skills/            # Installed agent skills (27 aif-* skills, project-local)
└── .ai-factory/           # AI Factory context (generated by /aif)
    ├── config.yaml        # Language, paths, git, rules configuration
    ├── DESCRIPTION.md     # Project spec (English, machine-readable)
    ├── rules/
    │   └── base.md        # Project base rules (placeholders — no code yet)
    └── ARCHITECTURE.md    # (pending — deferred until stack is chosen)
```

## Key Entry Points

| File | Purpose |
|------|---------|
| `vision.md` | Authoritative product vision: business requirements, MVP scope, roadmap, monetization |
| `specs/requirements/REQ-NFR-security.compliance.role-catalog.md` | RBAC role catalog (T8): universal archetypes (contract) + personas → role seed instances, Community/Enterprise catalogs |
| `specs/requirements/REQ-NFR-security.compliance.permission-matrix.md` | Role-permission matrix (T8): archetype × functional area × permission level (CRUD) × scope |
| `specs/requirements/REQ-NFR-api.compliance.ownership-boundary.md` | VEDO Hub ↔ EduTrack responsibility boundary (T10): data/computation/API/event/deployment ownership, ontology-port contract |
| `specs/requirements/` | Formalized requirements: 60 FR + 54 NFR, each with measurable acceptance criteria |
| `specs/adr/` | Architecture decision records (ADR-DES/IMPL.*): stack, storage, comm patterns, RBAC, CLI interface |
| `backend/cmd/vedo-edutrack/` | Единственный бинарник `vedo-edutrack`: тонкий entry (composition root, wire) → `cli.Execute()` |
| `backend/internal/cli/` | Cobra-дерево команд (входной адаптер): server, mcp, migrate, seed, ontology sync, route compute, plan get, gap diagnose, report — `ADR-DES.API.cli-interface` |
| `deploy/ci/gates.yaml` | Единый манифест гейтов качества (T16): tier fast/delivery, phase-регрессия, severity, trigger, group — единственный источник команд гейтов |
| `deploy/ci/run-gates.sh` | Двухуровневый раннер гейтов (T16): `--tier fast|delivery`, `--trigger`, `--group`, `--out-format table|json` (aif-gate-result v1); exit 1 при blocking-фейле |
| `Makefile` | Build automation (T9): up/down/dev/build/test/lint/format/migrate/hooks (lefthook install) + dev-check (fast) / check (delivery) / ci |
| `.github/workflows/ci.yml` | CI-пайплайн (T12): тонкие обёртки над run-gates.sh (lint → typecheck → test → coverage → security → build) |
| `backend/Dockerfile` | distroless multi-stage образ бэкенда (T4) — SPA embed встроен в backend-бинарник (M0.3) |
| `frontend/Dockerfile` | nginx-образ SPA для SaaS/CDN (T5) |
| `tests/e2e/gui/` | Playwright E2E (T15): config, placeholder.spec.ts (red), mvp-must-scenarios.spec.ts (M1–M10 стабы), fixtures/auth.ts |
| `frontend/src/features/*/__tests__/` | Компонентные тест-скаффолды (T14, Vitest + RTL, red) — по одному на фичу (зеркало bounded contexts) |
| `tests/` | System tests: `e2e/gui` (Playwright browser, M1–M10), `e2e/api` (Playwright API flows), `integration` (cross-layer, compose stack) |
| `specs/ddd/` | DDD artifacts: context map, aggregates, domain events |
| `specs/c4/` | C4 diagrams: System Context, Container, Component (Mermaid + legend) |
| `traceability.ttl` | OWL 2 DL traceability model linking US → UC → FR → ADR/C4 → COMP → TEST |
| `.ai-factory.json` | AI Factory manifest: installed skills and MCP server flags |
| `.ai-factory/DESCRIPTION.md` | Machine-readable project specification (English) |
| `.ai-factory/config.yaml` | AI Factory configuration: `language.ui: ru`, `language.artifacts: en`, git settings |

## Documentation

| Document | Path | Description |
|----------|------|-------------|
| README | `README.md` | Project landing page: what it is, why not an LMS, getting started, make targets (English) |
| Product vision | `specs/vision.md` | Business requirements, scope, MVP, roadmap (Russian) |
| Requirements | `specs/requirements/` | FR + NFR requirements with measurable acceptance criteria, README with conventions (Russian) |
| Use cases | `specs/use-cases/` | 42 UC covering route building, execution, viz, auth, API integration (Russian) |
| User stories | `specs/user-stories/` | 47 US in Gherkin with @UC/@FR tags (Russian) |
| Architecture decisions | `specs/adr/` | ADR records: stack, storage, comm patterns, repo structure, RBAC, CLI interface (Russian) |
| CLI interface | `specs/adr/ADR-DES.API.cli-interface.md` | Single binary with cobra subcommands; CLI as input adapter over Application layer; dev/support/testing tooling (Russian) |
| Domain model | `specs/ddd/` | Context map, aggregates, domain events (Russian) |
| RBAC matrix | `specs/requirements/REQ-NFR-security.compliance.*` | Role catalog + permission matrix + ops/admin separation (T8), Community/Enterprise differences (Russian) |
| Responsibility boundary | `specs/requirements/REQ-NFR-api.compliance.*` | EduTrack vs VEDO Hub ownership (T10): data, computation, API, events, deployment, ontology-port contract (Russian) |
| C4 diagrams | `specs/c4/` | Context/Container/Component diagrams with legend, context, F0–F6 links (Russian) |
| Traceability | `traceability.ttl` | OWL 2 DL: US → UC → FR → ADR/C4 → COMP → TEST chains |
| Project description | `.ai-factory/DESCRIPTION.md` | Spec for AI tooling (English) |

## AI Context Files

| File | Purpose |
|------|---------|
| AGENTS.md | This file — structural map for AI agents |
| .ai-factory/DESCRIPTION.md | Project specification (stack, features, NFRs) |
| .ai-factory/ARCHITECTURE.md | Architecture guidelines (pending — stack deferred) |
| vision.md | Full product vision & business context |

## Agent Rules

- Decompose shell commands into single steps rather than combining them.
  - Incorrect: `git checkout main && git pull`
  - Correct: First `git checkout main`, then `git pull origin main`
- **Auxiliary/tooling tasks run through pnpm** (root `package.json` scripts + `scripts/`): `pnpm install`, `pnpm validate:mermaid`, etc. Do **not** run ad-hoc `npm install` in the repo root — add tooling as devDependencies and install with pnpm (`pnpm add -D <pkg>`). Node version is pinned in `.nvmrc`.
- **Python (dev environment, Windows):** managed via `uv` — venv lives at `.venv` (CPython 3.13). Use the venv interpreter directly or after `uv venv` + activation:
  - `python` works inside the activated venv (`.venv\Scripts\activate`), **`python3` is NOT available on Windows** — always invoke `python`, never `python3`.
  - Run scripts: `.venv/Scripts/python.exe <script.py>` (or `python <script.py>` after activation).
  - Syntax check: `.venv/Scripts/python.exe -m py_compile <file.py>`.
  - Add packages with `uv pip install <pkg>` (venv must exist; create with `uv venv`). Example: `uv pip install requests` for `docs/integration/examples/python-example.py`.
  - The example client scripts under `docs/integration/examples/` are documented extras, not part of the Go backend build.
- **traceability.ttl is Python-only:** the traceability model (`traceability.ttl`, OWL 2 DL Turtle) MUST be read/edited/validated only through Python tooling — never by hand-editing without tooling, and never edited via the Go backend. After ANY edit, validate with the rdflib gate: `pnpm validate:traceability` (or `uv run --with rdflib python scripts/validate_traceability.py`, or `make validate-traceability`). The gate is enforced in CI as `traceability-validate` (fast tier, group validate). Use rdflib to inspect or modify the graph (e.g. `.venv/Scripts/python.exe -c "import rdflib; g=rdflib.Graph(); g.parse('traceability.ttl', format='turtle'); ..."`).
- The stack is **selected (2026-08-02)**: Go + chi + oapi-codegen backend, React + TS SPA frontend, PostgreSQL, VEDO Hub external. See `specs/adr/` (T3–T5) for the decision records.
- Product decisions must trace back to `vision.md` (authoritative) and `.ai-factory/DESCRIPTION.md` (summary).
- Communicate in Russian (`language.ui: ru`); write artifacts in English (`language.artifacts: en`).
