#!/usr/bin/env bash
# Unified colored target header — VEDO EduTrack
#
# Usage: header.sh <icon-key> <message...>
#
# Prints a bold cyan "▶ <icon> <message>" line (plain when non-TTY/NO_COLOR).
#
# IMPORTANT: emojis live inside this script, resolved from an ASCII key.
# Native GNU Make on Windows mangles non-ASCII characters that appear in
# recipe command lines (child-process arguments), while bytes emitted by
# child stdout pass through unchanged — so all decorative glyphs must come
# from scripts, never from Makefile recipe text.
#
# Colors auto-disable when stdout is not a TTY or NO_COLOR is set (CI-safe).
set -u

KEY="${1:?header.sh: icon key required (up|down|dev|build|test|test-e2e|lint|format|gen|dev-check|check|migrate|migrate-down|hooks|ci|docker|gates|clean)}"
shift
TEXT="$*"

case "$KEY" in
  up)           ICON="🚀" ;;
  down)         ICON="🛑" ;;
  dev)          ICON="🔄" ;;
  build)        ICON="🔨" ;;
  test)         ICON="🧪" ;;
  test-e2e)     ICON="🎯" ;;
  lint)         ICON="🧹" ;;
  format)       ICON="✨" ;;
  gen)          ICON="⚙️" ;;
  dev-check)    ICON="⚡" ;;
  check)        ICON="🏁" ;;
  migrate)      ICON="🗄️" ;;
  migrate-down) ICON="⏪" ;;
  hooks)        ICON="🪝" ;;
  ci)           ICON="🤖" ;;
  docker)       ICON="🐳" ;;
  gates)        ICON="🚦" ;;
  clean)        ICON="🗑️" ;;
  *)            ICON="❔" ;;
esac

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
  printf "\033[1;36m▶ %s %s\033[0m\n" "$ICON" "$TEXT"
else
  printf "▶ %s %s\n" "$ICON" "$TEXT"
fi
