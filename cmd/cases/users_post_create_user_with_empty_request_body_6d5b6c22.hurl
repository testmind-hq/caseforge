# ── Create user with empty request body ─────
# case_id=TC-6d5b6c22
# case_name=Create user with empty request body
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2

POST {{base_url}}/users
```json
{}
```

HTTP 400

