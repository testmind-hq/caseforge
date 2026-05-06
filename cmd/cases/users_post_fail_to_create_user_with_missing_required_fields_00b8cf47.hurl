# ── Fail to create user with missing required fields ──
# case_id=TC-00b8cf47
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0
# title=Fail to create user with missing required fields

POST {{base_url}}/users
```json
{
  "name": "John Doe",
  "password": "SecurePass123!"
}
```

HTTP 400

