#!/usr/bin/env bash
# Unified colored verdict printer — VEDO EduTrack
#
# Usage: verdict.sh <STATUS> <message...>
#   STATUS: PASS | FAIL | WARN | SKIP | INFO | OK | ERROR
#
# Prints a single, well-distinguishable verdict line with light indication:
#   PASS/OK     green
#   FAIL/ERROR  red
#   WARN        yellow
#   SKIP        cyan
#   INFO        blue
#
# Colors auto-disable when stdout is not a TTY or NO_COLOR is set (CI-safe).
# Used by the Makefile and the gate runner (deploy/ci/run-gates.sh) so every
# command/test result speaks the same visual language.
set -u

STATUS="${1:?verdict.sh: status required (PASS|FAIL|WARN|SKIP|INFO|OK|ERROR)}"
shift
MSG="$*"

case "$STATUS" in
  PASS | OK)     CODE="\033[1;32m" ;; # green
  FAIL | ERROR)  CODE="\033[1;31m" ;; # red
  WARN)          CODE="\033[1;33m" ;; # yellow
  SKIP)          CODE="\033[1;36m" ;; # cyan
  INFO)          CODE="\033[1;34m" ;; # blue
  *)             CODE="\033[1;37m" ;; # white (fallback)
esac

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  printf "  [%b%s%b] %s\n" "$CODE" "$STATUS" "\033[0m" "$MSG"
else
  printf "  [%s] %s\n" "$STATUS" "$MSG"
fi
