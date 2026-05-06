# ── POST /api/upload - [schema_violation] branch_too_short ──
# case_id=TC-76d8b912
# case_name=POST /api/upload - [schema_violation] branch_too_short
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "",
  "commitSha": "Brazilian",
  "service": "intimidate",
  "specContent": "tonight"
}
```

HTTP 422

