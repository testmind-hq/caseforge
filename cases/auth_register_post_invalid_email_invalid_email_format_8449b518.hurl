# ── POST /auth/register - invalid email: invalid email format ──
# case_id=TC-8449b518
# case_name=POST /auth/register - invalid email: invalid email format
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "not-an-email",
  "password": "this"
}
```

HTTP 422

