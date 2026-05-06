# ── POST /auth/login - mutation: email invalid email format ──
# case_id=TC-6926df81
# case_name=POST /auth/login - mutation: email invalid email format
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "not-an-email",
  "password": "staff"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

