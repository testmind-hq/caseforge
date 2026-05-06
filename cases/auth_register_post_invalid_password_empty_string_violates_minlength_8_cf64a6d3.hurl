# ── POST /auth/register - invalid password: empty string violates minLength 8 ──
# case_id=TC-cf64a6d3
# case_name=POST /auth/register - invalid password: empty string violates minLength 8
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "mayragrant@nichols.name",
  "password": ""
}
```

HTTP 422

