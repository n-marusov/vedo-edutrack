# Research

Updated: 2026-08-02 17:00
Status: active

## Active Summary (input for /aif-plan)
<!-- aif:active-summary:start -->
Topic: Roadmap update — integrating 21-step engineering algorithm (DDD → first `make up`) with existing product roadmap

Goal: Evolve ROADMAP.md to include a Phase 0 covering requirements elaboration, DDD modeling, architecture decisions, and engineering foundation — the work that must happen before the first line of production code.

Constraints:
- Artifact naming MUST follow existing semantic conventions from READMEs: US = `US-<domain>.<subdomain>.<action>`, UC = `UC-<L1>.<L2>.<L3>`, ADR = `ADR-<LEVEL>.<AREA>.<semantic-tag>`
- FR and NFR conventions do NOT exist yet — must be created following the same hierarchical pattern
- ui_language = ru, artifact_language = en, technical_terms = keep
- The roadmap is a product document; engineering phases should be added without diluting the product focus

Decisions:
1. Phase 0 with 4 milestones: M0.0 (Requirements), M0.1 (Domain Model & Architecture), M0.2 (Engineering Platform), M0.3 (App Scaffold)
2. Requirements elaboration produces 4 artifact types: US, UC, FR, NFR — all with semantic IDs
3. FR naming convention: `FR-<domain>.<subdomain>.<capability>` (by analogy with US)
4. NFR naming convention: `NFR-<category>.<subcategory>.<semantic-tag>` (by analogy with ADR/UC)
5. README files needed in `specs/requirements/` for FR and NFR conventions
6. UC L1-domains for EduTrack need to be defined (current UC README is VEDO Core-specific)

Open questions:
- Exact L1-domains for EduTrack UC (candidates: route, plan, viz, auth, api, resource, practice)
- FR/NFR category taxonomy detail
- Phase 0 timebox: 3-4 weeks for a team of 2-3

Success signals:
- ROADMAP.md includes Phase 0 with M0.0–M0.3
- `specs/requirements/` has README for FR and NFR conventions
- Traceability chain US → UC → FR → COMP → TEST is formalized

Next step: `/aif-plan` to create implementation plan for ROADMAP update with Phase 0
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
<!-- aif:sessions:end -->
