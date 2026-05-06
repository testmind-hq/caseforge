# ── POST /auth/login - [schema_violation] password_missing_required ──
# case_id=TC-5bddd51c
# case_name=POST /auth/login - [schema_violation] password_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "email": "clovissoto@clay.io"
}
```

HTTP 422

