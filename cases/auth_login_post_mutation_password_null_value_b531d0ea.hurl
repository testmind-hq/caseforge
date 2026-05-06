# ── POST /auth/login - mutation: password null value ──
# case_id=TC-b531d0ea
# case_name=POST /auth/login - mutation: password null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "naomipierce@lewis.biz",
  "password": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

