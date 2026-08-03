#!/usr/bin/env bash
# Secrets detection gate — check that no secrets are accidentally committed.
#
# Lightweight grep-based scan (no external tools). Complements gitleaks
# (which needs the binary installed — often skipped locally).
#
# Patterns checked:
#   - Private keys (RSA, EC, Ed25519, OpenSSH, PGP)
#   - AWS access keys (AKIA*)
#   - Generic "password=" / "secret=" assignments
#   - .env files outside .gitignore exemptions
#
# Usage: bash deploy/ci/secrets-check.sh
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
FAIL=0

echo "=== secrets-check: scanning for committed secrets ==="

# ------------------------------------------------------------------
# 1. Private keys — never commit these.
# ------------------------------------------------------------------
# The scanner itself declares the patterns below, so its own source file
# must be excluded — otherwise the gate false-positives on self-scanning.
SELF_EXCLUDE=':!deploy/ci/secrets-check.sh'
PRIVATE_KEY_PATTERNS=(
  "BEGIN RSA PRIVATE KEY"
  "BEGIN EC PRIVATE KEY"
  "BEGIN OPENSSH PRIVATE KEY"
  "BEGIN PRIVATE KEY"
  "BEGIN DSA PRIVATE KEY"
  "BEGIN PGP PRIVATE KEY BLOCK"
)

echo ""
echo "[1/4] private keys..."
for pattern in "${PRIVATE_KEY_PATTERNS[@]}"; do
  hits="$(git grep -n "$pattern" -- "$SELF_EXCLUDE" || true)"
  if [[ -n "$hits" ]]; then
    echo "  [FAIL] $pattern found:"
    echo "$hits" | sed 's/^/    /'
    FAIL=1
  fi
done

# ------------------------------------------------------------------
# 2. AWS access keys — AKIA prefix is AWS IAM, should never be in code.
# ------------------------------------------------------------------
echo ""
echo "[2/4] AWS access keys..."
AWS_HITS="$(git grep -nP 'AKIA[0-9A-Z]{16}' -- "$SELF_EXCLUDE" || true)"
if [[ -n "$AWS_HITS" ]]; then
  echo "  [FAIL] AWS access keys found:"
  echo "$AWS_HITS" | sed 's/^/    /'
  FAIL=1
fi

# ------------------------------------------------------------------
# 3. Hardcoded credentials — password=/secret=/token= in source files.
# ------------------------------------------------------------------
echo ""
echo "[3/4] hardcoded credentials..."
CRED_HITS="$(git grep -n -iP '(password|secret|token|api_key|apikey)\s*[:=]\s*["\x27][^"'"'"'\x27]{8,}' -- ':!*.md' ':!*.yaml' ':!*.yml' ':!*.json' ':!*.lock' ':!specs/' ':!pnpm-lock.yaml' ':!go.sum' ':!vendor/' "$SELF_EXCLUDE" || true)"
if [[ -n "$CRED_HITS" ]]; then
  echo "  [FAIL] hardcoded credentials found:"
  echo "$CRED_HITS" | sed 's/^/    /'
  FAIL=1
fi

# ------------------------------------------------------------------
# 4. .env files committed — should only be .env.example in the repo.
# ------------------------------------------------------------------
echo ""
echo "[4/4] .env files..."
ENV_HITS="$(git ls-files '*.env' '.env*' ':!:.env.example' ':!:deploy/.env.example' ':!:**/.env.example' ':!:**/.env.test' ':!:**/.env.dev' 2>/dev/null || true)"
if [[ -n "$ENV_HITS" ]]; then
  echo "  [FAIL] .env files committed (only .env.example should be in the repo):"
  echo "$ENV_HITS" | sed 's/^/    /'
  FAIL=1
fi

# ------------------------------------------------------------------
# Summary
# ------------------------------------------------------------------
echo ""
if [[ "$FAIL" -eq 0 ]]; then
  echo "[PASS] no secrets detected"
  exit 0
else
  echo "[FAIL] secrets detected — see above"
  exit 1
fi
