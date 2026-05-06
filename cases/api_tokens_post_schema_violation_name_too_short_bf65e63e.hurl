# ── POST /api/tokens - [schema_violation] name_too_short ──
# case_id=TC-bf65e63e
# case_name=POST /api/tokens - [schema_violation] name_too_short
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "",
  "scope": "read"
}
```

HTTP 422

