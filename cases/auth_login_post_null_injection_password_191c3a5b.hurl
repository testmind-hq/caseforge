# ── POST /auth/login - null injection: password ──
# case_id=TC-191c3a5b
# case_name=POST /auth/login - null injection: password
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "ottonorris@sullivan.com",
  "password": null
}
```

HTTP 422

