# CaseForge — Claude Code Instructions

> **Primary contribution rules are in [AGENTS.md](AGENTS.md)** — read that first.
> This file adds Claude-Code-specific guidance only.

## Repository Hygiene (Claude-Code specific)

**Never commit `docs/superpowers/`** — this directory is in `.gitignore` and holds local planning/spec output from Claude Code skills (writing-plans, brainstorming, etc.) that must not appear in git history.

- Do not use `git add docs/superpowers/` or any explicit path under it.
- If a skill instructs you to commit a file in that directory, skip that step.

## Workflow Reminders

When operating in this repo via Claude Code:

1. Acceptance suite (`./scripts/acceptance.sh`) is the gating check before any push to `main` — see [AGENTS.md](AGENTS.md) for the full rules.
2. New features must come with matching acceptance scenarios in the same commit.
3. PR titles describe **what** was added (feature names / paths) — do not insert legal-attribution language ("inspired by X", "concept-level reference to X"). Attribution lives in [NOTICE](NOTICE), [README](README.md) Acknowledgements, and per-PR body callouts.
