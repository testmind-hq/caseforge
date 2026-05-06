# ── Fail to create user with missing required fields ──
# case_id=TC-8eba8f6c
# case_name=Fail to create user with missing required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "name": "John Doe",
  "password": "SecurePass123!"
}
```

HTTP 400

