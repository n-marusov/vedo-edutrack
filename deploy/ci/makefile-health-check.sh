#!/usr/bin/env bash
# Makefile health gate — regression guard for two Windows GNU Make bugs:
#
#   Bug 1: parse-time `$(shell sed ...)` spawned sed via CreateProcess and
#          failed when sed wasn't in the make process PATH, polluting EVERY
#          make invocation (e.g. `make help`) with:
#            process_begin: CreateProcess(NULL, sed ..., ...) failed.
#            Makefile:N: pipe: No error
#   Bug 2: `help` grepped `$(MAKEFILE_LIST)` (now "Makefile deploy/.env.dev"
#          because -include appends the committed env file); grep over multiple
#          files prefixes matches with the filename, so `IFS=':' read` took
#          "Makefile" as every target name.
#
# Fixes: $(file < ...) instead of $(shell sed ...); help greps only
# $(firstword $(MAKEFILE_LIST)). This gate prevents regression on both.
#
# Usage: bash deploy/ci/makefile-health-check.sh
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
FAIL=0

echo "=== makefile-health: verifying 'make help' output ==="

cd "$ROOT" || exit 2

HELP_OUT="$(make help 2>&1)"
HELP_RC=$?

# 1. make help must succeed.
if [ "$HELP_RC" -ne 0 ]; then
  echo "  [FAIL] 'make help' exited $HELP_RC"
  FAIL=1
else
  echo "  [PASS] 'make help' exits 0"
fi

# 2. No parse-time spawn errors (Windows GNU Make $(shell) CreateProcess fast-path).
SPAWN_ERR="$(printf '%s\n' "$HELP_OUT" | grep -E 'process_begin|CreateProcess|pipe: No error' || true)"
if [ -n "$SPAWN_ERR" ]; then
  echo "  [FAIL] parse-time spawn error present in make output:"
  echo "$SPAWN_ERR" | sed 's/^/    /'
  FAIL=1
else
  echo "  [PASS] no process_begin/CreateProcess errors"
fi

# 3. Target names must display (not the makefile filename).
#    A broken help shows "Makefile" as every target name; a correct one shows
#    real target names like "help", "up", "test-up".
NAME_ERR="$(printf '%s\n' "$HELP_OUT" | grep -E '^  (Makefile|makefile)[[:space:]]' || true)"
if [ -n "$NAME_ERR" ]; then
  echo "  [FAIL] help lists 'Makefile' as target name (multi-file grep prefix bug):"
  echo "$NAME_ERR" | sed 's/^/    /'
  FAIL=1
else
  echo "  [PASS] target names displayed (no 'Makefile' placeholder)"
fi

# 4. Key documented targets must be present in the listing.
MISSING=0
for t in help up down test-up test-down dev build test bench lint format gen migrate clean; do
  if ! printf '%s\n' "$HELP_OUT" | grep -qE "^  $t[[:space:]]"; then
    echo "  [FAIL] target '$t' missing from 'make help'"
    MISSING=1
  fi
done
if [ "$MISSING" -eq 0 ]; then
  echo "  [PASS] all key targets listed (help/up/down/test-up/test-down/dev/build/test/bench/lint/format/gen/migrate/clean)"
else
  FAIL=1
fi

echo ""
if [ "$FAIL" -eq 0 ]; then
  echo "[PASS] Makefile health OK"
  exit 0
else
  echo "[FAIL] Makefile health checks failed — see above"
  exit 1
fi
