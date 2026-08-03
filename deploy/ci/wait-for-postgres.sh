#!/usr/bin/env bash
# Wait for PostgreSQL readiness (CI integration tests / local compose).
# Usage: wait-for-postgres.sh [host] [port] [user] [db] [timeout_sec]
# See ADR-IMPL.PROCESS.development-tooling §7 (integration tests with testcontainers/service).
set -euo pipefail

HOST="${1:-localhost}"
PORT="${2:-5432}"
USER="${3:-edutrack}"
DB="${4:-edutrack}"
TIMEOUT="${5:-60}"

echo "wait-for-postgres: waiting for $HOST:$PORT (db=$DB user=$USER, timeout=${TIMEOUT}s)"

for i in $(seq 1 "$TIMEOUT"); do
  if (exec 3<>"/dev/tcp/$HOST/$PORT") 2>/dev/null; then
    exec 3>&- 3<&-
    echo "wait-for-postgres: postgres ready at $HOST:$PORT after ${i}s"
    exit 0
  fi
  sleep 1
done

echo "wait-for-postgres: postgres NOT ready at $HOST:$PORT after ${TIMEOUT}s" >&2
exit 1
