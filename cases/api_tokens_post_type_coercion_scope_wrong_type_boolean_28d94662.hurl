# ── POST /api/tokens - [type_coercion] scope wrong_type_boolean ──
# case_id=TC-28d94662
# case_name=POST /api/tokens - [type_coercion] scope wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Nathanael Connelly",
  "scope": true
}
```

HTTP 422

