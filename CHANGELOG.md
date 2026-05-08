# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [Unreleased]

## [0.11.0] - 2026-05-08

### Added
- **AWS Bedrock provider**: configure `provider: bedrock` with any Claude model via AWS credentials; `caseforge doctor` validates region and credentials
- **Batch LLM annotation** (`--annotation-batch N`): group N operations per LLM call to cut annotation time dramatically on large specs (36-op spec: ~9 min → ~1–2 min with `--annotation-batch 10`)
- **Spec-hash deduplication**: `gen` exits early with `✓ Spec unchanged` when the spec file hash matches previously generated output; bypass with `--force`
- **`# case_name=` header in Hurl output**: every generated `.hurl` file now includes both `# case_id=` and `# case_name=` comment headers for traceability
- **Rate-limit backoff**: 429 responses from LLM providers trigger exponential backoff (5 s → 15 s → 30 s → 60 s) before retrying
- **MCP `generate_test_cases` tool**: exposes `annotation_batch` parameter so AI clients can control batch size
- **Onboarding wizard redesign**: provider-specific sub-prompts, MCP/skill multi-select, config written to `~/.caseforge.yaml`
- **AI config validation**: `config.Load()` validates `api_key`/`base_url` at startup with actionable error messages

### Changed
- Generated output filenames are now path-first (e.g. `pets_get_200.hurl`) for better filesystem sorting and readability
- OWASP methodology labels translated from Chinese to English in all generated output
- TUI annotation progress now visible during LLM pre-pass

### Fixed
- Makefile build output moved to `bin/` and added to `.gitignore`
- `gen` warn lines route through TUI instead of raw stderr during interactive runs
- `cases/` directory added to `.gitignore` — generated test output no longer accidentally committed

### Documentation
- NOTICE expanded to enumerate all 11 referenced projects and academic papers with explicit attribution wording
- README updated with `--force`, `--annotation-batch`, dedup behavior, rate-limit backoff, and Hurl header fields
- Bump `google.golang.org/grpc` to v1.80.0 (resolves GHSA-p77j-4mvh-x3m3)

## [0.10.2] - 2026-05-03

### Fixed
- `caseforge doctor` no longer checks for the `tree-sitter` CLI binary — RBT uses the `gotreesitter` Go library with embedded grammars, not the CLI
- `caseforge doctor` AI-disabled warning now only fires when all three provider keys (Anthropic, OpenAI, Gemini) are absent; individual provider warnings are per-provider

### Documentation
- Reframe README around spec-driven workflow, clarifying that LLM is optional

## [0.10.1] - 2026-04-27

### Documentation
- Correct technique flag names and complete the SKILL.md technique list

## [0.10.0] - 2026-04-27

### Added — six concept-level sprints (see NOTICE for full attribution)

- **Tcases-style techniques** (concept-level reference): `isolated_negative`, `schema_violation`, `variable_irrelevance` techniques; `--tuple-level N` for N-way coverage; `--seed` for deterministic generation; cross-variable constraint filtering in pairwise
- **RESTler-style stateful testing** (concept-level reference): dependency graph (`internal/methodology/depgraph.go`), N-step chain technique, mutation, auth-chain; `--max-cases-per-op` cap; `caseforge chain` BFS command
- **EvoMaster-style features** (concept-level reference): OpenAPI Links wiring through dep-graph and BFS chain; DELETE teardown for chain commands; `DataPool` for response-derived seeding; `--export-pool` / `--data-pool`; per-(endpoint × status-code) coverage metric; Postman v2.1 seeding via `--seed-postman`; `--prioritize-uncovered` two-pass exploration
- **Schemathesis-style features** (concept-level reference): operation filtering, response schema validation, `constraint_mutation` technique, named coverage scenarios in scoring, `--max-failures`, rule deduplication
- **CATS-style fuzzing** (concept-level reference): `type_coercion`, `unicode_fuzzing`, `mass_assignment`, `idor` techniques; explore extensions; pattern-aware datagen
- **Portman / Microcks-style techniques** (concept-level reference): `semantic_annotation` (nullable / readOnly / writeOnly), `field_boundary`, `required_omission`, `positive_parameter_examples` techniques; CRUD flow scenario tagging; `caseforge import har` command; `score --min-score` CI gate with conformance trend tracking
- **Paper-inspired features** (concept-level reference to RBCTest, AutoRestTest, RESTifAI): `--auth-bootstrap` for secured-endpoint chain wrapping; run failure classification (`server_error` / `missing_validation` / `auth_failure` / `security_regression`); `score --fill-gaps`; `--with-oracles` Observation-Confirmation oracle mining; `business_rule_violation`, `chain_sequence` techniques; `caseforge conformance` command

## [0.9.0] - 2026-04-05

### Added
- TUI improvements: progress list, dynamic flag completion
- `caseforge watch`, `caseforge stats`, `caseforge ci init` commands
- DEA edge-case hypotheses: array constraints, required query param, format violation
- MBT-style methodology techniques: Classification Tree, Orthogonal Array
- MCP tools: `lint_spec`, `ask_test_cases`
- Webhook event push: `on_generate` / `on_run_complete`
- `caseforge score` command — multi-dimensional test case quality scoring

### Changed
- Replaced `tree-sitter` CLI dependency with `gotreesitter` pure-Go library

### Fixed
- `score`: use `CaseSource.SpecPath` for operation grouping
- ask: harden LLM intent parsing (13 new test cases)

## [0.8.0] - 2026-04-04

### Added
- Risk-Based Testing (RBT) — `--generate` for high-risk uncovered operations; LLM parser improvements with concurrency, content cache; `caseforge rbt index` Embed Phase
- `caseforge export` — Allure / Xray / TestRail adapters
- `caseforge diff --gen-cases` — auto-generate test cases for breaking spec operations
- Example extraction & validation from OpenAPI specs (PH2-15)
- TestSuite cross-case DAG orchestration (`caseforge suite create / validate`)
- `caseforge ask` — generate test cases from natural-language descriptions
- Semantic field inference and cross-field consistency in datagen
- Hurl exit codes 3 & 4 (P1-14, P1-15, P1-16) and Appendix B structured comments
- `caseforge gen` flags: `--technique`, `--priority`, `--operations`, `--concurrency`
- Assertion operators: `matches`, `is_iso8601`, `is_uuid`, `gt`

### Fixed
- RBT V3 `--depth` semantics aligned with V2

## [0.7.0] - 2026-03-23

### Added
- Engineering infrastructure: LICENSE (Apache 2.0), CONTRIBUTING.md, CI workflow, GoReleaser configuration, initial CHANGELOG seeded from git history
- Multi-provider LLM support: OpenAI, OpenAI-compat (DeepSeek, Qwen, Moonshot, Azure), Google Gemini alongside Anthropic
- `NewProviderWithConfig` factory accepting `baseURL` for OpenAI-compatible APIs
- `ai.base_url` config field for OpenAI-compat providers
- `caseforge doctor` checks `OPENAI_API_KEY` and `GEMINI_API_KEY` / `GOOGLE_API_KEY`

### Fixed
- OpenAI provider: guard `MaxTokens=0` to avoid sending `max_tokens: 0`
- Gemini provider: guard empty `req.Messages` to prevent index panic

## [0.6.0] - 2026-03-23

### Added
- Multi-provider LLM support: OpenAI, OpenAI-compat (DeepSeek, Qwen, Moonshot, Azure), and Google Gemini alongside Anthropic
- `NewProviderWithConfig` factory accepting `baseURL` for OpenAI-compatible APIs
- `ai.base_url` config field for OpenAI-compat providers
- `caseforge doctor` now checks `OPENAI_API_KEY` and `GEMINI_API_KEY` / `GOOGLE_API_KEY`

### Fixed
- OpenAI provider: guard `MaxTokens=0` to avoid sending `max_tokens: 0` to API
- Gemini provider: guard empty `req.Messages` to prevent index panic

## [0.5.0] - 2026-03-23

### Added
- MCP Sampling mode (`caseforge mcp` command) — exposes CaseForge as an MCP server over stdio
- `generate_cases` MCP tool backed by the AI generation engine
- Integration with the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk`)

## [0.4.0] - 2026-03-23

### Added
- k6 JavaScript test renderer (`--format k6`)
- `caseforge run` command with Hurl and k6 runner backends
- Per-case `RunResult` in runner output

### Fixed
- k6 renderer: sorted headers, correct HTTP API arity, removed dead `capturedVars`

## [0.3.0] - 2026-03-23

### Added
- OWASP security testing methodology and test-case generation
- Full lint rule set with `--fail-on` threshold support
- Postman Collection v2.1 renderer (`--format postman`)
- Spec diff command (`caseforge diff`)

## [0.2.0] - 2026-03-22

### Added
- Chain test cases (multi-step flows with variable capture)
- Interactive TUI progress view
- CSV output format
- Event bus for async runner notifications
- URL encoding helpers and response report parsing

## [0.1.0] - 2026-03-22

### Added
- Initial release: `caseforge gen` command — AI-powered HTTP test-case generation from OpenAPI specs
- Hurl output renderer
- `caseforge lint` command — validates generated test cases against spec
- `caseforge doctor` command — environment dependency checker
- Apache 2.0 license
