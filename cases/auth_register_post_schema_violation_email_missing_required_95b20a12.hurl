# ── POST /auth/register - [schema_violation] email_missing_required ──
# case_id=TC-95b20a12
# case_name=POST /auth/register - [schema_violation] email_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "password": "these"
}
```

HTTP 422

