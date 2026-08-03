# pnpm Reference

> Source: pnpm official docs (pnpm.io/settings, pnpm.io/continuous-integration, pnpm.io/supply-chain-security) + vedo-edutrack project experience (2026-08-03)
> Created: 2026-08-03
> Updated: 2026-08-03

## Overview

pnpm is a fast, disk-space-efficient Node.js package manager with a content-addressable store. It enforces strict dependency resolution by default (no phantom dependencies) and has a first-class workspace (monorepo) model. This reference covers the **pnpm v11 configuration model**, which changed significantly from v10: settings moved out of `package.json`/`.npmrc` into `pnpm-workspace.yaml`. The vedo-edutrack project hit exactly these v11 changes during the 2026-08-03 infrastructure session — this reference captures both the official model and the project-specific lessons.

**Project context:** vedo-edutrack uses pnpm as a workspace (`pnpm-workspace.yaml` at repo root, `frontend/` as a workspace member). `packageManager: pnpm@11.18.0` is pinned in the root `package.json`.

## Core Concepts

**Config model (v11) — the single most important concept:** pnpm settings now live in `pnpm-workspace.yaml` using **camelCase** keys. The `pnpm` field in `package.json` is **no longer read at all**. `.npmrc` is used **only** for auth/registry credentials (and is gitignored).

| Category | Stored in | Format |
|----------|-----------|--------|
| All pnpm/install settings (`nodeLinker`, `allowBuilds`, `overrides`, `catalog`, …) | `pnpm-workspace.yaml` (project) + `config.yaml` (global) | YAML, camelCase |
| Auth & registry credentials | `.npmrc` (gitignored) + global `rc` | INI |

**allowBuilds (build-script approval):** by default pnpm does **not** run dependency lifecycle scripts (`preinstall`/`install`/`postinstall`). Packages must be explicitly approved in one `allowBuilds` map. Unreviewed packages with build scripts cause install to exit non-zero with `ERR_PNPM_IGNORED_BUILDS` (when `strictDepBuilds` is true, the default).

**Workspaces:** monorepo support via `packages:` globs in `pnpm-workspace.yaml`; `--filter` runs commands in specific packages; one shared `pnpm-lock.yaml` at the root.

**Frozen lockfile:** CI must use `--frozen-lockfile` (or `pnpm ci`) — fails if the lockfile needs updates, guaranteeing reproducibility.

## Configuration (vedo-edutrack)

Current project config in `pnpm-workspace.yaml` (repo root):

```yaml
packages:
  - "frontend"
  - "apps/*"
  - "tools/*"
onlyBuiltDependencies:   # ⚠️ see pitfall below — v11 uses allowBuilds, not this
  - '@biomejs/biome'
  - esbuild
  - lefthook
```

Root `.npmrc` (kept minimal — auth only per v11 model):

```
pnpm.onlyBuiltDependencies[]=@biomejs/biome   # ⚠️ legacy key — see pitfall
pnpm.onlyBuiltDependencies[]=esbuild
pnpm.onlyBuiltDependencies[]=lefthook
```

**Important project fact:** the current config uses the **legacy** v10 keys (`onlyBuiltDependencies` in `.npmrc`/workspace-yaml). Per the pnpm v11 skill, these were renamed to `allowBuilds: { name: true }` in `pnpm-workspace.yaml`. The project should migrate:

```yaml
# pnpm-workspace.yaml (v11 canonical)
allowBuilds:
  '@biomejs/biome': true
  esbuild: true
  lefthook: true
```

## Usage Patterns

### Install (host)

```bash
pnpm install                          # dev
pnpm install --frozen-lockfile        # CI / reproducible
pnpm ci                               # clean + frozen install
```

### Approve build scripts

```bash
pnpm approve-builds                   # interactive prompt
pnpm approve-builds --all             # approve all pending
pnpm approve-builds esbuild fsevents  # approve specific, !name = deny
pnpm add --allow-build=esbuild pkg    # approve while adding
```

### Workspace filtering

```bash
pnpm --filter @vedo-edutrack/frontend build   # build one package
pnpm --filter "...[origin/main]" test         # changed packages only
pnpm --filter @vedo-edutrack/frontend ...     # run from workspace root
```

### Docker build stage (v11 PATH change)

```dockerfile
FROM node:24-alpine AS spa-build
WORKDIR /app
# Copy workspace root files before sources (layer caching)
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc ./
COPY frontend/package.json ./frontend/package.json
RUN --mount=type=cache,id=pnpm-store,target=/root/.local/share/pnpm/store \
    corepack enable && pnpm install --frozen-lockfile
COPY frontend/ ./frontend/
RUN pnpm --filter @vedo-edutrack/frontend build
```

Note: in pnpm v11, global binaries live at `$PNPM_HOME/bin` — set `PATH="$PNPM_HOME/bin:$PATH"` in Docker (not `$PNPM_HOME`).

## Best Practices

1. **Put all settings in `pnpm-workspace.yaml` (camelCase)** — never in `package.json#pnpm` (ignored in v11) or `.npmrc` (auth only). `pnpm config get nodeLinker` reads it.
2. **Always use `--frozen-lockfile` in CI and Docker builds** — the lockfile is the single source of truth; fail instead of silently updating.
3. **Approve build scripts explicitly via `allowBuilds`** — never use `dangerouslyAllowAllBuilds: true`. Unreviewed builds fail with `ERR_PNPM_IGNORED_BUILDS`.
4. **Pin pnpm via `packageManager` in `package.json`** (exact) or `devEngines.packageManager` (ranges). Match the CI major to the lockfile writer — v11 fails on incompatible lockfiles.
5. **Cache the pnpm store in CI** (`pnpm store path --silent`) across trusted jobs only.
6. **Use `--mount=type=cache` for the pnpm store in Docker** to speed up rebuilds.
7. **Copy lockfile + workspace manifest before sources in Dockerfile** — maximizes layer cache hits.
8. **Keep auth tokens out of the repo** — project `.npmrc` is gitignored; use `${NPM_TOKEN}` env expansion or user-level auth file.
9. **In a workspace, run member commands from the root** with `--filter`, not by `cd`-ing into the package (keeps workspace resolution correct).
10. **Never mix npm and pnpm artifacts** — a stray `package-lock.json` (from `npm install`) breaks pnpm's deps check.

## Common Pitfalls

1. **`ERR_PNPM_IGNORED_BUILDS` on install** — package with a postinstall script isn't in `allowBuilds`. Fix: `pnpm approve-builds` or add to `allowBuilds` in `pnpm-workspace.yaml`. (Hit in this project: `@biomejs/biome`, `esbuild`.)
2. **`pnpm` field in `package.json` silently ignored** — v11 removed it entirely. Config written there (e.g. `pnpm.onlyBuiltDependencies`) has **no effect** — this bit the project on 2026-08-03 (frontend Docker build failed repeatedly).
3. **Legacy config keys still "work" partially** — `onlyBuiltDependencies` in `.npmrc`/workspace yaml is the v10 spelling. v11 renamed it to `allowBuilds`. The old key may be read or ignored depending on exact version — migrate to `allowBuilds`.
4. **`pnpm config get onlyBuiltDependencies` returns `undefined` from a workspace member** — config resolution from a subdirectory can miss root settings. Always run config-sensitive commands from the workspace root.
5. **Docker: frozen-lockfile fails because lockfile isn't visible** — when the build context doesn't include `pnpm-lock.yaml` (it lives at the workspace root). Copy it explicitly in the Dockerfile, or the container can't verify the lockfile.
6. **npm ↔ pnpm artifact mixing** — `npm install` in a directory with pnpm creates `package-lock.json` and a different `node_modules` layout, breaking pnpm's deps check. Use one package manager per project.
7. **v11 Docker PATH** — global binaries at `$PNPM_HOME/bin`, not `$PNPM_HOME`. Getting this wrong → `command not found: pnpm` in containers.
8. **Version mismatch CI vs lockfile** — v11 refuses to rewrite a lockfile written by a newer major. Keep CI pnpm in sync.

## Version Notes

- **pnpm v10 → v11 migration** (relevant to this project's `pnpm@11.18.0`):
  - `package.json#pnpm` field → removed; settings → `pnpm-workspace.yaml`
  - `onlyBuiltDependencies`/`neverBuiltDependencies`/`ignoredBuiltDependencies`/`onlyBuiltDependenciesFile` → `allowBuilds: { name: true|false }`
  - `managePackageManagerVersions`/`packageManagerStrict`/`COREPACK_ENABLE_STRICT` → `pmOnFail: download|ignore|warn|error`
  - Global binaries dir → `$PNPM_HOME/bin`
  - CI: frozen-lockfile auto-enabled; incompatible lockfile = hard error
  - `npm_config_*` env vars no longer read → use `pnpm_config_*`/`PNPM_CONFIG_*`
  - `pnpm config get`/`list` output JSON (not INI); `--location=project` writes to `pnpm-workspace.yaml`
- **Supply-chain defaults (v11):** `minimumReleaseAge: 1440` (1 day), `blockExoticSubdeps: true`, tarball integrity mismatch = hard error.

## Project To-Do (from 2026-08-03 session)

- [ ] Migrate `pnpm-workspace.yaml` from `onlyBuiltDependencies` (legacy list) to `allowBuilds` (v11 canonical map)
- [ ] Remove stale `pnpm.onlyBuiltDependencies[]` entries from root `.npmrc` (legacy key; `.npmrc` should be auth-only)
- [ ] Align frontend Docker build with v11: ensure `pnpm-lock.yaml` is visible in the container build context

## Source References

- pnpm skill (global, v2026.6.22): `core-config.md`, `best-practices-ci.md`, `features-supply-chain-security.md`
- https://pnpm.io/settings
- https://pnpm.io/configuring
- https://pnpm.io/pnpm-workspace_yaml
- https://pnpm.io/continuous-integration
- https://pnpm.io/supply-chain-security
- https://pnpm.io/cli/approve-builds
