# Project Base Rules — VEDO EduTrack

> Auto-detected conventions from codebase analysis. The project is greenfield (no code yet),
> so these rules are placeholders to be refined once the stack is chosen and code lands.
> Edit as needed.

## Naming Conventions

- Files: *TBD — to be filled when the stack is chosen*
- Variables: *TBD*
- Functions: *TBD*
- Classes: *TBD*

## Module Structure

- *TBD — no modules exist yet.* Expected boundaries (from vision): service layer reading VEDO Hub API
  (route engine, learning plan, gap diagnosis), REST API for EdTech, web dashboards, integrations.

## Error Handling

- *TBD — no code yet.* Expected: structured error responses (see DESCRIPTION.md, Non-Functional Requirements).

## Logging

- *TBD — no code yet.* Expected: configurable via `LOG_LEVEL` (see DESCRIPTION.md).

## Tooling (Auxiliary Tasks)

- Auxiliary/tooling tasks run through **pnpm** (root `package.json` scripts + `scripts/`), never ad-hoc `npm install` in the repo root.
- Tooling dependencies are added as devDependencies: `pnpm add -D <pkg>`; lockfile `pnpm-lock.yaml` is committed and is the single source of truth.
- `pnpm validate:mermaid` — validates mermaid blocks in `specs/c4/*.md`; `pnpm validate:mermaid:all` — across all specs (used in CI / before committing C4 or DDD diagrams).
