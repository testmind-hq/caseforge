# ── Create user with all required fields ────
# case_id=TC-ca607f38
# case_name=Create user with all required fields
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

