# Research

Updated: 2026-08-03 02:43
Status: active

## Active Summary (input for /aif-plan)
<!-- aif:active-summary:start -->
Topic: Gate automation for the agent pipeline — unified gate manifest + two-tier fast/delivery model

Goal: Make the "task ready for delivery" decision executable: one gate manifest consumed by the dev loop (fast tier), the delivery handoff (full tier), CI, and agent skills; the agent must run the full tier before declaring a task done.

Constraints:
- ui_language = ru, artifact_language = en, technical_terms = keep
- Single source of truth: deploy/ci/gates.yaml; consumers (Makefile, CI T12, /aif-verify, pre-commit later) call the runner and never duplicate commands
- Shared output contract: aif-gate-result JSON (schema v1: gate, status pass|warn|fail, blocking, blockers, affected_files, suggested_next)
- Fast tier (dev loop): ≤ 2–3 min, no Postgres/Docker/E2E — go build, tsc --noEmit, gofmt/biome, golangci-lint, biome ci, unit of touched modules, gen-consistency (openapi/sqlc), validate:mermaid, gitleaks
- Delivery tier (handoff): fast + full unit -race, integration (Postgres), e2e (phase <= current — regression: current + all previous phases), docker build ≤ 20 MB, gosec, pnpm audit, syft (ci-main), atlas validate, coverage-check, agent gates (TQS/RCS, docs/env/drift)
- Severity: blocking = NFR acceptance criteria (lint 0 errors, format 0 diff); advisory escalates (coverage advisory on M0.2 → blocking on M1)
- Traceability: gates reference NFR IDs (REQ-NFR-process.dev.engineering-gates P0, REQ-NFR-process.dev.test-coverage P1)

Decisions:
1. Two-tier model: tier fast (feedback loop during implementation) and tier delivery (full set, mandatory before "ready for delivery"); fast ⊆ delivery
2. Gate manifest + runner: deploy/ci/gates.yaml + deploy/ci/run-gates.sh; CI YAML becomes a skeleton calling the runner with --group; Makefile gains make dev-check (fast) and make check (delivery)
3. Decision rule: task ready ⇔ delivery tier completes with zero blocking fails; then /aif-verify (semantic: drift, docs, context gates) → /aif-review → /aif-security-checklist → /aif-commit
4. Wire into aif-implement: Step 3.4 → --tier fast after each task; Final Step → --tier delivery mandatory before /aif-verify
5. Phase-aware regression: run all gates with phase <= current_phase (current phase from ROADMAP)

Open questions:
- Plan task placement in m0-2-engineering-platform.md (new Phase 7 "Gate Automation" + T16 vs extending T12)
- Manifest location: deploy/ci/gates.yaml vs .ai-factory (project convention: deploy/ci for infra)
- Pre-commit hook scope (later, --trigger precommit, fast gates only)
- Mutation testing (REQ-NFR-process.dev.test-coverage ≤ 15%): phase m1, advisory initially

Success signals:
- RESEARCH saved with the gate-automation design ✅
- M0.2 plan updated with the gate-automation task (pending: /aif-improve)

Next step: /aif-improve to add the gate-automation task (Phase 7, T16) to m0-2-engineering-platform.md
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

### 2026-08-03 02:43 — Gate automation: unified manifest + two-tier fast/delivery model

What changed:
Designed how to unify delivery gates into automation: a single gate manifest + runner consumed by the dev loop, delivery handoff, CI, and agent skills. Mapped the user's gate list (build, lint/format, docker, unit, integration+e2e for current and all previous phases) to M0.2 T12 and closed the gaps: typecheck, generated-code consistency (oapi-codegen/sqlc), test quality (TQS/RCS), security (gosec/gitleaks/syft/pnpm audit), DB migrations (atlas), traceability, artifact consistency (docs/env/drift).

Key notes:
- Two tiers: fast (dev feedback loop: compile/types/lint/format/touched-unit/gen-consistency/mermaid/secrets, ≤ 2–3 min, no external services) and delivery (full set, phase <= current for regression; mandatory before declaring a task ready)
- Single source of truth: deploy/ci/gates.yaml; runner deploy/ci/run-gates.sh --tier fast|delivery [--phase] [--out-format table|json|github]; aggregated aif-gate-result JSON; exit 1 on blocking fail
- Decision rule: ready for delivery ⇔ delivery tier zero blocking fails → /aif-verify (semantic gates) → /aif-review → /aif-security-checklist → /aif-commit
- Wiring: aif-implement Step 3.4 → fast tier per task; Final Step → delivery tier before /aif-verify; CI T12 jobs call the runner with --group (no command duplication); Makefile gains make dev-check / make check
- NFR grounding: REQ-NFR-process.dev.engineering-gates (P0: lint 0 errors, format 0 diff, review ≥ 1 approve, runbooks, make up ≤ 30 min), REQ-NFR-process.dev.test-coverage (P1: unit ≥ 90% core / ≥ 80% rest, integration ≥ 70% API contracts, e2e 100% MVP Must, mutation ≤ 15%, regression = 0)
- Honest boundaries: code review, mutation testing, runbook coverage, onboarding — declared in the manifest (runner: agent / phase m1) but not executed as shell commands

Links (paths):
- specs/requirements/REQ-NFR-process.dev.engineering-gates.md
- specs/requirements/REQ-NFR-process.dev.test-coverage.md
- .ai-factory/plans/m0-2-engineering-platform.md (T12 CI pipeline, T9 Makefile)
- .agents/skills/aif-implement/SKILL.md (Step 3.4, Final Step — Verify or Commit)
- .agents/skills/aif-verify/SKILL.md (Step 2/3, aif-gate-result contract)
- .agents/skills/aif-test-quality/SKILL.md (TQS/RCS/B1–B7)
- .agents/skills/aif-security-checklist/SKILL.md

<!-- aif:sessions:end -->
