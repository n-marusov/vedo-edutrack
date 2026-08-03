#!/usr/bin/env bash
# Atlas migration validate gate — drift = 0 (REQ-NFR-ops.release.deployment-verification).
# Trivially passes when no migrations exist yet (migrations/ is empty at M0.2).
# Usage: atlas-migrate-validate.sh
set -euo pipefail

if ls backend/migrations/*.sql >/dev/null 2>&1; then
  if ! command -v atlas >/dev/null 2>&1; then
    echo "atlas-migrate-validate: atlas not installed — skipped (CI installs it)"
    exit 0
  fi
  (cd backend && atlas migrate validate --dir "file://migrations")
else
  echo "atlas-migrate-validate: no migration files yet — drift = 0 trivially (M0.2)"
fi
