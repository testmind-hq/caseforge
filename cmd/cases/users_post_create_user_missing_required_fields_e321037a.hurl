# ── Create user missing required fields ─────
# case_id=TC-e321037a
# case_name=Create user missing required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "name": "Jane Doe",
  "password": "SecurePass123!"
}
```

HTTP 400

