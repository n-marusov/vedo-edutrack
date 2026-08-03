#!/usr/bin/env bash
# Two-tier quality-gate runner — VEDO EduTrack (T16)
#
# Single executable over deploy/ci/gates.yaml (single source of truth).
# Consumers (Makefile dev-check/check/ci, CI jobs, agent skills) call this
# script and never duplicate commands.
#
# Usage:
#   run-gates.sh [--tier fast|delivery] [--trigger local|ci|ci-main|precommit]
#                [--group <group>] [--out-format table|json] [--list]
#
# Decision rule: a task is ready for delivery iff `--tier delivery` completes
# with zero blocking fails.
#
# Exit codes: 0 = no blocking failures (advisory/skipped may be present);
#             1 = at least one blocking gate failed.
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GATES_YAML="${GATES_YAML:-$(dirname "${BASH_SOURCE[0]}")/gates.yaml}"

TIER="fast"
TRIGGER="local"
GROUP=""
OUT_FORMAT="table"

# ---- args ----
while [[ $# -gt 0 ]]; do
  case "$1" in
    --tier) TIER="$2"; shift 2 ;;
    --trigger) TRIGGER="$2"; shift 2 ;;
    --group) GROUP="$2"; shift 2 ;;
    --out-format) OUT_FORMAT="$2"; shift 2 ;;
    --list) OUT_FORMAT="list" ; shift ;;
    *) echo "run-gates.sh: unknown argument: $1" >&2; exit 2 ;;
  esac
done

if [[ "$TIER" != "fast" && "$TIER" != "delivery" ]]; then
  echo "run-gates.sh: --tier must be fast|delivery (got: $TIER)" >&2
  exit 2
fi

# ---- colored output (unified verdict scheme with scripts/verdict.sh) ----
# Colors auto-disable when stdout is not a TTY or NO_COLOR is set (CI-safe).
colorize_status() {
  local st="$1"
  [[ -t 1 ]] || { printf '%s' "$st"; return; }
  [[ -z "${NO_COLOR:-}" ]] || { printf '%s' "$st"; return; }
  case "$st" in
    pass) printf '\033[1;32m%s\033[0m' "$st" ;;
    fail) printf '\033[1;31m%s\033[0m' "$st" ;;
    warn) printf '\033[1;33m%s\033[0m' "$st" ;;
    skip) printf '\033[1;36m%s\033[0m' "$st" ;;
    *)    printf '%s' "$st" ;;
  esac
}

# ---- phase scoring (m0.0..m10, post-mvp) ----
phase_score() {
  case "$1" in
    m0.0) echo 0 ;; m0.1) echo 1 ;; m0.2) echo 2 ;; m0.3) echo 3 ;;
    m1) echo 10 ;; m2) echo 20 ;; m3) echo 30 ;; m4) echo 40 ;;
    m5) echo 50 ;; m6) echo 60 ;; m7) echo 70 ;; m8) echo 80 ;;
    m9) echo 90 ;; m10) echo 100 ;; post-mvp) echo 999 ;;
    *) echo 0 ;;
  esac
}

CUR_PHASE="$(awk -F': *' '/^current_phase:/ {gsub(/[ "\r]/, "", $2); print $2; exit}' "$GATES_YAML")"
CUR_PHASE="${CUR_PHASE:-m0.2}"
CUR_SCORE="$(phase_score "$CUR_PHASE")"

# ---- parse gates.yaml into records ----
# Output: TSV lines: id \t command \t tier \t phase \t severity \t group \t trigger \t needs \t tool \t runner \t nfr
# NOTE: bash `read` collapses consecutive IFS delimiters, so empty fields are
# emitted as the @@E@@ sentinel (values never contain it) and stripped after read.
parse_gates() {
  awk '
    BEGIN { n=0; in_gates=0 }
    /^gates:/ { in_gates=1; next }
    in_gates && /^  - id:/ {
      if (n>0) emit()
      n++
      id=$3; cmd=""; tier=""; phase=""; sev=""; grp=""; trig=""; needs=""; tool=""; runner=""; nfr=""
      next
    }
    in_gates && n>0 && /^    [a-zA-Z-]+:/ {
      key=$1; sub(/:$/, "", key)
      val=substr($0, index($0, ": ")+2)
      # strip inline YAML comments (# ...) from values
      sub(/[ \t]*#.*$/, "", val)
      sub(/[ \t]+$/, "", val)
      if (key=="command") cmd=val
      else if (key=="tier") tier=val
      else if (key=="phase") phase=val
      else if (key=="severity") sev=val
      else if (key=="group") grp=val
      else if (key=="trigger") trig=val
      else if (key=="needs") needs=val
      else if (key=="tool") tool=val
      else if (key=="runner") runner=val
      else if (key=="nfr") nfr=val
    }
    function f(v) { return (v == "" ? "@@E@@" : v) }
    function emit() {
      printf "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", f(id), f(cmd), f(tier), f(phase), f(sev), f(grp), f(trig), f(needs), f(tool), f(runner), f(nfr)
    }
    END { if (n>0) emit() }
  ' "$GATES_YAML"
}

# ---- selection ----
SELECTED=()
select_gates() {
  local line
  while IFS=$'\t' read -r id cmd tier phase sev grp trig needs tool runner nfr; do
    [[ -z "$id" ]] && continue
    # restore empty fields collapsed by bash read (consecutive IFS delimiters)
    trig="${trig//@@E@@/}"; needs="${needs//@@E@@/}"; tool="${tool//@@E@@/}"; runner="${runner//@@E@@/}"; nfr="${nfr//@@E@@/}"
    # tier: delivery includes fast
    if [[ "$TIER" == "fast" && "$tier" != "fast" ]]; then continue; fi
    # phase regression: gate.phase <= current_phase
    local gs
    gs="$(phase_score "$phase")"
    if [[ "$gs" -gt "$CUR_SCORE" ]]; then continue; fi
    # trigger: ci-main gates run only with --trigger ci-main; empty = all triggers
    if [[ "$trig" == "ci-main" && "$TRIGGER" != "ci-main" ]]; then continue; fi
    if [[ -n "$trig" && "$trig" != "$TRIGGER" && "$TRIGGER" != "local" ]]; then continue; fi
    # group filter
    if [[ -n "$GROUP" && "$grp" != "$GROUP" ]]; then continue; fi
    SELECTED+=("$id|$cmd|$tier|$phase|$sev|$grp|$trig|$needs|$tool|$runner|$nfr")
  done < <(parse_gates)
}

# ---- internal gates ----
unit_touched() {
  local base_ref="${CI_BASE_REF:-}"
  if [[ -z "$base_ref" ]]; then
    if git rev-parse --verify origin/main >/dev/null 2>&1; then base_ref="origin/main"
    elif git rev-parse --verify HEAD~1 >/dev/null 2>&1; then base_ref="HEAD~1"
    else base_ref="HEAD"; fi
  fi
  local changed
  changed="$(git diff --name-only "$base_ref"...HEAD 2>/dev/null || git diff --name-only HEAD~1 2>/dev/null || true)"
  local pkgs
  pkgs="$(printf '%s\n' "$changed" | grep '^backend/' | grep '\.go$' | grep -v '_test\.go$' \
          | sed -e 's|^backend/||' -e 's|/[^/]*\.go$||' | sort -u)"
  if [[ -z "$pkgs" ]]; then
    echo "no touched backend packages — gate skipped"
    return 0
  fi
  local pkg_args=()
  while IFS= read -r p; do pkg_args+=("./$p"); done <<< "$pkgs"
  (cd backend && go test -count=1 "${pkg_args[@]}")
}

# ---- execution ----
START_TS="$(date +%s)"
RESULTS=()       # "id|status|duration_ms|severity|group"
BLOCKERS=()
AFFECTED=()

run_one() {
  local rec="$1"
  IFS='|' read -r id cmd tier phase sev grp trig needs tool runner nfr <<< "$rec"

  local g_start g_end dur status=pass
  g_start="$(date +%s%N)"

  if [[ "$runner" == "agent" ]]; then
    status=skip
    echo "[agent] $id — declared gate ($cmd; executed by the agent via /aif-verify or /aif-test-quality)"
  elif [[ "$runner" == "internal" && "$id" == "unit-touched" ]]; then
    if unit_touched >/tmp/unit-touched.out 2>&1; then
      cat /tmp/unit-touched.out | sed 's/^/  /'
      status=pass
    else
      cat /tmp/unit-touched.out | sed 's/^/  /'
      status=fail
    fi
  elif [[ -n "$tool" ]] && ! command -v "$tool" >/dev/null 2>&1; then
    echo "[skip] $id — tool '$tool' not installed (local run; CI installs it)"
    status=skip
  else
    if bash -c "$cmd" >/tmp/gate.out 2>&1; then
      status=pass
    else
      status=fail
    fi
    if [[ -s /tmp/gate.out ]]; then sed 's/^/  /' /tmp/gate.out | tail -40; fi
  fi

  g_end="$(date +%s%N)"
  local dur_ms=$(( (g_end - g_start) / 1000000 ))

  # severity folding: advisory fail -> warn; skip -> warn
  local out_status="$status"
  if [[ "$status" == "fail" && "$sev" == "advisory" ]]; then out_status=warn; fi
  if [[ "$status" == "skip" ]]; then out_status=warn; fi

  RESULTS+=("$id|$out_status|$dur_ms|$sev|$grp")
  if [[ "$out_status" == "fail" ]]; then BLOCKERS+=("$id"); fi
  if [[ "$id" == "unit-touched" && -f /tmp/unit-touched.out ]]; then
    : # affected files tracked below
  fi
  printf '%-24s %-8s %-8s %s\n' "$id" "$(colorize_status "$out_status")" "${dur_ms}ms" "$sev/$grp"
}

# changed files for affected_files (best effort)
collect_affected() {
  git diff --name-only HEAD 2>/dev/null | head -20 || true
}

# ---- output ----
print_summary_table() {
  echo ""
  echo "=== Gate run summary ($TIER tier, phase $CUR_PHASE, trigger $TRIGGER) ==="
  printf '%-24s %-8s %s\n' "GATE" "STATUS" "SEVERITY/GROUP"
  printf '%s\n' "----------------------------------------------------------"
  local r
  for r in "${RESULTS[@]}"; do
    IFS='|' read -r id st dur sev grp <<< "$r"
    printf '%-24s %-8s %s\n' "$id" "$(colorize_status "$st")" "$sev/$grp (${dur}ms)"
  done
  echo ""
  if [[ ${#BLOCKERS[@]} -gt 0 ]]; then
    echo "BLOCKING FAILURES: $(colorize_status fail) — ${BLOCKERS[*]}"
  else
    echo "No blocking failures."
  fi
}

print_json() {
  local status=pass
  if [[ ${#BLOCKERS[@]} -gt 0 ]]; then status=fail; fi
  if [[ ${#RESULTS[@]} -eq 0 ]]; then status=warn; fi

  local gates_json="["
  local first=1
  local r
  for r in "${RESULTS[@]}"; do
    IFS='|' read -r id st dur sev grp <<< "$r"
    [[ $first -eq 0 ]] && gates_json+=","
    first=0
    gates_json+="{\"id\":\"$id\",\"status\":\"$st\",\"severity\":\"$sev\",\"group\":\"$grp\",\"duration_ms\":$dur}"
  done
  gates_json+="]"

  local blockers_json="["
  first=1
  if [[ ${#BLOCKERS[@]} -gt 0 ]]; then
    for b in "${BLOCKERS[@]}"; do
      [[ $first -eq 0 ]] && blockers_json+=","
      first=0
      blockers_json+="\"$b\""
    done
  fi
  blockers_json+="]"

  local affected_json="["
  first=1
  for a in $(collect_affected); do
    [[ $first -eq 0 ]] && affected_json+=","
    first=0
    affected_json+="\"$a\""
  done
  affected_json+="]"

  local suggested=""
  if [[ ${#BLOCKERS[@]} -gt 0 ]]; then
    suggested="fix blockers, then rerun: bash deploy/ci/run-gates.sh --tier $TIER"
    [[ -n "$GROUP" ]] && suggested+=" --group $GROUP"
  fi

  cat <<EOF
{
  "schema": "aif-gate-result",
  "schema_version": 1,
  "phase": "$CUR_PHASE",
  "tier": "$TIER",
  "trigger": "$TRIGGER",
  "group": ${GROUP:-null},
  "status": "$status",
  "blocking": $( [[ ${#BLOCKERS[@]} -gt 0 ]] && echo true || echo false ),
  "gates": $gates_json,
  "blockers": $blockers_json,
  "affected_files": $affected_json,
  "suggested_next": "$suggested"
}
EOF
}

# ---- main ----
select_gates

if [[ "$OUT_FORMAT" == "list" ]]; then
  echo "Gates selected for --tier $TIER --trigger $TRIGGER (phase <= $CUR_PHASE):"
  for rec in "${SELECTED[@]:-}"; do
    IFS='|' read -r id cmd tier phase sev grp trig needs tool runner nfr <<< "$rec"
    printf '  %-24s %-10s phase=%-6s %-8s group=%-12s' "$id" "$tier" "$phase" "$sev" "$grp"
    [[ -n "$tool" ]] && printf ' tool=%s' "$tool"
    printf '\n'
  done
  exit 0
fi

echo "=== Gate runner: tier=$TIER trigger=$TRIGGER phase=$CUR_PHASE group=${GROUP:-all} ==="
if [[ ${#SELECTED[@]} -eq 0 ]]; then
  echo "No gates selected."
  exit 0
fi

for rec in "${SELECTED[@]}"; do
  run_one "$rec"
done

END_TS="$(date +%s)"
TOTAL_MS=$(( (END_TS - START_TS) * 1000 ))

if [[ "$OUT_FORMAT" == "json" ]]; then
  print_json
else
  print_summary_table
  echo "Total duration: ${TOTAL_MS}ms"
fi

if [[ ${#BLOCKERS[@]} -gt 0 ]]; then
  exit 1
fi
exit 0
