# ── Fail to create user with invalid email format ──
# case_id=TC-1ba1acf6
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1
# title=Fail to create user with invalid email format

POST {{base_url}}/users
```json
{
  "email": "invalid-email",
  "name": "John Doe",
  "password": "SecurePass123!"
}
```

HTTP 400

