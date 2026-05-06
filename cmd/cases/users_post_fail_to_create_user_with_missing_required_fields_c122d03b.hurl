# ── Fail to create user with missing required fields ──
# case_id=TC-c122d03b
# case_name=Fail to create user with missing required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{}
```

HTTP 400

