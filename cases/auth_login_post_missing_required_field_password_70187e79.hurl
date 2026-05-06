# ── POST /auth/login - missing required field "password" ──
# case_id=TC-70187e79
# case_name=POST /auth/login - missing required field "password"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "montemendez@campbell.name"
}
```

HTTP 422

