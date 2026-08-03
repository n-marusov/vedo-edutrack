#!/usr/bin/env bash
# Unified colored verdict printer — VEDO EduTrack
#
# Usage: verdict.sh <STATUS> <message...>
#   STATUS: PASS | FAIL | WARN | SKIP | INFO | OK | ERROR
#
# Prints a single, well-distinguishable verdict line with emoji + light indication:
#   PASS/OK     ✅ green
#   FAIL/ERROR  ❌ white-on-red badge (high contrast, unmissable)
#   WARN        ⚠️  yellow
#   SKIP        ⏭️  cyan
#   INFO        ℹ️  blue
#
# Colors auto-disable when stdout is not a TTY or NO_COLOR is set (CI-safe);
# emojis are kept in CI logs so verdicts stay scannable without color.
# Used by the Makefile and the gate runner (deploy/ci/run-gates.sh) so every
# command/test result speaks the same visual language.
set -u

STATUS="${1:?verdict.sh: status required (PASS|FAIL|WARN|SKIP|INFO|OK|ERROR)}"
shift
MSG="$*"

case "$STATUS" in
  PASS | OK)     ICON="✅"; CODE="\033[1;32m"   ;; # green
  FAIL | ERROR)  ICON="❌"; CODE="\033[1;37;41m" ;; # white-on-red badge
  WARN)          ICON="⚠️";  CODE="\033[1;33m"   ;; # yellow
  SKIP)          ICON="⏭️";  CODE="\033[1;36m"   ;; # cyan
  INFO)          ICON="ℹ️";  CODE="\033[1;34m"   ;; # blue
  *)             ICON="❔"; CODE="\033[1;37m"   ;; # white (fallback)
esac

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  printf "  %s [%b%s%b] %s\n" "$ICON" "$CODE" "$STATUS" "\033[0m" "$MSG"
else
  printf "  %s [%s] %s\n" "$ICON" "$STATUS" "$MSG"
fi
