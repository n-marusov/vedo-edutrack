# Project Base Rules — VEDO EduTrack

> Project-specific base conventions derived from the chosen stack (2026-08-02, ADR T3–T5)
> and the M1 implementation. Loaded after `.ai-factory/RULES.md` (axioms — IDs, traceability,
> Hub/EduTrack boundary, pnpm). Keep this file in sync with the actual codebase.

## Naming Conventions

- **Go backend:** files `snake_case.go` (`sync_service.go`, `compute_service.go`); exported types/functions `PascalCase`; unexported helpers `camelCase`; package names single lowercase word.
- **Frontend:** directories `kebab-case` (`execution-progress`, `ontology-port`); component files and exports `PascalCase.tsx` (`OntologyBrowser.tsx`, `RouteBuilder.tsx`); variables/functions `camelCase`; single quotes, width 100 (biome.json).
- **SQL migrations:** `<NNNNNN>_<schema>_<desc>.sql` (e.g. `000002_planmanagement_init.sql`), one file per schema, embedded via `go:embed`.
- **OpenAPI:** `backend/api/openapi/v1.yaml` is the single source of truth; JSON fields `snake_case` (`learner_id`, `goal_topic_id`); operation IDs `camelCase`; generated code is committed via `make gen` and never hand-edited.
- **Environment variables:** `UPPER_SNAKE_CASE` with a `VEDO_` prefix for Hub-related config; document every new var in `.env.dev.example` (root) and `deploy/.env.example` in the same change.

## Module Structure

- **Bounded contexts** live under `backend/internal/modules/<context>/` with three layers: `domain/` (pure, no I/O, no logging), `application/` (use cases: `commands/`, `queries/`, `*_service.go`), `adapters/` (`hub/`, `repository/`, `handler/`, `stub/`).
- **Shared infrastructure** lives in `backend/internal/platform/` (config, logger, postgres, telemetry, auth, wire.go) and is imported by adapters only, never by domain.
- **CLI** is an input adapter (ADR-DES.API.cli-interface): thin cobra commands in `backend/internal/cli/` that call application services; one command per file (`route_compute.go`, `gap_diagnose.go`).
- **Repositories** are hand-written pgx adapters over the SQL contract in `sqlc/queries.sql` until the sqlc toolchain is installed; row models live in `sqlc/models.go`.
- **Frontend features** mirror bounded contexts: `frontend/src/features/<feature>/` with `api.ts`, `use*` hooks, components, and `__tests__/`; shared UI in `frontend/src/shared/components/`; routes registered in `frontend/src/routes.tsx`.
- **API versioning:** the REST API is served under the `/api/v1` prefix; the frontend `apiBaseUrl` must match it (`/api/v1`).

## Error Handling

- **REST errors** use the OpenAPI `ErrorResponse` shape: `{error: <machine-code>, message: <detail>, endpoint: <path>}`; status codes map to error classes (400 invalid, 401/403 auth, 404 not found, 501 stub).
- **Domain errors** are returned as values, not panics; use typed errors for structured handling (e.g. `*UnreachableGoalError` in route planning, `ErrNotFound` in repositories).
- **SQL errors** are mapped in the repository `errors.go` (`pgx.ErrNoRows` → `ErrNotFound`), never leaked raw to the API layer.
- **Input validation** happens at the adapter boundary (HTTP handler / CLI) before application calls.

## Logging

- **Structured JSON logging** via zap (platform/logger), level controlled by `LOG_LEVEL` (debug|info|warn|error).
- **Component prefixes** identify the module and area: `[ontologyport.hub]`, `[route.compute]`, `[plan.fixation]`, `[execution.progress]`, `[gap.diagnose]`, `[coverage]`, `[ontology.sync]`, `[repository.<Module>]`.
- **Levels:** `DEBUG` for GraphQL queries/DB calls with params; `INFO` for service boundaries and key events (sync started/completed, path found, plan fixed); `WARN` for anomalies (unreachable goal, at-risk forecast, incremental changes); `ERROR` with full context (`reason`, `status`, `err`).
- **Config changes** (`config.go`) log a startup line via `[config.Load]`.

## Testing

- **TDD:** tests are written before implementation (RED), then made GREEN; domain test suites mirror FR acceptance criteria.
- **Go:** table-driven tests with `-count=1`; testify for assertions; testcontainers for Postgres-backed integration tests behind the `integration` build tag (`go test -tags=integration`).
- **NFR-critical paths** get `testing.B` benchmarks (pathfinder 5k, gap 1k, catalog 10k, hub copy 10k) with thresholds from FR acceptance criteria; run via `make bench`.
- **Frontend:** Vitest + React Testing Library in `__tests__/*.test.tsx` (jsdom); mock API modules (`vi.mock`) instead of hitting the network.
- **E2E:** Playwright — GUI specs in `tests/e2e/gui` (auth via the demo login UI), API contract specs in `tests/e2e/api` (absolute `API_BASE` URLs against the compose stack).
- **Contract tests** protect external boundaries (VEDO Hub GraphQL) in `backend/tests/contract/`.

## Tooling (Auxiliary Tasks)

- Auxiliary/tooling tasks run through **pnpm** (root `package.json` scripts + `scripts/`), never ad-hoc `npm install` in the repo root (see RULES.md axiom).
- **Format/lint:** gofmt + golangci-lint (backend), Biome (frontend), enforced by Lefthook git hooks.
- **Code generation:** `make gen` regenerates OpenAPI code via oapi-codegen (`types` + `chi-server`, excluding `issueToken`); committed generated files must match the spec (gen-consistency gate).
- **Database:** Atlas-style migrations via the embedded runner (`vedo-edutrack migrate up/down/validate`); drift = 0 is the target.
- **Dev env:** the compose stack reads `.env.dev` (root, git-ignored; template `.env.dev.example`); `make up` passes it with `--env-file .env.dev`.
- `pnpm validate:mermaid` — validates mermaid blocks in `specs/c4/*.md`; `pnpm validate:mermaid:all` — across all specs (used in CI / before committing C4 or DDD diagrams).
