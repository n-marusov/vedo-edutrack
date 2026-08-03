#!/usr/bin/env bash
# e2e-run.sh <gui|api> — run an E2E suite against the dev stack with automatic
# stack lifecycle: bring the stack up (make up) if not running, tear it down
# (make down) when this script started it. Disables Playwright webServer
# auto-start via E2E_STACK_MANAGED=1 so the stack is managed here, not by Playwright.
set -uo pipefail

SUITE="${1:?usage: e2e-run.sh <gui|api>}"
case "$SUITE" in
  gui) PROBE="http://localhost:5173" ;;
  api) PROBE="http://localhost:8080/healthz" ;;
  *) echo "e2e-run.sh: unknown suite '$SUITE' (want gui|api)" >&2; exit 2 ;;
esac

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

started=0
if ! curl -fsS -o /dev/null "$PROBE" 2>/dev/null; then
  echo "[e2e-run] $PROBE not responding — starting dev stack"
  if ! (cd "$ROOT" && make up); then
    echo "[e2e-run] failed to start dev stack" >&2
    exit 1
  fi
  started=1
fi

cleanup() {
  if [ "$started" -eq 1 ]; then
    echo "[e2e-run] tearing down dev stack"
    (cd "$ROOT" && make down) >/dev/null 2>&1 || true
  else
    echo "[e2e-run] stack was already up — leaving it running"
  fi
}
trap cleanup EXIT

export E2E_STACK_MANAGED=1
cd "$ROOT/tests/e2e/$SUITE" || exit 2
pnpm install --ignore-workspace >/dev/null 2>&1
pnpm test
