# ── POST /api/tokens - [type_coercion] scope wrong_type_integer ──
# case_id=TC-9bf5d669
# case_name=POST /api/tokens - [type_coercion] scope wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Nathanael Connelly",
  "scope": 123
}
```

HTTP 422

