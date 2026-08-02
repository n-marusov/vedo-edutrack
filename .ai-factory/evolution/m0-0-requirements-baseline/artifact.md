# M0.0 Requirements Baseline — Final Artifact

**Run:** `m0-0-requirements-baseline-20260802-100000`
**Status:** ✅ completed (`threshold_reached`)
**Date:** 2026-08-02
**Method:** aif-loop (quality-matrix skill v2.2 → explore → refine) × 2 iterations

---

## 1. Deliverables

| Artifact | Count | Location |
|----------|-------|----------|
| Use Cases (UC) | 42 | `specs/use-cases/UC-*.md` |
| Functional Requirements (FR) | 56 | `specs/requirements/REQ-FR-*.md` |
| Non-Functional Requirements (NFR) | 13 | `specs/requirements/REQ-NFR-*.md` |
| User Stories (US, Gherkin) | 47 | `specs/user-stories/US-*.md` |
| MVP Acceptance Criteria (MoSCoW) | 1 | `specs/requirements/MVP-ACCEPTANCE-CRITERIA.md` |
| Quality Matrix config | 1 | `.ai-factory/quality-matrix.yaml` |
| Quality Matrix audit | 1 | `.ai-factory/quality-matrix.md` |
| Traceability graph (OWL Turtle) | 1 | `traceability.ttl` (193 nodes, populated ABox) |

**Total artifacts:** 158 + 3 support documents.

## 2. Domain coverage

| Domain | UC | FR | US | P0 present |
|--------|----|----|----|------------|
| plan (F1) | 10 | 13 | 12 | ✅ (20) |
| execute (F2) | 9 | 11 | 10 | ✅ (12) |
| resource (F3) | 4 | 8 | 4 | ✅ (3) |
| viz (F4) | 7 | 8 | 8 | ✅ (14) |
| practice (F5) | 3 | 4 | 3 | ✅ (2) |
| api (F6) | 7 | 9 | 7 | ✅ (4) |
| a11y (WCAG) | 2 | 3 | 3 | ✅ (2) |

## 3. NFR areas (13 files)

- **api** ×3: latency-p95 (≤200ms @1000 conc), webhook-idempotency, hub-dependency-sla
- **security** ×3: role-based-access, pii-152-fz, owasp-application-security
- **data** ×2: backup-rpo (RPO≤1h, RTO≤4h), pii-export-deletion (right to forget)
- **infra** ×2: community-enterprise-isolation, cicd-supply-chain-security (SBOM, secrets)
- **ui** ×2: wcag-level (2.1 AA), supported-browsers
- **ops** ×1: log-level-config (LOG_LEVEL)

## 4. Quality matrix verdict (final)

- 🔴 Lacunae: **0** | 🟠 Not measurable: **0** | 🟡 Incomplete: **2** (debt M0.1) | 🟢 Closed: **18**
- Hierarchy violations: **0**
- **Verdict: 🟢 objectively confirmed for MVP scope** (0 critical + 0 high gaps)

### Iteration log
- **Iteration 1:** 7 gaps found (priority imbalance, env SLA, decommissioning, supply chain) → 5 NFRs added, 7 priorities promoted to P0, TTL updated
- **Iteration 2:** verification pass — 0 broken traceability, 0 bad Gherkin, 0 non-measurable requirements, 0 broken MVP refs

## 5. Traceability

`traceability.ttl` populated with 42 UC + 56 FR + 13 NFR + 47 US instances:
- US → UC: `tr:partOf` (47 links)
- US → FR: `tr:isSourceOf` (49 links)
- FR → UC: `tr:refines` (52 links)
- FR → US: `tr:derivesFrom` (52 links)
- All references resolve (193 declared, 0 missing)

## 6. Milestone

`ROADMAP.md`: **M0.0: Requirements Baseline** marked ✅ complete (2026-08-02).
Next: **M0.1: Domain Model & Architecture Baseline**.
