# ── Create user missing required name field ──
# case_id=TC-20f71db2
# case_name=Create user missing required name field
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "email": "missing.name@example.com",
  "password": "SecurePass123!"
}
```

HTTP 400

