# ── Create user with password too short ─────
# case_id=TC-6585f31e
# case_name=Create user with password too short
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2

POST {{base_url}}/users
```json
{
  "email": "weakpass@example.com",
  "name": "Weak Pass User",
  "password": "123"
}
```

HTTP 400

