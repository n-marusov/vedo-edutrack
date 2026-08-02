# Traceability Verification Reference

Cross-validation between `traceability.ttl` and filesystem artifacts.

## Overview

The traceability graph (`.ai-factory/traceability/traceability.ttl`) is the **single source of truth** for test→requirement links.
This reference documents how to verify consistency and compute RCS.

## Phase 3A: Parse TTL

Parse `traceability.ttl` (RDF Turtle) to extract:

### Entities

```turtle
# Requirement declarations
vdo:REQ-FUN-API-create-class a vdo:Requirement ;
    vdo:priority "P0" ;
    vdo:area "API" .

# Test suite declarations
vdo:test-api-gateway a vdo:TestSuite ;
    vdo:service "api-gateway" .

# User story declarations
vdo:US-101 a vdo:UserStory ;
    vdo:priority "P0" .
```

### Validation Links

```turtle
# Test validates requirement
vdo:test-api-gateway vdo:validates vdo:REQ-FUN-API-create-class .

# Requirement verified by test
vdo:REQ-FUN-API-create-class vdo:verifiedBy vdo:test-api-gateway .
```

### Internal Model

Build an in-memory model:
```
requirements: Map<id, {priority, area, exists_on_fs}>
tests: Map<id, {service, file_path, exists_on_fs}>
validations: List<{test_id, requirement_id}>
user_stories: Map<id, {priority, exists_on_fs}>
```

## Phase 3B: Filesystem Scan

Inventory all artifacts on disk:

### Requirements
```
specs/requirements/REQ-*.md → extract ID from filename
```

### User Stories
```
specs/user-stories/US-*.md → extract ID from filename
```

### Use Cases
```
specs/use-cases/UC-*.md → extract ID from filename
```

### Test Files
```
src/services/*/tests/*_test.*
src/services/*/*_test.*
src/services/*/test_*
src/services/*/tests/test_*
tests/**/*.test.*
tests/**/*.spec.*
```

## Phase 3C: Cross-Validation

Three validation directions:

### 1. TTL→FS (Stale References)

TTL references an artifact that doesn't exist on the filesystem.

```
For each {test_id} in validations:
    if test file not found → STALE_REF
    log: "TTL references <test_id> but file <path> not found"

For each {req_id} in validations:
    if REQ file not found → STALE_REF
    log: "TTL references <req_id> but specs/requirements/<req_id>.md not found"
```

**Severity:** HIGH — TTL has phantom links.

### 2. FS→TTL (Untracked Artifacts)

Filesystem artifact exists but has no TTL entry.

**Test files:** Check for `// Validates: REQ-...` annotations in test source code.
```
grep -rn '// Validates:\|# Validates:' <test-file>
```

If annotation exists but no TTL entry → UNTRACKED_TEST

For requirements/user stories:
```
If REQ-*.md exists but no vdo:REQ-* in TTL → UNTRACKED_REQUIREMENT
If US-*.md exists but no vdo:US-* in TTL → UNTRACKED_USER_STORY
```

**Severity:** MEDIUM — missing traceability.

### 3. Structural Integrity

TTL entity type doesn't match filesystem artifact.

```
If TTL declares vdo:TestSuite for path X, but X is not a test file → TYPE_MISMATCH
If TTL declares vdo:Requirement for path X, but X is not a REQ file → TYPE_MISMATCH
```

**Severity:** LOW — metadata inconsistency.

## Phase 3D: RCS Calculation

Requirements Coverage Score (RCS 0–10).

### Algorithm

```
1. Collect all requirements from TTL + filesystem
2. For each requirement:
   - Determine priority (P0, P1, P2, P3) from TTL or REQ-*.md header
   - Check if ANY test has vdo:validates link to this requirement
   - Check if any test file has annotation referencing this requirement
3. Compute weighted coverage:

   rcs = (Σ(priority_weight × is_covered)) / (Σ(priority_weight)) × 10

   Where:
     P0 weight = 3  (critical — must have tests)
     P1 weight = 2  (high — should have tests)
     P2 weight = 1  (medium — nice to have)
     P3 weight = 0  (low — excluded from score)

   is_covered = 1 if at least one test covers this requirement, 0 otherwise
```

### RCS Grades

| Grade | RCS Range | Meaning |
|-------|-----------|---------|
| 🟢 excellent | ≥ 9.0 | Most requirements covered |
| 🔵 good | ≥ 7.0 | Solid coverage, some gaps |
| 🟡 partial | ≥ 5.0 | Significant gaps |
| 🔴 poor | < 5.0 | Major coverage gaps |
| ⚪ na | N/A | No requirements or no tests |

### Orphan P0 Detection

**Critical flag:** P0 requirements with zero test coverage.

```
orphan_p0 = count of P0 requirements where is_covered == 0

if orphan_p0 > 0:
    → WARN in quality report
    → listed as highest-priority recommendation
```

## Phase 3E: TTL Sync Suggestions

Generate suggested additions/removals for traceability.ttl:

### For UNTRACKED_TEST:
```turtle
# Suggested addition:
vdo:test-<service>-<name> a vdo:TestSuite ;
    vdo:service "<service>" ;
    vdo:sourceFile "src/services/<service>/tests/<file>" .

vdo:test-<service>-<name> vdo:validates vdo:<REQ-ID> .
```

### For STALE_REF:
```
Suggested removal of triple:
vdo:<test> vdo:validates vdo:<requirement> .
```

### For UNTRACKED_REQUIREMENT:
```
Suggested addition:
vdo:<REQ-ID> a vdo:Requirement ;
    vdo:priority "<priority>" ;
    vdo:area "<area>" .
```

## Current State (from research)

Based on initial analysis of VEDO Hub:

| Metric | Value | Gap |
|--------|-------|-----|
| TTL entries | 4 NFR + 2 Spec | vs 200+ REQ files |
| Orphan test files | 41 | test files without TTL entry |
| Untracked user stories | 16 | US-* in TTL but not declared |
| Untracked REQ files | 198 | REQ-*.md without TTL entry |

**Recommendation:** Start with P0 requirements — ensure all P0 REQ files have TTL entries and at least one test link.
