# ── POST /api/upload - [schema_violation] branch_missing_required ──
# case_id=TC-381d4381
# case_name=POST /api/upload - [schema_violation] branch_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "commitSha": "Brazilian",
  "service": "intimidate",
  "specContent": "tonight"
}
```

HTTP 422

