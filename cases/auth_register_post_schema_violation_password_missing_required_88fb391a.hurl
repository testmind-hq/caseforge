# ── POST /auth/register - [schema_violation] password_missing_required ──
# case_id=TC-88fb391a
# case_name=POST /auth/register - [schema_violation] password_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "jadonrobertson@wu.org"
}
```

HTTP 422

