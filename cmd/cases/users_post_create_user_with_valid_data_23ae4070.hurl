# ── Create user with valid data ─────────────
# case_id=TC-23ae4070
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0
# title=Create user with valid data

POST {{base_url}}/users
```json
{
  "email": "john.doe@example.com",
  "name": "John Doe",
  "password": "SecurePass123!"
}
```

HTTP 201

