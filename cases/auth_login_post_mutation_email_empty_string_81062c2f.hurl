# ── POST /auth/login - mutation: email empty string ──
# case_id=TC-81062c2f
# case_name=POST /auth/login - mutation: email empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "",
  "password": "staff"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

