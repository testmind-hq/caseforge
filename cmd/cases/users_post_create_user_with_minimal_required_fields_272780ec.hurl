# ── Create user with minimal required fields ──
# case_id=TC-272780ec
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1
# title=Create user with minimal required fields

POST {{base_url}}/users
```json
{
  "email": "minimal@example.com",
  "password": "Password123"
}
```

HTTP 201

