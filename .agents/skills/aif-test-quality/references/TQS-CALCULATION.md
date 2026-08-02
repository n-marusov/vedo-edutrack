# TQS Calculation Reference

Complete scoring formula and criteria for the Test Quality Score (TQS 0–10).

## Formula

```
TQS = Σ(group_score × group_weight)

Where:
  group_score = Σ(criterion_score) / criterion_count × 10
  (each criterion is 0.0–1.0, group_score is 0–10)

  TQS_final = clamp(TQS × multiplier, 0, 10)
  multiplier = 1.0 (default)
  multiplier = 0.0 (if any B1-B7 block triggers)
```

## Group 1: Structural Integrity (weight: 0.30)

| # | Criterion | Score Logic | Max |
|---|-----------|-------------|-----|
| 1.1 | File location | Correct test file naming for language → 1.0; wrong location → 0.0 | 1.0 |
| 1.2 | BDD naming | Function/test name follows `[Condition]_[Action]_[ExpectedResult]` → 1.0; partial → 0.5; none → 0.0 | 1.0 |
| 1.3 | Module/package structure | Tests organized in logical groups → 1.0; flat/unstructured → 0.3 | 1.0 |
| 1.4 | Test class/file purpose | Single test file per class/feature → 1.0; monolithic test file → 0.3 | 1.0 |

**Group score** = average(1.1, 1.2, 1.3, 1.4) × 10

## Group 2: Dependencies (weight: 0.25)

| # | Criterion | Score Logic | Max |
|---|-----------|-------------|-----|
| 2.1 | Mock count | 0–2 mocks → 1.0; 3–5 → 0.7; 6+ → 0.3 | 1.0 |
| 2.2 | Fixture usage | Shared fixtures/helpers → 1.0; duplicated setup → 0.3 | 1.0 |
| 2.3 | Test isolation | No shared mutable state → 1.0; shared state detected → 0.0 | 1.0 |
| 2.4 | External dependency | No real network/DB calls (mocked) → 1.0; real calls → 0.0 | 1.0 |
| 2.5 | Location sanity | Test file next to source or in tests/ dir → 1.0; misplaced → 0.0 | 1.0 |

**Group score** = average(2.1–2.5) × 10

## Group 3: Readability (weight: 0.20)

| # | Criterion | Score Logic | Max |
|---|-----------|-------------|-----|
| 3.1 | given/when/then structure | Clear arrange-act-assert with comments → 1.0; implicit AAA → 0.5; no structure → 0.0 | 1.0 |
| 3.2 | Descriptive comments | Comments explain intent (not just "test X") → 1.0; some comments → 0.5; none → 0.0 | 1.0 |
| 3.3 | Test length | < 50 lines → 1.0; 50–100 → 0.7; 100+ → 0.3 | 1.0 |
| 3.4 | Parameterized tests | Uses table-driven/subtest where applicable → 1.0; repetitive tests → 0.3 | 1.0 |

**Group score** = average(3.1–3.4) × 10

## Group 4: Safety (weight: 0.15)

| # | Criterion | Score Logic | Max |
|---|-----------|-------------|-----|
| 4.1 | No anti-patterns | Zero B1-B7 violations → 1.0; any violation → 0.0 (also triggers BLOCK) | 1.0 |
| 4.2 | No sleep/async waits | No `sleep` or `time.Sleep` → 1.0; present → 0.0 | 1.0 |
| 4.3 | No flaky markers | No `#[ignore]`, `t.Skip`, `@skip` without reason → 1.0; unexplained skips → 0.5 | 1.0 |
| 4.4 | Error handling | Tests check error returns → 1.0; ignore errors → 0.3 | 1.0 |

**Group score** = average(4.1–4.4) × 10

## Group 5: Coverage (weight: 0.10)

| # | Criterion | Score Logic | Max |
|---|-----------|-------------|-----|
| 5.1 | Assertions per test | ≥ 1 meaningful assertion → 1.0; 0 → triggers B5 | 1.0 |
| 5.2 | Multiple scenarios | Happy path + error path → 1.0; happy path only → 0.5 | 1.0 |
| 5.3 | Edge cases | Boundary values tested → 1.0; no edge cases → 0.3 | 1.0 |
| 5.4 | Negative tests | Invalid input/error cases covered → 1.0; none → 0.5 | 1.0 |

**Group score** = average(5.1–5.4) × 10

## Grade Thresholds

| Grade | TQS Range | Meaning |
|-------|-----------|---------|
| 🟢 gold | ≥ 9.5 | Excellent quality — exemplary tests |
| 🔵 silver | ≥ 8.0 | Good quality — meets high standards |
| 🟡 bronze | ≥ 6.0 | Acceptable — minimum for production |
| 🔴 fail | < 6.0 | Below minimum — needs improvement |
| ⚪ na | N/A | No tests found |

## Service-Level TQS

```
service_TQS = Σ(file_TQS × file_weight) / Σ(file_weight)

Where:
  file_weight = line_count (or 1.0 if equal weighting)
  If only 1 test file: service_TQS = file_TQS
```

## Example: api-gateway

Based on 6 test files found in `src/services/api-gateway/`:

| File | Structural | Dependencies | Readability | Safety | Coverage | File TQS |
|------|-----------|-------------|-------------|--------|----------|----------|
| main_test.go | 8.5 | 7.0 | 7.5 | 10.0 | 6.0 | 7.88 |
| helpers_test.go | 7.0 | 8.0 | 7.0 | 10.0 | 5.0 | 7.48 |
| auth_test.go | 9.0 | 7.5 | 8.0 | 10.0 | 7.0 | 8.38 |
| routes_test.go | 8.0 | 6.0 | 6.5 | 10.0 | 6.0 | 7.28 |
| middleware_test.go | 7.5 | 7.0 | 7.0 | 10.0 | 5.5 | 7.33 |
| grpc_test.go | 8.0 | 6.5 | 6.0 | 10.0 | 5.0 | 7.08 |

**Service TQS** = weighted average = **7.25** → 🔵 silver
