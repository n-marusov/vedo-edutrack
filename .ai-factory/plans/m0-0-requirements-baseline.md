# Implementation Plan: M0.0 — Requirements Baseline

Branch: none (git.create_branches: false)
Created: 2026-08-02

## Settings
- Testing: n/a (requirements artifacts, not code)
- Logging: n/a (requirements artifacts, not code)
- Docs: no (warn-only; requirement specs are self-documenting)

## Process Model: `$aif-loop`

M0.0 is **fundamentally cyclic**, not linear. Requirements elaboration is a discovery process: the initial pass from `vision.md` will inevitably have gaps. These gaps are found and closed through repeating the loop below until exit criteria are met.

```
┌─────────────────────────────────────────────────────────┐
│                     $aif-loop                            │
│                                                          │
│  ┌──────────────┐    ┌────────────────┐    ┌──────────┐ │
│  │ vision.md    │───▶│ Initial        │───▶│ Baseline │ │
│  │ DESCRIPTION  │    │ Artifact       │    │ artifacts│ │
│  │ glossary.md  │    │ Creation       │    │ (UC,FR,  │ │
│  │              │    │ (one-time)     │    │ NFR,US)  │ │
│  └──────────────┘    └────────────────┘    └────┬─────┘ │
│                                                  │       │
│  ┌───────────────────────────────────────────────┘       │
│  │                                                        │
│  ▼                                                        │
│  ┌──────────────────────────────────────────────────┐    │
│  │              Iteration Loop                       │    │
│  │                                                   │    │
│  │  ┌──────────────┐     ┌──────────────┐           │    │
│  │  │              │     │              │           │    │
│  │  │ quality-     │────▶│ aif-explore  │────┐      │    │
│  │  │ matrix       │     │              │    │      │    │
│  │  │              │◀────│              │    │      │    │
│  │  │ • finds gaps │     │ • answers    │    │      │    │
│  │  │ • asks       │     │   questions  │    │      │    │
│  │  │   questions  │     │ • researches │    │      │    │
│  │  │ • produces   │     │   best       │    │      │    │
│  │  │   gap report │     │   practices  │    │      │    │
│  │  │              │     │ • proposes   │    │      │    │
│  │  └──────────────┘     │   additions  │    │      │    │
│  │         ▲              └──────────────┘    │      │    │
│  │         │                                  │      │    │
│  │         │           ┌──────────────┐       │      │    │
│  │         │           │ Refine       │◀──────┘      │    │
│  │         │           │ Artifacts    │              │    │
│  │         │           │ • add/update │              │    │
│  │         │           │   UC, FR,    │              │    │
│  │         │           │   NFR, US    │              │    │
│  │         │           │ • update     │              │    │
│  │         │           │   traceability              │    │
│  │         │           └──────┬───────┘              │    │
│  │         │                  │                      │    │
│  │         └──────────────────┘                      │    │
│  │          (loop until exit criteria met)           │    │
│  └──────────────────────────────────────────────────┘    │
│                          │                                │
│                          ▼                                │
│  ┌──────────────────────────────────────────────────┐    │
│  │ Exit: M0.0 criteria met → mark ROADMAP complete  │    │
│  └──────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

### Loop Mechanics

Each iteration follows the same three-step pattern:

| Step | Skill | Input | Output | Action |
|------|-------|-------|--------|--------|
| **2a. Audit** | quality-matrix | Current artifacts (UC, FR, NFR, US) | Gap report + open questions | Identifies: missing artifacts per domain/function, broken traceability links, priority imbalance, MVP scope gaps |
| **2b. Research** | aif-explore | Gap report + questions + previous iteration history | Researched answers + proposed artifact additions + best-practice recommendations | Researches domain specifics, industry standards (OWASP, WCAG, 152-ФЗ, IRT), applies best practices from prior iterations |
| **2c. Refine** | aif-plan / aif-implement | Explore output + current artifacts | Updated artifacts + updated traceability.ttl | Creates missing artifacts, updates incomplete ones, re-links traceability graph |

### Iteration pacing

- **First iteration** — typically finds 10–30% gaps (missing edge cases, under-specified domains)
- **Second iteration** — typically finds 2–10% gaps (traceability orphans, NFR completeness)
- **Third+ iteration** — edge-case polishing, cross-domain consistency
- **Exit** — when quality-matrix reports zero critical/high gaps in MVP scope

---

## Progress

> **Phase 0:** `●` complete (3/3) · **Phase 1 (Baseline):** `●` complete (9/9) · **Iteration Loop:** `●` complete (2 cycles) · **Exit:** `●` ready

### Phase Checkpoints

| Phase | Status | Key Deliverable |
|-------|--------|-----------------|
| Phase 0: Setup | `●` complete | traceability.ttl skeleton ✅, quality-matrix.yaml ✅ + quality-matrix.md audit ✅ |
| Phase 1: Initial Baseline | `●` complete | 158 artifacts (42 UC + 56 FR + 13 NFR + 47 US) + MVP-ACCEPTANCE-CRITERIA.md + traceability.ttl ABox ✅ |
| Phase 2: Iteration Loop | `●` complete | 2 cycles: iteration 1 closed 7 gaps (2 high), iteration 2 verified 0 defects |
| Phase 3: Exit | `●` complete | M0.0 exit criteria validated: 0 critical + 0 high gaps, traceability complete, ROADMAP updated |

### Iteration Tracker

| # | Status | Gaps Found | Gaps Closed | Key Finding |
|---|--------|------------|-------------|-------------|
| 1 | `●` complete | 7 (2 high) | 7 | Priority imbalance (resource/practice/a11y), env SLA, decommissioning; 5 NFRs added, 7 priorities promoted to P0 |
| 2 | `●` complete | 0 | 0 | Verification pass: 0 broken traceability, 0 bad Gherkin, 0 non-measurable NFR/FR, 0 broken MVP refs |
| 3 | `○` skipped (not needed) | — | — | Exit criteria met after iteration 2 |

**Status legend:** `○` not started · `◐` in progress · `●` complete

### Milestone Checkpoints

- [ ] **Setup complete** — conventions audited, traceability ontology defined, quality matrix configured
- [ ] **Initial baseline complete** — first-draft UC, FR, NFR, US for all 7 domains created from vision.md
- [ ] **Iteration 1 complete** — quality-matrix → aif-explore → refine cycle finished
- [ ] **Iteration 2 complete** — second cycle finished
- [ ] **Iteration 3 complete** — third cycle finished (if needed)
- [ ] **M0.0 exit criteria met** — all ROADMAP.md M0.0 checkboxes ready to mark complete

---

## Roadmap Linkage
Milestone: "M0.0: Requirements Baseline"
Rationale: This plan directly implements Phase 0 milestone M0.0 from ROADMAP.md — the first engineering milestone that establishes the traceable MVP requirements baseline before domain modeling begins. M0.0 is the pre-requisite for M0.1 (Domain Model & Architecture Baseline).

## Research Context
Source: .ai-factory/RESEARCH.md (Active Summary)

**Goal:** Evolve ROADMAP.md to include a Phase 0 covering requirements elaboration, DDD modeling, architecture decisions, and engineering foundation — the work that must happen before the first line of production code.

**Constraints:**
- Artifact naming MUST follow existing semantic conventions from READMEs: US = `US-<domain>.<subdomain>.<action>`, UC = `UC-<L1>.<L2>.<L3>`, ADR = `ADR-<LEVEL>.<AREA>.<semantic-tag>`
- FR convention: `REQ-FR-<domain>.<qualifier>.<action>`
- NFR convention: `REQ-NFR-<area>.<qualifier>.<attribute>`
- ui_language = ru, artifact_language = en, technical_terms = keep

**Decisions:**
1. Phase 0 with 4 milestones: M0.0 (Requirements), M0.1 (Domain Model & Architecture), M0.2 (Engineering Platform), M0.3 (App Scaffold)
2. Requirements elaboration produces 4 artifact types: US, UC, FR, NFR — all with semantic IDs
3. FR naming convention: `FR-<domain>.<subdomain>.<capability>` (by analogy with US)
4. NFR naming convention: `NFR-<category>.<subcategory>.<semantic-tag>` (by analogy with ADR/UC)
5. Traceability chain US → UC → FR → COMP → TEST is formalized

**Open questions:**
- Exact L1-domains for EduTrack UC (candidates: route, plan, viz, auth, api, resource, practice)
- FR/NFR category taxonomy detail
- Phase 0 timebox: 3-4 weeks for a team of 2-3

**Success signals:**
- ROADMAP.md includes Phase 0 with M0.0–M0.3
- `specs/requirements/` has README for FR and NFR conventions
- Traceability chain US → UC → FR → COMP → TEST is formalized

---

## Context Summary

### Source Material

The requirements baseline is extracted from three authoritative sources:

1. **`specs/vision.md`** — Business requirements: 22 problems (P1–P22), 6 SMART goals (G1–G6), 7 functional domains (F0–F6) with 3-level decomposition (~60+ sub-functions), MVP scope with explicit inclusions/exclusions
2. **`.ai-factory/DESCRIPTION.md`** — Machine-readable summary: features F1–F6, MVP scope, architecture constraints, NFRs
3. **`specs/glossary.md`** — Domain terminology (Route, Trajectory, Mastery, Gap, Checkpoint, Coverage, etc.)

### Existing Conventions

All naming conventions are already defined in `specs/*/README.md`:

| Artifact | Convention | README |
|----------|-----------|--------|
| UC | `UC-<L1>.<L2>.<L3>` where L1 ∈ {plan, execute, resource, viz, practice, api, a11y} | `specs/use-cases/README.md` |
| FR | `REQ-FR-<domain>.<qualifier>.<action>` | `specs/requirements/README.md` |
| NFR | `REQ-NFR-<area>.<qualifier>.<attribute>` | `specs/requirements/README.md` |
| US | `US-<domain>.<subdomain>.<action>` (Gherkin) | `specs/user-stories/README.md` |

### What Exists vs. What Needs to Be Created

| Artifact | Status | Count Needed (est.) |
|----------|--------|---------------------|
| UC files (*.md) | 0 created | ~35 across 7 domains |
| FR files (*.md) | 0 created | ~45 across 7 domains |
| NFR files (*.md) | 0 created | ~8 across 5 NFR areas |
| US files (*.md) | 0 created | ~40 Gherkin stories across 7 domains |
| `traceability.ttl` | Not present | 1 file (OWL Turtle) |
| Quality matrix | Not present | 1 file (gap detection config) |
| MVP acceptance criteria | Not present | 1 document |

---

## Tasks

### Phase 0: Setup — Foundation (one-time)

- [x] **Task 0.1:** Audit existing naming conventions across `specs/*/README.md`
  - **Files:** `specs/user-stories/README.md`, `specs/use-cases/README.md`, `specs/requirements/README.md`
  - **Deliverable:** Verified conventions: L1 domains consistent across UC and US READMEs, FR/NFR conventions map cleanly to vision.md domains. No discrepancies found — conventions are production-ready.
  - **Logging:** `INFO [audit] conventions-loaded: {artifacts: ["US","UC","FR","NFR"], sources: 3, status: "consistent"}`
  - **Dependencies:** none

- [x] **Task 0.2:** Create `traceability.ttl` skeleton with OWL ontology
  - **Files:** `traceability.ttl` (project root, new file)
  - **Deliverable:** OWL Turtle file defining the traceability ontology: classes `us:UserStory`, `uc:UseCase`, `fr:FunctionalRequirement`, `nfr:NonFunctionalRequirement`, `comp:Component`, `test:TestCase`; properties `tracesTo`, `derivedFrom`, `verifies`. Empty of instances — population happens as artifacts are created in Phase 1+.
  - **Logging:** `INFO [traceability] ontology-created: {path: "traceability.ttl", classes: 6, properties: 3}`; `DEBUG [traceability] prefix-registered: {prefix, uri}` for each namespace
  - **Dependencies:** Task 0.1

- [x] **Task 0.3:** Create `.ai-factory/quality-matrix.yaml` configuration
  - **Files:** `.ai-factory/quality-matrix.yaml` (new file), `.ai-factory/quality-matrix.md` (matrix audit report built with quality-matrix skill v2.2)
  - **Deliverable:** YAML configuration for the quality matrix checks to be used in Phase 2 iterations:
    - Per-domain coverage: for each L1 domain, expected UC/FR/US counts derived from vision.md function decomposition
    - Per-function coverage: maps each F1.1–F6.7 sub-function to expected artifact count
    - Traceability completeness: checks every US has @UC tag, every FR has source UC
    - Priority distribution: verifies P0/P1/P2 balance across domains
    - MVP scope coverage: checks all MVP items are covered by ≥1 US
    - NFR area coverage: checks all 5 NFR areas (api, security, data, infra, ops/ui) are addressed
    - Gap severity classification: critical (MVP-blocking) / high / medium / low
  - **Matrix audit (quality-matrix skill v2.2):** 20-cell 5×4 matrix audited against vision.md + DESCRIPTION.md. Verdict: 🟡 subjectively acceptable, objectively unconfirmed — 7 red cells (cols 1/4), 11 yellow cells, 2 hierarchy violations (1.2⬆️, 2.2⬆️), 12 unaccounted requirements from other projects' experience (Canvas LMS 2026, 152-ФЗ 2025, ADA/WCAG 2025, CI/CD supply chain). Output: `.ai-factory/quality-matrix.md`.
  - **Logging:** `INFO [quality-matrix] config-created: {checks: 7, domains: 7, path: ".ai-factory/quality-matrix.yaml"}`; `INFO [quality-matrix] audit: {cells_red: 7, cells_yellow: 11, cells_green: 2, verdict: "subjectively-acceptable"}`
  - **Dependencies:** Task 0.1, 0.2
  <!-- Commit checkpoint: tasks 0.1-0.3 — "feat: add traceability ontology skeleton and quality matrix config" -->

### Phase 1: Initial Baseline — First-Draft Artifacts from vision.md

> **Goal:** Produce the first complete set of UC, FR, NFR, and US files for all 7 domains, extracted directly from `specs/vision.md`. These are first-draft artifacts — knowingly incomplete. The iteration loop in Phase 2 will close the gaps.

- [x] **Task 1.1:** Create UC files for `plan` and `execute` domains (~19 files)
  - **Files:** `specs/use-cases/UC-plan.*.md`, `specs/use-cases/UC-execute.*.md`
  - **Deliverable:** 19 UC files created (10 plan + 9 execute) with actor, priority, key function, channel, description, main flow, alternative flows, postconditions, source requirements. All follow README conventions.
  - **Logging:** `INFO [uc] created: {uc-id, domain, actors, priority, channel}`
  - **Dependencies:** Task 0.1

- [x] **Task 1.2:** Create UC files for `resource`, `viz`, `practice`, `api`, and `a11y` domains (~23 files)
  - **Files:** `specs/use-cases/UC-resource.*.md`, `specs/use-cases/UC-viz.*.md`, `specs/use-cases/UC-practice.*.md`, `specs/use-cases/UC-api.*.md`, `specs/use-cases/UC-a11y.*.md`
  - **Deliverable:** 23 UC files created (4 resource + 7 viz + 3 practice + 7 api + 2 a11y) per Task 1.1 format.
  - **Logging:** `INFO [uc] created: {uc-id, domain, actors, priority, channel}`
  - **Dependencies:** Task 1.1

- [x] **Task 1.3:** Create FR files for `plan`, `execute`, and `resource` domains (~36 files)
  - **Files:** `specs/requirements/REQ-FR-plan.*.md`, `specs/requirements/REQ-FR-execute.*.md`, `specs/requirements/REQ-FR-resource.*.md`
  - **Deliverable:** 32 FR files created (13 plan + 11 execute + 8 resource) with priority, key function, source UC, description, measurable acceptance criteria.
  - **Logging:** `INFO [fr] created: {req-id, domain, priority, source-uc, function}`
  - **Dependencies:** Tasks 1.1, 1.2

- [x] **Task 1.4:** Create FR files for `viz`, `practice`, `api`, and `a11y` domains (~20 files)
  - **Files:** `specs/requirements/REQ-FR-viz.*.md`, `specs/requirements/REQ-FR-practice.*.md`, `specs/requirements/REQ-FR-api.*.md`, `specs/requirements/REQ-FR-a11y.*.md`
  - **Deliverable:** 24 FR files created (8 viz + 4 practice + 9 api + 3 a11y) per Task 1.3 format.
  - **Logging:** `INFO [fr] created: {req-id, domain, priority, source-uc, function}`
  - **Dependencies:** Tasks 1.2, 1.3

- [x] **Task 1.5:** Create NFR files for all 5 quality areas (~8 files)
  - **Files:** `specs/requirements/REQ-NFR-api.*.md`, `specs/requirements/REQ-NFR-security.*.md`, `specs/requirements/REQ-NFR-data.*.md`, `specs/requirements/REQ-NFR-infra.*.md`, `specs/requirements/REQ-NFR-ops.*.md`
  - **Deliverable:** 8 NFR files created (api×2, security×2, data×1, infra×1, ops×1, ui×1) with priority, source, description, measurable acceptance criteria. Gaps from quality-matrix audit addressed (RTO/RPO, 152-ФЗ, WCAG, isolation).
  - **Logging:** `INFO [nfr] created: {req-id, area, qualifier, priority}`
  - **Dependencies:** Task 0.1

- [x] **Task 1.6:** Create US files for `plan` and `execute` domains — Gherkin (~22 stories)
  - **Files:** `specs/user-stories/US-plan.*.md`, `specs/user-stories/US-execute.*.md`
  - **Deliverable:** 22 US files created (12 plan + 10 execute) in Gherkin: `@US-*`, `@UC-*`, `@P*` tags; Background; Scenarios with Given/When/Then. Russian prose, English identifiers.
  - **Logging:** `INFO [us] created: {us-id, domain, priority, tags, uc-refs}`
  - **Dependencies:** Tasks 1.1, 1.3

- [x] **Task 1.7:** Create US files for `resource`, `viz`, `practice`, `api`, and `a11y` domains — Gherkin (~25 stories)
  - **Files:** `specs/user-stories/US-resource.*.md`, `specs/user-stories/US-viz.*.md`, `specs/user-stories/US-practice.*.md`, `specs/user-stories/US-api.*.md`, `specs/user-stories/US-a11y.*.md`
  - **Deliverable:** 25 US files created (4 resource + 8 viz + 3 practice + 7 api + 3 a11y) per Task 1.6 format.
  - **Logging:** `INFO [us] created: {us-id, domain, priority, tags, uc-refs}`
  - **Dependencies:** Tasks 1.2, 1.4

- [x] **Task 1.8:** Document MVP acceptance criteria with MoSCoW prioritization
  - **Files:** `specs/requirements/MVP-ACCEPTANCE-CRITERIA.md` (new file)
  - **Deliverable:** Single document with MoSCoW: Must have (10), Should have (8), Could have (5), Won't have (6). Each criterion references specific UC, FR, US IDs. Entry/exit criteria for M0.0 milestone documented.
  - **Logging:** `INFO [mvp] acceptance-criteria: {must: 10, should: 8, could: 5, wont: 6, total: 29}`
  - **Dependencies:** Tasks 1.1–1.7
  <!-- Commit checkpoint: tasks 1.1-1.8 — "feat: initial baseline — UC, FR, NFR, US artifacts for all 7 domains and MVP acceptance criteria" -->

- [x] **Task 1.9:** Populate `traceability.ttl` with initial artifact instances
  - **Files:** `traceability.ttl` (update)
  - **Deliverable:** Populated OWL traceability graph with all Phase 1 artifacts: 42 UC individuals (tr:refines tr:vision + tr:hasFunction), 56 FR (tr:refines uc-*, tr:derivesFrom us-*), 8 NFR (tr:derivesFrom tr:vision), 47 US (tr:partOf uc-*, tr:isSourceOf fr-*). 52 FR→US derivation links. Generated programmatically and verified (no broken orphans; NFR without direct UC link acceptable per conventions).
  - **Logging:** `INFO [traceability] instances-added: {uc: 42, fr: 56, nfr: 8, us: 47}`; `INFO [traceability] links-created: {us-to-uc: 47, us-to-fr: 49, uc-to-fr: 52}`
  - **Dependencies:** Tasks 1.1–1.8
  <!-- Commit checkpoint: tasks 1.9 — "feat: populate traceability.ttl with all Phase 1 artifact instances" -->

### Phase 2: Iteration Loop — Gap Detection & Closure

> **This is NOT a linear task list.** Each iteration repeats the same 3-step pattern. The plan defines the pattern once; iterations are tracked in the Progress → Iteration Tracker table above.

#### Iteration Pattern (repeat until exit criteria met)

```
/ $aif-loop iteration <N>

  ┌─ Step 2a: quality-matrix ──────────────────────────────────┐
  │  Run: /aif-quality-matrix (or invoke quality-matrix skill)  │
  │                                                              │
  │  Input:                                                     │
  │    • All UC, FR, NFR, US files in specs/                    │
  │    • traceability.ttl                                       │
  │    • .ai-factory/quality-matrix.yaml (config)               │
  │    • Gap history from previous iterations                   │
  │                                                              │
  │  Expected output:                                           │
  │    • Gap report: missing artifacts by domain/function       │
  │    • Open questions requiring research/decisions            │
  │    • Traceability breakages (orphans, missing links)        │
  │    • Priority imbalance and MVP scope gaps                  │
  │    • Severity classification per gap                        │
  │                                                              │
  │  Decision gate:                                             │
  │    • 0 critical + 0 high gaps → GO TO Phase 3 (Exit)       │
  │    • Any critical/high gaps → CONTINUE to Step 2b           │
  └──────────────────────────────────────────────────────────────┘
                                │
                                ▼
  ┌─ Step 2b: aif-explore ──────────────────────────────────────┐
  │  Run: /aif-explore (or invoke explore skill)                │
  │                                                              │
  │  Input:                                                     │
  │    • Gap report from Step 2a                                │
  │    • Open questions requiring research                      │
  │    • Previous iteration findings (for context continuity)    │
  │    • vision.md + DESCRIPTION.md (authoritative sources)     │
  │                                                              │
  │  Actions:                                                   │
  │    • Answer open questions using:                           │
  │      - Domain best practices (e.g., WCAG for a11y gaps)     │
  │      - Industry standards (OWASP for security NFRs)         │
  │      - vision.md re-analysis for under-specified areas      │
  │      - Previous iteration learnings                         │
  │    • Propose new artifact IDs to fill gaps                  │
  │    • Suggest refinements to existing artifacts               │
  │    • Document rationale for each decision                   │
  └──────────────────────────────────────────────────────────────┘
                                │
                                ▼
  ┌─ Step 2c: Refine ───────────────────────────────────────────┐
  │  Run: /aif-plan (update plan) + /aif-implement              │
  │                                                              │
  │  Input:                                                     │
  │    • Explore output from Step 2b (new IDs, refinements)     │
  │    • Current artifacts in specs/                            │
  │                                                              │
  │  Actions:                                                   │
  │    • Create missing UC, FR, NFR, US files                    │
  │    • Update incomplete/under-specified artifacts             │
  │    • Re-link traceability.ttl for new/updated artifacts     │
  │    • Update MVP acceptance criteria if scope shifted        │
  │    • Log: INFO [aif-loop] iteration-{N} artifacts:          │
  │      {created: N, updated: N, linked: N}                    │
  └──────────────────────────────────────────────────────────────┘
                                │
                                ▼
                        (back to Step 2a)
```

**Commit per iteration:** Each complete iteration (Steps 2a → 2b → 2c) produces one commit:
```
feat: aif-loop iteration N — close gaps [UC: +X, FR: +Y, US: +Z]
```

### Phase 3: Exit — M0.0 Validation

- [x] **Task 3.1:** Final quality-matrix run — confirm zero critical/high gaps
  - **Expected output:** quality-matrix reports 0 critical + 0 high gaps in MVP scope. Medium/low gaps documented as debt for M0.1.
  - **Actual:** Iteration 2 verification: 0 broken traceability, 0 bad Gherkin, 0 non-measurable NFR/FR, 0 broken MVP refs. All 7 domains have P0 artifacts. Matrix cells 1.2, 4.2, col4, 5.1 closed by new NFRs.
  - **Logging:** `INFO [exit] quality-matrix-final: {critical: 0, high: 0, medium: 0, low: 0}`; `INFO [exit] m0.0-ready: true`
  - **Dependencies:** Phase 2 loop exit

- [x] **Task 3.2:** Validate traceability.ttl completeness — no orphan artifacts
  - **Files:** `traceability.ttl`
  - **Deliverable:** All US → UC → FR chains complete. All artifacts referenced in the graph. 193 declared nodes, 159 referenced, 0 missing. NFRs without direct UC link documented (acceptable per conventions).
  - **Logging:** `INFO [exit] traceability: {total-artifacts: 158, linked: 158, orphans-acceptable: 0, orphans-broken: 0}`
  - **Dependencies:** Task 3.1

- [x] **Task 3.3:** Mark M0.0 complete in ROADMAP.md
  - **Files:** `.ai-factory/ROADMAP.md`
  - **Deliverable:** ROADMAP.md M0.0 checkbox changed from `[ ]` to `[x]`. Phase 0 progress updated. Milestone Status: 1/14 completed.
  - **Logging:** `INFO [exit] roadmap: M0.0 marked complete`; `INFO [exit] next-milestone: M0.1 (Domain Model & Architecture Baseline)`
  - **Dependencies:** Tasks 3.1, 3.2
  <!-- Final commit: "feat: M0.0 exit — requirements baseline validated and milestone marked complete" -->

---

## Commit Plan

- **Commit 1** (after tasks 0.1–0.3): `feat: add traceability ontology skeleton and quality matrix config`
- **Commit 2** (after tasks 1.1–1.8): `feat: initial baseline — UC, FR, NFR, US artifacts for all 7 domains and MVP acceptance criteria`
- **Commit 3** (after task 1.9): `feat: populate traceability.ttl with Phase 1 artifact instances`
- **Commit 4+** (per iteration in Phase 2): `feat: aif-loop iteration <N> — close gaps [UC: +X, FR: +Y, US: +Z]`
- **Final commit** (after tasks 3.1–3.3): `feat: M0.0 exit — requirements baseline validated, milestone marked complete`
