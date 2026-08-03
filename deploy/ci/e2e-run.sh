#!/usr/bin/env bash
# e2e-run.sh <gui|api> — run an E2E suite against the TEST stack with automatic
# stack lifecycle: bring the stack up (make test-up) if not running, tear it
# down (make test-down) when this script started it. Disables Playwright
# webServer auto-start via E2E_STACK_MANAGED=1 so the stack is managed here.
#
# The test stack (deploy/docker-compose.test.yml, ADR-IMPL.INFRA.dev-test-
# compose-separation) is a minimal postgres+backend+hub-mock+frontend set with
# an isolated compose project (vedo-edutrack-test) — no observability/traefik,
# and its volumes can never collide with the dev stack.
set -uo pipefail

SUITE="${1:?usage: e2e-run.sh <gui|api>}"
# Test stack host ports are offset +50000 from dev (see docker-compose.test.yml):
# backend 58080, frontend 55173, hub-mock 58081, postgres 55432 — dev and
# test stacks can run side by side without port clashes.
case "$SUITE" in
  gui) PROBE="http://localhost:55173" ;;
  api) PROBE="http://localhost:58080/healthz" ;;
  *) echo "e2e-run.sh: unknown suite '$SUITE' (want gui|api)" >&2; exit 2 ;;
esac

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

started=0
if ! curl -fsS -o /dev/null "$PROBE" 2>/dev/null; then
  echo "[e2e-run] $PROBE not responding — starting test stack"
  if ! (cd "$ROOT" && make test-up); then
    echo "[e2e-run] failed to start test stack" >&2
    exit 1
  fi
  started=1
fi

cleanup() {
  if [ "$started" -eq 1 ]; then
    echo "[e2e-run] tearing down test stack"
    (cd "$ROOT" && make test-down) >/dev/null 2>&1 || true
  else
    echo "[e2e-run] stack was already up — leaving it running"
  fi
}
trap cleanup EXIT

export E2E_STACK_MANAGED=1
# Point the Playwright suites at the test stack's offset host ports.
export E2E_API_URL="http://localhost:58080/api/v1"
export E2E_ROOT_URL="http://localhost:58080"
export E2E_BASE_URL="http://localhost:55173"
cd "$ROOT/tests/e2e/$SUITE" || exit 2
pnpm install --ignore-workspace >/dev/null 2>&1
pnpm test
