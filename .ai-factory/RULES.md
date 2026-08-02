# Project Rules

> Short, actionable rules and conventions for this project. Loaded automatically by /aif-implement.
> Compiled from `specs/` (ADR, C4, requirements, use-cases, user-stories READMEs and glossary).

## Rules

- Artifact identifiers are semantic and hierarchical; follow the established ID formats (ADR, UC, REQ, US) and never use numeric indexes.
- ADR ID format: `ADR-<LEVEL>.<AREA>.<semantic-tag>` where LEVEL is `BIZ|DES|IMPL` and AREA is `API|DATA|INFRA|SECURITY|UI|PROCESS|INTEGRATION|STACK|OPS|DOC`.
- ADR semantic-tag is kebab-case, 2–5 words, max 40 characters, no spaces; use the documented patterns (`-vs-`, `-or-`, `-tradeoff`, `-adoption`, `-mandate`, `-strategy`/`-approach`/`-pattern`, `-evolution`/`-migration`, `-scope`/`-boundary`).
- Each ADR includes: status, date, context, source requirements, decision, considered alternatives, and consequences.
- ADR IDs are unique; check for topic overlap before creating a new ADR.
- Changes to accepted ADRs are recorded by marking the old one superseded/replaced and creating a new ADR, never by rewriting the accepted one.
- UC ID format: `UC-<L1>.<L2>.<L3>` where L1 is a domain (`plan|execute|resource|viz|practice|api|a11y`) and L3 is a free semantic kebab-case tag; no digits or camelCase.
- Each UC includes: actor(s), priority, key function (F1–F6), channel, description, main flow, alternative flows, postconditions, and source requirements.
- UC channel must be explicit: `GUI|API|CLI|Webhook|Schedule|Mixed`; automatic scenarios (recomputation, webhooks) use Webhook or Schedule, not GUI.
- REQ ID format: `REQ-<TYPE>-<domain-or-area>.<qualifier>.<action-or-attribute>` with TYPE being `FR` or `NFR`.
- FR domains are `plan|execute|resource|viz|practice|api|a11y`; NFR areas are `api|security|data|infra|ui|ops|integration|process|doc`.
- NFR qualifiers are `performance|availability|observability|compliance|maintainability`.
- One requirement file contains one atomic requirement; the filename matches the REQ ID.
- Each requirement records: priority (P0/P1/P2), key function (F1–F6), source (UC and problem P1–P22), description, and acceptance criteria.
- US ID format: `US-<domain>.<subdomain>.<action>` where domain mirrors the UC L1 domain.
- User stories are written in Gherkin with tags `@US-*`, `@UC-*`, and `@P*`; one story per file, filename matches the US ID, anchor `<a id="...">` precedes the heading.
- Spec content (UC, US, FR) is authored in Russian; entity names, technical identifiers, and system error messages stay in English.
- Traceability chain US → UC → FR → COMP → TEST is maintained in `traceability.ttl` (OWL Turtle at the project root).
- C4 diagrams cover System Context, Container, and Component levels only (Code level is not used).
- C4 filenames follow `<LEVEL>-<name>.md` where LEVEL is `context|container|component` and name is kebab-case.
- EduTrack is a service layer over VEDO Hub: it reads ontologies via Hub REST API/MCP and never stores or edits them.
- Never duplicate VEDO Hub responsibilities: ontology editing, versioning, ABox, Git model, social hub, and LLM extraction belong to Hub.
- Respect both deployment contours (Community and Enterprise) in INFRA/SECURITY ADRs and C4 diagrams.
- Route is a function, not a document: describe recomputation on triggers, never CRUD on routes; C4 must not depict route as data storage.
- Domain events `module.mastered`, `plan.deviated`, `route.recalculated` trigger cascade recomputations; describe these cascades explicitly.
- Route computation is deterministic; LLM output (stories, project ideas, assessment items) is generated with validation.
- Stack is TBD: STACK/IMPL ADRs are deferred until the stack is chosen; stack-independent invariants (boundaries, data model, API contract) are decided at DES level.
- Use the established domain terms from `specs/glossary.md` (Route, Trajectory, Mastery, Gap, Checkpoint, Coverage, etc.) consistently.
- Auxiliary and tooling tasks (validation, generation, doc checks) run through **pnpm** from the root `package.json` (`pnpm install`, `pnpm <script>`); tooling lives in `scripts/` and is added as devDependencies via `pnpm add -D`. Never run ad-hoc `npm install` in the repo root — pnpm is the single package manager (lockfile `pnpm-lock.yaml`), Node version pinned in `.nvmrc`.
- Any change to an artifact requires a traceability pass in `traceability.ttl`: verify and update the entire affected chain, both upstream (back to `specs/vision.md`) and downstream (to code, tests, and documentation).
- Every implementation plan must include a `## Progress` section with checkboxes for each task, grouped by phase, placed between `## Roadmap Linkage` and the `## Tasks` section; this serves as a quick progress overview independent of per-task checkboxes in the task body.
