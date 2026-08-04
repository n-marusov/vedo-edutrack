# Implementation Plan: M4 — Integration & Webhook Layer (F6)

Branch: none
Created: 2026-08-04

## Settings
- Testing: yes
- Logging: verbose
- Docs: yes

## Roadmap Linkage
Milestone: "M4: Integration & Webhook Layer (F6)"
Rationale: Establishes the MVP integration surface (REST API contracts, SPARQL endpoint, webhooks, MCP server) as the next milestone after M1 core infrastructure completion. M4 is the clear next step per ROADMAP network plan: M1 → M4 → M6.

## Research Context
Source: .ai-factory/RESEARCH.md (Active Summary) — not directly applicable; RESEARCH currently focuses on hub-mock strategy. M4 implementation builds on the hub-mock for VEDO Hub connectivity and the existing ontology-port adapter.

Decisions carried over:
- ADR-DES.API.communication-patterns: REST (OpenAPI-first) + async events (in-process bus + outbox + webhooks), SPARQL read-only, MCP read-oriented
- ADR-DES.API.cli-interface: Single binary with cobra subcommands; MCP runs as `vedo-edutrack mcp` over stdio
- ADR-DES.INFRA.modular-monolith-approach: In-process event bus for cascades, outbox for external delivery
- ADR-DES.DATA.storage-strategy: PostgreSQL for event_dedup and outbox tables

## Commit Plan
- **Commit 1** (after tasks 1-3): "feat(integrations): add domain models, migrations, and outbox infrastructure"
- **Commit 2** (after tasks 4-8): "feat(integrations): implement SPARQL endpoint and webhook subscription system"
- **Commit 3** (after tasks 9-11): "feat(integrations): wire domain events and implement MCP server"
- **Commit 4** (after tasks 12-15): "feat(integrations): complete API docs, sandbox, and validation tests"

## Tasks

### Phase 1: Domain & Infrastructure Foundation

- [x] **Task 1: Define domain models and PostgreSQL migrations for integrations context**
  Implement the domain layer for the `integrations` bounded context and create Atlas migrations.

  **Domain models** (`backend/internal/modules/integrations/domain/integrations.go`):
  - `WebhookSubscription`: subscriber URL, event types (module.mastered, plan.deviated, route.recalculated), signing secret, active/inactive status, created/updated timestamps
  - `WebhookDelivery`: delivery record with subscription ID, event ID, attempt count, status (pending/sent/failed/permanent_failure), last attempt timestamp, HTTP status, response body
  - `EventDedup`: (source, event_id) unique pair for idempotent webhook receive; per ADR-§3
  - `SPARQLQueryResult`: result set with head/variables, bindings, truncated flag
  - Value objects: `EventType` (enum), `DeliveryStatus` (enum), `SubscriptionID`

  **Migrations** (`backend/migrations/`):
  - `<timestamp>_integrations_webhook_subscriptions.sql`: webhook_subscriptions table
  - `<timestamp>_integrations_webhook_deliveries.sql`: webhook_deliveries table
  - `<timestamp>_integrations_event_dedup.sql`: event_dedup table with unique (source, event_id) constraint

  **Domain tests** (`backend/internal/modules/integrations/domain/domain_test.go`):
  - Replace T13 placeholder with real tests for subscription validation, event type enum, delivery state machine

  **Logging:**
  - `DEBUG [integrations.domain] SubscriptionCreated {subscriptionID, eventTypes, url}`
  - `DEBUG [integrations.domain] DeliveryRecorded {deliveryID, subscriptionID, eventID, attempt}`
  - `ERROR [integrations.domain] SubscriptionValidationFailed {reason}`

  **Files:** `backend/internal/modules/integrations/domain/integrations.go`, `backend/internal/modules/integrations/domain/domain_test.go`, `backend/migrations/<ts>_integrations_*.sql`

- [x] **Task 2: Implement PostgreSQL-backed outbox replacing in-memory prototype**
  Replace the in-memory outbox (`adapters/webhook/outbox.go`) with a PostgreSQL-backed implementation.

  **Outbox repository** (`backend/internal/modules/integrations/adapters/webhook/outbox.go`):
  - Keep existing `Event`, `EventType` types and `ValidateReadOnly` guard
  - Add `OutboxRepository` interface: `Enqueue(ctx, Event) error`, `Dequeue(ctx, limit) ([]Event, error)`, `MarkDelivered(ctx, eventID) error`, `MarkFailed(ctx, eventID, reason) error`
  - Implement `PostgresOutbox`: enqueue inserts into outbox table transactionally with business changes; dequeue polls pending events ordered by created_at; idempotent on duplicate event_id (unique constraint)
  - Worker: goroutine polling every 1s with configurable batch size (default 10), graceful shutdown via context

  **Outbox table** (in migration from Task 1):
  - `outbox_events`: id (UUID PK), event_type (text), payload (jsonb), created_at, delivered_at (nullable), retry_count (int, default 0), status (pending/failed/delivered)

  **Tests** (`backend/internal/modules/integrations/adapters/webhook/outbox_test.go`):
  - Extend existing tests: test PostgresOutbox with testcontainers; verify dedup on duplicate event_id; verify worker polls and delivers; verify retry backoff

  **Logging:**
  - `DEBUG [webhook.outbox] Enqueue {eventID, eventType}`
  - `DEBUG [webhook.outbox] Dequeue {batchSize, count}`
  - `INFO [webhook.outbox] Delivered {eventID, subscriptionID, attempt, httpStatus}`
  - `WARN [webhook.outbox] DeliveryFailed {eventID, subscriptionID, attempt, error, httpStatus}`
  - `ERROR [webhook.outbox] PermanentFailure {eventID, subscriptionID, attempts, reason}`

  **Files:** `backend/internal/modules/integrations/adapters/webhook/outbox.go`, `backend/internal/modules/integrations/adapters/webhook/outbox_test.go`

- [x] **Task 3: Implement rate limiter and circuit breaker infrastructure**
  Add shared infrastructure for API rate limiting and VEDO Hub circuit breaker.

  **Rate limiter** (`backend/internal/modules/integrations/adapters/` or `backend/internal/platform/`):
  - Token-bucket rate limiter per API key / user: configurable rate (e.g., 10 req/min for SPARQL, 100 req/min for REST)
  - Return `429 Too Many Requests` with `Retry-After` header
  - Middleware form factor: `chi` middleware that wraps the handler

  **Circuit breaker** (`backend/internal/platform/circuitbreaker/`):
  - Circuit breaker for VEDO Hub calls (SPARQL proxy): closed → half-open → open states
  - Configurable failure threshold (e.g., 5 failures → open), timeout (e.g., 30s → half-open), success threshold (e.g., 3 successes → closed)
  - Per ADR-§5: timeout ≤ 3s, retry with exponential backoff

  **Logging:**
  - `WARN [ratelimit] RateLimited {userID, endpoint, currentRate}`
  - `WARN [circuitbreaker] CircuitOpened {service, failures}`
  - `INFO [circuitbreaker] CircuitHalfOpen {service}`
  - `INFO [circuitbreaker] CircuitClosed {service}`

  **Files:** `backend/internal/platform/circuitbreaker/breaker.go`, `backend/internal/platform/ratelimit/limiter.go`, `backend/internal/platform/ratelimit/middleware.go`

<!-- Commit checkpoint: tasks 1-3 -->

### Phase 2: SPARQL Endpoint & Webhook System

- [x] **Task 4: Implement production SPARQL endpoint handler**
  Replace the existing stub `SparqlQuery` handler with a real implementation that proxies SPARQL queries to VEDO Hub.

  **Handler** (`backend/internal/modules/integrations/adapters/sparql/handler.go`):
  - Reuse existing `ValidateReadOnly` and `IsTooManyResults` guard functions
  - `SPARQLHandler`: accepts query via GET query param, validates read-only, forwards to VEDO Hub SPARQL endpoint via `ontology-port` adapter
  - Circuit breaker wrapping Hub calls (from Task 3)
  - Rate limiting per authenticated user (from Task 3)
  - Result truncation at 10 000 rows with `truncated: true`
  - Timeout: 30s execution → `504 Gateway Timeout`
  - Mutating queries → `403 Forbidden` (existing guard)
  - Missing auth → `401 Unauthorized`
  - Rate limit → `429 Too Many Requests` with `Retry-After`

  **Wire into existing router** (`backend/internal/cli/server_http.go`):
  - Replace the `SparqlQuery` stub on the StubHandler with real implementation
  - The `/api/v1/sparql` route is already defined in OpenAPI spec and registered via `api.HandlerWithOptions`

  **Logging:**
  - `DEBUG [sparql] QueryReceived {userID, queryLength}`
  - `INFO [sparql] QueryExecuted {userID, executionTime, resultRows, truncated}`
  - `WARN [sparql] QueryRejected {userID, reason: "mutation"|"rate_limit"|"timeout"}`
  - `ERROR [sparql] HubUnreachable {error}`

  **Files:** `backend/internal/modules/integrations/adapters/sparql/handler.go`, `backend/internal/api/handler.go` (update SparqlQuery call)

- [x] **Task 5: SPARQL endpoint contract and integration tests**
  Write comprehensive tests for the SPARQL endpoint.

  **Contract tests** (`backend/tests/contract/` or `backend/internal/modules/integrations/adapters/sparql/handler_test.go`):
  - Valid SELECT/ASK/DESCRIBE/CONSTRUCT queries → 200
  - Mutating queries (INSERT/DELETE/CREATE/DROP/LOAD/CLEAR) → 403
  - Empty query → 400
  - Missing auth → 401
  - Rate limit exceeded → 429
  - Result truncation > 10 000 rows → 200 with truncated flag
  - Timeout → 504
  - Hub down → 503 (circuit breaker open)

  **Integration tests** (`tests/integration/`):
  - End-to-end SPARQL flow with hub-mock container: query → validate response format (application/sparql-results+json)
  - Test with minimal ontology data from hub-mock seed

  **Logging:** per standard test practices — assert log output in contract tests

  **Files:** `backend/internal/modules/integrations/adapters/sparql/handler_test.go`, `tests/integration/sparql_test.go`

- [x] **Task 6: Implement webhook subscription domain and application layer**
  Implement the domain service and application layer for webhook subscription management.

  **Domain service** (`backend/internal/modules/integrations/domain/integrations.go` — extend):
  - `SubscriptionService`: `CreateSubscription`, `UpdateSubscription`, `DeleteSubscription`, `ListSubscriptions`, `GetSubscription`
  - Validation: URL format, event types must be from known set, signing secret min 32 chars
  - Business rules: max 10 subscriptions per tenant, deactivation on 5 consecutive delivery failures

  **Application commands** (`backend/internal/modules/integrations/application/commands/commands.go`):
  - `CreateWebhookSubscription`: validate input, create subscription via domain service, enqueue verification ping event
  - `UpdateWebhookSubscription`: update URL/events/secret, reset failure count
  - `DeleteWebhookSubscription`: soft-delete

  **Application queries** (`backend/internal/modules/integrations/application/queries/queries.go`):
  - `ListWebhookSubscriptions`: by tenant, with status filter
  - `GetWebhookDeliveryHistory`: for a subscription, paginated

  **Logging:**
  - `INFO [integrations.app] SubscriptionCreated {subscriptionID, url, events}`
  - `INFO [integrations.app] SubscriptionUpdated {subscriptionID, changes}`
  - `INFO [integrations.app] SubscriptionDeleted {subscriptionID}`
  - `WARN [integrations.app] SubscriptionDeactivated {subscriptionID, reason: "consecutive_failures"}`

  **Files:** `backend/internal/modules/integrations/domain/integrations.go`, `backend/internal/modules/integrations/application/commands/commands.go`, `backend/internal/modules/integrations/application/queries/queries.go`

- [x] **Task 7: Add webhook subscription endpoints to OpenAPI spec and implement handlers**
  Extend the OpenAPI spec with webhook subscription management endpoints and implement the HTTP handlers.

  **OpenAPI spec** (`backend/api/openapi/v1.yaml`):
  - `GET /api/v1/webhooks/subscriptions` — list subscriptions for authenticated tenant
  - `POST /api/v1/webhooks/subscriptions` — create subscription {url, events[], secret}
  - `GET /api/v1/webhooks/subscriptions/{id}` — get subscription details
  - `PUT /api/v1/webhooks/subscriptions/{id}` — update subscription
  - `DELETE /api/v1/webhooks/subscriptions/{id}` — delete subscription
  - `GET /api/v1/webhooks/subscriptions/{id}/deliveries` — delivery history
  - `POST /api/v1/webhooks/subscriptions/{id}/ping` — manual ping to verify endpoint
  - Schemas: `WebhookSubscription`, `WebhookSubscriptionCreate`, `WebhookDelivery`

  **Handler** (`backend/internal/api/handler.go`):
  - Wire new endpoints to application commands/queries from Task 6
  - Auth: require authenticated user (JWT middleware already in place)

  **Logging:**
  - `DEBUG [api.webhooks] ListSubscriptions {userID}`
  - `INFO [api.webhooks] CreateSubscription {userID, url, events}`

  **Files:** `backend/api/openapi/v1.yaml`, `backend/internal/api/handler.go`, `backend/internal/api/server.gen.go` (regenerated)

- [x] **Task 8: Implement webhook delivery worker with retry and HMAC signing**
  Build the outbox polling worker that delivers webhook events to subscriber URLs.

  **Worker** (`backend/internal/modules/integrations/adapters/webhook/outbox.go` — extend):
  - `DeliveryWorker`: goroutine started on server boot, stopped on graceful shutdown
  - Poll `outbox_events` every 1s for pending events, batch size 10
  - For each pending event, find active subscriptions matching the event type
  - Deliver: HTTP POST to subscriber URL with JSON payload `{event_id, event_type, timestamp, data}`
  - HMAC-SHA256 signature header: `X-Vedo-Signature: t=<timestamp>,v1=<hex(signature)>` (signed with subscription secret)
  - Retry: exponential backoff (1s, 2s, 4s, 8s, 16s), max 5 attempts
  - Permanent failure after 5 attempts → deactivate subscription
  - Dedup: skip events already delivered to this subscription (check webhook_deliveries table)

  **Worker lifecycle** (`backend/internal/cli/server_http.go`):
  - Start `DeliveryWorker` after database pool is ready
  - Stop on server shutdown signal

  **Logging:**
  - `DEBUG [webhook.worker] Polling {batchSize}`
  - `DEBUG [webhook.worker] Delivering {eventID, subscriptionID, url, attempt}`
  - `INFO [webhook.worker] Delivered {eventID, subscriptionID, httpStatus, attempt}`
  - `WARN [webhook.worker] DeliveryRetry {eventID, subscriptionID, attempt, backoff}`
  - `ERROR [webhook.worker] PermanentFailure {eventID, subscriptionID, attempts}`

  **Files:** `backend/internal/modules/integrations/adapters/webhook/outbox.go`, `backend/internal/cli/server_http.go`

<!-- Commit checkpoint: tasks 4-8 -->

### Phase 3: Event Integration & MCP Server

- [x] **Task 9: Wire domain events from other bounded contexts into the webhook outbox**
  Connect the in-process event bus from other modules to the integrations outbox so webhook subscribers receive events.

  **Event subscriber** (`backend/internal/modules/integrations/adapters/webhook/` or new file):
  - Subscribe to in-process events: `ModuleMastered` → outbox `module.mastered`, `PlanDeviationDetected` → outbox `plan.deviated`, `RouteRecalculated` → outbox `route.recalculated`, `StandardDeficitDetected` → outbox `standard.risk_detected`
  - Map domain events to webhook payload format per `specs/ddd/domain-events.md` webhook representations table
  - Enqueue into outbox transactionally (outbox write within the same DB transaction as the event)

  **Wire into bus** (check `backend/internal/modules/` for in-process bus registration):
  - Register integrations event subscriber during server startup
  - Handle events asynchronously to not block the publisher

  **Logging:**
  - `DEBUG [webhook.events] EventReceived {eventType, eventID}`
  - `INFO [webhook.events] EventEnqueued {eventType, eventID, subscriptions}`
  - `WARN [webhook.events] NoSubscriptions {eventType}`

  **Files:** `backend/internal/modules/integrations/adapters/webhook/event_subscriber.go` (new), `backend/internal/cli/server_http.go` (wire subscriber)

- [x] **Task 10: Implement MCP server over stdio**
  Implement the `vedo-edutrack mcp` subcommand as a production-ready MCP server over stdio.

  **MCP server** (`backend/internal/cli/mcp.go` — replace stub):
  - Implement MCP protocol (JSON-RPC over stdio: `initialize`, `tools/list`, `tools/call`)
  - Use a lightweight MCP library or implement protocol directly (MCP spec is simple — JSON-RPC 2.0 with specific method names)
  - Read-oriented tools (per ADR-§7):
    - `get_route`: returns current route for a learner (horizons, steps)
    - `get_progress`: returns plan-vs-actual progress
    - `get_coverage`: returns FGOS/profstandard coverage
    - `get_gaps`: returns diagnosed gaps with root causes
    - `get_resources`: returns resources for a module
  - Auth: API key passed via MCP initialization or environment variable
  - Reuse Application layer queries from existing bounded contexts (routeplanning, executionprogress, gapcoverage)
  - Graceful shutdown on stdin close

  **CLI wiring** (`backend/internal/cli/mcp.go`):
  - Remove stub `fmt.Println("MCP server not yet implemented")`
  - Real implementation reading from stdin, writing to stdout (stderr for logs)

  **Logging:**
  - All MCP logs go to stderr (not stdout — stdio is for JSON-RPC protocol)
  - `DEBUG [mcp] ToolCalled {tool, args}`
  - `INFO [mcp] SessionStarted {clientInfo}`
  - `INFO [mcp] SessionEnded`
  - `ERROR [mcp] ToolError {tool, error}`

  **Files:** `backend/internal/cli/mcp.go`, `backend/internal/modules/integrations/adapters/mcp/` (new: mcp adapter)

- [x] **Task 11: MCP server tests**
  Write tests for the MCP server.

  **Tests** (`backend/internal/cli/mcp_test.go` or `backend/internal/modules/integrations/adapters/mcp/mcp_test.go`):
  - `initialize` → returns server capabilities and protocol version
  - `tools/list` → returns tool definitions with schemas
  - `tools/call` for each tool → returns expected response format
  - Missing auth → error response
  - Invalid tool name → error response
  - Malformed JSON-RPC → error response
  - Stdin close → graceful shutdown

  **Logging:** assert log output to stderr in tests

  **Files:** `backend/internal/cli/mcp_test.go` (new), `backend/internal/modules/integrations/adapters/mcp/mcp_test.go` (new)

<!-- Commit checkpoint: tasks 9-11 -->

### Phase 4: API Documentation & Validation

- [ ] **Task 12: Complete OpenAPI spec and generate API documentation**
  Finalize the OpenAPI 3.1 specification and set up API documentation.

  **OpenAPI spec** (`backend/api/openapi/v1.yaml`):
  - Review and finalize all endpoints: ensure error responses (400/401/403/404/429/500/503) are documented for every endpoint
  - Add `description` and `example` fields for all schemas and parameters
  - Add `tags` grouping: SPARQL, Webhooks, Routes, Progress, Coverage, Resources, Practice, Auth
  - Add `security` section documenting JWT Bearer auth scheme
  - Document rate limiting behavior in endpoint descriptions
  - Add `x-readme` and `x-codeSamples` extensions for documentation generation

  **Swagger UI** (served at `/api/v1/docs` during development):
  - Already available via oapi-codegen generated server; verify it's accessible

  **Logging:** N/A (spec work)

  **Files:** `backend/api/openapi/v1.yaml`

- [ ] **Task 13: Create integration examples and sandbox demo data**
  Create ready-to-use integration examples and demo data for early API partners.

  **Integration examples** (`docs/integration/`):
  - `quickstart.md`: 5-minute guide — get token, compute a route, check progress
  - `webhook-guide.md`: how to subscribe, receive events, verify signatures, handle dedup
  - `sparql-guide.md`: example SPARQL queries for common use cases
  - `mcp-guide.md`: how to connect AI agents to EduTrack MCP server

  **Code examples** (`docs/integration/examples/`):
  - `curl-examples.sh`: curl commands for all endpoints
  - `python-example.py`: Python client using requests + JWT auth
  - `webhook-receiver-example.js`: Node.js webhook receiver with signature verification

  **Sandbox seed data** (`backend/cmd/vedo-edutrack/` or seed):
  - Extend existing `seed` command to create demo webhook subscriptions and test data
  - Add `--integration-demo` flag to `vedo-edutrack seed` for integration sandbox setup

  **Logging:** N/A (docs work)

  **Files:** `docs/integration/`, `backend/internal/cli/seed.go` (extend)

- [ ] **Task 14: Webhook E2E contract tests**
  Write end-to-end contract tests for the webhook system.

  **Tests** (`tests/e2e/api/webhook_contract_test.go` or `tests/integration/webhook_e2e_test.go`):
  - **Idempotency:** Send duplicate events with same event_id → 200 OK both times, only one delivery
  - **Delivery signing:** Verify HMAC-SHA256 signature header on delivered webhook
  - **Retry:** Simulate subscriber downtime (HTTP 500) → verify retry with exponential backoff
  - **Permanent failure:** Simulate 5 consecutive failures → verify subscription deactivated, no more deliveries
  - **Event type filtering:** Subscriber for `module.mastered` only → should not receive `plan.deviated`
  - **Multiple subscriptions:** Verify event delivered to all matching subscriptions
  - **Subscription lifecycle:** Create → ping → update → delete → verify no more deliveries

  **Test infrastructure:**
  - Use `httptest` server as mock webhook receiver
  - Use testcontainers for PostgreSQL outbox
  - Verify dedup constraint at database level

  **Logging:** Log test assertions for debugging

  **Files:** `tests/integration/webhook_e2e_test.go` (new), `backend/internal/modules/integrations/adapters/webhook/outbox_test.go` (extend)

- [ ] **Task 15: Integration hardening — security, rate limiting, and input validation**
  Final hardening pass across all integration surfaces.

  **Security** (`backend/internal/`):
  - SPARQL injection: verify parameterized queries (no string concatenation), add injection test suite
  - Webhook signing: ensure HMAC secret is never logged, rate-limit webhook subscription creation (max 5/hour per tenant)
  - Input validation: URL allowlist (only https://), event type enum validation, JSON payload size limits (max 64KB)
  - Auth: verify all endpoints require valid JWT (no anonymous access to SPARQL/webhook/API)
  - Permissions: RBAC matrix — only admin/methodologist roles can manage webhook subscriptions

  **Rate limiting hardening** (extend Task 3):
  - SPARQL: 10 req/min per user (configurable), burst allowance 2
  - REST API: 100 req/min per user (configurable)
  - Webhook delivery: max 10 concurrent deliveries per worker

  **Monitoring** (`backend/internal/platform/telemetry/`):
  - Add OTel spans for SPARQL queries, webhook deliveries, MCP tool calls
  - Prometheus metrics: `webhook_delivery_total`, `webhook_delivery_duration_seconds`, `sparql_query_duration_seconds`, `mcp_tool_calls_total`

  **Logging:**
  - `WARN [security] InvalidWebhookURL {url, reason}` 
  - `INFO [gate] SubscriptionRateLimited {userID}`
  - `DEBUG [gate] PermissionCheck {userID, role, action, resource}`

  **Files:** `backend/internal/platform/`, `backend/internal/modules/integrations/adapters/`, `backend/internal/modules/identityaccess/` (RBAC integration)

<!-- Commit checkpoint: tasks 12-15 -->
