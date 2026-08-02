# Research

Updated: 2026-08-03 12:00
Status: active

## Active Summary (input for /aif-plan)
<!-- aif:active-summary:start -->
Topic: CLI interface for the backend executable — embedded into the single binary (decision taken, ADR created)

Goal: Fixed the CLI-embedding decision as ADR-DES.API.cli-interface; updated folder structure (cmd/vedo-edutrack + internal/cli) and all corresponding documentation. Next: carry the CLI skeleton into M0.2/M0.3 plans (T1 structure + T9 Makefile already updated).

Constraints:
- ui_language = ru, artifact_language = en, technical_terms = keep
- Single binary `vedo-edutrack` (Go, cobra); CLI = input adapter over Application layer (same pattern as MCP)
- CLI is dev/support/testing tooling (trusted operator, bypasses JWT by design); audit logging via zap; scriptable (no prompts); --output json|table|csv
- Per-command lazy wire: each command builds its minimal graph (migrate = config + Postgres; server = full graph)
- MCP stdio (REQ-FR-api.mcp.server) requires a spawnable mode of the same binary — CLI embedding is required, not optional

Decisions:
1. Single binary with cobra subcommands: server / mcp / migrate / seed / ontology sync / route compute / plan get / gap diagnose / report — ADR-DES.API.cli-interface (ПРИНЯТО, 2026-08-03)
2. cmd/server renamed to cmd/vedo-edutrack; CLI command tree lives in internal/cli (not cmd/, not platform/)
3. CLI also used for testing: route compute --stub (route engine TDD without Postgres/Hub); CLI commands unit-tested with mocked ports
4. Makefile becomes a thin wrapper: make migrate → vedo-edutrack migrate up
5. Traceability: new ADR instance tr:adr-des-api-cli-interface in traceability.ttl

Open questions:
- Exact M0.3 scope: minimal skeleton (server + migrate + seed + version) in M0.3 vs route compute --stub deferred to M1 with the engine
- --output formats finalized (json/table/csv) once commands are implemented
- Audit-log schema for CLI invocations (actor, command, args, result) — design in M0.3/M1

Success signals:
- ADR-DES.API.cli-interface.md created with full alternative analysis ✅
- repository-structure + development-tooling ADRs, AGENTS.md, DESCRIPTION.md, M0.2 plan, traceability.ttl updated ✅
- go build ./... + pnpm validate:mermaid pass ✅

Next step: /aif-plan to fold the CLI skeleton into M0.2/M0.3 tasks (T1 structure already updated; add cobra wiring + command stubs + audit logging)
<!-- aif:active-summary:end -->

## Sessions
<!-- aif:sessions:start -->
### 2026-08-02 17:00 — Initial exploration: 21-step algorithm vs current ROADMAP

What changed:
Mapped the user's 21-step DDD-to-`make up` algorithm against the existing 10-milestone product roadmap to identify gaps.

Key findings:
1. **The current ROADMAP has no engineering foundation phase.** It assumes architecture, stack, repo structure, containers, CI/CD, Keycloak, and test framework already exist by M1. They don't.

2. **The algorithm covers 5 phases of engineering work:**
   - Phase 0: DDD modeling (steps 1-2)
   - Phase 1: Architecture basis — ADR, C4, RBAC, traceability (steps 3-6)
   - Phase 2: Engineering foundation — repo, Docker, Makefile, test scaffold (steps 7-10)
   - Phase 3: Implementation — scaffold, Keycloak, landing page (steps 11-16)
   - Phase 4: Quality — TDD, CI gates, observability, README, first launch (steps 17-21)

3. **Only phases 3-4 roughly correspond to roadmap M1-M4.** Phases 0-2 have no roadmap representation.

4. **A requirements elaboration phase is missing from the algorithm itself.** The algorithm jumps to DDD without formalizing US, UC, FR, NFR artifacts. These MUST be produced before DDD begins.

5. **Semantic naming conventions exist** for US (`US-<domain>.<subdomain>.<action>`), UC (`UC-<L1>.<L2>.<L3>`), and ADR (`ADR-<LEVEL>.<AREA>.<semantic-tag>`). FR and NFR conventions do NOT exist — `specs/requirements/` is empty.

6. **Proposed Phase 0 structure** (to insert before the existing Phase 1):
   - M0.0: Requirements elaboration (US, UC, FR, NFR with semantic IDs)
   - M0.1: Domain model & architecture basis (DDD, C4, ADR, RBAC, traceability)
   - M0.2: Engineering platform (repo structure, Docker, Makefile, CI gates)
   - M0.3: App scaffold (healthcheck, Keycloak, stubs, landing page)

Links (paths):
- `vedo-edutrack/specs/user-stories/README.md` — US naming convention
- `vedo-edutrack/specs/use-cases/README.md` — UC naming convention (VEDO Core; EduTrack needs its own L1 domains)
- `vedo-edutrack/specs/adr/README.md` — ADR naming convention
- `vedo-edutrack/.ai-factory/ROADMAP.md` — current roadmap (M1–M10, 4 product phases)
- `vedo-edutrack/specs/vision.md` — business requirements (authoritative)
- `vedo-edutrack/specs/glossary.md` — domain terminology
- `vedo-edutrack/specs/requirements/` — empty; needs README for FR and NFR conventions
### 2026-08-02 18:30 — Quality-matrix final audit + recommendatory NFR triage

What changed:
Final quality-matrix audit (2026-08-02) confirmed 20/20 green cells: 107 requirements (60 FR + 47 NFR), 0 lacunas, 0 hierarchy violations, all metrologically valid. Three recommendatory open questions from audit Step 7 were triaged with the user.

Key notes:
- A/B testing framework: deferred — no MVP product hypotheses; overlaps with canary-kill-switch
- Vertical certification (ISO 27001/PCI DSS): deferred — no regulated vertical selected; PCI scope likely with acquirer
- Localization: approved — i18n-readiness NFR (P2, cell 2.2): RU+EN, easy language addition (translations + config only, ≤ 1 day, 0 code changes), RTL explicitly out of scope
- Requirement corpus grew from 79 (66 FR + 13 NFR) to 107 (60 FR + 47 NFR) across 8 loop iterations (commits 0fe934d..2975ae3)

Links (paths):
- specs/requirements/ (60 FR + 47 NFR)
- traceability.ttl (107 requirement instances, 0 orphans)
- specs/vision.md (authoritative)

### 2026-08-03 — CLI interface decision: single binary with cobra subcommands

What changed:
Explored embedding a CLI into the backend executable. Decided: single binary `vedo-edutrack` with cobra subcommands (server / mcp / migrate / seed / ontology sync / route compute / plan get / gap diagnose / report). CLI is an input adapter over the Application layer (same pattern as MCP server). Created ADR-DES.API.cli-interface (ПРИНЯТО) and updated structure + documentation.

Key notes:
- MCP stdio requirement (REQ-FR-api.mcp.server) forces a spawnable binary mode — CLI embedding is required, not optional
- cmd/server renamed to cmd/vedo-edutrack; CLI tree in internal/cli (not cmd/ — minimal entry point; not platform/ — platform must not import modules)
- Per-command lazy wire: migrate = config + Postgres only; server = full graph
- CLI = trusted operator tooling (bypasses JWT by design); audit logging via zap; scriptable (no prompts); --output json|table|csv; queries + targeted commands (no event cascades)
- CLI also for testing: route compute --stub enables route-engine TDD without Postgres/Hub
- Makefile becomes thin wrapper: make migrate → vedo-edutrack migrate up
- Alternatives rejected: two binaries (cmd/server+cmd/cli), thin REST CLI client (kubectl pattern), Makefile-only, CLI in main.go, CLI in platform/

Links (paths):
- specs/adr/ADR-DES.API.cli-interface.md (new ADR)
- specs/adr/ADR-IMPL.PROCESS.repository-structure.md (tree §4, principle 4, alternatives, consequences)
- specs/adr/ADR-IMPL.PROCESS.development-tooling.md (§4 cobra, §11 Makefile wrappers)
- specs/c4/component-api-server.md (Composition Root references)
- AGENTS.md (structure + entry points + docs)
- .ai-factory/DESCRIPTION.md (stack: CLI line)
- .ai-factory/plans/m0-2-engineering-platform.md (T1, T5, T9, T10, T12)
- traceability.ttl (tr:adr-des-api-cli-interface)
- backend/cmd/vedo-edutrack/main.go, backend/internal/cli/cli.go

<!-- aif:sessions:end -->
