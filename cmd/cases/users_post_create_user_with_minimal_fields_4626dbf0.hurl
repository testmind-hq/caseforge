# ── Create user with minimal fields ─────────
# case_id=TC-4626dbf0
# case_name=Create user with minimal fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{
  "email": "minimal@example.com"
}
```

HTTP 201

