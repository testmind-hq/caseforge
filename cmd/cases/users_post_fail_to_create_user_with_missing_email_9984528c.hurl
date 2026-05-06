# ── Fail to create user with missing email ──
# case_id=TC-9984528c
# case_name=Fail to create user with missing email
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{
  "name": "Jane Doe",
  "password": "SecurePass123!"
}
```

HTTP 400

