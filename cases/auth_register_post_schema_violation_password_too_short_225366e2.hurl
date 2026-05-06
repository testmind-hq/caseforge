# ── POST /auth/register - [schema_violation] password_too_short ──
# case_id=TC-225366e2
# case_name=POST /auth/register - [schema_violation] password_too_short
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "jadonrobertson@wu.org",
  "password": ""
}
```

HTTP 422

