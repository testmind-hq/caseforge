# ── POST /auth/register - [type_coercion] email wrong_type_integer ──
# case_id=TC-c40fa64f
# case_name=POST /auth/register - [type_coercion] email wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": 123,
  "password": "it"
}
```

HTTP 422

