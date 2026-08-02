# Project Roadmap

> VEDO EduTrack — an educational route service on top of VEDO Hub that builds personalized, cross-subject learning paths from a shared knowledge ontology.

## Milestones

### Phase 1 — MVP (Q2)

- [ ] **M1: Core Infrastructure (F0 + F1 + F2 + F3 + F6 scaffold)**
  Ontology port via VEDO Hub REST API/MCP (F0.1, F0.2). Route engine: shortest path with strict links, gap diagnosis (root-cause climb), three horizons, recompute on progress/goal change (F1.1–F1.7 except F1.8). Learning plan: checkpoint snapshots, plan-vs-actual, binary readiness forecast (F2.1–F2.3). Basic resource catalog (F3.1). REST API scaffold (`GET /routes/compute`, `/progress`, `/fgos/coverage`). Starter ontology in VEDO Hub: 1000+ topics (grades 5–11, 9 subjects), 500+ cross-subject links, FGOS mapping, 3–5 pedagogy concepts.

- [ ] **M2: Family Education — «Дай пять» (F4 + F2 FGOS + F5)**
  Knowledge map (F4.1): 2D graph with progress/gap color-coding, gap map view (F4.6). Route builder (F4.3): select goal → visualize route. Learner dashboard (F4.2), parent dashboard (F4.4), group panel for multiple children (F4.7). FGOS coverage in real time (F2.8), deficit list, attestation readiness report (F2.9). 50+ stories and 30+ project ideas (F5.1–F5.3). Demo editor for quick feature showcase. Target: parent builds route in ≤5 min, NPS ≥ 40.

- [ ] **M3: Corporate Onboarding — «Вектор Компетенций» (MVP corporate)**
  Corporate onboarding plan: compute personalized route for new hires covering professional skills + domain context + values + internal services + socialization. Time-to-productivity metric. Basic corporate dashboard (F4.4 adapted). Target: 2 weeks to productivity (down from 4), one pilot company.

- [ ] **M4: Integration & Webhook Layer (F6 production-ready)**
  Production REST API (F6.1): route computation, progress, FGOS coverage. Read-only SPARQL endpoint (F6.2). Webhooks: `module.mastered`, `plan.deviated` (F6.4). MCP server for AI agents (F6.6).

### Phase 2 — Enrichment & B2B (Q3)

- [ ] **M5: Ontology & Route Enrichment**
  All 5 cross-subject link types: add `enriches`, `isAnalogousTo` (F1.1). Pedagogy-concept-aware routes (F1.8): spiral learning, project immersion, practice-to-theory — route adapts to chosen concept. Cascade recompute on ontology updates (F0.3, F1.2). Learning-style resource matching: visual/text/audio (F1.4, F3.2). Complex risk-level forecast: green/yellow/red with ±10% accuracy (F2.3 upgraded). 200+ stories, 100+ project ideas (F5).

- [ ] **M6: EdTech Platform Integration**
  EdTech partner onboarding: fork public ontology + API integration in ≤1 week. Production API with rate limiting, documentation, SLA foundation. Out-of-the-box LMS connectors: WebTutor, iSpring (F6.3). Enterprise SSO/SAML via Keycloak (F6.5). Content-agnostic integration: platform retains content & audience ownership (F6.7). Target: 20+ platforms using public ontology.

### Phase 3 — Verticals & Community (Q3–Q4)

- [ ] **M7: Corporate Application & Compliance**
  Full «Вектор Компетенций» application on shared service mechanics. Career tracks: gap analysis to target role, learning path from current position (F1.5 expanded). Regulatory compliance module: bind regulatory requirements to modules and real work scenarios. Corporate group dashboard: onboarding funnel, time-to-productivity per department, compliance coverage, ROI analytics. Target: 50+ employees on simultaneous onboarding under control, ROI 16:1.

- [ ] **M8: Community & Network Effect**
  Contributor program via VEDO Hub: fork + merge request flow, contributor profiles. Semantic merge comparison (diff visualization). Partner program for EdTech platforms. Community growth targets: 5000+ topics, 8000+ links, 200+ active contributors. Each contributor's work reaches thousands of learners through platform forks.

### Phase 4 — Enterprise & Scale (Q4+)

- [ ] **M9: Enterprise Deployment & Compliance**
  On-premise / private cloud deployment: VEDO Hub Enterprise + EduTrack in corporate perimeter. Dedicated API endpoints with SLAs. SAP SuccessFactors integration (F6.3). Predictive analytics: deficit forecasting, churn prediction. Isolated corporate ontology with private forks.

- [ ] **M10: Multilingual & Global Readiness**
  Multilingual ontology: Russian + English (expandable). Localization framework for UI and ontology content. Foundation for international expansion of the public knowledge graph.

## Completed

| Milestone | Date |
|-----------|------|
| *(none yet — implementation pending stack selection)* | |

---

## Milestone ↔ Business Goals Traceability

| Milestone | Business Goals |
|-----------|---------------|
| M1: Core Infrastructure | G1, G3 |
| M2: Family Education | G1, G3, G4 |
| M3: Corporate Onboarding | G6 |
| M4: Integration Layer | G2, G5 |
| M5: Ontology & Route Enrichment | G1, G3 |
| M6: EdTech Integration | G2 |
| M7: Corporate & Compliance | G4, G6 |
| M8: Community & Network Effect | G5 |
| M9: Enterprise & Scale | G5, G6 |
| M10: Multilingual | G5 |
