# ── Fail to create user with empty request body ──
# case_id=TC-9787221a
# case_name=Fail to create user with empty request body
# step_id=step-1
# step_type=test
# technique=ask
# priority=P2

POST {{base_url}}/users
```json
{}
```

HTTP 400

