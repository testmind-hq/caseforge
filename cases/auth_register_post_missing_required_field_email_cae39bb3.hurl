# ── POST /auth/register - missing required field "email" ──
# case_id=TC-cae39bb3
# case_name=POST /auth/register - missing required field "email"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "password": "this"
}
```

HTTP 422

