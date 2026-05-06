# ── POST /auth/register - missing required field "password" ──
# case_id=TC-31707ae5
# case_name=POST /auth/register - missing required field "password"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "email": "leahawkins@white.io"
}
```

HTTP 422

