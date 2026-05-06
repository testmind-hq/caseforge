# ── Fail to create user with weak password ──
# case_id=TC-3cf31478
# case_name=Fail to create user with weak password
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "charlie@example.com",
  "name": "Charlie Brown",
  "password": "123"
}
```

HTTP 400

