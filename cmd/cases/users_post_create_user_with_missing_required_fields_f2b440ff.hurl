# ── Create user with missing required fields ──
# case_id=TC-f2b440ff
# case_name=Create user with missing required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{
  "name": "John Doe"
}
```

HTTP 400

