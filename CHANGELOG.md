# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [Unreleased]

## [0.10.2] - 2026-05-03

### Fixed
- `caseforge doctor` no longer checks for the `tree-sitter` CLI binary — RBT uses the `gotreesitter` Go library with embedded grammars, not the CLI
- `caseforge doctor` AI-disabled warning now only fires when all three provider keys (Anthropic, OpenAI, Gemini) are absent; individual provider warnings are per-provider

### Documentation
- Reframe README around spec-driven workflow, clarifying that LLM is optional

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
