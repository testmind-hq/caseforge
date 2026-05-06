# ── POST /auth/login - mutation: email null value ──
# case_id=TC-b5693707
# case_name=POST /auth/login - mutation: email null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": null,
  "password": "staff"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

