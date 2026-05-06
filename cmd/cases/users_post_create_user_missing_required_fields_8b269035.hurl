# ── Create user missing required fields ─────
# case_id=TC-8b269035
# case_name=Create user missing required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{
  "name": "Jane Doe"
}
```

HTTP 400

