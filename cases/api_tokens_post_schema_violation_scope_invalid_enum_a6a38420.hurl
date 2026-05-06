# ── POST /api/tokens - [schema_violation] scope_invalid_enum ──
# case_id=TC-a6a38420
# case_name=POST /api/tokens - [schema_violation] scope_invalid_enum
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Bonita Hermann",
  "scope": "__invalid__"
}
```

HTTP 422

