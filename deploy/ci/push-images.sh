#!/usr/bin/env bash
# GHCR image push gate (ci-main trigger only).
# Builds backend image with :latest and :sha-<commit> tags and pushes to GHCR.
# Requires REGISTRY, GITHUB_REPOSITORY, GITHUB_SHA and a prior docker login.
# Usage: push-images.sh
set -euo pipefail

REGISTRY="${REGISTRY:-ghcr.io}"
REPO="${GITHUB_REPOSITORY:-}"
SHA="${GITHUB_SHA:-}"

if [[ -z "$REPO" || -z "$SHA" ]]; then
  echo "push-images: GITHUB_REPOSITORY/GITHUB_SHA not set — ci-main gate skipped locally"
  exit 0
fi

IMG="$REGISTRY/$REPO/vedo-edutrack"
SHORT_SHA="${SHA:0:12}"

echo "push-images: building and pushing $IMG (:latest, :sha-$SHORT_SHA)"
docker build -q -t "$IMG:latest" -t "$IMG:sha-$SHORT_SHA" -f backend/Dockerfile backend/
docker push "$IMG:latest"
docker push "$IMG:sha-$SHORT_SHA"
echo "push-images: done"
