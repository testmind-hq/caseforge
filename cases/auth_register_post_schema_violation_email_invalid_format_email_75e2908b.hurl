# ── POST /auth/register - [schema_violation] email_invalid_format_email ──
# case_id=TC-75e2908b
# case_name=POST /auth/register - [schema_violation] email_invalid_format_email
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "not-an-email",
  "password": "these"
}
```

HTTP 422

