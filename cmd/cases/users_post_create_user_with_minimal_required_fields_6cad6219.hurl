# ── Create user with minimal required fields ──
# case_id=TC-6cad6219
# case_name=Create user with minimal required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{
  "email": "minimal@example.com",
  "password": "Password123"
}
```

HTTP 201

