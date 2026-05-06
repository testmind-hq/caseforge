# ── POST /auth/login - missing required field "email" ──
# case_id=TC-9b253ab6
# case_name=POST /auth/login - missing required field "email"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "password": "sigh"
}
```

HTTP 422

