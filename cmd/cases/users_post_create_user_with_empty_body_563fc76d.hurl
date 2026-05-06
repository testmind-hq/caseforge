# ── Create user with empty body ─────────────
# case_id=TC-563fc76d
# case_name=Create user with empty body
# step_id=step-1
# step_type=test
# technique=ask
# priority=P1

POST {{base_url}}/users
```json
{}
```

HTTP 400

