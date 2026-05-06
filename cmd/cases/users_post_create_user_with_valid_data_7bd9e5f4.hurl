# ── Create user with valid data ─────────────
# case_id=TC-7bd9e5f4
# case_name=Create user with valid data
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{
  "email": "john.doe@example.com",
  "name": "John Doe",
  "password": "SecurePass123!"
}
```

HTTP 201

