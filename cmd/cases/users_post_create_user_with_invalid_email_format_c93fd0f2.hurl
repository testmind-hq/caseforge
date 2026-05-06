# ── Create user with invalid email format ───
# case_id=TC-c93fd0f2
# case_name=Create user with invalid email format
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "invalid-email",
  "name": "Jane Doe",
  "password": "SecurePass123!"
}
```

HTTP 400

