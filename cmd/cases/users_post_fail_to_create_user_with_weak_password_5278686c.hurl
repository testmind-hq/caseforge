# ── Fail to create user with weak password ──
# case_id=TC-5278686c
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2
# title=Fail to create user with weak password

POST {{base_url}}/users
```json
{
  "email": "john.weak@example.com",
  "name": "John Doe",
  "password": "123"
}
```

HTTP 400

