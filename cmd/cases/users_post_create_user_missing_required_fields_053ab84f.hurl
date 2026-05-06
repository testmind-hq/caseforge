# ── Create user missing required fields ─────
# case_id=TC-053ab84f
# case_name=Create user missing required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{
  "name": "John Doe"
}
```

HTTP 400

