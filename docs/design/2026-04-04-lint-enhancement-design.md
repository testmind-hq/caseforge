# Lint Enhancement — Design Spec

**Date:** 2026-04-04
**Status:** Approved
**Scope:** `caseforge lint` — new rules (A), output formats (B), rule management (C)

---

## Overview

Three orthogonal enhancements to the existing `caseforge lint` command:

- **A) New rules** — 9 additional lint rules (L013–L021) across completeness, consistency, and security categories
- **B) Output formats** — `--format json` for CI consumption + `--output <dir>` to write `lint-report.json`
- **C) Rule management** — `.caseforgelint.yaml` project-level config + `--skip-rules` flag for one-off overrides

All changes follow the existing architecture: rules registered via `init()` in category files, `RunAll` in `runner.go`, command wiring in `cmd/lint.go`.

---

## A) New Rules (L013–L021)

### completeness.go — L013, L014, L015

| ID | Severity | Check |
|----|----------|-------|
| L013 | warning | `parameter` missing `type` (query/path/header parameter has no type declaration) |
| L014 | warning | operation missing any 4xx error response definition |
| L015 | warning | 2xx response schema properties missing `example` |

### consistency.go — L016, L017, L018

| ID | Severity | Check |
|----|----------|-------|
| L016 | error | duplicate `operationId` — two operations share the same ID |
| L017 | warning | inconsistent path versioning — `/v1/` and `/v2/` (or other versions) mixed in same spec |
| L018 | warning | inconsistent response Content-Type — some operations return `application/json`, others return different types |

### security_rules.go — L019, L020, L021

| ID | Severity | Check |
|----|----------|-------|
| L019 | warning | GET operation missing security scheme (L011 only covers non-GET) |
| L020 | error | sensitive field in query parameter (e.g. `?token=`, `?password=`) |
| L021 | warning | no global `securitySchemes` defined anywhere in the spec |

**Total rules after enhancement:** 21 (L001–L021)

---

## B) Output Formats

### New flags on `caseforge lint`

```
--format terminal|json    Output format (default: terminal)
--output <dir>            Write lint-report.json to this directory
```

### JSON stdout format

```json
{
  "score": 72,
  "error_count": 2,
  "warning_count": 4,
  "issues": [
    {
      "rule_id": "L004",
      "severity": "error",
      "path": "POST /users",
      "message": "no 2xx response defined"
    }
  ]
}
```

### Behaviour matrix

| `--format` | `--output` | stdout | file written |
|-----------|-----------|--------|--------------|
| terminal (default) | — | coloured text | — |
| terminal | `./reports` | coloured text | `reports/lint-report.json` |
| json | — | JSON | — |
| json | `./reports` | JSON | `reports/lint-report.json` |

`--output` and `--format` are independent. File output uses the same JSON structure regardless of `--format`.

### Implementation

Add a `LintReport` struct to `internal/lint/report.go` (new file):

```go
type LintReport struct {
    Score        int         `json:"score"`
    ErrorCount   int         `json:"error_count"`
    WarningCount int         `json:"warning_count"`
    Issues       []LintIssue `json:"issues"`
}
```

`cmd/lint.go` constructs `LintReport` from `[]LintIssue` + `Score()`, then renders terminal or JSON accordingly, and optionally writes the file.

---

## C) Rule Management

### `.caseforgelint.yaml` (project-level, committed to repo)

```yaml
# .caseforgelint.yaml
skip_rules:
  - L003
  - L015
fail_on: warning   # optional — overrides caseforge.yaml lint.fail_on
```

File is loaded from the working directory. Missing file is silently ignored.

### `--skip-rules` flag

```bash
caseforge lint --spec openapi.yaml --skip-rules L001,L003
```

Accepts comma-separated rule IDs. Merged with (not replacing) `.caseforgelint.yaml`.

### Priority (three-layer merge)

| Layer | `skip_rules` | `fail_on` |
|-------|-------------|-----------|
| `caseforge.yaml` `lint.skip_rules` | baseline | lowest |
| `.caseforgelint.yaml` | merged (union) | overrides caseforge.yaml |
| `--skip-rules` flag | merged (union) | — (no flag for fail_on) |

`skip_rules` is always a **union** across all layers — a rule skipped at any layer is skipped.
`fail_on` follows highest-priority wins: `.caseforgelint.yaml` > `caseforge.yaml`.

### Implementation

`RunAll` signature changes:

```go
// Before
func RunAll(ps *spec.ParsedSpec) []LintIssue

// After
func RunAll(ps *spec.ParsedSpec, skip map[string]bool) []LintIssue
```

`cmd/lint.go` builds the skip set by merging all three layers before calling `RunAll`.

A new `internal/lint/lintconfig.go` handles loading `.caseforgelint.yaml`:

```go
type LintFileConfig struct {
    SkipRules []string `yaml:"skip_rules"`
    FailOn    string   `yaml:"fail_on"`
}

func LoadLintFileConfig(dir string) (LintFileConfig, error)
```

---

## File Changes

```
internal/lint/
  completeness.go       — add L013, L014, L015
  consistency.go        — add L016, L017, L018
  security_rules.go     — add L019, L020, L021
  report.go             — NEW: LintReport struct + JSON marshal helper
  lintconfig.go         — NEW: LoadLintFileConfig(.caseforgelint.yaml)
  runner.go             — RunAll signature: add skip map[string]bool param
  runner_test.go        — update existing call sites + add skip tests
  lint_test.go          — update call sites

cmd/lint.go             — add --format, --output flags; merge skip layers; render report

docs/acceptance/acceptance-tests.md  — AT-03x new scenarios (see below)
scripts/acceptance.sh                — matching checks
CLAUDE.md                            — update scenario count
```

---

## Acceptance Tests

| ID | Scenario | Command | Expected |
|----|----------|---------|----------|
| AT-039 | lint --format json outputs valid JSON | `caseforge lint --spec petstore.yaml --format json` | parseable JSON with `score`, `issues` fields |
| AT-040 | lint --output writes lint-report.json | `caseforge lint --spec petstore.yaml --output /tmp/r` | `/tmp/r/lint-report.json` created |
| AT-041 | lint --skip-rules suppresses rule | `caseforge lint --spec petstore.yaml --skip-rules L001` | L001 issues absent from output |
| AT-042 | .caseforgelint.yaml skip_rules respected | file with `skip_rules: [L002]`, run lint | L002 issues absent |
| AT-043 | L016 duplicate operationId detected | spec with duplicate IDs | error L016 reported |
| AT-044 | L020 sensitive query param detected | spec with `?token` query param | error L020 reported |

Total acceptance scenarios: 38 → **44**

---

## Out of Scope (V1)

- `--fix` auto-remediation
- Per-path rule suppression (inline spec comments)
- AI-assisted lint suggestions
- Rule severity override via config
