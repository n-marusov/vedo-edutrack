#!/usr/bin/env bash
# Docker build gate — build the backend distroless image and check its size.
# Usage: docker-build-check.sh [--max-mb <mb>] [--tag <tag>]
# Enforces the image-size NFR (distroless backend <= 20 MB, development-tooling §8).
set -euo pipefail

MAX_MB=20
TAG="vedo-edutrack:gate-check"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --max-mb) MAX_MB="$2"; shift 2 ;;
    --tag) TAG="$2"; shift 2 ;;
    *) echo "docker-build-check: unknown arg: $1" >&2; exit 2 ;;
  esac
done

if ! command -v docker >/dev/null 2>&1; then
  echo "docker-build-check: docker not installed — skipped (CI installs it)"
  exit 0
fi

echo "docker-build-check: building backend image (tag=$TAG)..."
# M0.3: the unified backend Dockerfile builds the SPA in a node stage and
# embeds it — the build context is the REPO ROOT (pnpm workspace + frontend/),
# not backend/.
docker build -q -t "$TAG" -f backend/Dockerfile . >/dev/null

SIZE_BYTES="$(docker image inspect --format='{{.Size}}' "$TAG")"
SIZE_MB="$(awk -v b="$SIZE_BYTES" 'BEGIN { printf "%.1f", b / 1048576 }')"
MAX_BYTES="$(awk -v m="$MAX_MB" 'BEGIN { printf "%d", m * 1048576 }')"

echo "docker-build-check: image size ${SIZE_MB} MB (max ${MAX_MB} MB)"
docker rmi "$TAG" >/dev/null 2>&1 || true

if [[ "$SIZE_BYTES" -gt "$MAX_BYTES" ]]; then
  echo "docker-build-check: FAIL — image ${SIZE_MB} MB exceeds ${MAX_MB} MB" >&2
  exit 1
fi
echo "docker-build-check: OK"
