# ── POST /auth/login - missing required field "email" ──
# case_id=TC-4cc99b0c
# case_name=POST /auth/login - missing required field "email"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "password": "fuel"
}
```

HTTP 422

