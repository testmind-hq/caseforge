# ── POST /auth/login - mutation: password oversized string (300 chars) ──
# case_id=TC-acbb9354
# case_name=POST /auth/login - mutation: password oversized string (300 chars)
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "naomipierce@lewis.biz",
  "password": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

