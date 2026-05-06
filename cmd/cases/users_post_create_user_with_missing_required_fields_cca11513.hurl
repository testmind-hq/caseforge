# ── Create user with missing required fields ──
# case_id=TC-cca11513
# case_name=Create user with missing required fields
# step_id=step-1
# step_type=test
# technique=ask
# priority=P0

POST {{base_url}}/users
```json
{}
```

HTTP 400

