# Research

Updated: 2026-08-02 18:30
Status: active

## Active Summary (input for /aif-plan)
<!-- aif:active-summary:start -->
Topic: i18n-readiness NFR — close quality-matrix cell 2.2 recommendation; deferred decisions recorded (A/B testing, vertical certification, RTL)

Goal: Create an NFR requirement for localization readiness (RU+EN supported, easy addition of new languages) following REQ-NFR conventions, plus record deferred decisions (A/B framework, industry vertical certification, RTL support) for roadmap review.

Constraints:
- NFR naming: `REQ-NFR-<area>.<qualifier>.<attribute>` (specs/requirements/README.md)
- Mandatory fields: Приоритет, Ключевая функция, Источник, Описание, Критерии приёмки (measurable)
- Priority: P2 (recommendatory); cell 2.2 (Эргономика/Эксплуатация)
- traceability.ttl must get a new NFR instance (tr:nfr-<area>-<qualifier>-<attribute>)
- ui_language = ru, artifact_language = en, technical_terms = keep
- Phase: M10 (Phase 4) — i18n hygiene NFR is set now (cheap now, expensive later)

Decisions:
1. Write i18n-readiness NFR now (P2): RU (primary) + EN (M10); "easy language addition" = adding a language requires only translations + config, 0 code changes, ≤ 1 day for one developer
2. A/B testing framework — DEFERRED: no product hypotheses for MVP; overlaps with REQ-NFR-ops.release.canary-kill-switch
3. Industry vertical certification (ISO 27001, PCI DSS) — DEFERRED: no regulated vertical selected (M7); PCI scope likely sits with payment acquirer
4. RTL support — DEFERRED: no RTL-market decision; SNG expansion (Moscow, Kazan, Vladivostok) does not require RTL

Open questions:
- Exact acceptance metric for "easy addition": ≤ 1 day per language? 0 code changes?
- Should user documentation (REQ-NFR-ops.compliance.user-documentation) be localized in the same NFR or cross-referenced?

Success signals:
- REQ-NFR-ops.compliance.i18n-readiness.md created with measurable criteria
- traceability.ttl updated with new instance
- Deferred decisions visible in RESEARCH for roadmap review

Next step: `/aif-plan` to create implementation plan for the i18n-readiness NFR (single file + traceability + commit)
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

<!-- aif:sessions:end -->
