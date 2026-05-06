# ── POST /api/tokens - [schema_violation] name_missing_required ──
# case_id=TC-c2cef5a1
# case_name=POST /api/tokens - [schema_violation] name_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "scope": "read"
}
```

HTTP 422

