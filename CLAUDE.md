# CaseForge — Claude Code Instructions

## Git Rules

**Never commit `docs/superpowers/`** — this directory is in `.gitignore` and contains internal planning/spec files (writing-plans, brainstorming skill outputs) that must not appear in git history. Do not use `git add docs/superpowers/` or any explicit path under it. If a skill instructs you to commit a file in this directory, skip that step.

---

## Acceptance Tests

Every development session must run the acceptance suite before committing and after completing any feature or fix:

```bash
./scripts/acceptance.sh
```

**Rules:**
- All 92 scenarios must pass before pushing to `main`.
- When adding a new command or feature, add the corresponding scenario(s) to `docs/acceptance/acceptance-tests.md` and a matching `contains` / `run` call in `scripts/acceptance.sh` **in the same commit**.
- If a scenario fails, determine root cause before proceeding:
  - **Code bug** → fix the code, re-run the suite.
  - **Environment issue** (no API key, no running server) → document in the Expected column as "expected failure" and ensure the scenario tests the error path, not the happy path.
