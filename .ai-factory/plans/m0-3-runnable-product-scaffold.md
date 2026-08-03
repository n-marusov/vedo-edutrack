# Implementation Plan: M0.3 — Runnable Product Scaffold

Branch: main
Created: 2026-08-03

## Settings
- Testing: no
- Logging: standard (INFO level — key events only)
- Docs: yes (mandatory docs checkpoint at completion)

## Roadmap Linkage
Milestone: "M0.3: Runnable Product Scaffold"
Rationale: Delivers the minimal authenticated product skeleton proving the chosen stack, local environment, API shell, ontology stub, route stub, and role-aware UI shell work together — the last Phase 0 milestone before MVP feature development.

---

## Overview

M0.3 transforms the M0.2 engineering scaffold (empty directories, platform stubs, red tests) into a runnable product skeleton. Every platform stub becomes a real implementation; the SPA embeds into the backend binary; health checks, auth, RBAC, ontology/route stubs, and role-aware dashboard shells are wired end-to-end.

**Boundary:** this plan covers the scaffold skeleton only — real domain logic, database schemas, VEDO Hub integration, and production hardening belong to Phase I (M1+).

**Mock VEDO Hub (Phase 8):** the plan additionally scaffolds a minimal GraphQL stand-in for the VEDO Hub ontology platform — a separate Docker container (`hub-mock`) serving the ontology-service schema from an in-memory `.ttl` ontology (first seed `traceability.ttl`). It is dev/test infrastructure (not Hub integration): the same handler runs as `httptest.Server` for M1 contract tests, and the container closes the dead `VEDO_HUB_API_URL=http://localhost:8081` gap in the dev stack and CI (see `ADR-DES.INFRA.mock-hub-strategy`, T20).

---

## Commit Plan

| # | After Tasks | Conventional Commit |
|---|-------------|---------------------|
| 1 | T1–T2 | `feat: add Go dependencies and wire platform services (zap, OTel, pgx)` |
| 2 | T3–T4 | `feat: wire cobra CLI and chi HTTP server with graceful shutdown` |
| 3 | T5–T6 | `feat: implement JWT auth middleware and RBAC engine with seed` |
| 4 | T7–T9 | `feat: expand OpenAPI spec and implement ontology + route stubs` |
| 5 | T10–T12 | `feat: add frontend dependencies, auth context, API client, and routing` |
| 6 | T13–T16 | `feat: build landing page, login, role-aware dashboards, and protected routes` |
| 7 | T17–T19 | `feat: embed SPA into backend binary with multi-arch Docker` |
| 8 | T20–T25 | `feat: add mock VEDO Hub GraphQL container for dev/test/CI` |
| 9 | T26–T28 | `feat: integration verification — health checks, end-to-end smoke test, gate pass` |

---

## Tasks

### Phase 1: Backend Dependency Wiring

- [ ] **T1: Add Go dependencies (go.mod)**

  Add all Go external dependencies required for M0.3. This is the first time `go.mod` gains a `require` block and `go.sum`.

  Dependencies to add:
  ```
  github.com/go-chi/chi/v5          # HTTP router
  github.com/go-chi/cors/v2         # CORS middleware (SPA)
  github.com/spf13/cobra             # CLI tree
  github.com/jackc/pgx/v5            # PostgreSQL driver + pool
  github.com/google/wire             # Compile-time DI
  github.com/lestrrat-go/jwx/v2      # JWT RS256 + JWKS
  github.com/vektah/gqlparser/v2    # GraphQL parsing (mock VEDO Hub, test-only — Phase 8)
  go.uber.org/zap                    # Structured logging
  go.opentelemetry.io/otel           # OTel SDK
  go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc  # OTLP gRPC trace exporter
  go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc # OTLP gRPC metric exporter
  go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp      # HTTP auto-instrumentation
  github.com/prometheus/client_golang  # /metrics endpoint
  ```

  Run `go mod tidy` after adding imports. Verify `go build ./...` compiles (stubs still ok).

  **Files:** `backend/go.mod`, `backend/go.sum`
  **Logging:** INFO — log each dependency group added; log `go mod tidy` result.

- [ ] **T2: Wire real platform services (zap, OTel, pgx)**

  Replace all four platform stubs with real implementations. Each retains the same function signature used in `platform/wire.go` so the call site is unchanged.

  **2a. Logger** (`backend/internal/platform/logger/logger.go`)
  - `New(level string) (*zap.Logger, error)` — JSON encoder, level from config, caller skip=1
  - `NewNop() *zap.Logger` for tests
  - `Sync()` helper for graceful shutdown
  - Env: `LOG_LEVEL` (default: `info`)

  **2b. OTel Tracer** (`backend/internal/platform/telemetry/tracer.go`)
  - `InitTracer(ctx, serviceName string) (*sdktrace.TracerProvider, error)`
  - OTLP gRPC exporter to `OTEL_EXPORTER_OTLP_ENDPOINT` (default: `http://otel-collector:4317`)
  - Resource attributes: `service.name`, `service.version`, `deployment.environment`
  - Sampler: `OTEL_SAMPLING_RATIO` env (default `0.1`) — 100% error + 10% success
  - `Shutdown(ctx)` helper for graceful shutdown

  **2c. OTel Meter** (`backend/internal/platform/telemetry/metrics.go`)
  - `InitMeter(ctx string) (*sdkmetric.MeterProvider, error)`
  - OTLP gRPC exporter + Prometheus `/metrics` endpoint (port `:9090` in dev)
  - Same resource attributes as tracer

  **2d. PostgreSQL pool** (`backend/internal/platform/postgres/connection.go`)
  - `Connect(ctx, databaseURL string) (*pgxpool.Pool, error)`
  - Parse config from `DATABASE_URL`, set max conns = 20, min conns = 2
  - Health check via `pool.Ping(ctx)`
  - `Close()` helper for graceful shutdown

  **2e. Platform wire** (`backend/internal/platform/wire.go`)
  - Update `InitPlatform` to call real constructors (no-op → real)
  - Return cleanup functions for graceful shutdown
  - Add `Shutdown(ctx)` that calls Sync on logger and Shutdown on TracerProvider/MeterProvider

  **Files:** `backend/internal/platform/logger/logger.go`, `backend/internal/platform/telemetry/tracer.go`, `backend/internal/platform/telemetry/metrics.go`, `backend/internal/platform/postgres/connection.go`, `backend/internal/platform/wire.go`
  **Logging:** INFO — log "platform initialized" with service name and env; log each subsystem status (logger/OTel/pgx ready); WARN on exporter connection failure (non-fatal in dev).

---

### Phase 2: CLI & HTTP Server

<!-- Commit checkpoint: T1–T2 -->

- [ ] **T3: Wire cobra CLI tree**

  Replace the current `switch` in `main.go` with a cobra command tree. CLI is an input adapter over the Application layer (per ADR-DES.API.cli-interface).

  **3a. Root command** (`backend/internal/cli/cli.go`)
  - `NewRoot() *cobra.Command`
  - Persistent pre-run: init zap logger (INFO), set `log.Printf` → zap bridge
  - Persistent flags: `--output json|table|csv` (default: `table`), `--yes` (skip confirmations)
  - Subcommands registered via `AddCommand`

  **3b. Subcommand stubs** (9 commands, all in `backend/internal/cli/`)
  - `server` → starts HTTP server (T4), production-ready
  - `mcp` → stub: `fmt.Println("MCP server not yet implemented")`
  - `migrate` → delegates to Atlas CLI or runs embedded migration engine. Subcommands: `up`, `down`, `validate` — all stubs printing "not yet implemented"
  - `seed` → inserts RBAC role catalog + demo data (wired in T6)
  - `ontology sync` → stub
  - `route compute` → stub with `--stub` flag (uses in-memory graph later, T9)
  - `plan get` → stub
  - `gap diagnose` → stub
  - `report` → stub

  **3c. Entry point** (`backend/cmd/vedo-edutrack/main.go`)
  - Replace `switch args[0]` with `cli.Execute()`
  - Bare `vedo-edutrack` → help text (cobra default)
  - `vedo-edutrack version` → ldflags-injected version (keep existing logic)
  - `vedo-edutrack health` → self-probe (keep existing logic, move to cobra)

  **Files:** `backend/internal/cli/cli.go` (rewrite from doc-only to real), `backend/internal/cli/server.go`, `backend/internal/cli/mcp.go`, `backend/internal/cli/migrate.go`, `backend/internal/cli/seed.go`, `backend/internal/cli/ontology_sync.go`, `backend/internal/cli/route_compute.go`, `backend/internal/cli/plan_get.go`, `backend/internal/cli/gap_diagnose.go`, `backend/internal/cli/report.go`, `backend/cmd/vedo-edutrack/main.go`
  **Logging:** INFO — log subcommand invocation with args; log "platform initialized" from persistent pre-run; WARN on unknown subcommand (cobra handles, but zap bridge should capture).

- [ ] **T4: Upgrade server to chi router with graceful shutdown**

  Replace `http.ServeMux` with `chi.NewRouter()`. The server becomes the production HTTP frontend serving API + metrics + health + (later) embedded SPA.

  **4a. Router setup** (`backend/cmd/vedo-edutrack/server.go`)
  - `chi.NewRouter()` with middleware stack:
    - `middleware.RequestID` — injects request ID into context and response header
    - `middleware.RealIP` — trust Traefik X-Forwarded-For
    - `middleware.Logger` — chi built-in (structured, request/response summary)
    - `middleware.Recoverer` — panic recovery with stack trace
    - `middleware.Timeout(30 * time.Second)` — global timeout
  - Route groups:
    - `GET /healthz` → liveness (always 200)
    - `GET /readyz` → readiness (DB + JWKS, see 4b)
    - `GET /metrics` → Prometheus handler
    - `/.well-known/jwks.json` → JWKS endpoint (T5)
    - `/api/v1/` → API group (mounted later, T8–T9)
    - `/*` → SPA fallback (T17)
  - CORS: `cors.Handler` allowing `http://localhost:5173` (dev), `*` in dev mode
  - Graceful shutdown: trap `SIGINT`/`SIGTERM`, 30s timeout, close DB pool, flush OTel, sync logger
  - `http.Server` with `ReadTimeout`, `WriteTimeout`, `IdleTimeout` (all 30s)

  **4b. Readiness with JWKS** (`backend/cmd/vedo-edutrack/server.go` readyz handler)
  - Extend `/readyz` JSON response with `checks.identity_provider` field
  - TCP dial check to JWKS URL host (similar to existing DB check pattern)
  - New env: `JWKS_URL` (default: `http://localhost:8080/.well-known/jwks.json` for dev)
  - Status: `up` if reachable, `down` with error message otherwise
  - Overall status: `ok` only if all checks pass

  **Files:** `backend/cmd/vedo-edutrack/server.go` (rewrite), `backend/cmd/vedo-edutrack/health.go` (keep, move to cobra subcommand in T3)
  **Logging:** INFO — log "server starting" with port, "server stopping" with reason; log each middleware init; WARN on readyz check failure with reason.

---

### Phase 3: Authentication & RBAC

<!-- Commit checkpoint: T3–T4 -->

- [ ] **T5: Implement JWT authentication middleware**

  Implement local JWT auth with self-signed RS256 key pair for dev environments. Keycloak is post-MVP (ADR-DES.SECURITY.rbac-model); local RS256/JWKS is the MVP path.

  **5a. Key pair generation** (`backend/internal/platform/auth/keys.go`)
  - `GenerateKeyPair()` → `*rsa.PrivateKey` (2048-bit)
  - `MarshalJWKS(publicKey)` → JWKS JSON
  - Auto-generate on startup if no key file exists
  - `JWKS_PRIVATE_KEY_PATH` env (default: `./tmp/jwt_key.pem` — gitignored)

  **5b. JWKS endpoint**
  - Mount `GET /.well-known/jwks.json` on chi router
  - Returns the public key as JWKS (kid, kty=RSA, n, e, alg=RS256)

  **5c. Token issuance** (`backend/internal/platform/auth/token.go`)
  - `POST /api/v1/auth/token` — dev-only endpoint
  - Accepts `{"user_id": "...", "roles": ["learner", ...]}`
  - Returns signed JWT with claims: `sub`, `iat`, `exp` (24h), `roles`
  - `--dev-mode` flag or `ENV=development` guard — WARN if used in production

  **5d. Auth middleware** (`backend/internal/platform/auth/middleware.go`)
  - Chi middleware: extracts `Authorization: Bearer <token>`
  - Validates via JWKS (local or remote `JWKS_URL`), RS256
  - Validates `iss`, `aud`, `exp` claims
  - Injects into `context.Context`: `UserID`, `Roles`, `TokenClaims`
  - 401 on missing/invalid token (JSON body with error code)
  - Skips validation for `/healthz`, `/readyz`, `/.well-known/jwks.json`

  **5e. Auth helpers** (`backend/internal/platform/auth/context.go`)
  - `GetUserID(ctx)` → `string`
  - `GetRoles(ctx)` → `[]string`

  **Files:** `backend/internal/platform/auth/keys.go` (new), `backend/internal/platform/auth/token.go` (new), `backend/internal/platform/auth/middleware.go` (new), `backend/internal/platform/auth/context.go` (new)
  **Logging:** INFO — log token issuance with user_id (no secret); log token validation result; WARN on invalid/expired tokens; WARN if dev-mode token endpoint used with ENV=production.

- [ ] **T6: Implement RBAC engine and seed command**

  Implement a deny-by-default RBAC engine matching the role catalog and permission matrix from the RBAC ADR. This is the server-side enforcement layer; UI role visibility is cosmetic (UX only).

  **6a. RBAC engine** (`backend/internal/platform/auth/rbac.go`)
  - Archetype enum: `self`, `dependents_owner`, `staff`, `management`, `integration`, `admin`, `ops`
  - Permission enum per functional area × action: `route:read`, `route:write`, `plan:read`, `plan:manage`, `progress:read`, `gap:read`, `coverage:read`, `ontology:read`, `resource:read`, `resource:manage`, `user:read`, `user:manage`, `webhook:configure`
  - Permission matrix lookup table (hardcoded from REQ-NFR-security.compliance.permission-matrix.md)
  - `HasPermission(roles []string, perm Permission) bool` — O(1) map lookup
  - `RequireRole(perm Permission)` → chi middleware: reads roles from context (T5), calls HasPermission, 403 on deny

  **6b. Seed command** (`backend/internal/cli/seed.go`)
  - Replace stub with real implementation
  - Inserts role catalog: 10 role instances mapping to 7 archetypes
  - Inserts permission matrix rows
  - Creates default admin user (`admin@edutrack.local`) with `admin` role
  - Idempotent: `ON CONFLICT DO NOTHING` for re-runs
  - Table schema: `identity_access.roles(id, name, archetype, description)`, `identity_access.role_permissions(role_id, permission)`

  **6c. Migration** (`backend/migrations/000001_identity_access_init.sql`)
  - First Atlas migration: creates `identity_access` schema with `roles` and `role_permissions` tables
  - Minimal columns needed for M0.3 RBAC seed

  **Files:** `backend/internal/platform/auth/rbac.go` (new), `backend/internal/cli/seed.go` (rewrite), `backend/migrations/000001_identity_access_init.sql` (new), `backend/migrations/atlas.sum` (new)
  **Logging:** INFO — log seed start/completion with role count; INFO — log permission check results in RequireRole middleware (user_id, permission, allowed/denied).

---

### Phase 4: API Contract & Domain Stubs

<!-- Commit checkpoint: T5–T6 -->

- [ ] **T7: Expand OpenAPI specification**

  Expand `backend/api/openapi/v1.yaml` from the current minimal `/healthz` spec to include the M0.3 API surface.

  Add schemas and paths:
  - `POST /api/v1/routes/compute` — request: `{learner_id, goal_topic_id}`, response: `{route: [{topic_id, order, link_type}], metadata: {...}}`
  - `GET /api/v1/ontology/concepts` — query: `?topic_id`, response: `{concept: {id, title, description, links: [...]}}`
  - `POST /api/v1/auth/token` — request: `{user_id, roles}`, response: `{access_token, token_type, expires_in}`
  - `GET /api/v1/me` — response: `{user_id, roles}`
  - `GET /readyz` — response: `{status, version, env, uptime, checks: {database, identity_provider}}`
  - `GET /.well-known/jwks.json` — JWKS response
  - Security scheme: `bearerAuth` (Bearer JWT)
  - Apply security globally, with exceptions for `/healthz`, `/readyz`, `/.well-known/jwks.json`, `/api/v1/auth/token`

  **Files:** `backend/api/openapi/v1.yaml`
  **Logging:** N/A (spec file change only).

- [ ] **T8: Run oapi-codegen and scaffold server interface**

  Generate Go types and server interface from the expanded OpenAPI spec.

  **8a. Codegen config** (`backend/api/openapi/codegen.yaml`)
  - `oapi-codegen` config for chi server generation
  - Output: `backend/internal/api/server.gen.go` (server interface), `backend/internal/api/types.gen.go` (request/response types)

  **8b. Generated code** (`backend/internal/api/`)
  - Run `oapi-codegen -config codegen.yaml v1.yaml`
  - Generated `ServerInterface` with method stubs for each endpoint

  **8c. Stub handler** (`backend/internal/api/handler.go`)
  - `StubHandler` struct implementing `ServerInterface`
  - All methods return 501 Not Implemented with JSON `{"error": "not_implemented", "endpoint": "..."}`
  - Mounted on chi router at `/api/v1/`

  **8d. Makefile gen target** (`Makefile`)
  - Wire `make gen`: runs `oapi-codegen` + (future) `sqlc generate`

  **Files:** `backend/api/openapi/codegen.yaml` (new), `backend/internal/api/server.gen.go` (generated), `backend/internal/api/types.gen.go` (generated), `backend/internal/api/handler.go` (new), `Makefile` (update `gen` target)
  **Logging:** INFO — log each generated file path; INFO — log "API handler mounted" on server start.

- [ ] **T9: Implement ontology stub and route computation stub**

  Replace the two critical stub endpoints with working in-memory implementations. These are hardcoded demo stubs — real VEDO Hub integration comes in M1.

  **9a. Ontology stub** (`backend/internal/modules/ontologyport/adapters/stub/graph.go`)
  - Hardcoded small module graph (10–15 topics covering math 5th grade)
  - 5 link types represented: `hasStrictPrerequisite`, `hasSoftPrerequisite`, `enriches`, `appliesTo`, `isAnalogousTo`
  - `GetConcept(topicID)` → concept with title, description, linked topics
  - `GetAllConcepts()` → list of all concepts
  - Wire into `GET /api/v1/ontology/concepts` handler

  **9b. Route stub** (`backend/internal/modules/routeplanning/adapters/stub/computer.go`)
  - `ComputeRoute(learnerID, goalTopicID)` → fixed pre-computed route
  - Uses ontology stub graph to walk `hasStrictPrerequisite` links from goal back to root
  - Returns ordered list of topics with link type
  - Wire into `POST /api/v1/routes/compute` handler

  **9c. Handler wiring** (`backend/internal/api/`)
  - Replace stub `StubHandler` methods for `PostRoutesCompute` and `GetOntologyConcepts` with real implementations
  - Other endpoints remain 501 stubs

  **9d. Makefile** (`Makefile`)
  - Add `route compute --stub` to cobra: runs `ComputeRoute` from CLI and prints result (table format)

  **Files:** `backend/internal/modules/ontologyport/adapters/stub/graph.go` (new), `backend/internal/modules/routeplanning/adapters/stub/computer.go` (new), `backend/internal/api/handler.go` (update — replace stubs)
  **Logging:** INFO — log ontology query with topic_id; INFO — log route computation with learner_id, goal_topic_id, and result length; WARN if topic not found in stub graph.

---

### Phase 5: Frontend Foundation

<!-- Commit checkpoint: T7–T9 -->

- [ ] **T10: Add frontend dependencies (routing, auth)**

  Add npm packages required for M0.3 UI scaffolding.

  ```bash
  pnpm --filter frontend add react-router-dom@7
  pnpm --filter frontend add jose   # JWT RS256 verification (lightweight, no Node.js crypto deps)
  ```

  `jose` is chosen over `jwx` (Go lib) or `jsonwebtoken` — it works in the browser, supports JWKS, and is maintained.

  **Files:** `frontend/package.json`, `pnpm-lock.yaml`
  **Logging:** N/A (dependency installation).

- [ ] **T11: Implement auth context and API client**

  Build the frontend auth layer: JWT storage, silent refresh (stub), and a typed fetch wrapper.

  **11a. Auth store** (`frontend/src/store/authStore.ts`)
  - Zustand store with `persist` middleware (sessionStorage — XSS-resistant)
  - State: `token`, `user`, `roles`, `isAuthenticated`
  - Actions: `login(token, user)`, `logout()`, `refreshToken()` (stub — returns same token)
  - `useAuth()` hook: `{ user, roles, isAuthenticated, login, logout }`

  **11b. Auth provider** (`frontend/src/features/identity-access/AuthProvider.tsx`)
  - React context provider wrapping the Zustand store
  - On mount: checks for stored token, validates expiry (JWT `exp` claim via `jose`)
  - Exposes `useAuth()` hook

  **11c. API client** (`frontend/src/shared/api/client.ts`)
  - Thin `fetch` wrapper:
    - Base URL from `appConfig.apiBaseUrl`
    - Injects `Authorization: Bearer <token>` from auth store
    - JSON request/response serialization
    - Typed error handling: `ApiError` class with status code + message
    - `api.get<T>(path)`, `api.post<T>(path, body)` helpers
  - Auto-redirect to `/login` on 401

  **11d. Web OTel** (`frontend/src/shared/telemetry/`)
  - Initialize `@opentelemetry/sdk-trace-web` with OTLP exporter
  - Export to collector (same endpoint as backend)
  - Configuration via `VITE_OTEL_EXPORTER_OTLP_ENDPOINT`
  - Resource attributes: `service.name=vedo-edutrack-web`
  - **Note:** the `@opentelemetry/*` web packages are not yet added to package.json; add them in this task:
    ```bash
    pnpm --filter frontend add @opentelemetry/sdk-trace-web @opentelemetry/exporter-trace-otlp-http @opentelemetry/resources @opentelemetry/semantic-conventions
    ```

  **Files:** `frontend/src/store/authStore.ts` (new), `frontend/src/features/identity-access/AuthProvider.tsx` (new), `frontend/src/features/identity-access/index.ts` (new), `frontend/src/shared/api/client.ts` (new), `frontend/src/shared/api/index.ts` (new), `frontend/src/shared/telemetry/tracer.ts` (new), `frontend/src/shared/telemetry/index.ts` (new), `frontend/package.json`
  **Logging:** INFO (frontend console) — log auth state changes (login/logout/token refresh); log API errors with status code and endpoint; WARN on token expiry.

- [ ] **T12: Set up client-side routing**

  Define the SPA route tree and mount it in `App.tsx`.

  **12a. Route definition** (`frontend/src/routes.tsx`)
  - `react-router-dom` v7 with JSX route config
  - Route tree:
    ```
    /                   → LandingPage (public)
    /login              → LoginPage (public)
    /dashboard          → DashboardLayout (protected, role-aware)
    /dashboard/route    → RouteView (protected)
    /dashboard/plan     → PlanView (protected)
    /dashboard/progress → ProgressView (protected)
    /*                  → NotFound
    ```
  - Lazy-loaded feature pages (`React.lazy` + `Suspense`)

  **12b. App root** (`frontend/src/App.tsx`)
  - Replace minimal landing banner with `<RouterProvider>` + `<AuthProvider>`
  - Wrap in `<ErrorBoundary>` (catch-all React error UI)
  - `<Suspense fallback={<LoadingSpinner />}>` for lazy pages

  **12c. Shared loading/error** (`frontend/src/shared/components/`)
  - `LoadingSpinner.tsx` — centered spinner with "Loading…" label
  - `ErrorFallback.tsx` — "Something went wrong" with retry button

  **Files:** `frontend/src/routes.tsx` (new), `frontend/src/App.tsx` (rewrite), `frontend/src/shared/components/LoadingSpinner.tsx` (new), `frontend/src/shared/components/ErrorFallback.tsx` (new), `frontend/src/shared/components/index.ts` (new)
  **Logging:** INFO (frontend console) — log route transitions (pathname); ERROR — log React render errors caught by ErrorBoundary.

---

### Phase 6: Frontend UI Pages

<!-- Commit checkpoint: T10–T12 -->

- [ ] **T13: Build shared component primitives and layout shell**

  Create a minimal design system of reusable components — no UI library, pure Tailwind CSS v4.

  **13a. Component primitives** (`frontend/src/shared/components/`)
  - `Button` — variants: `primary`, `secondary`, `ghost`; sizes: `sm`, `md`, `lg`; loading state
  - `Card` — padding, shadow, rounded corners, optional header/footer slots
  - `Input` — label, error message, disabled state
  - `Badge` — color variants (success/warning/danger/info)
  - `Avatar` — initials fallback when no image

  **13b. Layout shell** (`frontend/src/shared/layouts/`)
  - `MainLayout.tsx`:
    - Sidebar: collapsible, navigation links filtered by user roles
    - Header: app logo, user avatar + name, logout button
    - Content: `<Outlet />` from react-router
    - Responsive: sidebar collapses to icon-only on < 768px
  - `AuthLayout.tsx`: centered card layout for login/register pages
  - `LandingLayout.tsx`: full-width no-sidebar layout for landing page

  **13c. Role-based navigation** (`frontend/src/shared/guards/`)
  - `ProtectedRoute.tsx` — redirects to `/login` if unauthenticated
  - `RoleGate.tsx` — renders children only if user has required role; shows 403 message otherwise

  **Files:** `frontend/src/shared/components/Button.tsx` (new), `frontend/src/shared/components/Card.tsx` (new), `frontend/src/shared/components/Input.tsx` (new), `frontend/src/shared/components/Badge.tsx` (new), `frontend/src/shared/components/Avatar.tsx` (new), `frontend/src/shared/layouts/MainLayout.tsx` (new), `frontend/src/shared/layouts/AuthLayout.tsx` (new), `frontend/src/shared/layouts/LandingLayout.tsx` (new), `frontend/src/shared/guards/ProtectedRoute.tsx` (new), `frontend/src/shared/guards/RoleGate.tsx` (new)
  **Logging:** INFO — log layout mount; WARN — log access denied in RoleGate with required role.

- [ ] **T14: Build landing page**

  Public unauthenticated page (`/`) communicating the product value proposition.

  **14a. Landing page** (`frontend/src/pages/Landing.tsx`)
  - Hero section: "VEDO EduTrack" title, tagline: "Персональная образовательная траектория на основе графа знаний" (or English equivalent)
  - Value propositions (3 cards):
    - "Knowledge Graph Routes" — personalized learning paths
    - "Progress Tracking" — plan vs actual comparison
    - "Gap Diagnosis" — find root causes of learning gaps
  - CTA buttons: "Sign In" → `/login`, "Learn More" → anchor to features section
  - Footer: copyright, links

  **14b. Wire in routes**

  **Files:** `frontend/src/pages/Landing.tsx` (new), `frontend/src/pages/index.ts` (new)
  **Logging:** N/A (static page).

- [ ] **T15: Build login page**

  Simple login form for dev/demo auth flow.

  **15a. Login page** (`frontend/src/pages/Login.tsx`)
  - Form: user ID input + role selector dropdown
  - Calls `POST /api/v1/auth/token` via API client
  - On success: stores token via `authStore.login()`, redirects to `/dashboard`
  - On error: displays inline error message
  - Loading state on submit button
  - Note: this is a dev login flow (user_id + role → JWT). Real login with credentials comes post-MVP.

  **15b. Token types** (`frontend/src/shared/api/types.ts`)
  - `TokenResponse`: `{access_token, token_type, expires_in}`
  - `LoginRequest`: `{user_id, roles}`
  - `UserInfo`: `{user_id, roles}`

  **Files:** `frontend/src/pages/Login.tsx` (new), `frontend/src/shared/api/types.ts` (new)
  **Logging:** INFO — log login attempt (user_id only, no secrets); WARN — log login failure with error.

- [ ] **T16: Build role-aware dashboard shells**

  Create dashboard placeholder pages for each MVP persona. All protected behind `<ProtectedRoute>` and `<RoleGate>`.

  **16a. Dashboard layout** (`frontend/src/pages/Dashboard.tsx`)
  - Role-aware routing: renders different content based on user's primary role
  - Welcomes user by name: "Welcome, {user_id}"
  - Quick stats cards (stub data): "Active Plans", "Completed Modules", "Coverage %"

  **16b. Feature page placeholders** (one per bounded context, under `/dashboard/*`)
  - `RouteView.tsx` — "Route Builder" placeholder with coming-soon message
  - `PlanView.tsx` — "My Plan" placeholder
  - `ProgressView.tsx` — "Progress" placeholder
  - Each page wrapped in `<RoleGate requiredRole="self">` (learner default)

  **16c. Role variants**
  - Learner dashboard (role `learner`): my route, plan, progress, gaps
  - Parent dashboard (role `parent`): children's progress overview
  - Teacher dashboard (role `teacher`): class/group management placeholder
  - Methodologist dashboard (role `methodologist`): FGOS coverage placeholder
  - Admin dashboard (role `admin`): user management placeholder

  **16d. NotFound page** (`frontend/src/pages/NotFound.tsx`)
  - 404 message with link back to `/` or `/dashboard`

  **Files:** `frontend/src/pages/Dashboard.tsx` (new), `frontend/src/pages/RouteView.tsx` (new), `frontend/src/pages/PlanView.tsx` (new), `frontend/src/pages/ProgressView.tsx` (new), `frontend/src/pages/NotFound.tsx` (new)
  **Logging:** INFO — log dashboard mount with primary role; WARN — log unauthorized access attempts in RoleGate.

---

### Phase 7: SPA Embed & Docker Unification

<!-- Commit checkpoint: T13–T16 -->

- [ ] **T17: Embed SPA into backend binary**

  Merge the standalone `frontend/Dockerfile.embed` build stage into the backend binary. The result: a single `go build` produces one binary serving both API and SPA from one port.

  **17a. Embed filesystem** (`backend/internal/platform/spa/embed.go`)
  - `//go:embed all:frontend/dist` directive
  - `SPAFileSystem()` → `http.FileSystem` (subdirectory `frontend/dist`)
  - Graceful: if no embedded files (dev without build), serve a JSON message: `{"spa": "not embedded (dev mode)"}`
  - Build tag or env guard: embed only works after `make build-frontend`

  **17b. SPA fallback on chi router**
  - `/*` catch-all route after all API routes
  - `FileServer` from embedded filesystem
  - SPA fallback: serve `index.html` for unknown paths (react-router handles client-side)

  **17c. Frontend build integration**
  - `Makefile` target: `build-frontend` — runs `pnpm --filter frontend build`
  - Copy `frontend/dist/` to `backend/frontend/dist/` before `go build`
  - `.gitignore`: `backend/frontend/dist/` (generated)

  **Files:** `backend/internal/platform/spa/embed.go` (new), `backend/cmd/vedo-edutrack/server.go` (add SPA fallback route), `Makefile` (add `build-frontend` target, update `build` to include frontend build), `backend/.gitignore` (add frontend/dist/)
  **Logging:** INFO — log "SPA files embedded" with file count on server start; WARN — log "SPA not embedded (dev mode)" when running without build.

- [ ] **T18: Update Dockerfiles for multi-arch and unified binary**

  Update the backend Dockerfile to include the frontend embed build stage. Add multi-architecture support.

  **18a. Backend Dockerfile** (`backend/Dockerfile`)
  - Stage 1: `golang:1.26-alpine` — Go build with `CGO_ENABLED=0`
  - Stage 2: `node:24-alpine` — pnpm + Vite build (frontend)
  - Stage 3: copy frontend dist to backend embed path, run `go build`
  - Stage 4: `gcr.io/distroless/static:nonroot` — final image
  - Multi-arch: `FROM --platform=$BUILDPLATFORM` for build stages, `FROM --platform=$TARGETPLATFORM` for runtime
  - `ARG TARGETARCH` in build args for Go compiler (`GOARCH=$TARGETARCH`)
  - HEALTHCHECK: `vedo-edutrack health`
  - Labels: OCI annotations (`org.opencontainers.image.*`)

  **18b. Standalone Dockerfiles consolidated**
  - Retire `frontend/Dockerfile.embed` (or convert to comment-only doc reference)
  - Keep `frontend/Dockerfile.nginx` for SaaS/CDN deployment (post-MVP use case)
  - Update `deploy/README.md` container strategy

  **18c. Makefile docker targets** (`Makefile`)
  - `docker-build`: `docker buildx build --platform linux/amd64,linux/arm64 -f backend/Dockerfile -t vedo-edutrack:latest .`
  - `docker-build-local`: single-arch for local dev (`--load`)
  - `docker-build-all` already exists — update to use unified Dockerfile

  **Files:** `backend/Dockerfile` (rewrite), `deploy/README.md` (update container strategy), `Makefile` (update docker targets)
  **Logging:** N/A (Docker/CI).

- [ ] **T19: Update docker-compose and Traefik for unified binary**

  Adapt the dev stack to use the unified backend binary (API + SPA on one port).

  **19a. docker-compose.yml** (`deploy/docker-compose.yml`)
  - Replace separate `frontend` service with SPA served from backend
  - Uncomment and enable `spa-embed` profile or update `backend` service notes
  - Add `JWKS_URL`, `JWT_AUDIENCE`, `JWT_ISSUER` env vars to backend
  - Add `DATABASE_URL` with credentials to backend (was already there)
  - Keep `frontend` Vite dev server as optional dev profile (`profiles: [frontend-dev]`) for hot-reload workflow

  **19b. Traefik dynamic config** (`deploy/traefik/dynamic.yml`)
  - Update `spa` router to point to `backend:8080` (unified binary)
  - Add `api` router for backend (already exists: `PathPrefix(/api)`)
  - Remove or comment out separate frontend router

  **Files:** `deploy/docker-compose.yml`, `deploy/traefik/dynamic.yml`
  **Logging:** INFO — log compose service discovery on startup.

---

### Phase 8: Mock VEDO Hub (GraphQL test container)

<!-- Commit checkpoint: T17–T19 (Phase 7); Phase 8 commits after T25 -->

- [ ] **T20: Mock-hub design docs — ADR, C4 deployment diagrams, contract channel update**

  Materialize the design decision (RESEARCH session 2026-08-03) into architecture artifacts before implementation.

  **20a. ADR** (`specs/adr/ADR-DES.INFRA.mock-hub-strategy.md`, new, Russian per `specs/adr/README.md`)
  - Status: ПРЕДЛОЖЕНО · Date: 2026-08-03
  - Context: no Hub in the dev stack (`VEDO_HUB_API_URL=http://localhost:8081` points nowhere; in-container `localhost` ≠ host); contract tests of the Hub boundary (M1) and CI need a controllable stand-in; the Hub's ontology-service exposes a read-only GraphQL interface purpose-built for graph navigation (`graphNeighborhood`, `classDescendants`), while the REST API has no traversal endpoint
  - Source requirements: `REQ-FR-api.hub.read-ontology` (F0.1), `REQ-FR-api.hub.copy-subgraph` (F0.2), `REQ-NFR-api.availability.hub-dependency-sla`, `REQ-NFR-process.dev.test-coverage`, `ADR-DES.API.communication-patterns` §6
  - Decision: separate Docker container (`hub-mock`), same Go module, GraphQL-only (`POST /graphql`), ontology in memory from an arbitrary `.ttl` (`ONTOLOGY_FILE`; first seed `traceability.ttl`), `gqlparser/v2` resolvers, the same handler exposed as `httptest.Server`; documented exception to the single-binary ADR (test-only, not shipped)
  - Alternatives (brief): REST-only mock, WireMock/Prism/MockServer, real triplestore (Fuseki/GraphDB), gqlgen, hand-rolled GraphQL executor, in-process fake only
  - Consequences + mitigations (dependency pinning, contract-spec sync, second `cmd/` entry)

  **20b. Stack ADR cross-reference** (`specs/adr/ADR-DES.STACK.framework-vs-vs.md`)
  - Add `ADR-DES.INFRA.mock-hub-strategy` to «Связанные артефакты»; note `gqlparser/v2` as the only new (test-only) dependency

  **20c. C4 deployment diagrams**
  - `specs/c4/deployment-dev.md` (update): add `hub-mock` node in `edutrack-net` (:8081, healthcheck `/healthz`, volume `../traceability.ttl` → `/data/ontology.ttl`), `Rel(api, hub-mock, "Read ontology (F0), GraphQL", "HTTP :8081")`, legend/context note (`VEDO_HUB_API_URL=http://hub-mock:8081`)
  - `specs/c4/deployment-test.md` (new): CI (GitHub Actions) — runner + service containers `postgres` + `hub-mock`; contract/e2e tests hit `hub-mock`

  **20d. Contract channel update (GraphQL as F0 channel)**
  - `specs/requirements/REQ-FR-api.hub.read-ontology.md` — add GraphQL to the channel list
  - `specs/adr/ADR-DES.API.communication-patterns.md` §6 — add GraphQL to the ontology-port boundary
  - `traceability.ttl` — add instances: `tr:adr-des-infra-mock-hub-strategy`, `tr:c4-deployment-test`, TEST (contract tests); keep 0 orphans

  **Files:** `specs/adr/ADR-DES.INFRA.mock-hub-strategy.md` (new), `specs/adr/ADR-DES.STACK.framework-vs-vs.md`, `specs/c4/deployment-dev.md`, `specs/c4/deployment-test.md` (new), `specs/requirements/REQ-FR-api.hub.read-ontology.md`, `specs/adr/ADR-DES.API.communication-patterns.md`, `traceability.ttl`
  **Logging:** INFO — log each ADR/C4/requirement artifact written.
  **Dependencies:** T19 (after Docker unification; design anchor for the phase)

- [ ] **T21: Implement ontology loader — mini-Turtle parser + in-memory model**

  `backend/internal/testing/mockhub/` (test-tooling package; `internal/` scopes it to the backend module; not part of the product binary).

  **21a. Parser** (`parser.go`)
  - Turtle subset: prefixes, `s a owl:Class | owl:ObjectProperty | owl:DatatypeProperty`, `rdfs:label` / `rdfs:comment` / `rdfs:subClassOf` / `rdfs:domain` / `rdfs:range`, `owl:FunctionalProperty`; ignore comments and blank nodes
  - 0 new dependencies (hand-rolled for the known TBox shape; `knakk/rdf` as fallback if the shape grows)
  - `Parse(reader io.Reader) (*Ontology, error)` — errors with line numbers

  **21b. Model** (`ontology.go`)
  - `Ontology{Classes map[string]*Class, Properties map[string]*Property}`; `Class{ID, Label, Comment, Parents, Children, IsAbstract, IsDeprecated}`; `Property{ID, Label, Comment, Type, Domains, Ranges, XSDType, Functional}`
  - Helpers: `Counts()`, `Breadcrumb(classID)`, `Tree()`, `Descendants(classID, maxDepth)`, `Autocomplete(q, limit)`
  - Loaded once at startup from `ONTOLOGY_FILE` (default: `traceability.ttl`); arbitrary `.ttl` supported; in-memory only (no persistence)

  **Files:** `backend/internal/testing/mockhub/parser.go`, `backend/internal/testing/mockhub/ontology.go`
  **Logging:** INFO — log ontology loaded: `{file, classes, properties}`.
  **Dependencies:** T1 (gqlparser/v2, used by T22), T20 (design anchor)

- [ ] **T22: Implement GraphQL server — gqlparser resolvers + HTTP handler**

  **22a. Handler** (`handler.go`)
  - `NewHandler(ont *Ontology, schemaSDL string) http.Handler` — routes: `POST /graphql` (gqlparser parse + execute against QueryRoot resolvers), `GET /healthz` (200)
  - Auth: any `Authorization: Bearer <token>` accepted; missing token → GraphQL error
  - Resolvers: `classes`, `class`, `classTree`, `classAncestors`, `classDescendants`, `properties`, `property`, `individuals`, `individual`, `graphNeighborhood`, `autocompleteClasses`, `_service { sdl }` (vedo-hub `schema.graphql` SDL embedded as a string)
  - Pagination: `{items, total, page, perPage}` (default 20); TBox-only ontology → empty individuals/connections, `graphNeighborhood` returns the node without edges (data-driven)
  - Errors: standard GraphQL `errors` array

  **22b. Test server** (`server.go`)
  - `NewTestServer(t testing.TB, ttlPath string) *httptest.Server` — in-process server for Go integration/contract tests (M1); same handler as the container

  **Files:** `backend/internal/testing/mockhub/handler.go`, `backend/internal/testing/mockhub/server.go`, `backend/internal/testing/mockhub/schema.graphql` (copy of vedo-hub SDL)
  **Logging:** INFO — log each GraphQL operation (operationName); ERROR — log execution errors.
  **Dependencies:** T21

- [ ] **T23: Add cmd/mockhub entry point**

  `backend/cmd/mockhub/main.go` — standalone dev/test binary (documented exception to the single-binary ADR `ADR-DES.API.cli-interface`: test-only, not shipped, not in SBOM).

  - Flags/env: `PORT` (default 8081), `ONTOLOGY_FILE` (default `../traceability.ttl`)
  - Load ontology at startup, serve `NewHandler` with graceful shutdown (SIGINT/SIGTERM)
  - `GET /healthz` used by the container healthcheck

  **Files:** `backend/cmd/mockhub/main.go`
  **Logging:** INFO — log startup: `{port, ontology_file, classes, properties}`.
  **Dependencies:** T22

- [ ] **T24: Add Dockerfile.mockhub**

  `backend/Dockerfile.mockhub` — multi-stage image for the mock hub (parallel to `backend/Dockerfile`).

  - Stage 1: `golang:1.26-alpine` — build `./cmd/mockhub` with `CGO_ENABLED=0`, `-ldflags="-s -w"`
  - Stage 2: `gcr.io/distroless/static:nonroot` (or alpine for wget healthcheck) — copy binary, `USER nonroot`, expose 8081
  - `ONTOLOGY_FILE` default `/data/ontology.ttl` (mounted at runtime)
  - OCI labels (source, version, description)

  **Files:** `backend/Dockerfile.mockhub` (new)
  **Logging:** INFO — log build stages and final image size.
  **Dependencies:** T23

- [ ] **T25: Wire hub-mock into compose + CI + smoke verification**

  **25a. Dev stack** (`deploy/docker-compose.yml`)
  - Add `hub-mock` service: build `backend/Dockerfile.mockhub`, port `8081:8081`, volume `../traceability.ttl:/data/ontology.ttl:ro`, `ONTOLOGY_FILE=/data/ontology.ttl`, healthcheck `wget -q -O /dev/null http://localhost:8081/healthz`, network `edutrack-net`
  - Fix backend env default: `VEDO_HUB_API_URL: ${VEDO_HUB_API_URL:-http://hub-mock:8081}` (container-to-container; the previous `localhost:8081` default was dead in-container)

  **25b. CI** (`.github/workflows/ci.yml`)
  - Build the mockhub image in the `build` job (alongside the backend image)
  - Smoke step (test job): run the built mockhub image, `curl -X POST localhost:8081/graphql` with a `classes` query → 200 + valid JSON

  **25c. Smoke verification (local)**
  - `make up` → `hub-mock` healthy; `curl -X POST localhost:8081/graphql -d '{"query":"{ classes(ontologyId:\"traceability\", perPage: 5) { total items { id label } } }"}' -H "Authorization: Bearer test"` → JSON with ~22 classes

  **Files:** `deploy/docker-compose.yml`, `.github/workflows/ci.yml`
  **Logging:** INFO — log compose service up, smoke result, CI build gate result.
  **Dependencies:** T24

### Phase 9: Integration & Verification

<!-- Commit checkpoint: T17–T19 -->

- [ ] **T26: Extend health checks and verify reachability**

  Finalize the health/readiness endpoints to verify all M0.3 components.

  **20a. /readyz checks**
  - `database`: pgx pool ping (replaces TCP dial from M0.2)
  - `identity_provider`: HTTP GET to JWKS URL (TCP dial is enough for M0.3)
  - All checks must pass for overall `status: "ok"`

  **20b. Verify local environment**
  - `make up` → all services healthy (`docker compose ps` shows healthy)
  - `curl localhost:8080/healthz` → 200 + JSON
  - `curl localhost:8080/readyz` → 200 + all checks up
  - `curl localhost:8080/.well-known/jwks.json` → valid JWKS
  - `curl -X POST localhost:8080/api/v1/auth/token -d '{"user_id":"demo","roles":["learner"]}'` → JWT returned
  - `curl -H "Authorization: Bearer <token>" localhost:8080/api/v1/me` → user info
  - `curl localhost:8080/api/v1/ontology/concepts?topic_id=math-5-1` → concept JSON
  - `curl -X POST localhost:8080/api/v1/routes/compute -d '{"learner_id":"demo","goal_topic_id":"math-5-10"}'` → route array
  - `curl -X POST localhost:8081/graphql -H "Authorization: Bearer test" -d '{"query":"{ classes(ontologyId:\"traceability\", perPage: 5) { total } }"}'` → 200 + valid GraphQL JSON (hub-mock reachable, Phase 8)
  - Browser: `http://localhost:8080/` → landing page (SPA served)
  - Browser: `http://localhost:8080/login` → login page
  - Browser: login → dashboard with role-aware UI

  **Files:** `backend/cmd/vedo-edutrack/server.go` (readyz handler update)
  **Logging:** INFO — log each health check result in readyz; INFO — log verification results.

- [ ] **T27: Seed database and verify RBAC**

  Run the seed command and verify the RBAC engine works end-to-end.

  **21a. Run seed**
  - `vedo-edutrack migrate up` → creates identity_access schema
  - `vedo-edutrack seed` → inserts roles + permissions + admin user
  - Verify: `SELECT * FROM identity_access.roles` → 10+ rows

  **21b. Verify RBAC enforcement**
  - Request with no token → 401
  - Request with learner token to admin endpoint → 403
  - Request with admin token to admin endpoint → allowed
  - UI: role-aware navigation only shows allowed links

  **Files:** N/A (verification task)
  **Logging:** INFO — log each RBAC verification step with result.

- [ ] **T28: Run gate tiers and confirm green**

  Run the gate runner to verify all automated checks pass.

  **22a. Fast tier** (`make dev-check`)
  - `go build ./...` — compiles
  - `golangci-lint run` — 0 errors
  - `gofmt -l .` — 0 diff
  - `biome ci frontend/src` — 0 errors
  - `tsc --noEmit` (frontend) — 0 errors
  - `pnpm validate:mermaid` — pass
  - `gitleaks detect` — 0 leaks

  **22b. Delivery tier** (`make check`)
  - All fast gates + integration, docker build, oapi-codegen consistency, coverage

  **22c. Scaffold verification**
  - All route endpoints return valid JSON (health, auth, ontology, routes, SPA)
  - Landing page loads in browser
  - Auth flow: login → token → protected route → dashboard
  - RBAC: role-based access enforced on both API and UI

  **Files:** N/A (verification task)
  **Logging:** INFO — log gate tier results summary; ERROR — log any gate failure with details.

---

## Dependencies Between Tasks

```
T1 ──→ T2 ──→ T3 ──→ T4
                │
                ├──→ T5 ──→ T6
                │         │
                │         └──→ T7 ──→ T8 ──→ T9
                │
                └──→ T10 ──→ T11 ──→ T12
                                     │
                                     └──→ T13 ──→ T14
                                                  │
                                                  ├──→ T15
                                                  └──→ T16
                                                         │
                                                         └──→ T17 ──→ T18 ──→ T19
                                                                                │
                                                                                ├──→ T20 ──→ T21 ──→ T22 ──→ T23 ──→ T24 ──→ T25
                                                                                │                                                    │
                                                                                └───────────────────────────────────────────────────┴──→ T26 ──→ T27 ──→ T28
```

- **T1→T2**: platform stubs need Go deps in go.mod (incl. `gqlparser/v2` for the mock hub, Phase 8)
- **T2→T3**: CLI needs zap logger from platform
- **T2→T4**: server needs OTel tracer + graceful shutdown patterns from platform
- **T3→T5**: auth middleware mounts on chi router
- **T5→T6**: RBAC reads roles from auth context
- **T3→T7**: OpenAPI spec is standalone (no code dependency)
- **T7→T8**: codegen runs from spec
- **T8→T9**: ontology/route stubs implement generated server interface
- **T3→T10**: frontend deps are standalone
- **T10→T11**: auth context uses jose + zustand
- **T11→T12**: routing uses AuthProvider
- **T12→T13**: layout uses routing (Outlet)
- **T13→T14,T15,T16**: pages use shared components + layout + guards
- **T4,T16→T17**: SPA embed needs server (T4) and built frontend (T16)
- **T17→T18**: Dockerfile includes embed build stage
- **T18→T19**: compose update references new Dockerfile
- **T19→T20**: mock-hub phase starts after Docker unification (design anchor)
- **T20→T21**: ADR fixes the design; parser implements it
- **T21→T22**: GraphQL resolvers need the in-memory model
- **T22→T23**: cmd/mockhub needs the handler
- **T23→T24**: image needs the binary
- **T24→T25**: compose/CI need the image
- **T19,T25→T26**: health checks verify the composed stack incl. hub-mock
- **T6,T26→T27**: seed needs DB migration + RBAC engine
- **T27→T28**: gate pass requires all components working

---

## Architecture References

| ADR | Relevance |
|-----|-----------|
| `ADR-DES.STACK.language-vs-vs` | Go backend, TypeScript frontend |
| `ADR-DES.STACK.framework-vs-vs` | chi + oapi-codegen backend, React SPA frontend |
| `ADR-DES.DATA.storage-strategy` | PostgreSQL + sqlc + Atlas |
| `ADR-DES.API.communication-patterns` | REST API (OpenAPI-first), SPARQL endpoint |
| `ADR-DES.API.cli-interface` | Single binary with cobra subcommands, CLI as input adapter |
| `ADR-DES.SECURITY.rbac-model` | Two-layer RBAC: archetypes + permission matrix |
| `ADR-IMPL.PROCESS.repository-structure` | 10 bounded contexts, Clean Architecture |
| `ADR-IMPL.PROCESS.development-tooling` | Biome, golangci-lint, Lefthook |
| `ADR-DES.INFRA.mock-hub-strategy` | Mock VEDO Hub GraphQL container (Phase 8, T20) |
| `specs/c4/deployment-test.md` | C4 deployment — CI test environment (Phase 8, T20) |
| `specs/requirements/REQ-NFR-security.compliance.role-catalog.md` | Role archetype definitions |
| `specs/requirements/REQ-NFR-security.compliance.permission-matrix.md` | Permission matrix per archetype |
| `specs/ddd/context-map.md` | Bounded context relationships |

## Non-Functional Requirements Addressed

| NFR | How |
|-----|-----|
| `REQ-NFR-ops.observability.structured-logging` | zap JSON logger (T2) |
| `REQ-NFR-ops.observability.distributed-tracing` | OTel OTLP gRPC exporter (T2) |
| `REQ-NFR-ops.observability.golden-signals-dashboards` | Prometheus /metrics + OTel metrics (T2); dashboards deferred to M1 |
| `REQ-NFR-security.access.jwt-auth` | JWT RS256/JWKS middleware (T5) |
| `REQ-NFR-security.access.rbac` | RBAC engine + permission matrix (T6) |
| `REQ-NFR-process.dev.engineering-gates` | make dev-check / make check pass (T28) |
| `REQ-NFR-ux.i18n-readiness` | Landing page in RU + EN; ICU deferred to M1 |
| `REQ-NFR-api.availability.hub-dependency-sla` | Mock VEDO Hub container + contract-test stand-in (T20–T25) |
| `REQ-NFR-process.dev.test-coverage` | `httptest.Server` from mockhub — basis for M1 contract tests (T22) |
