# ── POST /auth/login - invalid email: invalid email format ──
# case_id=TC-2286db52
# case_name=POST /auth/login - invalid email: invalid email format
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "not-an-email",
  "password": "sigh"
}
```

HTTP 422

