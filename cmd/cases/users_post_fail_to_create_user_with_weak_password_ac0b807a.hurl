# ── Fail to create user with weak password ──
# case_id=TC-ac0b807a
# case_name=Fail to create user with weak password
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2

POST {{base_url}}/users
```json
{
  "email": "weakpass@example.com",
  "name": "Weak Password User",
  "password": "123"
}
```

HTTP 400

