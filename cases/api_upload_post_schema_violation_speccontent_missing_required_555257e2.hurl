# ── POST /api/upload - [schema_violation] specContent_missing_required ──
# case_id=TC-555257e2
# case_name=POST /api/upload - [schema_violation] specContent_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "to",
  "commitSha": "Brazilian",
  "service": "intimidate"
}
```

HTTP 422

