# ── POST /auth/register - [type_coercion] password wrong_type_boolean ──
# case_id=TC-4af1b36a
# case_name=POST /auth/register - [type_coercion] password wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "bentonwoods@marsh.net",
  "password": true
}
```

HTTP 422

