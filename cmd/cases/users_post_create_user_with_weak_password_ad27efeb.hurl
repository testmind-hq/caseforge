# ── Create user with weak password ──────────
# case_id=TC-ad27efeb
# case_name=Create user with weak password
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2

POST {{base_url}}/users
```json
{
  "email": "charlie.wilson@example.com",
  "name": "Charlie Wilson",
  "password": "123"
}
```

HTTP 400

