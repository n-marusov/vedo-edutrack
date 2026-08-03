#!/usr/bin/env bash
# Compose stack health gate — start the full dev stack and verify every service is running.
# Usage: bash deploy/ci/compose-health-check.sh
# See gates.yaml and ADR-IMPL.PROCESS.development-tooling §8.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
COMPOSE_FILE="${COMPOSE_FILE:-deploy/docker-compose.yml}"

echo "=== compose-health: starting dev stack (10 services) ==="
if ! docker compose -f "$ROOT/$COMPOSE_FILE" up -d --wait 2>&1; then
  echo "[FAIL] docker compose up -d --wait failed"
  docker compose -f "$ROOT/$COMPOSE_FILE" ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}" 2>&1 || true
  exit 1
fi

echo ""
echo "=== compose-health: service status ==="
docker compose -f "$ROOT/$COMPOSE_FILE" ps --format "table {{.Name}}\t{{.Status}}\t{{.Health}}"

echo ""
FAIL=0
# M0.3: frontend (Vite) is an optional dev profile (frontend-dev); the
# production SPA is served by the unified backend binary. hub-mock is part of
# the core stack.
SERVICES=(postgres backend hub-mock otel-collector prometheus loki tempo grafana traefik)

for svc in "${SERVICES[@]}"; do
  state="$(docker compose -f "$ROOT/$COMPOSE_FILE" ps --format '{{.State}}' "$svc" 2>/dev/null || echo "missing")"
  health="$(docker compose -f "$ROOT/$COMPOSE_FILE" ps --format '{{.Health}}' "$svc" 2>/dev/null || echo "")"

  case "$state" in
    running)
      if [[ -n "$health" && "$health" != "healthy" ]]; then
        echo "  [FAIL] $svc — running but unhealthy: $health"
        FAIL=1
      elif [[ -n "$health" && "$health" == "healthy" ]]; then
        echo "  [OK]   $svc — healthy"
      else
        echo "  [OK]   $svc — running"
      fi
      ;;
    *)
      echo "  [FAIL] $svc — state: ${state:-unknown}"
      FAIL=1
      ;;
  esac
done

if [[ "$FAIL" -ne 0 ]]; then
  echo ""
  echo "[FAIL] one or more services are not healthy"
  exit 1
fi

echo ""
echo "[PASS] all ${#SERVICES[@]} services are running"
