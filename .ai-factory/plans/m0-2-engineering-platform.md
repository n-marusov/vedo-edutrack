# Implementation Plan: M0.2 — Engineering Platform

Branch: none
Created: 2026-08-03

## Settings
- Testing: no
- Logging: standard
- Docs: no

## Roadmap Linkage
Milestone: "M0.2: Engineering Platform"
Rationale: This plan implements the full M0.2 exit criteria from ROADMAP.md — repository layout, Docker configs, Makefile, CI pipeline, and test scaffolds.

## Progress

### Phase 1: Repository Structure & Module Initialization
- [x] T1: Initialize Go module and backend directory tree
- [x] T2: Initialize React + Vite frontend application
- [x] T3: Create deploy/ directory structure

### Phase 2: Containerization & Dev Environment
- [x] T4: Backend Dockerfile (distroless, multi-stage)
- [x] T5: Frontend Dockerfile (Go embed SPA)
- [x] T6: Docker Compose with full dev stack
- [x] T7: Traefik reverse-proxy configuration
- [x] T8: Container strategy documentation (deploy/README.md)

### Phase 3: Build Automation
- [ ] T9: Makefile (up / down / dev / build / test / lint / migrate / hooks / ci)

### Phase 4: Development Tooling
- [ ] T10: Pre-commit hooks + Go lint config + Go hot-reload
- [ ] T11: Frontend lint/formatter config (Biome, TypeScript)

### Phase 5: CI Pipeline
- [ ] T12: GitHub Actions CI workflow

### Phase 6: Test Scaffolds
- [ ] T13: Backend unit/integration test scaffolds (placeholder, red)
- [ ] T14: Frontend test scaffolds (Vitest + RTL, placeholder, red)
- [ ] T15: E2E test scaffold (Playwright, placeholder, red)

### Phase 7: Gate Automation
- [ ] T16: Unified gate manifest + two-tier runner

## Commit Plan
- **Commit 1** (after T1–T3): `chore: scaffold repository directory tree`
- **Commit 2** (after T4–T8): `feat: add Docker configs and container strategy`
- **Commit 3** (after T9): `feat: add Makefile build automation`
- **Commit 4** (after T10–T11): `chore: configure development tooling (lint, pre-commit)`
- **Commit 5** (after T12): `ci: add GitHub Actions CI pipeline`
- **Commit 6** (after T13–T15): `test: add placeholder test scaffolds`
- **Commit 7** (after T16): `ci: add unified gate manifest and two-tier runner`

## Tasks

### Phase 1: Repository Structure & Module Initialization

- [x] **T1: Initialize Go module and backend directory tree**

  Create the backend directory skeleton following ADR-IMPL.PROCESS.repository-structure.md (§2, §4) and ADR-DES.API.cli-interface.md:
  ```
  backend/
  ├── cmd/vedo-edutrack/main.go        # Thin entry: cli.Execute() (composition root, wire)
  ├── internal/
  │   ├── cli/                         # Cobra command tree (input adapter, ADR-DES.API.cli-interface)
  │   │   └── cli.go                   # package cli placeholder (root/server/mcp/migrate/… in M0.3+)
  │   ├── modules/                    # 10 bounded contexts (empty module stubs)
  │   │   ├── routeplanning/           # core: domain/, application/, adapters/
  │   │   ├── executionprogress/       # core
  │   │   ├── gapcoverage/             # core
  │   │   ├── planmanagement/          # supporting
  │   │   ├── ontologyport/            # supporting (ACL to VEDO Hub)
  │   │   ├── resources/               # supporting
  │   │   ├── practicelife/            # supporting
  │   │   ├── visualization/           # supporting (read models)
  │   │   ├── identityaccess/          # generic (auth, RBAC)
  │   │   └── integrations/            # generic (REST/SPARQL/webhooks/MCP/LMS/SSO)
  │   ├── pkg/                         # Shared utilities (infrastructure-neutral)
  │   └── platform/                    # Shared adapters
  │       ├── postgres/                # connection.go, migration.go (Atlas)
  │       ├── telemetry/               # tracer.go (OTel), metrics.go (Prometheus)
  │       ├── config/                  # env-based configuration loader
  │       ├── logger/                  # zap + otelzapbridge
  │       └── wire.go                  # shared DI providers
  ├── migrations/                      # Atlas: <timestamp>_<schema>_<desc>.sql
  ├── api/openapi/v1.yaml              # OpenAPI spec (source of truth)
  ├── tests/                           # E2E/integration (testcontainers, cross-module)
  └── go.mod                           # module vedo-edutrack/backend
  ```

  **Note:** binary is `vedo-edutrack` with cobra subcommands (server / mcp / migrate / seed / …) — renamed from `cmd/server` per ADR-DES.API.cli-interface.

  Each module directory gets a minimal Clean Architecture skeleton (ADR-DES.INFRA.clean-architecture-adoption):
  - `domain/` — package with empty `.go` file
  - `application/commands/` — package with empty `.go` file
  - `application/queries/` — package with empty `.go` file
  - `adapters/` — package with empty `.go` file

  `.gitkeep` files in empty directories where Go packages are not yet needed (migrations/, tests/).

  **Source:** ADR-IMPL.PROCESS.repository-structure.md §2, §4; ADR-DES.API.cli-interface.md; DDD context-map.md (10 contexts).
  **Logging:** INFO — log Go module initialization path and directory counts.

- [x] **T2: Initialize React + Vite frontend application**

  Create the frontend application skeleton (ADR-DES.STACK.framework-vs-vs, ADR-IMPL.PROCESS.development-tooling.md §11):
  ```
  frontend/
  ├── src/
  │   ├── index.tsx                    # React entry point
  │   ├── App.tsx                      # Root component (placeholder)
  │   ├── design/                      # Design tokens (Tailwind v4 @theme)
  │   ├── features/                    # Feature modules (mirror backend bounded contexts)
  │   │   ├── route-planning/
  │   │   ├── execution-progress/
  │   │   ├── gap-coverage/
  │   │   ├── plan-management/
  │   │   ├── ontology-port/
  │   │   ├── resources/
  │   │   ├── practice-life/
  │   │   ├── visualization/
  │   │   ├── identity-access/
  │   │   └── integrations/
  │   ├── shared/                      # Shared UI components
  │   ├── store/                       # Zustand stores
  │   └── styles/
  │       └── index.css                # Tailwind v4 entry (@import "tailwindcss")
  ├── design/                          # Pencil .pen files (design process)
  ├── package.json                     # Dependencies: react, react-dom, react-router-dom
  ├── vite.config.ts                   # Vite config with React plugin
  ├── tsconfig.json                    # Strict TypeScript config
  ├── index.html                       # Vite HTML entry point
  └── biome.json                       # (moved to T11)
  ```

  Package dependencies (added via `pnpm add` from workspace root; frontend lives under `apps/web` per pnpm-workspace.yaml):
  - `react`, `react-dom`, `react-router-dom` (core)
  - `tailwindcss`, `@tailwindcss/vite` (styling)
  - `zustand` (state management)
  - `@types/react`, `@types/react-dom` (dev)
  - `typescript`, `vite`, `@vitejs/plugin-react` (dev)

  **Note:** The frontend is a pnpm workspace package under `apps/web`, not a standalone directory. The `frontend/` root directory maps to `apps/web` in the pnpm workspace, keeping the mono-repo structure clean. Adjust `pnpm-workspace.yaml` to include `"frontend"` (or create `apps/web` and symlink/reference).

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §11; ADR-IMPL.PROCESS.repository-structure.md §3.
  **Logging:** INFO — log package installation output and directory structure creation.

- [x] **T3: Create deploy/ directory structure**

  Create the deploy/ layout (ADR-IMPL.PROCESS.repository-structure.md §4):
  ```
  deploy/
  ├── ci/                              # CI scripts and pipeline helpers (empty, for T12)
  ├── keycloak/                        # Keycloak realm, clients, role mapping (placeholder for M0.3)
  ├── observability/                   # OTel Collector, Prometheus, Loki, Tempo, Grafana configs
  │   ├── collector-config.yaml        # OTel Collector: receivers (OTLP), exporters (Prometheus, Loki, Tempo)
  │   ├── prometheus.yml               # Prometheus scrape config
  │   ├── loki-config.yaml             # Loki config
  │   ├── tempo-config.yaml            # Tempo config
  │   ├── grafana-datasources.yml      # Grafana datasource provisioning
  │   └── grafana-dashboards/          # Dashboard JSONs (empty, provisioned later)
  ├── postgres/                        # PostgreSQL init scripts
  │   └── init.sql                     # Placeholder: CREATE EXTENSION, initial schema flag
  ├── seeds/                           # Seed data (RBAC role catalog, demo data — placeholder for M0.3)
  ├── traefik/                         # Traefik configs
  │   ├── traefik.yml                  # Static config: entrypoints, providers, API
  │   └── dynamic.yml                  # Dynamic config: routers, services, middlewares
  ├── docker-compose.yml               # Full dev stack (T6)
  └── README.md                        # Container strategy documentation (T8)
  ```

  **Source:** ADR-IMPL.PROCESS.repository-structure.md §4; ADR-IMPL.PROCESS.development-tooling.md §8, §9.
  **Logging:** INFO — log directory creation and file counts.

### Phase 2: Containerization & Dev Environment

- [x] **T4: Backend Dockerfile (distroless, multi-stage)**

  Create `backend/Dockerfile` following ADR-IMPL.PROCESS.development-tooling.md §8:
  - Multi-stage build: `golang:1.24-alpine` (build) → `gcr.io/distroless/static-debian12:nonroot` (run)
  - Build stage: download deps, compile with `-ldflags="-s -w"`, `CGO_ENABLED=0`
  - Run stage: copy binary, set `USER nonroot:nonroot`, expose port 8080
  - Health check via binary subcommand or lightweight endpoint
  - Labels: org.opencontainers annotations (source, version, description)

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §8; DESCRIPTION.md (Go stack).
  **Logging:** INFO — log Docker build stage names and binary size.

- [x] **T5: Frontend Dockerfile (Go embed SPA)**

  Create `frontend/Dockerfile` — Go embed approach for single-artifact deployment (ADR modular-monolith §7, §8; ADR development-tooling §8):
  - Build stage: `node:24-alpine` — `pnpm install`, `pnpm build` (vite build → static assets)
  - Embed stage: Go binary reads embedded `dist/` and serves via `embed.FS`
  - The embed server is part of the backend binary (`backend/cmd/vedo-edutrack/main.go` wires the SPA handler), so this Dockerfile is used for standalone frontend development only
  - Production: SPA assets are embedded in the Go binary — single artifact for on-prem Enterprise simplicity

  **Alternative valid approach:** Standalone SPA Dockerfile with nginx (variant B from ADR). For M0.2, implement both:
  - `Dockerfile.embed` — Go embed variant (production, on-prem)
  - `Dockerfile.nginx` — nginx variant (SaaS/CDN-friendly)
  The docker-compose uses nginx for dev mode (hot-reload via Vite dev server, no Docker build needed).

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §8; ADR-DES.INFRA.modular-monolith-approach.md §7.
  **Logging:** INFO — log build approach used and output size.

- [x] **T6: Docker Compose with full dev stack**

  Create `deploy/docker-compose.yml` — full development environment (ADR-IMPL.PROCESS.development-tooling.md §8):

  Services:
  - `backend` — Go binary with hot-reload (air), mounts `./backend`, port 8080
  - `frontend` — Vite dev server (hot-reload), mounts `./frontend`, port 5173
  - `postgres` — PostgreSQL 16 with init scripts, port 5432
  - `otel-collector` — OTel Collector with config from `deploy/observability/`
  - `prometheus` — metrics scraping, port 9090
  - `loki` — log aggregation, port 3100
  - `tempo` — distributed tracing, port 3200
  - `grafana` — dashboards + alerting, port 3000, provisioned datasources
  - `traefik` — reverse proxy (T7), ports 80/443

  Networks: `edutrack-net` (internal), `edutrack-public` (Traefik ingress).
  Volumes: `postgres_data`, `grafana_data`, `loki_data`, `tempo_data`.
  Health checks on PostgreSQL, backend readiness endpoint.
  `.env.example` with defaults: ports, credentials, OTel sampling rates.

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §8, §9; ADR modular-monolith §5.
  **Logging:** INFO — log service startup order and health check status.

- [x] **T7: Traefik reverse-proxy configuration**

  Create Traefik configuration files (ADR-IMPL.PROCESS.development-tooling.md §8):
  - `deploy/traefik/traefik.yml` — static config:
    - Entrypoints: `web` (80 → redirect to 443), `websecure` (443)
    - Providers: Docker (watch), file (dynamic.yml)
    - API/Dashboard: enabled (dev only), TLS via Let's Encrypt staging (dev) / production (CI)
    - Access logs, metrics (Prometheus)
  - `deploy/traefik/dynamic.yml` — dynamic config:
    - Routers: `api` → backend:8080, `spa` → frontend:5173 (dev) / backend (embed, prod)
    - Middlewares: rate limiting, headers (CSP, HSTS per OWASP NFR), circuit breaker
    - Services: backend, frontend load-balanced
  - Blue-green deployment labels (post-MVP: two backend services, Traefik weighted round-robin)

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §7, §8; REQ-NFR-security.compliance.owasp-application-security.
  **Logging:** INFO — log entrypoints created, routers configured, TLS status.

- [x] **T8: Container strategy documentation**

  Write `deploy/README.md` — container strategy aligned with stack and deployment contours (ADR modular-monolith, ADR development-tooling, M0.2 exit criterion "container strategy documented and aligned with selected stack"):

  Sections:
  - **Architecture overview**: modular monolith → one deployable artifact; SPA embedded via Go embed for on-prem Enterprise simplicity; docker-compose for dev + staging
  - **Image strategy**:
    - Backend: distroless multi-stage (~15 MB), Go binary with embedded SPA
    - Frontend dev: Vite hot-reload (no Docker build needed in dev)
  - **Dev environment**: `docker compose up -d` starts all 9 services, one command
  - **SaaS deployment**: Traefik + compose, blue-green via weighted round-robin
  - **Enterprise on-prem**: single Go binary (embedded SPA) + PostgreSQL only; Redis optional; no orchestration required on MVP; K8s/Helm post-MVP
  - **Observability stack**: OTel Collector → Prometheus/Loki/Tempo + Grafana, provisioned as-code
  - **CI integration**: `make ci` mirrors CI pipeline; images built and pushed by GitHub Actions
  - **Security**: nonroot user, CSP/HSTS headers via Traefik, rate limiting, no secrets in images
  - **Contours**: Community (full stack, SaaS) vs Enterprise (minimal, on-prem, 242-ФЗ)

  **Source:** ADR-DES.INFRA.modular-monolith-approach.md; ADR-IMPL.PROCESS.development-tooling.md §7, §8; M0.2 exit criteria.
  **Logging:** N/A (documentation artifact).

### Phase 3: Build Automation

- [ ] **T9: Makefile (up / down / dev / build / test / lint / migrate / hooks / ci)**

  Create `Makefile` at repository root — single entry point for all workflows (ADR-IMPL.PROCESS.development-tooling.md §11):

  | Target | Command | Description |
  |--------|---------|-------------|
  | `up` | `docker compose -f deploy/docker-compose.yml up -d --wait` | Start dev environment |
  | `down` | `docker compose -f deploy/docker-compose.yml down --volumes` | Stop and clean up |
  | `dev` | `up` + `air` (Go hot-reload) in background | Development mode |
  | `build` | `cd backend && go build -ldflags="-s -w" -o bin/vedo-edutrack ./cmd/vedo-edutrack` + `cd frontend && pnpm build` | Production build check (single binary with cobra CLI) |
  | `test` | `cd backend && go test ./... -count=1` + `cd frontend && pnpm test` | All tests |
  | `lint` | `cd backend && golangci-lint run ./...` + `cd frontend && pnpm biome ci .` | Lint both ends |
  | `format` | `cd backend && gofmt -l -w .` + `cd frontend && pnpm biome check --write .` | Auto-format |
  | `dev-check` | `bash deploy/ci/run-gates.sh --tier fast` | Fast feedback loop: compile, types, lint, format, touched-unit, gen-consistency, mermaid (T16) |
  | `check` | `bash deploy/ci/run-gates.sh --tier delivery` | Full delivery gate set for task handoff (T16) |
  | `migrate` | `cd backend && go run ./cmd/vedo-edutrack migrate up` (wraps Atlas, ADR-DES.API.cli-interface) | Apply DB migrations |
  | `migrate-down` | `cd backend && go run ./cmd/vedo-edutrack migrate down` | Revert last migration |
  | `hooks` | `pre-commit install` | Install git hooks |
  | `ci` | `bash deploy/ci/run-gates.sh --tier delivery --trigger ci` | Full local CI run (delegates to gate runner, T16) |
  | `clean` | `down` + `rm -rf backend/bin/ frontend/dist/` | Full cleanup |

  Additional requirements:
  - `.PHONY` declarations for all targets
  - `SHELL := /bin/bash` for consistency
  - Default target: `help` — prints available targets with descriptions
  - `make ci` delegates to the gate runner (`run-gates.sh --tier delivery --trigger ci`, T16) — single source of truth is `deploy/ci/gates.yaml`, not the CI YAML
  - Environment variables read from `.env` (auto-loaded via `include .env` or `-include .env`)
  - `DATABASE_URL` default: `postgres://edutrack:edutrack@localhost:5432/edutrack?sslmode=disable`
  - Idempotent targets: running twice doesn't break state

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §11; M0.2 exit criterion "build automation covers up/down/build/test/lint".
  **Logging:** INFO — Makefile targets echo their step name on entry and exit status.

### Phase 4: Development Tooling

- [ ] **T10: Pre-commit hooks + Go lint config + Go hot-reload**

  Configure development tooling for the Go backend (ADR-IMPL.PROCESS.development-tooling.md §1, §2):

  **`.pre-commit-config.yaml`** (repo root):
  - `local-biome-check`: `npx @biomejs/biome check --write --files-ignore-unknown=true --no-errors-on-unmatched --staged` (frontend files)
  - `gofmt`: `gofmt -l -w` (Go files only)
  - `golangci-lint`: `golangci-lint run --new-from-rev=HEAD~1` (Go files only)

  **`.golangci.yml`** (backend/):
  - Linters: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `depguard` (ban cross-module imports per modular monolith rules)
  - `depguard` config: modules under `internal/modules/` must NOT import sibling modules; only `internal/pkg/` and `internal/platform/` allowed
  - `gocyclo` max complexity: 10 (per REQ-NFR-process.dev.code-complexity)
  - `goconst`, `misspell`, `nilerr`, `nolintlint`, `prealloc`, `revive`
  - Exclude: `vendor/`, `mocks/`, generated code (`wire_gen.go`, `sqlc/`)
  - Run timeout: 5 min

  **`.air.toml`** (backend/):
  - Watch: `.go` files in `backend/`
  - Exclude: `tmp/`, `vendor/`, `_test.go`, `mocks/`
  - Build command: `go build -o ./tmp/vedo-edutrack ./cmd/vedo-edutrack`
  - Run command: `./tmp/vedo-edutrack server`
  - Log file: `tmp/air.log`
  - Full bin rebuild on change (modular monolith = single binary)

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §1, §2; ADR modular-monolith §6 (dependency rule enforcement).
  **Logging:** INFO — log linter version at hook install, lint violations count.

- [ ] **T11: Frontend lint/formatter config (Biome, TypeScript)**

  Configure frontend tooling (ADR-IMPL.PROCESS.development-tooling.md §1, §11):

  **`frontend/biome.json`**:
  - Formatter: indent 2 spaces, semicolons always, single quotes, trailing commas, line width 100
  - Linter: recommended rules + `useExhaustiveDependencies`, `useHookAtTopLevel` (React hooks)
  - JavaScript: `react` framework, JSX runtime `react-jsx`
  - Organize imports: enabled
  - Ignore: `node_modules/`, `dist/`, `coverage/`, `vite.config.ts.timestamp-*`

  **`frontend/tsconfig.json`**:
  - `strict: true`, `noUncheckedIndexedAccess: true`, `noUnusedLocals: true`, `noUnusedParameters: true`
  - `module: "ESNext"`, `moduleResolution: "bundler"`, `target: "ES2022"`
  - `jsx: "react-jsx"`
  - Paths alias: `@/*` → `src/*`
  - Include: `src/`, `vite.config.ts`
  - Exclude: `node_modules/`, `dist/`

  **Tailwind v4 setup** (frontend/src/styles/index.css):
  ```css
  @import "tailwindcss";
  @theme {
    /* Design tokens from frontend/src/design/ (placeholder, populated later) */
  }
  ```

  **Vite config** (`frontend/vite.config.ts`):
  - `@vitejs/plugin-react`
  - `@tailwindcss/vite` plugin
  - Dev server: port 5173, proxy `/api` → `http://localhost:8080`
  - Alias: `@` → `./src`

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §1, §11; ADR DDD context-map (feature modules mirror bounded contexts).
  **Logging:** INFO — log Biome version, TypeScript compilation status.

### Phase 5: CI Pipeline

- [ ] **T12: GitHub Actions CI workflow**

  Create `.github/workflows/ci.yml` — complete CI pipeline (ADR-IMPL.PROCESS.development-tooling.md §7):

  Trigger: `push` to all branches, `pull_request` to `main`.
  Concurrency group: cancel in-progress runs for same branch/PR.

  Jobs (sequential where arrows indicate dependency). Each job is a thin wrapper over the gate runner (T16): `bash deploy/ci/run-gates.sh --trigger ci --group <job>`. The actual commands live only in `deploy/ci/gates.yaml` (single source of truth); this YAML stays a skeleton without command duplication.

  1. **`lint`** (parallel: backend + frontend):
     - Backend: `golangci-lint run --out-format=github-actions ./...` (with annotations)
     - Frontend: `biome ci . --reporter=github` (with annotations)
     - Go format check: `test -z "$(gofmt -l .)"`

  2. **`typecheck`** (frontend only):
     - `pnpm tsc --noEmit`

  3. **`test`** (runs after lint + typecheck pass):
     - Backend unit: `go test -race -count=1 -coverprofile=coverage.out ./...`
     - Backend integration: `go test -tags=integration -count=1 ./...` (requires PostgreSQL service container)
     - Frontend: `pnpm vitest --run --coverage`
     - **Services:** PostgreSQL 16 service container (for integration tests)

  4. **`coverage-gate`** (runs after tests, reports only — gate is advisory at M0.2):
     - Go: `go tool cover -func=coverage.out` → check core modules ≥ 90% (warning only)
     - Frontend: vitest coverage report (warning only)

  5. **`security`**:
     - SAST: `gosec ./...` (Go)
     - SBOM: `syft backend/cmd/vedo-edutrack -o spdx-json=backend-sbom.spdx.json`
     - Secret scan: `gitleaks detect --source . --verbose`
     - Dependency audit: `pnpm audit` (frontend)

  6. **`build`** (runs after tests + security pass):
     - Backend: Docker build `backend/` (distroless, verify image size ≤ 20 MB)
     - Frontend: `vite build` (verify output)
     - Push images to GHCR (on `main` branch only, with `:latest` and `:sha-<commit>` tags)

  Shared setup:
  - `actions/setup-go@v5` with `go-version-file: backend/go.mod`
  - `pnpm/action-setup@v4` with `version: 9` (from pnpm-workspace)
  - Caching: Go modules + pnpm store
  - Timeout: 15 min per job (MTR ≤ 2 h across all jobs)

  Deferred stages (declared in the manifest, not active at M0.2): `mutation` (gremlins, `phase: m1`, advisory — core F1/F2 lands on M1, REQ-NFR-process.dev.test-coverage) and `smoke` (post-MVP). Declared so DESCRIPTION.md's CI/CD chain (lint → test → mutation → coverage → security → build → deploy → smoke) and the pipeline never drift.

  **`deploy/ci/`** helpers:
  - `wait-for-postgres.sh` — health-check loop for PostgreSQL readiness in integration test job
  - `coverage-check.sh` — parse coverage output and fail below threshold (configurable `COVERAGE_MIN`)

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §7; REQ-NFR-process.dev.engineering-gates; REQ-NFR-infra.compliance.cicd-supply-chain-security.
  **Logging:** INFO — CI steps log job name on entry/exit; test commands output via `go test -v`; coverage gate logs percentages.

### Phase 6: Test Scaffolds

- [ ] **T13: Backend unit/integration test scaffolds (placeholder, red)**

  Create intentionally failing (red) test scaffolds across backend modules, following ADR repository-structure §5 (test structure by layer):

  **Per core module** (`routeplanning`, `executionprogress`, `gapcoverage`):
  - `domain/domain_test.go` — placeholder unit test with `t.Error("TODO: implement domain tests")`
  - `application/commands/*_test.go` — placeholder with testify mock skeleton
  - `application/queries/*_test.go` — placeholder with testify mock skeleton
  - `adapters/handler/*_test.go` — placeholder (skipped, `t.Skip("TODO: integration test")`)
  - `adapters/repository/*_test.go` — placeholder (skipped)

  **Per supporting context** (`planmanagement`, `ontologyport`, `resources`, `practicelife`, `visualization`):
  - `domain/domain_test.go` — placeholder red test

  **Generic contexts** (`identityaccess`, `integrations`):
  - `domain/domain_test.go` — placeholder red test

  **Integration tests** (`backend/tests/`):
  - `integration_test.go` — placeholder with `t.Skip("TODO: integration tests with testcontainers")`
  - `testhelper_test.go` — test helper utilities (setup/teardown, Docker container lifecycle)

  **Pattern:** Each red test includes a `TODO:` comment referencing the relevant UC/FR to guide implementation later. Use `testing.T` and standard library (no testify for scaffold tests — keep it minimal).

  **Source:** ADR-IMPL.PROCESS.repository-structure.md §5; REQ-NFR-process.dev.test-coverage (core ≥ 90% target).
  **Logging:** INFO — `go test` output shows intentional failures with TODO messages.

- [ ] **T14: Frontend test scaffolds (Vitest + RTL, placeholder, red)**

  Create intentionally failing test scaffolds for frontend (ADR repository-structure §3, ADR development-tooling §6):

  **`frontend/vitest.config.ts`**:
  - Environment: `jsdom`
  - Setup file: `src/test/setup.ts` (RTL cleanup, matchMedia mock)
  - Coverage: v8 provider, thresholds not enforced (advisory)

  **Placeholder tests** (`frontend/src/features/*/__tests__/`):
  - One `.test.tsx` per feature module:
    ```tsx
    import { describe, it, expect } from "vitest";
    import { render, screen } from "@testing-library/react";

    describe("<feature-name>", () => {
      it("should render the <feature> placeholder component", () => {
        expect.fail("TODO: implement <feature> component tests (US-*, UC-*)");
      });
    });
    ```

  **Shared component test** (`frontend/src/shared/__tests__/`):
  - Placeholder for design-system components

  **Source:** ADR-IMPL.PROCESS.development-tooling.md §6; REQ-NFR-process.dev.test-coverage.
  **Logging:** INFO — vitest output shows intentional failure with TODO markers.

- [ ] **T15: E2E test scaffold (Playwright, placeholder, red)**

  Create E2E test scaffold with intentionally failing test (M0.2 exit criterion: "E2E test scaffolds exist with intentionally red or placeholder tests where implementation is pending"):

  **`e2e/`** directory (separate from frontend, not a workspace package):
  ```
  e2e/
  ├── playwright.config.ts             # Playwright config: baseURL, browsers, reporter
  ├── package.json                     # playwright + @playwright/test
  ├── tests/
  │   └── placeholder.spec.ts          # Red placeholder test
  └── fixtures/
      └── auth.ts                      # Auth fixture skeleton (placeholder for M0.3 Keycloak)
  ```

  **`playwright.config.ts`**:
  - Projects: Chromium, Firefox, WebKit
  - `baseURL: "http://localhost:5173"` (Vite dev) / `http://localhost:8080` (prod, embedded SPA)
  - Reporter: `html` + `github` (CI annotations per REQ-NFR-ops.release.canary-kill-switch)
  - Retries: 0 (CI: 2)
  - Timeout: 30s per test
  - Web server: `make dev` (auto-start for local runs)
  - Screenshot on failure

  **`tests/placeholder.spec.ts`**:
  ```typescript
  import { test, expect } from "@playwright/test";

  test("placeholder — fails until M0.3 scaffold is ready", async ({ page }) => {
    test.fail(true, "TODO: implement E2E tests for MVP Must-scenarios (M0.3+)");
    await page.goto("/");
    await expect(page).toHaveTitle(/VEDO EduTrack/);
  });
  ```

  **10 MVP Must-scenario stubs** (ADAR repository-structure §5): outline as `test.skip()` blocks referencing UC/US IDs — these guide future implementation without running.

  **Source:** ADR-IMPL.PROCESS.repository-structure.md §5; M0.2 exit criterion "E2E test scaffolds exist with intentionally red/placeholder tests".
  **Logging:** INFO — Playwright outputs test status with TODO annotations visible in report.

### Phase 7: Gate Automation

- [ ] **T16: Unified gate manifest + two-tier runner**

  Create a single source of truth for quality gates and a runner that both the dev loop (fast tier) and the delivery handoff (full tier) execute. Implements the gate-automation design from RESEARCH (session 2026-08-03 02:43) and makes the "task ready for delivery" decision executable: the agent runs the full tier before declaring a task done.

  **`deploy/ci/gates.yaml`** — gate manifest (single source of truth):
  - Each gate: `id`, `command`, `tier: fast|delivery`, `phase` (m0.0…m10), `severity: blocking|advisory`, `group` (lint|typecheck|test|coverage|gen|db|validate|security|build), `trigger` (local|ci|ci-main|precommit), `needs` (postgres, stack-up), `nfr` (requirement ID), optional `runner: agent` for non-command gates
  - `current_phase` synchronized with ROADMAP (M0.2)
  - Deferred stages declared explicitly so DESCRIPTION.md and CI don't drift: `mutation` (gremlins, `phase: m1`, advisory — core F1/F2 lands on M1), `smoke` (post-MVP)

  **`deploy/ci/run-gates.sh`** — runner:
  - `--tier fast` (dev loop, ≤ 2–3 min, no Postgres/Docker/E2E): `go build ./...`, `pnpm tsc --noEmit`, gofmt/biome format check, `golangci-lint run ./...`, `biome ci .`, unit tests for touched modules (auto-scope from `git diff --name-only`), `make gen && git diff --exit-code` (openapi + sqlc), `pnpm validate:mermaid:all`, `gitleaks detect`
  - `--tier delivery` (task handoff): fast + `go test -race -count=1 -coverprofile=coverage.out ./...`, `go test -tags=integration -count=1 ./...` (Postgres service), e2e Playwright (phase ≤ current — regression: current + all previous phases), docker build (distroless ≤ 20 MB), `gosec ./...`, `pnpm audit`, `syft ... -o spdx-json` (ci-main), `atlas migrate validate`, `coverage-check.sh --min 90` (advisory at M0.2)
  - Phase filter: run all gates with `phase <= current_phase` — current + all previous phases always execute
  - Output: human table + `--out-format json` → aggregated `aif-gate-result` JSON (schema v1: status pass|warn|fail, blocking, blockers, affected_files, suggested_next); exit 1 on blocking fail
  - Agent-declared gates (`runner: agent`): TQS/RCS via `/aif-test-quality`, docs/env/drift via `/aif-verify` Step 3 — listed in the report, executed by the agent

  **Integration:**
  - Makefile (T9): `make dev-check` → `run-gates.sh --tier fast`; `make check` → `run-gates.sh --tier delivery`; `make ci` → `--tier delivery --trigger ci`
  - CI (T12): jobs call `run-gates.sh --trigger ci --group <job>`; commands removed from YAML
  - Agent hooks (documented for `aif-implement`): Step 3.4 → `--tier fast` after each task; Final Step → `--tier delivery` mandatory before `/aif-verify`

  **Decision rule:** task ready for delivery ⇔ `--tier delivery` completes with zero blocking fails.

  **Source:** REQ-NFR-process.dev.engineering-gates (P0); REQ-NFR-process.dev.test-coverage (P1); RESEARCH session 2026-08-03 02:43; `aif-gate-result` contract (aif-verify/aif-review/aif-rules-check/aif-security-checklist).
  **Logging:** INFO — runner logs each gate (id, status, duration); JSON summary at end.

---

## Dependencies Between Tasks

```
T1 ──┐
T2 ──┤
T3 ──┼── T4, T5 (Dockerfiles depend on directory structure)
     │
     ├── T6 (compose depends on T4, T5, T7)
     │      └── T7 (Traefik config depends on T3 deploy structure)
     │
     ├── T8 (docs — can be written anytime after T3)
     │
     ├── T9 (Makefile — references paths from T1, T2, T3, T6)
     │
     ├── T10, T11 (tooling configs — independent of Docker, depend on T1, T2)
     │
     ├── T12 (CI — depends on T9 Makefile targets + T10, T11 configs)
     │
     ├── T13, T14, T15 (test scaffolds — depend on T1, T2 directory structure)
     │
     └── T16 (gate manifest + runner — depends on T9 Makefile, T12 CI, T13–T15 scaffolds; unifies them)
```

**Execution order by phase:**
- Phase 1 (T1, T2, T3) — all parallel (disjoint write sets)
- Phase 2 (T4, T5, T7) — parallel after T1, T2, T3; T6 after T4, T5, T7; T8 anytime after T3
- Phase 3 (T9) — after T1, T2, T3, T6 (needs compose paths)
- Phase 4 (T10, T11) — parallel after T1, T2
- Phase 5 (T12) — after T9, T10, T11
- Phase 6 (T13, T14, T15) — parallel after T1, T2
- Phase 7 (T16) — after T9, T12, T13, T14, T15

## Architecture References

| ADR | Role in M0.2 |
|-----|-------------|
| `ADR-DES.INFRA.modular-monolith-approach` | Single artifact, Go embed SPA, in-process events, PostgreSQL only (no Redis on MVP) |
| `ADR-DES.INFRA.clean-architecture-adoption` | Module skeleton layout: domain / application / adapters |
| `ADR-IMPL.PROCESS.repository-structure` | Full repo layout, 10 bounded context modules, test structure by layer |
| `ADR-IMPL.PROCESS.development-tooling` | Makefile targets, CI pipeline, Docker configs, pre-commit, Biome, Go tools |
| `ADR-DES.STACK.language-vs-vs` | Go (backend) + TypeScript (frontend) |
| `ADR-DES.STACK.framework-vs-vs` | chi + oapi-codegen (backend), React + Vite (frontend) |
| `ADR-DES.DATA.storage-strategy` | PostgreSQL, sqlc + Atlas migrations |
| `specs/ddd/context-map` | 10 bounded contexts (core/supporting/generic) with exact slugs and relationships |
| `deploy/ci/gates.yaml` + `deploy/ci/run-gates.sh` | Unified gate manifest + two-tier runner (T16): single source of truth for fast/delivery gates |

## Non-Functional Requirements Addressed

| NFR | How M0.2 addresses it |
|-----|----------------------|
| `REQ-NFR-process.dev.engineering-gates` | CI pipeline with lint → test → coverage → security gates (§5, T12); unified gate manifest + runner makes the acceptance criteria executable (T16) |
| `REQ-NFR-process.dev.test-coverage` | Delivery gate tier: full unit (race) + integration + e2e (phase ≤ current), coverage-check (advisory at M0.2 → blocking at M1), mutation deferred to m1 (T16) |
| `REQ-NFR-process.dev.code-complexity` (CC ≤ 10) | golangci-lint gocyclo configured in T10 |
| `REQ-NFR-process.dev.developer-documentation` | Container strategy doc (T8), deploy README |
| `REQ-NFR-security.compliance.owasp-application-security` | CSP/HSTS headers via Traefik (T7), nonroot user (T4), rate limiting (T7), secret scan in CI (T12) |
| `REQ-NFR-infra.compliance.cicd-supply-chain-security` | SBOM via syft, gitleaks secret scan, pnpm audit in CI (T12) |
| `REQ-NFR-ops.compliance.support-sla` | Single artifact for on-prem Enterprise simplifies support (T5, T8) |
| `REQ-NFR-infra.compliance.environment-isolation` | Docker Compose networks + dev-only configs (T6, T7) |
| `REQ-NFR-ops.observability.*` | OTel stack provisioned as-code in docker-compose (T6, T3) |
| `REQ-NFR-ops.release.deployment-verification` (drift = 0) | Atlas migrations in Makefile `migrate` target (T9) |
| `REQ-NFR-ops.release.cicd-resilience` (MTR ≤ 2 h) | CI pipeline caching + parallel jobs (T12) |
| `REQ-NFR-data.availability.migration-rollback` (≤ 15 min) | Atlas `migrate-down` target (T9) |
| `REQ-NFR-infra.availability.multi-az-geo-dr` (≥ 2 zones) | Traefik load balancing (T7), stateless replicas (T6) |
