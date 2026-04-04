# RBT (Risk-Based Testing) — Design Spec

**Date:** 2026-04-04
**Status:** Approved
**Commands:** `caseforge rbt`, `caseforge rbt index`

---

## Overview

RBT (Risk-Based Testing) connects git changes to API test coverage. Given a git diff, source code, an OpenAPI spec, and a directory of generated test cases, it:

1. Identifies which source files changed
2. Extracts the API routes those files implement (via tree-sitter call graph + LLM fallback)
3. Maps those routes to existing test cases in `./cases/`
4. Produces a risk report: which operations are untested, what risk level, what to do

Primary use case: CI gate — "did my changes break untested API operations?"

---

## Command

```bash
caseforge rbt \
  --spec openapi.yaml \
  [--cases ./cases] \
  [--src ./] \
  [--base HEAD~1] \
  [--head HEAD] \
  [--generate] \
  [--output ./reports] \
  [--format terminal|json] \
  [--fail-on none|low|medium|high] \
  [--map caseforge-map.yaml] \
  [--dry-run]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--spec` | (required) | OpenAPI spec file or URL |
| `--cases` | `./cases` | Directory of generated test cases |
| `--src` | `./` | Source code root to scan for routes |
| `--base` | `HEAD~1` | Git base ref for diff |
| `--head` | `HEAD` | Git head ref (or working tree) |
| `--generate` | false | Auto-generate tests for uncovered HIGH-risk operations |
| `--output` | `./reports` | Output directory for `rbt-report.json` |
| `--format` | `terminal` | Output format: `terminal` or `json` |
| `--fail-on` | `high` | Exit 1 if any operation has risk ≥ this level; `none` = never exit 1 |
| `--map` | (none) | Explicit `caseforge-map.yaml` override file |
| `--dry-run` | false | Return empty RiskReport (0 affected ops) without running git/tree-sitter; always exit 0 |

---

## Architecture & Data Flow

```
git diff (--base..--head)
    └─→ differ.go → []ChangedFile{Path, ChangedLines, IsNew, IsDeleted}
                         │
              mapper.go (SourceParser chain)
              ├─ 1. MapFileParser      — caseforge-map.yaml (explicit, highest priority)
              ├─ 2. TreeSitterParser   — tree-sitter CLI subprocess
              │       ├─ ExtractRoutes: route registration nodes
              │       └─ BuildCallGraph: fn→callers map for service layer tracing
              ├─ 3. RegexParser        — fallback when tree-sitter not installed
              └─ 4. LLMParser          — when call graph trace is uncertain
                         │
                    []RouteMapping{SourceFile, Line, Method, RoutePath}
                         │
              spec path matching (fuzzy + exact)
                         │
                    []AffectedOperation
                         │
              assessor.go
              ├─ ScanCases(casesDir): scan *.json (schema.TestCase) → op→[]TestCaseRef index
              └─ Assess(ops, affected, index) → RiskReport
                         │
           ┌─────────────┼──────────────┐
      terminal         rbt-report.json  generator.go
      report           + exit code      (--generate: call gen pipeline)
```

---

## Core Data Types

```go
// internal/rbt/types.go

type ChangedFile struct {
    Path         string
    ChangedLines []int
    IsNew        bool
    IsDeleted    bool
}

type RouteMapping struct {
    SourceFile string
    Line       int
    Method     string    // "GET", "POST", ...
    RoutePath  string    // "/users/{id}"
    Via        string    // "mapfile"|"treesitter"|"regex"|"llm"
    Confidence float64   // 0.0–1.0; <0.5 → "uncertain"
}

type RiskLevel string

const (
    RiskNone      RiskLevel = "none"
    RiskLow       RiskLevel = "low"
    RiskMedium    RiskLevel = "medium"
    RiskHigh      RiskLevel = "high"
    RiskUncertain RiskLevel = "uncertain" // call graph trace inconclusive
)

type TestCaseRef struct {
    File   string // e.g. "cases/POST_users_201.json"
    CaseID string
    Title  string
}

type OperationCoverage struct {
    OperationID string
    Method      string
    Path        string
    Affected    bool          // touched by git diff
    SourceRefs  []RouteMapping
    TestCases   []TestCaseRef
    Risk        RiskLevel
}

type RiskReport struct {
    DiffBase       string
    DiffHead       string
    ChangedFiles   []ChangedFile
    Operations     []OperationCoverage
    TotalAffected  int
    TotalCovered   int
    TotalUncovered int
    RiskScore      float64   // uncovered/affected, 0.0–1.0
    GeneratedAt    time.Time
}
```

---

## Layer Implementations

### differ.go — Git Diff Parser

```
git diff --name-only --diff-filter=AMD <base>..<head>  → changed file list
git diff -U0 <base>..<head> <file>                     → changed line numbers
```

Returns `[]ChangedFile`. If not in a git repo, returns empty list (no operations affected).

### mapper.go — SourceParser Chain

```go
type SourceParser interface {
    ExtractRoutes(srcDir string, files []ChangedFile) ([]RouteMapping, error)
    BuildCallGraph(srcDir string) (CallGraph, error) // CallGraph = map[string][]string
}
```

**Execution:** Try parsers in order (MapFile → TreeSitter → Regex → LLM). Each parser receives only the files not claimed by a higher-priority parser. A file is "claimed" once at least one RouteMapping is produced for it.

**TreeSitterParser:**
- Detects language from file extension (`.go`, `.py`, `.ts`, `.java`, `.rb`, `.rs`)
- Runs `tree-sitter query <query-file> <source-file>` subprocess
- Route extraction query (example for Go/gin):
  ```scheme
  (call_expression
    function: (selector_expression
      field: (field_identifier) @method
      (#match? @method "^(GET|POST|PUT|DELETE|PATCH)$"))
    arguments: (argument_list
      (interpreted_string_literal) @path))
  ```
- Call graph query: extracts all `function_declaration` + `call_expression` pairs
- Service layer changes: traverse CallGraph upward until hitting a known route handler

**LLMParser (fallback):**
```
Prompt: "Given this git diff and OpenAPI spec, which operations are likely affected?
Return JSON: [{"method":"POST","path":"/users","confidence":0.9}]"
```
Results with `confidence < 0.5` → `RiskUncertain`.

### assessor.go — Coverage Analysis

**ScanCases:** Scans `*.json` files only (format-agnostic `schema.TestCase`). Reads `source.spec_path` (e.g. `"POST /users"`). Falls back to filename-based inference if no `.json` files exist.

**Risk levels:**

| Condition | Risk |
|-----------|------|
| Operation not in affected set | `none` |
| Affected, ≥2 test cases | `low` |
| Affected, 1 test case | `medium` |
| Affected, 0 test cases | `high` |
| Call graph trace inconclusive | `uncertain` |

**Risk score:** `uncovered_affected / total_affected` (0.0–1.0)

### generator.go — Gap Filling (--generate)

For all `high`-risk operations, invokes the existing `internal/methodology` pipeline (same as `caseforge gen`) targeted at those specific operations. Writes new test cases to `--cases` directory.

---

## Output

### Terminal (default)

```
┌──────────────────────┬──────────┬──────────┬───────┬──────────────────────┐
│ Operation            │ Risk     │ Affected │ Cases │ Source               │
├──────────────────────┼──────────┼──────────┼───────┼──────────────────────┤
│ POST /users          │ HIGH     │ ✓        │ 0     │ service/user.go:42 (llm) │
│ GET  /users/{id}     │ MEDIUM   │ ✓        │ 1     │ handler/user.go:18   │
│ GET  /users          │ LOW      │ ✓        │ 3     │ handler/user.go:10   │
│ DELETE /users/{id}   │ NONE     │ -        │ 2     │ -                    │
└──────────────────────┴──────────┴──────────┴───────┴──────────────────────┘
Risk Score: 0.33  (1 uncovered / 3 affected)
```

### rbt-report.json (--output)

Full `RiskReport` serialized as JSON, written to `<output>/rbt-report.json`.

### Exit codes

| Exit | Condition |
|------|-----------|
| 0 | No operation has risk ≥ `--fail-on` level |
| 1 | At least one operation has risk ≥ `--fail-on` level |

---

## File Structure

```
internal/rbt/
  types.go              — core data types
  differ.go             — git diff → []ChangedFile
  differ_test.go
  mapper.go             — SourceParser interface + chain orchestration
  mapper_test.go
  parser_mapfile.go     — explicit caseforge-map.yaml parser
  parser_treesitter.go  — tree-sitter CLI subprocess parser
  parser_regex.go       — regex fallback parser
  parser_llm.go         — LLM-assisted inference
  assessor.go           — coverage analysis + risk scoring
  assessor_test.go
  generator.go          — gap-filling test generation (--generate)
  report.go             — terminal + JSON + file output
  report_test.go

cmd/rbt.go              — cobra command
cmd/rbt_test.go         — command registration + flag tests

docs/acceptance/acceptance-tests.md  — AT-039 through AT-044 (43 total)
scripts/acceptance.sh                — AT-039 through AT-044 checks
CLAUDE.md                            — update "37 scenarios" → "43 scenarios"
```

---

## caseforge-map.yaml Format

Optional explicit mapping file for cases where automatic extraction is insufficient (e.g. dynamic routing, languages without tree-sitter grammar):

```yaml
# caseforge-map.yaml
mappings:
  - source: internal/user/service.go
    operations:
      - POST /users
      - GET /users/{id}
      - PUT /users/{id}
      - DELETE /users/{id}
  - source: internal/order/handler.go
    operations:
      - POST /orders
      - GET /orders/{id}
```

---

## caseforge doctor Integration

Add tree-sitter to doctor checks:

```
✓ hurl 6.0.0
✓ git 2.43.0
⚠ tree-sitter not found — install with: brew install tree-sitter
  (RBT will use regex fallback for route extraction)
```

---

## Acceptance Tests (AT-039 – AT-044)

| ID | Scenario | Command | Expected |
|----|----------|---------|----------|
| AT-039 | rbt command registered | `caseforge --help` | `rbt` listed |
| AT-040 | missing --spec returns error | `caseforge rbt` | error: --spec is required |
| AT-041 | --format json produces valid JSON | `caseforge rbt --spec p.yaml --format json --dry-run` | parseable JSON |
| AT-042 | --fail-on high, no risk → exit 0 | `caseforge rbt --spec p.yaml --dry-run --fail-on high` | exit 0 |
| AT-043 | --dry-run skips git/tree-sitter | `caseforge rbt --spec p.yaml --dry-run` | no git/tree-sitter calls |
| AT-044 | doctor shows tree-sitter status | `caseforge doctor` | `tree-sitter` line present |

Total acceptance scenarios: 37 → **43**

---

---

## `caseforge rbt index` — LLM + Embedding Index Generation

Generates `caseforge-map.yaml` automatically by analyzing the project, eliminating manual maintenance.

### Command

```bash
caseforge rbt index \
  --spec openapi.yaml \
  --src ./ \
  [--out caseforge-map.yaml] \
  [--strategy llm|embed|hybrid] \
  [--overwrite]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--spec` | (required) | OpenAPI spec |
| `--src` | `./` | Source code root |
| `--out` | `caseforge-map.yaml` | Output map file path |
| `--strategy` | `llm` | Indexing strategy (see below) |
| `--overwrite` | false | Overwrite existing map file |

### Strategies

| Strategy | How | Best For | Extra Deps |
|----------|-----|---------|-----------|
| `llm` | Send files + spec to LLM, infer mapping | Small projects (<50 files), one-off | None (existing LLM) |
| `embed` | Generate embeddings for code chunks + spec ops, cosine similarity | Large projects, CI-frequent | Embedding API |
| `hybrid` (**recommended**) | embed narrows candidates (top-k), LLM confirms | General use, best accuracy | Embedding API |

### Hybrid Strategy Flow

```
All source files
    │
    ├─ tree-sitter: files with direct route registration
    │     → map entry (via: treesitter, confidence: 1.0)
    │
    └─ remaining files (service / repo / util)
          │
          ▼
    Embedding phase:
      1. Chunk each file by function (tree-sitter or line-window fallback)
      2. Generate embedding per chunk via LLM embedding API
      3. Generate embedding per spec operation (method + path + description)
      4. Cosine similarity → top-k candidates per operation
      Cache: .caseforge-index/index.json (skip re-embedding unchanged files)
          │
          ▼
    LLM confirmation phase (only top-k candidates, not all files):
      Prompt: "Does this function likely implement or call logic for [POST /users]?
               Return {match: true/false, confidence: 0.0-1.0}"
          │
          ▼
    map entry (via: hybrid, confidence: <score>)
```

### Local Index Storage (no vector DB required)

```
.caseforge-index/
  index.json         — per-chunk embeddings + file content hash
  meta.json          — index version, spec hash, created_at
```

`index.json` schema:
```json
{
  "chunks": [
    {
      "file": "internal/user/service.go",
      "fn":   "UserService.Create",
      "hash": "sha256:abc123",
      "embedding": [0.12, -0.34, ...]
    }
  ],
  "spec_ops": [
    {
      "operation": "POST /users",
      "description": "Create a new user account",
      "embedding": [0.15, -0.28, ...]
    }
  ]
}
```

**Incremental update:** on re-index, skip chunks whose `hash` matches — only re-embed changed functions.

### Generated `caseforge-map.yaml`

```yaml
# caseforge-map.yaml — generated by `caseforge rbt index`
# Strategy: hybrid | Spec: openapi.yaml | Indexed: 2026-04-04T10:00:00Z
# Review entries marked ⚠ (confidence < 0.7) before committing.
mappings:
  - source: internal/handler/user.go       # via: treesitter
    operations:
      - POST /users
      - GET /users/{id}
      - DELETE /users/{id}

  - source: internal/service/user.go       # via: hybrid (confidence: 0.92)
    operations:
      - POST /users
      - GET /users/{id}
      - PUT /users/{id}
      - DELETE /users/{id}

  - source: internal/service/auth.go       # via: hybrid (confidence: 0.61) ⚠ review
    operations:
      - POST /auth/login
      - POST /auth/refresh
```

### Recommended Workflow

```bash
# One-time project setup
caseforge rbt index --spec openapi.yaml --strategy hybrid
git add caseforge-map.yaml .caseforge-index/
git commit -m "chore: add caseforge RBT index"

# Daily CI usage (uses committed map, zero LLM calls)
caseforge rbt --spec openapi.yaml --fail-on high

# Re-index after significant refactor
caseforge rbt index --spec openapi.yaml --strategy hybrid --overwrite
```

`.caseforge-index/` should be committed for reproducibility, or added to `.gitignore` to always re-generate (slower CI but no committed binary data).

### New Files (index subcommand)

```
internal/rbt/
  indexer.go           — orchestrates index generation (tree-sitter + embed + LLM)
  indexer_test.go
  embedder.go          — embedding API calls + cosine similarity
  embedder_test.go
  index_store.go       — .caseforge-index/ read/write + incremental hash check
  index_store_test.go

cmd/rbt_index.go       — `caseforge rbt index` subcommand
cmd/rbt_index_test.go
```

### Acceptance Tests (AT-045 – AT-047)

| ID | Scenario | Command | Expected |
|----|----------|---------|----------|
| AT-045 | rbt index command registered | `caseforge rbt --help` | `index` listed |
| AT-046 | rbt index --strategy llm writes map file | `caseforge rbt index --spec p.yaml --strategy llm --out /tmp/map.yaml` | map.yaml created |
| AT-047 | rbt index incremental: unchanged files skipped | re-run with same files | index.json unchanged for unmodified files |

Total acceptance scenarios: 43 → **46**

---

## Out of Scope (V1)

- Multi-repo analysis
- Non-git VCS (SVN, Mercurial)
- Runtime coverage (instrumentation hooks)
- PR-level comment posting (future CI integration)
- Non-Go AST via native bindings (use tree-sitter subprocess instead)
- Vector database (Qdrant/Chroma) — `.caseforge-index/index.json` flat file is sufficient for V1
