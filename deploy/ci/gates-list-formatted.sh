#!/usr/bin/env bash
# Pretty-print `run-gates.sh --list` output grouped by category with counts.
# Usage: run-gates.sh --tier delivery --list | gates-list-formatted.sh
set -euo pipefail

awk '
  /^Gates selected/ { printf "%s\n\n", $0; next }
  /^  / {
    id = $1; grp = ""; sev = ""; tool = ""
    for (i = 2; i <= NF; i++) {
      if ($i ~ /^group=/)   grp = substr($i, 7)
      if ($i ~ /^tool=/)    tool = substr($i, 6)
      if ($i ~ /^blocking/ || $i ~ /^advisory/) sev = $i
    }
    extra = (tool ? sprintf(" (tool: %s)", tool) : "")
    gates[grp] = gates[grp] sprintf("    %-24s  %-8s%s\n", id, sev, extra)
    cnt[grp]++
    total++
    next
  }
  END {
    n = split("build typecheck lint test gen validate db security coverage", order)
    for (i = 1; i <= n; i++) {
      g = order[i]
      if (cnt[g] > 0) {
        printf "  \033[1;36m%s\033[0m (%d)\n", g, cnt[g]
        printf "%s", gates[g]
      }
    }
    printf "\n  \033[1mtotal: %d gates\033[0m\n", total
  }
' "${1:-/dev/stdin}"
