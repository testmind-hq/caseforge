# ── Create user with weak password ──────────
# case_id=TC-61182975
# case_name=Create user with weak password
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

