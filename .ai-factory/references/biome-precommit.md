# Biome Pre-commit / Git Hooks Reference

> Source: https://biomejs.dev/guides/getting-started/, https://biomejs.dev/guides/integrate-in-vcs/, https://biomejs.dev/reference/cli/, https://biomejs.dev/reference/configuration/, https://biomejs.dev/recipes/git-hooks/, https://biomejs.dev/recipes/continuous-integration/
> Created: 2026-08-02
> Updated: 2026-08-02

## Overview

Biome is a fast, all-in-one toolchain for JavaScript/TypeScript/JSON/CSS/HTML/GraphQL: formatter, linter, import sorter, and (experimental) GritQL-powered analyzer. It replaces ESLint + Prettier in one dependency written in Rust, with no plugins needed for the core workflow. It runs with zero configuration, works standalone (no Node.js required for the binary) or as an npm dev dependency, and is designed to be used in pre-commit hooks (format + lint only staged files) and CI (`biome ci`).

This reference focuses on integrating Biome with git pre-commit workflows: VCS integration options (`--staged`, `--changed`, `--since`), hook managers (Lefthook, Husky + lint-staged, git-format-staged, pre-commit framework), raw shell hooks, and CI usage.

## Core Concepts

**`biome check`**: runs formatter + linter + import sorting together (the "one command" workflow). With `--write` it applies safe fixes, formatting and import sorting in one pass. This is the command most pre-commit setups use.

**`--write` / `--fix`**: `--fix` is an alias of `--write`. Without it, commands only report diagnostics and exit non-zero on errors (read-only). `--unsafe` additionally applies unsafe fixes and should only be used with `--write`/`--fix`.

**`--staged`**: process only files added to the git index (files prepared to be committed). Intended for local pre-commit use. Not available on the `ci` command.

**`--changed`**: process only files changed relative to a base branch (`vcs.defaultBranch` in `biome.json`, or `--since=<ref>` which takes precedence). Intended for CI. On the `ci` command, VCS mode uses `--changed` instead of `--staged` (a remote repo has no "staged" concept).

**VCS integration is opt-in** — enabling `vcs.enabled` alone does nothing; you must enable the specific features you want.

**`--no-errors-on-unmatched`**: silence errors when no files are processed (essential in hooks — e.g. a commit touching only markdown files). Pre-commit recipes always include it.

**`--files-ignore-unknown=true`**: don't emit diagnostics for file types Biome doesn't know (lets hooks handle any future supported extension without updating the hook).

**`biome ci`**: CI variant of `check` — read-only (no `--write`), GitHub/GitLab-native reporters, `--threads` control, and `--changed` VCS mode. "Files won't be modified, the command is a read-only operation."

## API / Interface

### CLI commands (v2.x)

| Command | Purpose |
|---------|---------|
| `biome check [--write] [--unsafe] [--staged/--changed] [--since=REF] [--only=...] [--skip=...] [--watch] [PATH]...` | Formatter + linter + import sorting |
| `biome lint [--write] [--unsafe] [--suppress] [--reason=...] [--staged/--changed]` | Lint only, apply safe fixes |
| `biome format [--write] [--staged/--changed]` | Format only |
| `biome ci [--changed] [--since=REF] [--threads=N] [--only=...] [--skip=...]` | Read-only check for CI; no `--write` |
| `biome init [--jsonc]` | Bootstrap `biome.json` (or `biome.jsonc`) with defaults |
| `biome migrate [--write]` / `biome migrate prettier` / `biome migrate eslint [--include-inspired] [--include-nursery]` | Update config for breaking changes; map Prettier/ESLint config into Biome |
| `biome explain <NAME>` | Show docs for a rule or CLI aspect (e.g. `biome explain noDebugger`) |
| `biome upgrade`, `biome rage`, `biome start`, `biome stop`, `biome clean`, `biome version` | Tooling: upgrade standalone binary, debug info, daemon control |
| `biome search '`pattern`'` | EXPERIMENTAL GritQL search (quote patterns — shells treat backticks specially) |

### Key options (check/lint/format)

| Option | Behavior |
|--------|----------|
| `--write` / `--fix` | Write safe fixes / formatting / import sorting to disk |
| `--unsafe` | Apply unsafe fixes (needs `--write`) |
| `--staged` | Only files in the git index (local hooks) |
| `--changed` | Only files changed vs base branch (CI) |
| `--since=<REF>` | Base branch for `--changed`; overrides `vcs.defaultBranch` |
| `--only=<GROUP\|RULE\|DOMAIN\|ACTION\|PLUGIN>` | Run only given rules/groups/domains; `plugin` runs analyzer plugins |
| `--skip=<GROUP\|RULE\|DOMAIN\|ACTION\|PLUGIN>` | Skip given rules/groups (takes precedence over `--only`) |
| `--no-errors-on-unmatched` | Don't fail when no files were processed |
| `--error-on-warnings` | Exit with error code if any diagnostic is a warning |
| `--max-diagnostics=<N\|none>` | Cap diagnostics shown (default: 20) |
| `--reporter=<default\|json\|json-pretty\|github\|junit\|summary\|gitlab\|checkstyle\|rdjson\|sarif\|concise>` | Output format; `github` prints GitHub annotations, `gitlab` produces code-quality report |
| `--reporter-file=<PATH>` | Write reporter output to file |
| `--stdin-file-path=<PATH>` | Process code piped from stdin; extension determines language; virtual paths skip ignore checks |
| `--files-ignore-unknown=true` | Skip unknown file types without diagnostics |
| `--diagnostic-level=<info\|warn\|error>` | Minimum diagnostic level to show (default: info) |
| `--threads=<N>` | Limit threads (ci; env `BIOME_THREADS`) |
| `--watch` | Re-run automatically on file changes (not for hooks/CI) |

### Environment variables

`BIOME_CONFIG_PATH` (config file path), `BIOME_LOG_FILE`, `BIOME_LOG_PREFIX_NAME` (default `server.log`), `BIOME_LOG_PATH`, `BIOME_LOG_LEVEL`, `BIOME_LOG_KIND`, `BIOME_WATCHER_KIND`, `BIOME_WATCHER_POLLING_INTERVAL`, `BIOME_THREADS`, `BIOME_DISTRIBUTION` (npm|homebrew|standalone for `biome upgrade`).

## Usage Patterns

### Install (pin the version!)

`-E` pins the exact version — Biome docs recommend it because Biome changes its configuration and CLI across releases; unpinned installs break CI/hooks on new releases.

```bash
npm i -D -E @biomejs/biome
pnpm add -D -E @biomejs/biome
bun add -D -E @biomejs/biome
deno add -D npm:@biomejs/biome
yarn add -D -E @biomejs/biome
```

```bash
npx @biomejs/biome init            # creates biome.json
npx @biomejs/biome init --jsonc    # creates biome.json instead of biome.jsonc
```

### Enable VCS integration (biome.json)

```json
{
  "vcs": {
    "enabled": true,
    "clientKind": "git",
    "useIgnoreFile": true,
    "defaultBranch": "main"
  }
}
```

- `useIgnoreFile: true` — Biome reads `.gitignore` files, Git's local exclude file `.git/info/exclude`, and `.ignore` files (nested ignore files supported). In linked worktrees `.git/info/exclude` is read from the common git directory, matching Git's behavior.
- `defaultBranch` — base branch for `--changed` evaluation.
- `vcs.root` — needed when `biome.json` is not at the VCS root; path is relative to the config file, e.g. `"root": "../"` for `frontend/biome.json` in a repo-root VCS.

### Process only changed/staged files

```bash
biome check --changed                 # vs defaultBranch
biome check --changed --since=next    # vs an arbitrary branch
biome check --staged                  # only staged files (local hooks)
```

### Lefthook (fast, cross-platform, dependency-free hook manager)

`lefthook.yml` at repo root:

```yaml
pre-commit:
  commands:
    check:
      glob: "*.{js,ts,cjs,mjs,d.cts,d.mts,jsx,tsx,json,jsonc,css}"
      run: npx @biomejs/biome check --no-errors-on-unmatched --files-ignore-unknown=true {staged_files}
```

Format, lint, and apply safe code fixes before committing:

```yaml
pre-commit:
  commands:
    check:
      glob: "*.{js,ts,cjs,mjs,d.cts,d.mts,jsx,tsx,json,jsonc,css}"
      run: npx @biomejs/biome check --write --no-errors-on-unmatched --files-ignore-unknown=true {staged_files}
      stage_fixed: true
```

`stage_fixed: true` re-stages the fixed files. Then run `lefthook install`.

### Husky + lint-staged

`package.json`:

```json
{
  "scripts": {
    "prepare": "husky"
  }
}
```

`.husky/pre-commit`:

```sh
lint-staged
```

`lint-staged` config lives in `package.json`:

```json
{
  "lint-staged": {
    "*.{js,ts,cjs,mjs,d.cts,d.mts,jsx,tsx,json,jsonc,css}": [
      "biome check --files-ignore-unknown=true",
      "biome check --write --no-errors-on-unmatched",
      "biome check --write --organize-imports-enabled=false --no-errors-on-unmatched",
      "biome check --write --unsafe --no-errors-on-unmatched",
      "biome format --write --no-errors-on-unmatched",
      "biome lint --write --no-errors-on-unmatched"
    ],
    "*": [
      "biome check --no-errors-on-unmatched --files-ignore-unknown=true"
    ]
  }
}
```

The `"*"` entry handles every file type; unknown extensions are skipped via `--files-ignore-unknown=true`. Always add `--no-errors-on-unmatched` to silence errors when no matching files are staged.

### git-format-staged

Doesn't use `git stash` internally, so no manual intervention on conflicts between unstaged and updated staged changes.

`.husky/pre-commit` (check only):

```sh
git-format-staged --formatter 'biome check --files-ignore-unknown=true --no-errors-on-unmatched --stdin-file-path="{}"' '*'
```

`.husky/pre-commit` (format, lint, apply safe fixes):

```sh
git-format-staged --formatter 'biome check --write --files-ignore-unknown=true --no-errors-on-unmatched --stdin-file-path="{}"' '*'
```

### pre-commit framework (Python)

Four hooks are provided via the `biomejs/pre-commit` repository:

| hook `id` | description |
| --- | --- |
| `biome-ci` | Check formatting, check import organization, lint (read-only) |
| `biome-check` | Format, organize imports, lint, apply safe fixes to committed files |
| `biome-format` | Format the committed files |
| `biome-lint` | Lint and apply safe fixes to the committed files |

`.pre-commit-config.yaml`:

```yaml
repos:
-   repo: https://github.com/biomejs/pre-commit
    rev: "v2.0.6"  # Use the sha / tag you want to point at
    hooks:
    -   id: biome-check
        additional_dependencies: ["@biomejs/biome@2.1.1"]
```

You MUST specify the Biome version via `additional_dependencies` — pre-commit installs tools separately and needs to know which one to install. To avoid double maintenance, use a local hook that reuses the project's npm-installed Biome:

```yaml
repos:
  - repo: local
    hooks:
      - id: local-biome-check
        name: biome check
        entry: npx @biomejs/biome check --write --files-ignore-unknown=true --no-errors-on-unmatched
        language: system
        types: [text]
        files: "\\.(jsx?|tsx?|c(js|ts)|m(js|ts)|d\\.(ts|cts|mts)|jsonc?|css)$"
```

The `files` regex is optional — Biome skips unknown files with `--files-ignore-unknown=true`.

### Raw shell hook (no extra tools)

Check formatting and lint before committing:

```sh
#!/bin/sh
set -eu

npx @biomejs/biome check --staged --files-ignore-unknown=true --no-errors-on-unmatched
```

Format, lint, and apply safe fixes (fails if a staged file also has unstaged changes, then re-stages fixed files):

```sh
#!/bin/sh
set -eu

if git status --short | grep --quiet '^MM'; then
  printf '%s\n' "ERROR: Some staged files have unstaged changes" >&2
  exit 1;
fi

npx @biomejs/biome check --write --staged --files-ignore-unknown=true --no-errors-on-unmatched

git update-index --again
```

### CI

`biome check` vs `biome ci` — `ci` is for CI environments: no `--write` option, native GitHub annotations / GitLab reporters, `--threads` control, and `--changed` instead of `--staged` when VCS integration is enabled.

GitHub Actions (first-party action `biomejs/setup-biome@v2`):

```yaml
name: Code quality
on:
  push:
  pull_request:
jobs:
  quality:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    steps:
      - name: Checkout
        uses: actions/checkout@v5
        with:
          persist-credentials: false
      - name: Setup Biome
        uses: biomejs/setup-biome@v2
        with:
          version: latest
      - name: Run Biome
        run: biome ci .
```

If `biome.json` extends config from a package, install dependencies first (`actions/setup-node@v4` + `npm ci` etc.). Third-party: `mongolyy/reviewdog-action-biome@v1` posts review comments/suggestions on PRs.

GitLab CI (Docker image `ghcr.io/biomejs/biome:latest`; `entrypoint: [""]` is required):

```yaml
stages:
  - quality

lint:
    image:
      name: ghcr.io/biomejs/biome:latest
      entrypoint: [""]
    stage: quality
    script:
        - biome ci --reporter=gitlab --colors=off > /tmp/code-quality.json
        - cp /tmp/code-quality.json code-quality.json
    artifacts:
      reports:
        codequality:
          - code-quality.json
    rules:
        - if: $CI_COMMIT_BRANCH
        - if: $CI_MERGE_REQUEST_ID
```

## Configuration

### VCS section

| Key | Default | Values / notes |
|-----|---------|----------------|
| `vcs.enabled` | `false` | Opt-in; enabling alone does nothing |
| `vcs.clientKind` | — | `"git"` (only supported client) |
| `vcs.useIgnoreFile` | `false` | Read `.gitignore`, `.git/info/exclude`, `.ignore` (nested supported) |
| `vcs.root` | config dir | VCS root, relative to `biome.json`; for non-root configs |
| `vcs.defaultBranch` | — | Base branch for `--changed` evaluation |

### files section

| Key | Default | Notes |
|-----|---------|-------|
| `files.includes` | — | Glob list; `"**"`, negated `!` patterns (processed in order), force-ignore `!!**/dist` for output dirs; `node_modules/` always ignored; matching a folder requires `/**` suffix to include its files, `!dist` (no suffix) to ignore it |
| `files.ignoreUnknown` | `false` | Don't emit diagnostics for unknown file types |
| `files.maxSize` | `1048576` (1 MiB) | Files above this are ignored for performance |

### linter section

| Key | Default | Notes |
|-----|---------|-------|
| `linter.enabled` | `true` | |
| `linter.includes` | files.includes | Applied after `files.includes`; cannot match folders without `/**` |
| `linter.rules.preset` | `"recommended"` | `"recommended"` / `"all"` (except nursery) / `"none"`; `linter.rules.recommended` is deprecated |
| `linter.rules.[group]` | — | Groups: `a11y`, `complexity`, `correctness`, `nursery`, `performance`, `security`, `style`, `suspicious`. Value = severity string or per-rule object |
| severity values | — | `"on"` (rule default), `"off"`, `"info"`, `"warn"`, `"error"` |

Nursery rules require explicit opt-in on stable versions, are enabled by default on nightly, and can be promoted or removed without semver.

### formatter section (global defaults)

| Key | Default |
|-----|---------|
| `formatter.enabled` | `true` |
| `formatter.indentStyle` | `"tab"` |
| `formatter.indentWidth` | `2` (ignored when indentStyle is tab) |
| `formatter.lineEnding` | `"lf"` (`lf`/`crlf`/`cr`/`auto` — auto = CRLF on Windows, LF elsewhere) |
| `formatter.lineWidth` | `80` |
| `formatter.bracketSpacing` | `true` |
| `formatter.delimiterSpacing` | `false` |
| `formatter.expand` | `"auto"` (`auto`/`always`/`never`; `package.json` always uses `always`) |
| `formatter.trailingNewline` | `true` (setting `false` is highly discouraged) |
| `formatter.useEditorconfig` | `false` (biome.json overrides .editorconfig; nested .editorconfig not supported) |
| `formatter.formatWithErrors` | `false` |

### javascript.formatter section (key defaults)

| Key | Default |
|-----|---------|
| `quoteStyle` | `"double"` |
| `jsxQuoteStyle` | `"double"` |
| `quoteProperties` | `"asNeeded"` |
| `trailingCommas` | `"all"` |
| `semicolons` | `"always"` (`always`/`asNeeded` for ASI) |
| `arrowParentheses` | `"always"` |
| `bracketSameLine` | `false` |
| `operatorLinebreak` | `"after"` |

### overrides

`overrides` is a list of `{ "includes": [globs], <top-level section overrides> }`. First matching pattern wins. Sections: `formatter`, `linter`, `javascript`, `json`, or any language section (minus `includes`/`ignore`).

```json
{
  "formatter": { "lineWidth": 100 },
  "overrides": [
    { "includes": ["generated/**"], "formatter": { "lineWidth": 160, "indentStyle": "space" } }
  ]
}
```

### Glob syntax notes

- `*` — any chars except `/`; `**` — recursive (must be a full path component); `[...]`/`[!...]` — character classes; leading `!` — negated exception pattern (can't be used alone).
- Command-line globs are interpreted by the shell, which may not support `**` — prefer globs inside `biome.json`.

## Best Practices

1. **Pin the Biome version** (`npm i -D -E @biomejs/biome`). Biome changes CLI/config behavior between releases; an unpinned devDependency silently breaks hooks and CI after an update.
2. **Always add `--no-errors-on-unmatched` to hook commands** — otherwise a commit that touches no matching files (e.g. only markdown) fails with an error.
3. **Use `--files-ignore-unknown=true` instead of (or in addition to) a fixed glob list** in hooks — the hook then keeps working when Biome adds support for new file types.
4. **Use `biome check --write` in pre-commit** to format + lint + organize imports and apply safe fixes in one pass; use read-only `biome ci` (or `biome check` without `--write`) in CI so CI catches what hooks auto-fixed.
5. **Use `--staged` locally and `--changed`/`--since` in CI** — with VCS integration enabled, `biome ci` on a remote repo uses `--changed` automatically; `--staged` is not available on `ci`.
6. **Two-step safety with git-format-staged or lint-staged** — prefer tools that re-stage fixed files (`stage_fixed: true` in Lefthook, `git update-index --again` in shell hooks) so the commit contains the formatted result.
7. **Set `vcs.useIgnoreFile: true`** so `.gitignore`d files (build outputs, generated code) are never processed by hooks.
8. **In the `pre-commit` framework, always set `additional_dependencies`** to a concrete `@biomejs/biome` version, or use a `repo: local` system hook reusing the project's npm install to avoid version drift.
9. **Use force-ignore `!!**/dist` (or `!!**/build`) for output directories** in `files.includes` when project-domain rules are enabled, so the scanner doesn't index them.
10. **In CI, prefer `biome ci`** — native GitHub annotations/GitLab code-quality reports, no write access, thread control.

## Common Pitfalls

- **`--changed`/`--staged` only see the file list, not the diff content** — even adding a space or newline marks a file as "changed", so the whole file is checked, not just the changed lines.
- **`--changed` won't catch downstream errors** — if you change an exported type, files importing it are not checked; run a full `biome check`/`ci` periodically.
- **VCS integration enabled but no features opted in** — `vcs.enabled` alone does nothing; combine with `useIgnoreFile` and/or `--staged`/`--changed` flags.
- **Husky alone can't list staged files** — it doesn't hide unstaged changes, so combine with lint-staged or git-format-staged; otherwise Biome checks the working tree, not the index.
- **Lefthook < v2.1.6 shows raw ANSI escape sequences** — upgrade Lefthook, configure its colors, set `NO_COLOR`, or pass `--colors=off` to Biome.
- **Shell-expanded globs in CLI arguments** — shells may not support `**`; put globs in `biome.json` instead.
- **Staged file with unstaged edits** (`MM` state) — `biome check --write --staged` may rewrite the staged content from a mix of states; the recommended shell hook fails explicitly on `MM` files.
- **`files.includes` vs `linter.includes` overlap** — `linter.includes`/`formatter.includes`/`assist.includes` are applied after `files.includes`; files excluded by `files.includes` can never be re-included.
- **Ignoring directories** — `!test` (no `/**`) works only in `files.includes`; in other `includes` lists you must write `!/test/**`.
- **Scanner indexing** — with project-domain rules enabled, the scanner indexes imported files even if ignored in `files.includes`; use force-ignore `!!` for explicit exclusion.
- **No trailing newline** (`formatter.trailingNewline: false`) breaks other tools and git — highly discouraged by the docs.

## Version Notes

- Docs are for **Biome v2.x** (docs site has separate v1.x / v2.x / next versions). v2 changed configuration semantics (e.g. `extends` accepts `"//"` for monorepos, `linter.rules.preset` replaces `recommended`, scanner with force-ignore `!!` patterns, HTML formatter still experimental/disabled by default).
- `biome migrate` (optionally `biome migrate prettier` / `biome migrate eslint --include-inspired --include-nursery`) updates configuration on breaking changes — run it after upgrading Biome.
- Version pinning is required for the `pre-commit` framework hooks (`additional_dependencies`) and recommended for npm installs (`-E`).
- `files.experimentalScannerIgnores` is **deprecated** — use force-ignore syntax (`!!`) in `files.includes`.
- `linter.rules.recommended` is **deprecated** — use `linter.rules.preset: "recommended"`.
