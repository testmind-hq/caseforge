# ── POST /auth/register - [type_coercion] password wrong_type_integer ──
# case_id=TC-4a32c12b
# case_name=POST /auth/register - [type_coercion] password wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "bentonwoods@marsh.net",
  "password": 123
}
```

HTTP 422

