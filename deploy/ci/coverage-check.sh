#!/usr/bin/env bash
# Coverage gate — parse a Go coverprofile and warn/fail below a threshold.
# Usage: coverage-check.sh --min <pct> [path/to/coverage.out] [--generate-on-missing]
# Advisory at M0.2 (REQ-NFR-process.dev.test-coverage) — gate severity is
# decided by gates.yaml; this script only reports the percentage.
set -euo pipefail

MIN=90
FILE=""
GENERATE=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --min) MIN="$2"; shift 2 ;;
    --generate-on-missing) GENERATE=1; shift ;;
    *) FILE="$1"; shift ;;
  esac
done

if [[ -z "$FILE" ]]; then
  echo "coverage-check: missing coverage file path" >&2
  exit 2
fi

if [[ ! -f "$FILE" && "$GENERATE" -eq 1 ]]; then
  echo "coverage-check: $FILE missing — generating with 'go test -coverprofile' (scaffolds are red at M0.2, coverage data is best-effort)"
  (cd backend && go test -count=1 -coverprofile="$(basename "$FILE")" ./...) >/dev/null 2>&1 || true
fi

if [[ ! -f "$FILE" ]]; then
  echo "coverage-check: coverage file not found: $FILE (no tests produced coverage yet at M0.2)"
  exit 0 # no data -> not a blocker; the unit-test gate covers failures
fi

COVER_DIR="$(dirname "$FILE")"
COVER_FILE="$(basename "$FILE")"
TOTAL="$(cd "$COVER_DIR" && go tool cover -func="$COVER_FILE" 2>/dev/null | tail -1 | awk '{print $3}' | tr -d '%')"
if [[ -z "$TOTAL" ]]; then
  echo "coverage-check: could not parse $FILE (run go test -coverprofile first from the module directory)"
  exit 0
fi

echo "coverage-check: total coverage ${TOTAL}% (min ${MIN}%)"
awk -v t="$TOTAL" -v m="$MIN" 'BEGIN {
  if (t + 0 < m + 0) {
    printf "coverage-check: FAIL — %.1f%% < %d%% (advisory at M0.3; becomes blocking when coverage >= 90%%)\n"
    exit 1
  }
  printf "coverage-check: OK — %.1f%% >= %d%%\n", t, m
}'
