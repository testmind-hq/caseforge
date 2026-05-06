# ── POST /auth/login - mutation: email integer instead of string ──
# case_id=TC-d7ccf79e
# case_name=POST /auth/login - mutation: email integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": 12345,
  "password": "staff"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

