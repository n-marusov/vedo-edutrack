# Anti-Patterns Catalog

Complete catalog of test anti-patterns detected by $aif-test-quality, organized by severity.

## Blocking Rules (B1–B7)

Any blocking rule sets the **entire file's TQS to 0** and generates a BLOCK verdict.

### B1: Tautology Assertion

Assertion that always passes regardless of code behavior.

| Language | Pattern | Regex |
|----------|---------|-------|
| Go | `assert.True(t, true)` | `assert\.True\(t?,\s*true\)` |
| Go | `assert.False(t, false)` | `assert\.False\(t?,\s*false\)` |
| Go | `require.True(t, true)` | `require\.True\(t?,\s*true\)` |
| Rust | `assert!(true)` | `assert!\(true\)` |
| Rust | `assert_eq!(x, x)` (self-equality) | `assert_eq!\((\w+),\s*\1\)` |
| Python | `assert True` | `^\s*assert\s+True\s*$` |
| Python | `assert x == x` (self-equality) | `assert\s+(\w+)\s*==\s*\1\s*$` |
| TypeScript | `expect(true).toBe(true)` | `expect\(true\)\.toBe\(true\)` |
| TypeScript | `expect(x).toEqual(x)` (self-eq) | `expect\((\w+)\)\.toEqual\(\1\)` |

**Severity:** CRITICAL — test is meaningless.

### B2: Sleep in Test

Hard-coded sleep introduces flakiness and slows test suites.

| Language | Pattern | Regex |
|----------|---------|-------|
| Go | `time.Sleep(...)` | `time\.Sleep\(` |
| Rust | `std::thread::sleep(...)` | `std::thread::sleep\(` |
| Python | `time.sleep(...)` | `time\.sleep\(` |
| TypeScript | `setTimeout(...)` / `await new Promise` | `setTimeout\(` \| `new Promise.*resolve` |

**Severity:** HIGH — test is flaky and slow.

### B3: Empty Catch / Swallowed Error

Error handling that silently swallows failures.

| Language | Pattern | Regex |
|----------|---------|-------|
| Go | `if err != nil { }` (empty block) | `if\s+err\s*!=\s*nil\s*\{\s*\}` |
| Python | `except: pass` | `except.*:\s*pass` |
| Python | `except Exception: pass` | `except\s+Exception.*:\s*pass` |
| TypeScript | `catch(e) {}` (empty) | `catch\s*\(\w*\)\s*\{\s*\}` |

**Severity:** HIGH — errors are hidden, test passes incorrectly.

### B4: Process Exit in Test

Terminating the test process instead of reporting failure.

| Language | Pattern | Regex |
|----------|---------|-------|
| Go | `os.Exit(...)` | `os\.Exit\(` |
| Rust | `std::process::exit(...)` | `std::process::exit\(` |
| Python | `sys.exit(...)` | `sys\.exit\(` |
| TypeScript | `process.exit(...)` | `process\.exit\(` |

**Severity:** CRITICAL — kills test runner, masks failures.

### B5: Test Without Assertion

Test function that performs actions but never asserts the result.

**Detection heuristic:**
- Function matches `Test*` (Go/Rust) or `test_*` (Python) or `it()/test()` (TS)
- No `assert`, `require`, `expect`, `assert_eq`, `check` calls in the function body
- Function body is not empty (empty = different issue)

**Severity:** CRITICAL — test verifies nothing.

### B6: Test Code in Production File

Test logic embedded in non-test source files (non-idiomatic languages).

| Language | Rule | Exception |
|----------|------|-----------|
| Go | `func Test*` in non-`_test.go` file | None — always BLOCK |
| Rust | `#[cfg(test)]` module in source file | ✅ **Idiomatic** — NOT blocked |
| Python | `test_*` in module (not `test_*.py`) | WARN, not BLOCK |
| TypeScript | `describe()`/`it()` in non-spec file | None — always BLOCK |

**Severity:** CRITICAL for Go/TS — test code pollutes production bundle.

### B7: __main__ + Assert (Python Doctest Anti-Pattern)

Python file using `if __name__ == "__main__"` with assert statements as pseudo-tests.

```python
# BLOCKED pattern:
if __name__ == "__main__":
    assert some_function() == expected
    print("All tests passed")
```

**Severity:** HIGH — not a proper test framework, no reporting, no isolation.

## Warning Patterns (W1–W5)

Warnings reduce TQS but do NOT trigger BLOCK.

### W1: Excessive Mock Count

More than 5 mocks in a single test file suggest the code under test has too many dependencies (design issue).

**Detection:** Count `mock`, `Mock`, `mock_*`, `@patch`, `NewMock`, `vi.mock` calls per file.

### W2: Duplicated Test Setup

Multiple test functions with identical setup code instead of using shared fixtures.

**Detection:** Identical first 3+ lines in 3+ test functions within same file.

### W3: Overly Long Test Function

Single test function exceeding 100 lines.

**Detection:** Line count between `func Test*`/`def test_*`/`it(` and closing brace.

### W4: Missing Error Path Testing

All test functions follow happy-path only; no error/edge-case tests.

**Detection:** No test function name contains `error`, `fail`, `invalid`, `missing`, `null`, `empty`, `err`.

### W5: Unexplained Test Skip

`#[ignore]`, `t.Skip()`, `@skip`, `test.skip()` without a reason string.

**Detection:** Skip call without string argument.

## Polyglot Regex Quick Reference

```
# B1 — Tautologies
Go:     assert\.True\(t?,\s*true\)|require\.True\(t?,\s*true\)
Rust:   assert!\(true\)|assert_eq!\((\w+),\s*\1\)
Python: assert\s+True\s*$|assert\s+(\w+)\s*==\s*\1
TS:     expect\(true\)\.toBe\(true\)

# B2 — Sleep
Go:     time\.Sleep\(
Rust:   std::thread::sleep\(
Python: time\.sleep\(
TS:     setTimeout\(|new Promise.*resolve

# B3 — Empty catch
Go:     if\s+err\s*!=\s*nil\s*\{\s*\}
Python: except.*:\s*pass
TS:     catch\s*\(\w*\)\s*\{\s*\}

# B4 — Process exit
Go:     os\.Exit\(
Rust:   std::process::exit\(
Python: sys\.exit\(
TS:     process\.exit\(
```
