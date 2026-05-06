# ── POST /auth/register - missing required field "password" ──
# case_id=TC-72f7ecb7
# case_name=POST /auth/register - missing required field "password"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "mayragrant@nichols.name"
}
```

HTTP 422

