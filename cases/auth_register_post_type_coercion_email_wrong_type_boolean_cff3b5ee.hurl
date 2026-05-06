# ── POST /auth/register - [type_coercion] email wrong_type_boolean ──
# case_id=TC-cff3b5ee
# case_name=POST /auth/register - [type_coercion] email wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": true,
  "password": "it"
}
```

HTTP 422

