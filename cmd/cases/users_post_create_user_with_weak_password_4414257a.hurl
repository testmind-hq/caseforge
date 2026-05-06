# ── Create user with weak password ──────────
# case_id=TC-4414257a
# case_name=Create user with weak password
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2

POST {{base_url}}/users
```json
{
  "email": "john.weak@example.com",
  "name": "John Doe",
  "password": "123"
}
```

HTTP 400

