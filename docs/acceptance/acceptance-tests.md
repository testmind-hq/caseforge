# CaseForge Acceptance Test Suite

> **Usage:** Run `scripts/acceptance.sh` to execute all scenarios and regenerate results.
> Update this document whenever a new feature is added.

---

## How to Run

```bash
# Full acceptance run (builds binary + executes all scenarios)
./scripts/acceptance.sh

# Individual scenario (manual smoke test)
/tmp/caseforge <command> [flags]
```

---

## Acceptance Scenarios

### Core / CLI

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-001 | `--version` flag | `caseforge --version` | prints `caseforge version <tag>` | ✅ PASS |
| AT-002 | All commands registered | `caseforge --help` | lists ask, completion, config, diff, doctor, explore, fake, gen, init, lint, mcp, onboard, pairwise, run | ✅ PASS |
| AT-003 | `init` creates config | `caseforge init` in empty dir | `.caseforge.yaml` created | ✅ PASS |

---

### `gen` — Test Case Generation

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-004 | gen hurl format | `caseforge gen --spec petstore.yaml --format hurl --output ./cases` | 17 test cases, `.hurl` files written | ✅ PASS |
| AT-005 | gen json format | `caseforge gen --spec petstore.yaml --format json --output ./cases` | valid JSON with `$schema`, `version`, `test_cases[]` | ✅ PASS |
| AT-006 | gen postman format | `caseforge gen --spec petstore.yaml --format postman --output ./cases` | Postman collection file written | ✅ PASS |
| AT-007 | gen k6 format | `caseforge gen --spec petstore.yaml --format k6 --output ./cases` | k6 JS script written | ✅ PASS |
| AT-008 | gen csv format | `caseforge gen --spec petstore.yaml --format csv --output ./cases` | CSV file written | ✅ PASS |
| AT-009 | gen markdown format | `caseforge gen --spec petstore.yaml --format markdown --output ./cases` | Markdown file written | ✅ PASS |
| AT-010 | gen --no-ai flag | `caseforge gen --spec petstore.yaml --no-ai --format hurl` | generates without LLM, same count | ✅ PASS |
| AT-011 | gen invalid spec path | `caseforge gen --spec nonexistent.yaml` | error: file not found | ✅ PASS |

---

### `gen` — Technique Coverage

| ID | Scenario | Expected Techniques | Status |
|----|----------|---------------------|--------|
| AT-012 | Equivalence partitioning cases | `equivalence_partitioning` cases generated for GET/POST/DELETE | ✅ PASS |
| AT-013 | OWASP API Top 10 security cases | `owasp_api_top10` and `owasp_api_top10_spec` cases generated | ✅ PASS |
| AT-014 | Idempotency chain cases | `idempotency` technique → `kind: "chain"` with 2 steps (setup + test), `step-test.depends_on = ["step-setup"]` | ✅ PASS |
| AT-015 | CRUD chain cases | `chain_crud` technique generated for POST+GET+DELETE operations | ✅ PASS |

---

### `lint` — Spec Linting

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-016 | lint valid spec | `caseforge lint --spec petstore.yaml` | score reported, warnings for missing descriptions | ✅ PASS |
| AT-017 | lint missing operationId | `caseforge lint --spec bad.yaml` (no operationId) | `[L001]` warning reported | ✅ PASS |

---

### `diff` — Spec Diff

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-018 | diff identical specs | `caseforge diff --old spec.yaml --new spec.yaml` | `No changes detected.` | ✅ PASS |
| AT-019 | diff with breaking changes | `caseforge diff --old v1.yaml --new v2.yaml` | BREAKING: removed endpoints listed, NON-BREAKING: new endpoints listed | ✅ PASS |

---

### `doctor` — Environment Check

| ID | Scenario | Expected | Status |
|----|----------|----------|--------|
| AT-020 | doctor checks tools and API keys | reports hurl/k6 found/missing, API key set/not set per provider | ✅ PASS |

---

### `fake` — Data Generation

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-021 | fake from JSON schema | `caseforge fake --schema '{"type":"object","properties":{"name":{"type":"string"},"age":{"type":"integer"}}}'` | valid JSON object with name (string) and age (integer) | ✅ PASS |

---

### `pairwise` — Pairwise Combinations

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-022 | pairwise 3 params | `caseforge pairwise --params "browser:chrome,firefox os:win,mac lang:en,zh"` | 4 combinations (< full factorial 8), all pairs covered | ✅ PASS |

---

### `completion` — Shell Completion

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-023 | bash completion | `caseforge completion bash` | bash completion script output | ✅ PASS |
| AT-024 | zsh completion | `caseforge completion zsh` | zsh completion script output | ✅ PASS |
| AT-025 | fish completion | `caseforge completion fish` | fish completion script output | ✅ PASS |

---

### `config show` — Configuration Display

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-026 | config show defaults (no file) | `caseforge config show` in dir without `.caseforge.yaml` | shows defaults: provider=noop, format=hurl | ✅ PASS |
| AT-027 | config show with file | `caseforge config show` in dir with `.caseforge.yaml` | shows config values, API key masked as `sk-ant...` (first 6 chars + `...`) | ✅ PASS |

---

### `ask` — NL-to-Test-Case

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-028 | ask with no LLM (noop provider) | `caseforge ask "POST /users"` with noop config | error: AI provider unavailable | ✅ PASS |
| AT-029 | ask with live LLM (gemini) | `caseforge ask "POST /users create user with name and email"` | generates test cases, writes to `./cases` | ✅ PASS |

---

### `explore` — Dynamic Exploration Agent

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-035 | explore command registered | `caseforge --help` | `explore` listed | ✅ PASS |
| AT-036 | explore --dry-run produces report | `caseforge explore --spec petstore.yaml --dry-run --output ./reports` | dea-report.json written, planned rules listed | ✅ PASS |
| AT-037 | explore missing --spec returns error | `caseforge explore --target http://x` | error: --spec is required | ✅ PASS |
| AT-038 | explore missing --target returns error | `caseforge explore --spec petstore.yaml` | error: --target is required (or use --dry-run) | ✅ PASS |

---

### `rbt` — Regression-Based Testing

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-039 | rbt command registered | `caseforge --help` | `rbt` listed | ✅ PASS |
| AT-040 | missing --spec returns error | `caseforge rbt` | error message | ✅ PASS |
| AT-041 | --format json + --dry-run produces valid JSON | see script | rbt-report.json with diff_base field | ✅ PASS |
| AT-042 | --fail-on high, dry-run → exit 0 | see script | exit 0 | ✅ PASS |
| AT-043 | --dry-run skips git/tree-sitter | see script | no git errors | ✅ PASS |
| AT-044 | doctor shows tree-sitter status | `caseforge doctor` | tree-sitter line present | ✅ PASS |
| AT-045 | rbt index command registered | `caseforge rbt --help` | `index` listed | ✅ PASS |
| AT-046 | rbt index --strategy llm writes map file | see script | map.yaml created with mappings: | ✅ PASS |
| AT-044b | rbt index --out existing without --overwrite fails | see script | error: already exists | ✅ PASS |

---

### `dedupe` — Duplicate Test Case Detection

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-047 | dedupe command registered | `caseforge --help` | `dedupe` listed | ✅ PASS |
| AT-048 | no cases dir returns error | `caseforge dedupe --cases /nonexistent/xyz/cases` | error: cases | ✅ PASS |
| AT-049 | no duplicates exits 0 | `caseforge dedupe --cases <unique-cases-dir>` | exit 0 | ✅ PASS |
| AT-050 | exact duplicate reports group | `caseforge dedupe --cases <dup-cases-dir>` | output contains `Group 1` | ✅ PASS |
| AT-051 | --dry-run exits 0 and files still exist | `caseforge dedupe --cases <dup-cases-dir> --dry-run` | exit 0, both files present | ✅ PASS |
| AT-052 | --merge exits 0 and deletes lower-scoring file | `caseforge dedupe --cases <dup-cases-dir> --merge` | exit 0, lower-scoring file removed | ✅ PASS |

---

### `onboard` — Setup Wizard

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-030 | onboard --yes (non-interactive) | `ANTHROPIC_API_KEY=x caseforge onboard --yes` | auto-selects anthropic, writes `.caseforge.yaml` with `provider: anthropic`, prints Next steps | ✅ PASS |
| AT-031 | onboard skip existing config | `echo n \| caseforge onboard` with existing `.caseforge.yaml` | prompts overwrite, answers n → file unchanged | ✅ PASS |

---

### `run` — Test Execution

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-032 | run hurl (no server) | `caseforge run --cases ./cases --format hurl` | hurl exits with `base_url not set` error (expected without server) | ✅ PASS |
| AT-033 | run k6 (no server) | `caseforge run --cases ./cases --format k6` | k6 exits with connection refused (expected without server) | ✅ PASS |
| AT-034 | run non-existent dir | `caseforge run --cases /nonexistent --format k6` | error: file not found | ✅ PASS |
| AT-053 | run --target injects BASE_URL | `caseforge run --cases ./cases --target http://localhost:9999` | BASE_URL injected (hurl error mentions base_url) | ✅ PASS |
| AT-054 | run --output writes run-report.json | `caseforge run --cases ./cases --target http://localhost:9999 --output ./reports` | `run-report.json` created | ✅ PASS |

---

## Summary (last run: 2026-04-04)

| Category | Total | Pass | Fail |
|----------|-------|------|------|
| Core / CLI | 3 | 3 | 0 |
| gen — formats | 7 | 7 | 0 |
| gen — techniques | 4 | 4 | 0 |
| lint | 2 | 2 | 0 |
| diff | 2 | 2 | 0 |
| doctor | 1 | 1 | 0 |
| fake | 1 | 1 | 0 |
| pairwise | 1 | 1 | 0 |
| completion | 3 | 3 | 0 |
| config show | 2 | 2 | 0 |
| ask | 2 | 2 | 0 |
| explore | 4 | 4 | 0 |
| rbt | 9 | 9 | 0 |
| dedupe | 6 | 6 | 0 |
| onboard | 2 | 2 | 0 |
| run | 5 | 5 | 0 |
| **Total** | **54** | **54** | **0** |

---

## Issues Found During This Run

| Issue | Root Cause | Fix Applied |
|-------|-----------|-------------|
| `ask`, `config`, `completion` commands missing from main | PR #13 code lost when main was force-rebuilt | Restored from `feature/ask-config-completion` worktree, committed `07faf10` |
| Idempotency cases were `kind: "single"` with 1 step instead of `kind: "chain"` with 2 steps | PR #12 chain implementation lost in main rebuild | Restored from `feature/idempotent-chain` worktree, committed `bf5c25c` |

---

## Maintenance Guide

**When adding a new feature:**
1. Add scenario row(s) to the relevant table above
2. Run `./scripts/acceptance.sh` (or manually run the scenario)
3. Update the Summary table counts
4. Update the "last run" date

**When a scenario fails:**
- Code bug → fix code, update status to ✅ PASS after fix
- Environment issue (no server, no API key) → document in Expected column as "expected failure"
- Infrastructure issue → note in Issues Found table
