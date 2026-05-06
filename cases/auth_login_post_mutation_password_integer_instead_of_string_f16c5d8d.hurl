# ── POST /auth/login - mutation: password integer instead of string ──
# case_id=TC-f16c5d8d
# case_name=POST /auth/login - mutation: password integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "naomipierce@lewis.biz",
  "password": 12345
}
```

HTTP *

[Asserts]
status >= 400
status < 500

