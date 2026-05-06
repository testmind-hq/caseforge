# ── Fail to create user with invalid email format ──
# case_id=TC-984e56e9
# case_name=Fail to create user with invalid email format
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "not-an-email",
  "name": "Invalid User",
  "password": "SecurePass123!"
}
```

HTTP 400

