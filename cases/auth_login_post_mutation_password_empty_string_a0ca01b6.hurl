# ── POST /auth/login - mutation: password empty string ──
# case_id=TC-a0ca01b6
# case_name=POST /auth/login - mutation: password empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "naomipierce@lewis.biz",
  "password": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

