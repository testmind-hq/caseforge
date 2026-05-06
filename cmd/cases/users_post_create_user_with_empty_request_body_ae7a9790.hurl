# ── Create user with empty request body ─────
# case_id=TC-ae7a9790
# case_name=Create user with empty request body
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{}
```

HTTP 400

