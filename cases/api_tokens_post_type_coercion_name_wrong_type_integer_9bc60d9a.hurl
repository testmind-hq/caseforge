# ── POST /api/tokens - [type_coercion] name wrong_type_integer ──
# case_id=TC-9bc60d9a
# case_name=POST /api/tokens - [type_coercion] name wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": 123,
  "scope": "write"
}
```

HTTP 422

