# ── POST /auth/register - missing required field "email" ──
# case_id=TC-445d8b1f
# case_name=POST /auth/register - missing required field "email"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "password": "still"
}
```

HTTP 422

