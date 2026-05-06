# ── POST /api/tokens - [type_coercion] name wrong_type_boolean ──
# case_id=TC-bd1e61be
# case_name=POST /api/tokens - [type_coercion] name wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": true,
  "scope": "write"
}
```

HTTP 422

