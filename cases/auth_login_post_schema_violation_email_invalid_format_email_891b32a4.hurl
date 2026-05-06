# ── POST /auth/login - [schema_violation] email_invalid_format_email ──
# case_id=TC-891b32a4
# case_name=POST /auth/login - [schema_violation] email_invalid_format_email
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "not-an-email",
  "password": "eye"
}
```

HTTP 422

