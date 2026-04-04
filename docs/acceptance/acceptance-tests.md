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

### `gen` — CLI Flags (P1-1 to P1-4)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-067 | gen --technique filters output | `caseforge gen --spec petstore.yaml --no-ai --technique equivalence_partitioning --output ./cases` | only `equivalence_partitioning` cases in index.json | ✅ PASS |
| AT-068 | gen --priority filters output | `caseforge gen --spec petstore.yaml --no-ai --priority P1 --output ./cases` | index.json contains only P0/P1 cases | ✅ PASS |
| AT-069 | gen --operations filters spec | `caseforge gen --spec petstore.yaml --no-ai --operations listPets --output ./cases` | only cases for listPets operationId | ✅ PASS |
| AT-070 | gen --concurrency flag accepted | `caseforge gen --help` | `--concurrency` listed | ✅ PASS |

---

### `gen` — index.json Metadata (P1-6 to P1-10)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-071 | index.json contains meta object | `caseforge gen --spec petstore.yaml --no-ai --output ./cases` | `meta` key present in index.json | ✅ PASS |
| AT-072 | meta.spec_hash is non-empty SHA256 | same as AT-071 | `meta.spec_hash` is 64-char hex string | ✅ PASS |
| AT-073 | meta.caseforge_version present | same as AT-071 | `meta.caseforge_version` is non-empty | ✅ PASS |
| AT-074 | meta.by_technique counts match cases | same as AT-071 | `meta.by_technique` sums to total case count | ✅ PASS |
| AT-075 | meta.by_kind counts match cases | same as AT-071 | `meta.by_kind` sums to total case count | ✅ PASS |

---

### `gen` — Assertion Operators (P1-11 to P1-13)

| ID | Scenario | Expected | Status |
|----|----------|----------|--------|
| AT-076 | `exists` operator used in response assertions | Generated cases for an endpoint with object response include `exists` assertions | ✅ PASS |
| AT-077 | `is_uuid` operator used for uuid-format fields | Cases for endpoint with `format: uuid` response field use `is_uuid` operator | ✅ PASS |
| AT-078 | `is_iso8601` operator used for date-time fields | Cases for endpoint with `format: date-time` response field use `is_iso8601` operator | ✅ PASS |

---

### `rbt --generate` — High-Risk Auto-Generation (2.2)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-079 | `--generate` flag registered | `caseforge rbt --help` | `--generate` listed in help | ✅ PASS |
| AT-080 | `--generate --dry-run` prints "ignored" info message | `caseforge rbt --spec petstore.yaml --dry-run --generate` | output contains "ignored with" | ✅ PASS |
| AT-081 | `--generate` writes index.json for real high-risk op | git repo with changed handler.go mapped via caseforge-map.yaml, run `rbt --generate --no-ai` | `index.json` created in cases dir | ✅ PASS |

### `rbt index --strategy embed` — Embed Phase (2.3)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-082 | `rbt index --strategy embed` writes map file (regex fallback without API key) | `caseforge rbt index --spec openapi.yaml --src /tmp/src --strategy embed` with no `OPENAI_API_KEY` | map file written with `mappings:` | ✅ PASS |

### `caseforge export` — Platform Adapters (3.2)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-083 | `export` command registered | `caseforge --help` | `export` listed | ✅ PASS |
| AT-084 | `--format allure` creates Allure result file | `caseforge export --cases ./cases --format allure --output /tmp/out` | `*-result.json` in `/tmp/out/allure/` | ✅ PASS |
| AT-085 | `--format xray` creates xray-import.json | `caseforge export --cases ./cases --format xray --output /tmp/out` | `xray-import.json` in `/tmp/out/xray/` | ✅ PASS |
| AT-086 | `--format testrail` creates testrail-import.csv | `caseforge export --cases ./cases --format testrail --output /tmp/out` | `testrail-import.csv` in `/tmp/out/testrail/` | ✅ PASS |
| AT-087 | `--technique example_extraction` generates cases from spec examples | `caseforge gen --spec example-spec.yaml --no-ai --technique example_extraction` | Output contains `example_extraction` | ✅ PASS |
| AT-088 | Example extraction produces valid (P1) and invalid (P2) cases | Run `--technique example_extraction` on spec with named examples | `.hurl` output contains example name `valid_widget` | ✅ PASS |

---

### `caseforge diff --gen-cases` — Auto-generate for Breaking Changes (3.3)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-089 | `--gen-cases` flag registered | `caseforge diff --help` | `--gen-cases` listed | ✅ PASS |
| AT-090 | breaking changes → `index.json` written | `caseforge diff --old v1.yaml --new v2.yaml --gen-cases /tmp/gen` | `index.json` with `test_cases` array in gen dir | ✅ PASS |

---

<!-- AT-091–AT-092 reserved for future features -->

### `caseforge suite` — TestSuite Orchestration (3.6)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-093 | `suite` command registered | `caseforge --help` | `suite` listed | ✅ PASS |
| AT-094 | `suite create` writes valid suite.json | `caseforge suite create --id S --title T --kind chain --cases TC-001,TC-002` | `suite.json` with `$schema` and `cases` array | ✅ PASS |
| AT-095 | `suite validate` confirms valid suite | `caseforge suite validate --suite suite.json` | `valid ✓` output | ✅ PASS |

---

<!-- AT-091–AT-092 reserved for future features -->

### Assertion Operator Rendering (1.3 completeness)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-096 | `gen` produces index.json with assertions | `caseforge gen --no-ai` on numeric+uuid+datetime spec | `assertions` key present in index.json | ✅ PASS |
| AT-097 | Hurl output has no unrendered assertions | `caseforge gen --no-ai --format hurl` | No `# unrendered assertion` comment in any `.hurl` file | ✅ PASS |
| AT-098 | k6 output has no unrendered assertions | `caseforge gen --no-ai --format k6` | No `// unrendered:` comment in k6 output | ✅ PASS |

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

### `rbt` — Call Graph (V2)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-061 | --depth flag registered on rbt index | `caseforge rbt index --help` | `--depth` listed | ✅ PASS |
| AT-062 | rbt --dry-run exits 0 | `caseforge rbt --spec petstore.yaml --dry-run` | exit 0, report generated | ✅ PASS |
| AT-063 | --depth default is 0 on rbt index | `caseforge rbt index --help` | output contains `depth int` | ✅ PASS |

---

### `rbt` — Call Graph V3 (Go type-aware)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-064 | --algo flag registered on rbt index | `caseforge rbt index --help` | `--algo` listed | ✅ PASS |
| AT-065 | rbt index hybrid runs without error (no Go module) | `caseforge rbt index --spec petstore.yaml --strategy hybrid --src /tmp` | exit 0, map file written | ✅ PASS |
| AT-066 | --algo vta flag accepted | `caseforge rbt index --help` | `vta` mentioned in --algo description | ✅ PASS |

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

### `lint` — Enhancement (AT-055–AT-060)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-055 | lint --format json outputs valid JSON | `caseforge lint --spec petstore.yaml --format json` | parseable JSON with `score` and `issues` keys | ✅ PASS |
| AT-056 | lint --output writes lint-report.json | `caseforge lint --spec petstore.yaml --output /tmp/lr` | `/tmp/lr/lint-report.json` created | ✅ PASS |
| AT-057 | lint --skip-rules suppresses rule | `caseforge lint --spec petstore.yaml --skip-rules L014 --format json` | L014 absent from issues | ✅ PASS |
| AT-058 | .caseforgelint.yaml skip_rules respected | `.caseforgelint.yaml` with `skip_rules: [L014]`, run lint | L014 absent from output | ✅ PASS |
| AT-059 | L016 duplicate operationId detected | spec with two operations sharing same operationId | error L016 reported | ✅ PASS |
| AT-060 | L020 sensitive query param detected | spec with `?token` query parameter | error L020 reported | ✅ PASS |

---

### Exit Codes (P1-15, P1-16)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-071 | lint exits 3 when errors found | `caseforge lint --spec <spec-with-errors>` (spec with duplicate operationId) | exit code 3 | ✅ PASS |
| AT-072 | gen exits 4 when LLM unavailable without --no-ai | `caseforge gen --spec petstore.yaml` (no API key, no --no-ai) | exit code 4, error message about --no-ai | ✅ PASS |

---


## Summary (last run: 2026-04-04)

| Category | Total | Pass | Fail |
|----------|-------|------|------|
| Core / CLI | 3 | 3 | 0 |
| gen — formats | 7 | 7 | 0 |
| gen — techniques | 4 | 4 | 0 |
| gen — CLI flags | 4 | 4 | 0 |
| gen — metadata | 5 | 5 | 0 |
| gen — assertion operators | 3 | 3 | 0 |
| lint | 2 | 2 | 0 |
| lint enhancement | 6 | 6 | 0 |
| diff | 2 | 2 | 0 |
| doctor | 1 | 1 | 0 |
| fake | 1 | 1 | 0 |
| pairwise | 1 | 1 | 0 |
| completion | 3 | 3 | 0 |
| config show | 2 | 2 | 0 |
| ask | 2 | 2 | 0 |
| explore | 4 | 4 | 0 |
| rbt | 15 | 15 | 0 |
| dedupe | 6 | 6 | 0 |
| onboard | 2 | 2 | 0 |
| run | 5 | 5 | 0 |
| exit codes | 2 | 2 | 0 |
| example_extraction | 2 | 2 | 0 |
| **Total** | **90** | **90** | **0** |

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
