# ── POST /api/tokens - [schema_violation] scope_missing_required ──
# case_id=TC-ad285328
# case_name=POST /api/tokens - [schema_violation] scope_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Bonita Hermann"
}
```

HTTP 422

