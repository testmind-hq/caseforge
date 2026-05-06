# ── Fail to create user with empty request body ──
# case_id=TC-cea3990a
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1
# title=Fail to create user with empty request body

POST {{base_url}}/users
```json
{}
```

HTTP 400

