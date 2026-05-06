# ── Fail to create user with invalid email format ──
# case_id=TC-a2bd888d
# case_name=Fail to create user with invalid email format
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "invalid-email",
  "name": "Bob Smith",
  "password": "SecurePass123!"
}
```

HTTP 400

