# ── POST /auth/login - null injection: email ──
# case_id=TC-a1de0446
# case_name=POST /auth/login - null injection: email
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": null,
  "password": "float"
}
```

HTTP 422

