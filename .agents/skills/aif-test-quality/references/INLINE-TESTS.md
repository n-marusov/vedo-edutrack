# Inline Test Detection Reference

Rules for detecting test code embedded in non-test source files.

## Philosophy

Some languages have idiomatic patterns for inline tests (Rust `#[cfg(test)]`, Python `doctest`).
Others treat inline tests as anti-patterns (Go, TypeScript). This reference documents both.

## Language Rules

### Rust: `#[cfg(test)]` Modules — ✅ Idiomatic

Rust convention places unit tests directly in source files within `#[cfg(test)]` modules.
This is NOT an anti-pattern — analyze these modules as first-class tests.

**What to look for:**
```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn should_return_error_when_input_is_empty() {
        // ...
    }
}
```

**Analysis rules:**
1. Extract the `#[cfg(test)]` module contents
2. Apply full TQS scoring to extracted tests
3. Count these as part of the file's test suite
4. B6 does NOT apply to Rust — this is idiomatic

**Regex for detection:**
```
#\[cfg\(test\)\]\s*mod\s+tests\s*\{
```

### Go: `func Test*` in Non-`_test.go` — 🚫 BLOCK (B6)

Go strictly separates test and production code via file naming. Test code in production files is a bug.

**Detection:**
```
func\s+Test\w+\s*\(
```
in files NOT matching `*_test.go`

**Exception:** None. Go has no idiomatic inline test pattern.

**Fix:** Move test function to `<original>_test.go`.

### Python: Multiple Inline Patterns

#### `doctest` — ✅ Allowed

Python doctests are lightweight inline tests embedded in docstrings.

```python
def add(a, b):
    """
    >>> add(1, 2)
    3
    >>> add(-1, 1)
    0
    """
    return a + b
```

**Analysis:** Count as test coverage. Do not apply full TQS scoring (no fixtures, no structure).
Just verify the doctest assertions are meaningful (not `>>> True`).

**Detection:**
```
>>>\s+\S+
```
in docstring context (triple-quoted strings).

#### `test_*` functions in non-test module — ⚠️ WARN

```python
# In my_module.py (not test_my_module.py)
def test_something():
    assert my_function() == 42
```

**Verdict:** WARN — extract to `test_my_module.py`.
**Not BLOCK** — the code runs, it's just organizationally wrong.

#### `__main__` + assert — 🚫 BLOCK (B7)

```python
if __name__ == "__main__":
    assert func_a() == 1
    assert func_b() == 2
    print("OK")
```

**Verdict:** BLOCK — not a test framework, no reporting.

**Detection:**
```
if\s+__name__\s*==\s*['"]__main__['"]
```
followed by `assert` statements within the same block.

### TypeScript: `describe`/`it` in Non-Spec File — 🚫 BLOCK (B6)

TypeScript projects use `.spec.ts` or `.test.ts` files. Test code in other files is a bug.

**Detection:**
```
describe\(|it\(|test\(|expect\(
```
in files NOT matching `*.spec.ts`, `*.test.ts`, `__tests__/*.ts`

**Exception:** None for standard patterns.
**Note:** `.vue` files with `<script>` may contain `describe`/`it` — check if it's a Storybook story (allowed) or test code (BLOCK).

## Decision Matrix

| Language | Pattern | File | Verdict | Action |
|----------|---------|------|---------|--------|
| Rust | `#[cfg(test)]` | `src/*.rs` | ✅ Idiomatic | Analyze as tests |
| Rust | `#[cfg(test)]` | `tests/*.rs` | ✅ Standard | Standard test file |
| Go | `func Test*` | `*_test.go` | ✅ Standard | Standard test file |
| Go | `func Test*` | `*.go` (non-test) | 🚫 B6 | BLOCK, suggest extraction |
| Python | `>>>` doctest | `*.py` | ✅ Allowed | Count as coverage |
| Python | `test_*` | `test_*.py` | ✅ Standard | Standard test file |
| Python | `test_*` | `*.py` (non-test) | ⚠️ WARN | Suggest extraction |
| Python | `__main__`+assert | `*.py` | 🚫 B7 | BLOCK |
| TypeScript | `describe`/`it` | `*.spec.ts`/`*.test.ts` | ✅ Standard | Standard test file |
| TypeScript | `describe`/`it` | `*.ts` (non-spec) | 🚫 B6 | BLOCK, suggest extraction |

## Detection Strategy

### Phase 1: Standard Test File Discovery

Use language-specific glob patterns to find dedicated test files:
- Go: `*_test.go`
- Rust: `*_test.rs`, `tests/*.rs`
- Python: `test_*.py`, `*_test.py`, `tests/*.py`
- TS: `*.spec.ts`, `*.test.ts`, `__tests__/*.ts`

### Phase 2: Source File Scan

For all non-test source files, grep for inline test patterns:

```
# Rust — always grep (could be idiomatic)
grep -rn '#\[cfg(test)\]' src/ --include='*.rs'

# Go — grep for func Test in non-test files
grep -rn 'func Test' src/ --include='*.go' | grep -v '_test.go'

# Python — grep for doctest, test_*, __main__
grep -rn '>>> ' src/ --include='*.py'
grep -rn 'def test_' src/ --include='*.py' | grep -v 'test_'
grep -rn '__main__' src/ --include='*.py'

# TypeScript — grep for describe/it in non-spec files
grep -rn 'describe(\|it(' src/ --include='*.ts' | grep -v -E '\.spec\.|\.test\.'
```

### Phase 3: Classification

For each hit, classify per the decision matrix above.
Output inline test findings in the quality report.
