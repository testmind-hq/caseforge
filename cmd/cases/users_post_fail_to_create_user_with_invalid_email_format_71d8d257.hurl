# ── Fail to create user with invalid email format ──
# case_id=TC-71d8d257
# case_name=Fail to create user with invalid email format
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "not-an-email",
  "name": "Bad Email User",
  "password": "Password123"
}
```

HTTP 400

