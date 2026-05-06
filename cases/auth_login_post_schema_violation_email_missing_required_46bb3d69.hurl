# ── POST /auth/login - [schema_violation] email_missing_required ──
# case_id=TC-46bb3d69
# case_name=POST /auth/login - [schema_violation] email_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "password": "eye"
}
```

HTTP 422

