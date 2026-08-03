# GNU Make on Windows: env-export parse-time failure AND broken `make help` output

**Date:** 2026-08-03
**Files:** Makefile (env-file export block ~83-92; help target ~120-133)
**Severity:** medium

## Problem

Two distinct Makefile bugs surfaced together:

1. `make help` (and every other target) printed a scary error before running:

   ```
   process_begin: CreateProcess(NULL, sed s/=.*// .env.dev, ...) failed.
   Makefile:78: pipe: No error
   ```

2. `make help` listed every target with the name **`Makefile`** instead of the
   actual target name:

   ```
   VEDO EduTrack - available targets:
     Makefile       Print available targets
     Makefile       Start the full dev stack (9 services)
   ```

## Root Cause

### Bug 1 — parse-time `$(shell sed)` spawn on Windows

The Makefile exported env-file variable names with:

```make
export $(shell sed 's/=.*//' $(ENV_FILE))
```

This runs EAGERLY at parse time (even for `make help`). On Windows GNU Make,
a `$(shell ...)` command without shell metacharacters is executed directly via
`CreateProcess` (no shell), so it depends on the binary being in the Windows
PATH of make's process — `sed` lives in Git's `usr/bin`, only reliably
reachable through bash.

### Bug 2 — `help` target greps `$(MAKEFILE_LIST)`, which now has 2 entries

The help recipe ran `grep -E '^[a-zA-Z0-9_%-]+:.*## ' $(MAKEFILE_LIST)`.
`MAKEFILE_LIST` became `Makefile deploy/.env.dev` because `-include $(ENV_FILE)`
appends the env file, and `deploy/.env.dev` now exists (committed env file from
the dev/test stack split). grep over multiple files prefixes each match with
the filename — `Makefile:help: ## ...` — and the recipe parsed it with
`IFS=':' read -r name desc`, so `name` took the FIRST colon field = `Makefile`.
Before M1 the root `.env.dev` did not exist, `MAKEFILE_LIST` had one entry,
grep emitted no prefix, and the names displayed correctly.

## Solution

### Bug 1 — zero-child-process env read

Replaced `$(shell sed ...)` with make-native `$(file < FILE)` + text functions:

```make
export $(foreach _w,$(file < $(ENV_FILE)),$(if $(findstring =,$(_w)),$(firstword $(subst =, ,$(_w)))))
```

- `$(file < ...)` reads the file with no child process — cannot fail with
  CreateProcess on any platform (GNU Make ≥ 4.0; project uses 4.4.1).
- `$(findstring =,...)` keeps only `KEY=VALUE` words; comments/blanks skipped.
- Same construct applied to the legacy `.env` block.

### Bug 2 — grep only the primary makefile

Replaced `$(MAKEFILE_LIST)` with `$(firstword $(MAKEFILE_LIST))` in the help
recipe (both TTY and non-TTY branches) so grep searches only the main
Makefile and never emits a `filename:` prefix.

## Prevention

- **Never use `$(shell sed/grep/...)` at Makefile parse time** — Windows GNU
  Make spawns shell commands directly (CreateProcess fast-path) when there are
  no shell metacharacters; the make process PATH (not bash's) decides success.
- **Prefer `$(file < ...)` for reading files in make** — pure make, identical
  on all platforms.
- **`grep $(MAKEFILE_LIST)` breaks with >1 makefile** (env files via
  `-include`): either restrict to `$(firstword $(MAKEFILE_LIST))` or use
  `grep -h` to suppress filename prefixes.
- **Regression gate:** `deploy/ci/makefile-health-check.sh` (fast tier gate
  `makefile-health`) runs `make help` and asserts: exit 0, no
  `process_begin`/`CreateProcess` errors, target names shown (no `Makefile`
  placeholder), key targets present.
- During review: check parse-time `$(shell ...)` and multi-file `grep` in
  Makefile recipes — both are portability smells on Windows.

## Tags

`#makefile` `#windows` `#gnu-make` `#env-files` `#shell-function` `#portability` `#help-target`
