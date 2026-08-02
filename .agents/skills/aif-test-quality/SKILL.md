---
name: aif-test-quality
description: >-
  Analyzes test quality using TQS (Test Quality Score 0-10) and RCS (Requirements Coverage Score 0-10).
  Detects anti-patterns (sleep, assertTrue(true), empty catch, tautologies), inline tests in source files,
  and validates traceability.ttl cross-references. Use when reviewing test quality, before PR merge,
  during $aif-verify, or when user says "test quality", "check tests", "TQS", "RCS", or "test review".
argument-hint: "[path|service-name] [--strict] [--ttl-only] [--tqs-only] [--inline-only]"
allowed-tools: Read Grep Glob Bash(find *) Bash(wc *) Bash(grep *)
disable-model-invocation: false
metadata:
  author: AI Factory
  version: "1.0"
  category: quality
---

# Test Quality Analyzer

Static analysis of test quality across the VEDO Hub polyglot codebase (Go, Rust, Python, TypeScript).
Produces two independent metrics: **TQS** (test writing quality) and **RCS** (requirements coverage).

> Detailed scoring formulas: [references/TQS-CALCULATION.md](references/TQS-CALCULATION.md)
> Anti-pattern catalog: [references/ANTIPATTERNS.md](references/ANTIPATTERNS.md)
> Inline test rules: [references/INLINE-TESTS.md](references/INLINE-TESTS.md)
> Traceability verification: [references/TRACEABILITY-VERIFICATION.md](references/TRACEABILITY-VERIFICATION.md)

## Quick Start

```
$aif-test-quality api-gateway       # TQS + RCS + inline scan for api-gateway service
$aif-test-quality ontology-service  # Full analysis
$aif-test-quality --tqs-only src/   # TQS only on src directory
$aif-test-quality --ttl-only        # Traceability cross-validation only
$aif-test-quality --inline-only src/services/ # Inline test scan only
```

## Workflow

### Step 0: Load Configuration

1. Read `.ai-factory/config.yaml` for test-quality thresholds:
   ```yaml
   # Expected structure under test-quality:
   test-quality:
     thresholds:
       bronze: 6.0
       silver: 8.0
       gold: 9.5
     rcs:
       p0-weight: 3
       p1-weight: 2
       p2-weight: 1
     paths:
       traceability: ".ai-factory/traceability/traceability.ttl"
       requirements: "specs/requirements"
       user-stories: "specs/user-stories"
   ```
2. If config missing, use defaults (bronze=6.0, silver=8.0, gold=9.5).
3. Determine target path from `$ARGUMENTS`. If empty, scan entire `src/services/`.

### Step 1: TQS Analysis (Test Quality Score)

For each test file found via glob patterns:

**Discovery:**
- Go: `*_test.go`
- Rust: `*_test.rs`, `tests/*.rs`
- Python: `test_*.py`, `*_test.py`, `tests/*.py`
- TypeScript: `*.spec.ts`, `*.test.ts`, `__tests__/*.ts`

**For each test file, evaluate 5 groups** (see [references/TQS-CALCULATION.md](references/TQS-CALCULATION.md)):

| Group | Weight | Criteria |
|-------|--------|----------|
| Structural Integrity | 0.30 | File location, naming, BDD conventions |
| Dependencies | 0.25 | Mock count, fixture quality, isolation |
| Readability | 0.20 | given/when/then, comments, test length |
| Safety | 0.15 | Anti-patterns absent, no flakiness |
| Coverage | 0.10 | Multiple scenarios, edge cases, assertions per test |

**Blocking rules** — any of these instantly set TQS to 0:

| ID | Rule | Languages |
|----|------|-----------|
| B1 | Tautology assertion (`assertTrue(true)`, `assert!(true)`) | All |
| B2 | `sleep`/`thread::sleep` in test body | All |
| B3 | Empty catch block (`except: pass`) | Python, Go, TS |
| B4 | `System.exit`/`os.Exit`/`process.exit` in test | All |
| B5 | Test function with zero assertions | All |
| B6 | Test code in production source file (non-idiomatic) | Go, TS, Python |
| B7 | `__main__` + assert (doctest-in-main anti-pattern) | Python |

**Per-file TQS** → weighted average across files → **service-level TQS**.

**Threshold gates:**
- TQS ≥ 9.5 → 🟢 **gold**
- TQS ≥ 8.0 → 🔵 **silver**
- TQS ≥ 6.0 → 🟡 **bronze**
- TQS < 6.0 → 🔴 **below bronze** (BLOCK in $aif-verify)
- No tests found → TQS = N/A (recorded, not BLOCKED)

### Step 2: Inline Test Detection

Scan source files (non-test files) for embedded test code:

**Language-specific rules:**

| Language | What to find | Verdict |
|----------|-------------|---------|
| Rust | `#[cfg(test)]` modules | ✅ **Idiomatic** — analyze as tests |
| Go | `func Test*` in non-`_test.go` | 🚫 **BLOCK** (B6) |
| Python | `doctest` (`>>>`) | ✅ **Allowed** — lightweight test |
| Python | `test_*` functions in non-test module | ⚠️ **WARN** — should be separate |
| Python | `__main__` + assert pattern | 🚫 **BLOCK** (B7) |
| TypeScript | `describe`/`it` in non-spec/test file | 🚫 **BLOCK** (B6) |

For Rust `#[cfg(test)]` modules: apply full TQS analysis to the inline module.
For Go/TS/Python B6/B7 violations: report as BLOCK, suggest extraction.

### Step 3: Traceability Cross-Validation

Verify consistency between `traceability.ttl` and filesystem artifacts.

See [references/TRACEABILITY-VERIFICATION.md](references/TRACEABILITY-VERIFICATION.md) for full procedure.

**Quick version:**

1. **Parse TTL** → extract all `vdo:validates` (test→requirement) and `vdo:verifiedBy` (requirement→test) triples
2. **Scan filesystem** → inventory all test files, REQ-*, US-*, UC-* artifacts
3. **Cross-validate:**
   - **TTL→FS:** TTL references a file that doesn't exist → stale reference
   - **FS→TTL:** Test file exists with `// Validates: REQ-...` annotation but no TTL entry → untracked
   - **Structural:** TTL declares entity type that doesn't match filesystem artifact
4. **Compute RCS** = weighted coverage: `(Σ matched_requirements × weight) / (Σ all_requirements × weight) × 10`

RCS weights: P0=3, P1=2, P2=1, P3=0 (unweighted).

### Step 4: Generate Report

Output a structured quality report:

```
## Test Quality Report: <service-name>

### TQS: <score> (<grade>) | RCS: <score>

#### Anti-Patterns Found
| File | Line | Pattern | Severity | Rule |
|------|------|---------|----------|------|

#### Inline Tests
| File | Language | Verdict | Details |
|------|----------|---------|---------|

#### Traceability Issues
| Type | Detail | Recommendation |
|------|--------|----------------|

#### Recommendations
1. ...
2. ...
```

Use the template at [templates/QUALITY-REPORT.md](templates/QUALITY-REPORT.md).

### Step 5: Gate Decision

Determine if tests pass the quality gate:

```
if TQS == N/A (no tests):
  → Record "no tests" — do NOT block (informational)
if TQS < bronze (6.0):
  → BLOCK — tests are below minimum quality
if any B1-B7 violation found:
  → BLOCK — blocking anti-pattern detected
if orphan P0 requirements (RCS gap on P0):
  → WARN — high-priority requirement has no test coverage
otherwise:
  → PASS
```

Machine-readable output (for CI integration):
```
<!-- test-quality-gate: PASS|BLOCK|WARN -->
<!-- tqs: <score> | grade: <gold|silver|bronze|fail|na> -->
<!-- rcs: <score> | orphan-p0: <count> -->
```

## Integration Points

### Called by $aif-verify (Step 1 — Test Quality Gate)

When `$aif-verify` invokes this skill, it runs as **Step 1** (before task audit, code quality, and consistency checks). This provides fail-fast: bad tests block immediately without wasting time on other checks.

Output from this skill feeds into `$aif-verify`'s Step 4 report as the "Test Quality" section.

### Referenced by $aif-plan

When planning features with tests, `$aif-plan` should:
- Add TQS gate to acceptance criteria (if tests are planned)
- Add "Update traceability.ttl" to implementation steps
- Reference `.ai-factory/rules/test-quality.md` for quality thresholds

### Guided by $aif-implement

When implementing tests, `$aif-implement` should:
- Use BDD naming: `[Condition]_[Action]_[ExpectedResult]`
- Add `// Validates: REQ-<id>` annotations
- Include given/when/then structure
- Avoid anti-patterns from [references/ANTIPATTERNS.md](references/ANTIPATTERNS.md)
- Update `traceability.ttl` after test creation

## Error Handling

- **No test files found:** Report TQS = N/A, continue with other checks
- **traceability.ttl missing or malformed:** Skip RCS, report WARNING, continue TQS
- **Unrecognizable language:** Skip file with INFO note
- **Service with 0 tests (e.g. frontend):** TQS = N/A, no BLOCK, record for awareness

## Context

- Run standalone: `$aif-test-quality <service>`
- Called as subagent: via `$aif-verify` or `$aif-implement`
- No `context: fork` needed — runs locally with file reads only
