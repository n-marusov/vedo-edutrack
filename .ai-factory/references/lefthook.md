# Lefthook Reference

> Source: lefthook.dev (/, /usage/, /configuration/, /examples/), github.com/evilmartians/lefthook (README.md, docs/configuration/, docs/usage/, docs/examples/, docs/installation/)
> Created: 2026-08-03
> Updated: 2026-08-03
> Version documented: lefthook 2.1.10 (latest at time of writing)

## Overview

Lefthook is a Git hooks manager written in Go, maintained by Evil Martians. It is designed to be **fast** (parallel command execution), **powerful** (fine control over which files are passed to commands, conditional skip/only rules, local config overrides, remotes), and **simple** (a single dependency-free binary that works in Node.js, Ruby, Python, Go and other projects).

How it works: you create a `lefthook.yml` config file and run `lefthook install`. Lefthook writes simple hook scripts into `.git/hooks/`; each hook script calls `lefthook run <hook-name>` on execution. The config file is re-read every time a hook runs, so **re-install is not needed after editing `lefthook.yml`** (only hook *wiring* changes require `lefthook install`).

Key differentiators vs Husky/lint-staged: no per-hook JS runtime dependency, parallel execution by default opt-in, built-in file templates (`{staged_files}` etc.), glob/regex filters, and a `lefthook-local.yml` personal override mechanism.

## Core Concepts

- **Hook**: a Git hook name (`pre-commit`, `commit-msg`, `pre-push`, `post-merge`, ...) or any custom name. Custom names are runnable via `lefthook run <name>` but are not installed as Git hooks.
- **Config file**: `lefthook.yml` (or yaml/toml/json/jsonc variants). See Configuration section for accepted names.
- **Command** (`commands:`): a named map entry with a `run` string. Legacy but widely used.
- **Script** (`scripts:`): your own executable stored under `<source_dir>/<hook-name>/`, executed with a `runner`.
- **Job** (`jobs:`): unified task definition (added 1.10.0) that can be a `run` command or a `script`, optionally wrapped in a `group`. Modern, recommended way.
- **File templates**: placeholders in `run` substituted with file lists: `{files}`, `{staged_files}`, `{push_files}`, `{all_files}`, `{cmd}`, `{0}`..`{N}` (git args), `{lefthook_job_name}`.
- **Tags**: labels on commands/scripts used with `exclude_tags` and `LEFTHOOK_EXCLUDE` to filter groups of jobs.
- **skip / only**: conditional execution control (merge/rebase state, branch refs, custom commands).
- **stage_fixed**: auto `git add` of fixed files — works only for `pre-commit`.
- **parallel / piped**: concurrent execution vs stop-on-first-failure sequential pipeline.
- **extends / remotes / lefthook-local**: config merging hierarchy: `lefthook.yml` → `extends` → `remotes` → `lefthook-local.yml` (later overrides earlier).
- **Glob matcher**: gobwas (default; `**` matches 1+ dirs) or doublestar (standard behavior).

## API / Interface

### CLI commands

| Command | Purpose |
|---------|---------|
| `lefthook install [hook...]` | Create/update Git hooks from config; creates empty `lefthook.yml` if none exists. Optional hook names install only those hooks. `-f/--force` re-installs (needed after `rc` changes). `-d/--dirs` variant via `add`. |
| `lefthook add <hook> [--dirs]` | Install a single hook; `--dirs` creates `.git/hooks/<hook>/` + `.lefthook/<hook>/` script dirs. |
| `lefthook run <hook>` | Execute configured commands/scripts for a hook (Git hooks call this implicitly). Options: `--job <name>`, `--tag <tag>` (repeatable), `--all-files`, `--file <path>` (repeatable; overrides file templates; `--all-files` ignored if `--file` given). |
| `lefthook validate` | Validate config against JSON schema from the lefthook GitHub repo. |
| `lefthook dump` | Print the fully merged config (main + extends + remotes + local). |
| `lefthook check-install` | Exit `0` if hooks installed & synchronized, `1` otherwise. |
| `lefthook uninstall` | Clear Git hooks installed by lefthook. |
| `lefthook self-update` | Update binary from GitHub releases (only for source/release-binary installs). |
| `lefthook version [--full]` | Print version; `--full` adds commit hash. |
| `lefthook completion <shell>` | (standard) shell completion. |

### Environment variables

| Variable | Effect |
|----------|--------|
| `LEFTHOOK=0` / `LEFTHOOK=false` | Disable lefthook for a git command: `LEFTHOOK=0 git commit`. |
| `LEFTHOOK=1` | Force hook install in NPM postinstall even when `CI=true`. |
| `CI=true` | With NPM package: skip auto-install in postinstall (most CIs set it). |
| `LEFTHOOK_BIN=<path>` | Use a specific lefthook binary instead of PATH/package detection. |
| `LEFTHOOK_CONFIG=<path>` | Override main config file (local config, extends, remotes still load). |
| `LEFTHOOK_EXCLUDE=<tags-or-names>` | Comma-separated tags/command names to skip (e.g. `ruby,security,lint`). |
| `LEFTHOOK_OUTPUT=<list>` | Override output verbosity; `false` disables all output except errors. |
| `LEFTHOOK_VERBOSE=1` | Verbose printing. |
| `NO_COLOR=true` | Disable colored output (lefthook + subcommands). |
| `CLICOLOR_FORCE=true` | Force colored output (lefthook + subcommands). |

## Usage Patterns

### Install (typical flows)

```bash
# Node.js (npm) — hooks auto-install via postinstall
npm install --save-dev lefthook
# pnpm requires approval:
#   pnpm-workspace.yaml: onlyBuiltDependencies: [lefthook]
#   and package.json: "pnpm": { "onlyBuiltDependencies": ["lefthook"] }
pnpm add -D lefthook

# Ruby
gem install lefthook   # or gem "lefthook", require: false in Gemfile

# Python
pipx install lefthook  # or: python -m pip install --user lefthook; uv add --dev lefthook

# Go (>= 1.26)
go install github.com/evilmartians/lefthook/v2@v2.1.10
go get -tool github.com/evilmartians/lefthook/v2   # as project tool

# OS packages: brew, winget, scoop, snap, apt (deb), yum/dnf (rpm), apk (alpine), pacman (arch), mise, devbox
```

NPM package variants: `lefthook` (single executable for your OS, recommended), legacy `@evilmartians/lefthook` (all OS executables) and `@evilmartians/lefthook-installer` (fetches on install) — legacy ones still maintained but will be shut down.

### Minimal setup

```bash
vim lefthook.yml    # write config
lefthook install    # wire hooks into .git/hooks
git add -A && git commit -m '...'
```

### Run linters on pre-commit (jobs style, recommended)

```yml
# lefthook.yml
pre-commit:
  parallel: true
  jobs:
    - name: stylelint
      run: yarn stylelint --fix '{staged_files}'
      glob: "*.css"
      stage_fixed: true
    - name: eslint
      run: yarn eslint --fix '{staged_files}'
      glob:
        - "*.ts"
        - "*.js"
        - "*.tsx"
        - "*.jsx"
      stage_fixed: true
```

### Commands style (classic)

```yml
pre-commit:
  commands:
    lint:
      glob: "*.{js,ts,jsx,tsx}"
      run: yarn eslint {staged_files}
```

### Scripts (own executables)

```bash
lefthook add -d pre-commit          # creates .lefthook/pre-commit/
# edit .lefthook/pre-commit/my-script.sh
```

```yml
pre-commit:
  scripts:
    "my-script.sh":
      runner: bash
```

`runner` is invoked as `<runner> <path-to-script>` (e.g. `ruby .lefthook/pre-commit/lint.rb`, `node`, `go run`). Commit-message template checker example:

```bash
# .lefthook/commit-msg/template_checker
INPUT_FILE=$1
START_LINE=`head -n1 $INPUT_FILE`
PATTERN="^(TICKET)-[[:digit:]]+: "
if ! [[ "$START_LINE" =~ $PATTERN ]]; then
  echo "Bad commit message, see example: TICKET-123: some text"
  exit 1
fi
```

```yml
commit-msg:
  scripts:
    "template_checker":
      runner: bash
```

### File templates in `run`

| Template | Meaning |
|----------|---------|
| `{files}` | Result of custom `files` command (job-level `files`). |
| `{staged_files}` | Files staged for the commit. |
| `{push_files}` | Files committed but not pushed. |
| `{all_files}` | All git-tracked files. |
| `{cmd}` | The full command string from config (for wrapping, e.g. `docker run ... {cmd}`). |
| `{0}` | All git hook args joined by space; `{1}`, `{2}`... individual args (e.g. commit-msg file is `{1}`). |
| `{lefthook_job_name}` | Current job/command/script name. |

```yml
pre-commit:
  jobs:
    - name: govet
      files: git ls-files -m
      glob: "*.go"
      run: go vet -- {files}
pre-push:
  jobs:
    - name: eslint
      glob: "*.{js,ts,jsx,tsx}"
      run: yarn eslint {push_files}
```

Quoting: `run: yarn eslint "{staged_files}"` double-quotes every file (Windows-friendly); `'{staged_files}'` single-quotes; unquoted files are quoted only where needed (spaces in names). Line-length safety: lefthook splits long file lists into several sequential runs to fit the OS command-line limit.

### Custom git task groups

```yml
fixer:
  jobs:
    - run: bundle exec rubocop --force-exclusion --safe-auto-correct -- {staged_files}
    - run: yarn eslint --fix {staged_files}
```
```bash
lefthook run fixer
```

### Docker wrapping (local override)

```yml
# lefthook-local.yml
pre-commit:
  jobs:
    - name: lint
      run: docker run -it --rm <container_id_or_name> {cmd}
```

### Git LFS

Lefthook runs LFS hooks internally for `post-checkout`, `post-commit`, `post-merge`, `pre-push`. Errors suppressed if LFS not required; disable with global `skip_lfs: true`; see LFS output with `LEFTHOOK_VERBOSE=1`.

### CI / committing formatted code

```yml
pre-commit:
  parallel: true
  fail_on_changes: "ci"   # fail in CI if files were modified, frictionless locally
  jobs:
    - run: yarn lint
      stage_fixed: true
```

`fail_on_changes` values: `never` (default) / `always` / `ci` / `non-ci`.

### Commitlint + commitizen example

```yml
prepare-commit-msg:
  commands:
    commitzen:
      interactive: true
      run: yarn run cz --hook
      env:
        LEFTHOOK: 0
commit-msg:
  commands:
    "lint commit message":
      run: yarn run commitlint --edit {1}
```

## Configuration

### Config file names (use ONE format per project)

| Format | Acceptable names |
|--------|------------------|
| YAML | `lefthook.yml`, `lefthook.yaml`, `.lefthook.yml`, `.lefthook.yaml`, `.config/lefthook.yml`, `.config/lefthook.yaml` |
| TOML | `lefthook.toml`, `.lefthook.toml`, `.config/lefthook.toml` |
| JSON | `lefthook.json`, `.lefthook.json`, `.config/lefthook.json` |
| JSONC | `lefthook.jsonc`, `.lefthook.jsonc`, `.config/lefthook.jsonc` |

If multiple config files exist, only one is used (unspecified which). A `lefthook-local` variant merges on top (`lefthook-local.yml`; if main config starts with a dot, local must too, e.g. `.lefthook-local.json`). `-local` can be used standalone — put `lefthook-local.yml` in global `.gitignore` for personal configs.

### Top-level options

| Option | Default | Description |
|--------|---------|-------------|
| `assert_lefthook_installed` | `false` | Fail (exit 1) if lefthook binary can't be found (PATH/node_modules/gem/...), so hooks never silently skip. |
| `colors` | `auto` | `true`/`false` or map of `cyan`, `gray`, `green`, `red`, `yellow` to ANSI or hex codes. Overridable via `--colors`. |
| `extends` | — | List of extra YAML files to merge (globs supported: `projects/*/specific-lefthook-config.yml`). |
| `lefthook` | `null` | Full path or `sh` command to run lefthook (e.g. `bundle exec lefthook`, `pnpm lefthook`). NOT merged from remotes/extends (security); merged from local. |
| `min_version` | — | e.g. `min_version: 1.1.3` — fail if binary older. |
| `no_tty` | `false` | Ignore `interactive` option; useful in CI. |
| `output` | all enabled | List of output sections; `output: false` prints only errors. |
| `rc` | — | `sh` script path sourced before hooks (e.g. `~/.lefthookrc`) to set PATH/ENV for GUI-run hooks (nvm/fnm/rbenv). Requires `lefthook install -f` after changes. |
| `remotes` | — | Remote git repos with shared configs to fetch & merge (see below). |
| `source_dir` | `.lefthook/` | Directory containing `<hook-name>/` subfolders with scripts. |
| `source_dir_local` | — | Local-only scripts dir. |
| `skip_lfs` | `false` | Skip internal Git LFS hook execution. |
| `templates` | — | Custom `{name}` replacements in `run` (since 1.10.8); overridable via local config. |
| `install_non_git_hooks` | — | Install non-Git hooks into `.git/hooks` (since 2.0.17; e.g. git-flow). |
| `glob_matcher` | `gobwas` | `gobwas` or `doublestar` (standard `**` semantics). |
| `{hook-name}` | — | Per-hook config (see Hook options). |

### Hook-level options (e.g. `pre-commit:`)

| Option | Default | Description |
|--------|---------|-------------|
| `files` (global) | — | Custom command for files list (hook-wide). |
| `parallel` | `false` | Run commands/scripts concurrently. |
| `piped` | `false` | Run sequentially, stop on first failure. Error if both `piped` and `parallel` true. |
| `follow` | `false` | Stream STDOUT of running commands live (avoid with `parallel`). |
| `fail_on_changes` | `never` | Exit non-zero if files modified: `never`/`always`/`ci`/`non-ci`. |
| `fail_on_changes_diff` | — | Diff-based variant of fail_on_changes. |
| `exclude_tags` | — | Tags or command names to exclude for the hook (overridable via `LEFTHOOK_EXCLUDE`). |
| `exclude` | — | Glob patterns excluded from file templates. |
| `skip` | — | Skip whole hook: `true`, or list of conditions (see skip section). |
| `only` | — | Run hook only when conditions met (opposite of skip; skip wins on conflict). |
| `jobs` | — | Job definitions (see below). |
| `commands` | — | Command map (legacy). |
| `scripts` | — | Script map. |

### Job options (and command/script options, minus a few)

| Option | Applies to | Description |
|--------|-----------|-------------|
| `name` | jobs | Job name (used for merging across extends/local and `--job`). |
| `run` | commands/jobs | Shell command with templates. Mandatory for commands. |
| `script` | scripts/jobs | Script filename under `<source_dir>/<hook>/`. |
| `runner` | scripts/jobs(script) | Executor for script, e.g. `bash`, `node`, `ruby`, `go run`. |
| `args` | scripts/jobs | Extra args appended to command (since 2.0.5); templates allowed; omits git-provided args unless `{0}` given. |
| `group` | jobs | Nested `jobs` with own `parallel`/`piped`; `env`, `root`, `glob`, `exclude` on the group are inherited. |
| `skip` | all | See skip section. |
| `only` | all | See only section. |
| `tags` | commands/jobs/scripts | Tag list for filtering (`exclude_tags`, `LEFTHOOK_EXCLUDE`, `run --tag`). |
| `glob` | all | Glob filter for file templates (list allowed since 1.10.10). Without files template, filters `{staged_files}`/`{push_files}` and skips if none left. |
| `files` | commands/jobs | Custom command producing the files for `{files}`; empty result skips execution. Overrides hook-level `files`. |
| `file_types` | all | Filter by type: `text`, `binary`, `executable`, `not executable`, `symlink`, `not symlink`, or MIME (`text/html`, `text/xml`, `text/javascript`, `text/x-php`, `text/x-lua`, `text/x-perl`, `text/x-python`, `text/x-shellscript`, `text/x-sh`, `application/json`, ...). AND logic for the first six, OR for MIME. |
| `env` | all | ENV vars for the command (e.g. `RAILS_ENV: test`); can extend `PATH: $PATH:/home/me/bin`. |
| `root` | commands/jobs | Change CWD for the command (trailing slash); also used to filter files for pre-commit/pre-push. Globs still relative to git root. |
| `exclude` | all | Globs excluded from file templates (affected by `glob_matcher`). |
| `fail_text` | all | Custom message shown on failure. |
| `stage_fixed` | all | Auto `git add` after run — **pre-commit only**. Uses `files` command result if set, else `{staged_files}`; glob/exclude filters applied. |
| `interactive` | all | Connect TTY to command (executed after non-interactive unless piped; ignored with `no_tty`; Linux/Unix /dev/tty). |
| `use_stdin` | all | Pass OS stdin to command (needed for `pre-push` scripts reading stdin; avoids hang from pseudo-TTY). Only one command gets stdin. |
| `priority` | commands/scripts | Order in sequential/piped mode; `0` = +Infinity (run last); 1..N run first. |

### `skip` / `only` conditions

Possible values (scalar or list):
- `rebase` — when in rebase git state
- `merge` — when in merge git state
- `merge-commit` — when HEAD commit is a merge commit
- `ref: main` / `ref: dev/*` — branch name or glob
- `run: test ${SKIP_ME} -eq 1` — custom sh command; skipped if exit code 0 (`only` runs if exit code 0)

`skip` takes precedence over `only`. Typical uses: skip tests during rebase/merge, skip a linter on `main` pushes, `run: "! which aiautocommit"` to run only if tool exists, `skip: true` in `lefthook-local.yml` to disable a shared command locally.

### `remotes` (shared configs)

```yml
remotes:
  - git_url: git@github.com:evilmartians/lefthook
    ref: v1.0.0                # optional branch/tag; if set once, keep it set
    configs:                   # default: [lefthook.yml], relative to remote root
      - examples/ruby-linter.yml
    refetch_frequency: 24h     # always | never | duration; `refetch: true` overrides
```

Remote scripts require the `scripts` folder at the remote repo root. `extends` inside remotes is relative to the remote root. Failed fetch → warning; falls back to last successful fetch or ignores the remote.

### `output` values

`meta`, `summary`, `empty_summary`, `success`, `failure`, `execution`, `execution_out`, `execution_info`, `skips` — all enabled by default. `LEFTHOOK_OUTPUT` env extends/overrides.

### Config merge order

`lefthook.yml` → `extends` → `remotes` → `lefthook-local.yml` (each level overrides the previous). Named jobs merge by name; unnamed jobs append. `lefthook dump` shows the effective merged config.

## Best Practices

1. **Pin lefthook versions.** Install with a pinned version in package.json / Gemfile / go.mod. For remote configs, pin `ref` to a tag and set `refetch_frequency` if the ref is mutable.
2. **Prefer `jobs` over `commands`/`scripts` maps** — it's the modern unified format, supports `group`, and merges cleanly across extends/local.
3. **Use `{staged_files}` + `glob` for speed** so linters only touch changed files; combine with `stage_fixed: true` + `fail_on_changes: ci` for auto-fix locally and strict CI.
4. **Use `lefthook-local.yml` for personal overrides** (skip slow tests, wrap commands in docker); commit nothing but keep the file in `.gitignore`/global gitignore.
5. **Enable `assert_lefthook_installed: true`** to fail loudly when the binary is missing instead of silently skipping checks.
6. **Set `LEFTHOOK=1` explicitly in CI** when you *want* hooks installed despite `CI=true`; otherwise CI's auto-set `CI=true` prevents postinstall hook install.
7. **Use `rc` (e.g. `~/.lefthookrc`) for GUI tools** (VSCode, etc.) that don't source your shell rc — fixes "executable not found" for nvm/fnm/rbenv-managed tools.
8. **Remember `**` semantics**: with the default gobwas matcher, `**/*.js` does NOT match `app.js` (only 1+ dirs deep). Use `glob_matcher: doublestar` when migrating from other tools.
9. **Quote file templates for Windows** (`"{staged_files}"`) to handle paths with spaces.
10. **Use `--all-files` / `--file` in `lefthook run`** to force/replace file lists during debugging.
11. **Order sequential steps with `priority`** (1..N; `0` runs last) instead of relying on YAML order in `commands` maps.
12. **Run `lefthook validate` + `lefthook dump`** after non-trivial config changes to catch schema errors and see the effective merged config.

## Common Pitfalls

- **Multiple config formats in one project** — lefthook picks one silently; use a single format.
- **pnpm + lefthook**: hooks won't install unless `lefthook` is in `onlyBuiltDependencies` (pnpm-workspace.yaml and/or package.json `pnpm.onlyBuiltDependencies`).
- **Hooks not re-installed after config change** — usually fine: config is read on every run. Re-run `lefthook install` only when hook *wiring* must change (e.g. new hooks or `rc`/`lefthook` option changes).
- **`rc` changes not applied** — must run `lefthook install -f` to regenerate hooks.
- **`stage_fixed` in non-pre-commit hooks silently does nothing** — it works only for `pre-commit`.
- **`piped: true` + `parallel: true` together** → error.
- **Long file lists** — OS command-line limits apply; lefthook splits into sequential runs automatically (usually just a note, not an error).
- **`pre-push` scripts reading stdin hang** — add `use_stdin: true`, otherwise the pseudo-TTY never closes stdin.
- **`skip` vs `only` conflicts** — `skip` always wins.
- **`root` doesn't affect glob matching** — globs are computed from the git repo root.
- **RuboCop `--force-exclusion` needed** when passing `{all_files}` — otherwise RuboCop's own `Exclude` settings are ignored.
- **Legacy NPM packages** (`@evilmartians/lefthook*`) are deprecated and will be shut down — use the plain `lefthook` package.
- **Interactive commands in CI** — use `no_tty: true` to ignore `interactive`, or they'll misbehave without a TTY.

## Version Notes

- **1.10.0**: `jobs` introduced (unified run/script tasks).
- **1.10.5**: `lefthook` config option (explicit binary/command).
- **1.10.8**: `templates` custom replacements.
- **1.10.10**: `glob` accepts a list of patterns.
- **2.0.5**: `args` option for scripts/commands.
- **2.0.17**: `install_non_git_hooks`.
- **2.1.10**: latest documented release; added experimental AI hooks.
- Go install requires Go >= 1.26 (for `go install` route); `min_version` config guards against too-old binaries.
- Legacy distributions (`@evilmartians/lefthook`, `@evilmartians/lefthook-installer`) still maintained but planned for shutdown.

## Related References

- See `biome-precommit.md` in this directory for Lefthook + Biome.js integration examples (pre-commit linter, `stage_fixed`).
- Official docs: https://lefthook.dev/ — Configuration reference: https://lefthook.dev/configuration/ — Usage: https://lefthook.dev/usage/ — GitHub: https://github.com/evilmartians/lefthook
