# ── Create user with invalid email format ───
# case_id=TC-12d150e0
# case_name=Create user with invalid email format
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "invalid-email",
  "name": "Test User",
  "password": "Password123"
}
```

HTTP 400

