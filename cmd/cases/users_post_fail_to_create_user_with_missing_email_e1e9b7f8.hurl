# ── Fail to create user with missing email ──
# case_id=TC-e1e9b7f8
# case_name=Fail to create user with missing email
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "name": "No Email User",
  "password": "Password123"
}
```

HTTP 400

