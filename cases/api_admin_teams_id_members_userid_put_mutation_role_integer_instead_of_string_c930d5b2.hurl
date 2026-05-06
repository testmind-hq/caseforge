# ── PUT /api/admin/teams/{id}/members/{userId} - mutation: role integer instead of string ──
# case_id=TC-c930d5b2
# case_name=PUT /api/admin/teams/{id}/members/{userId} - mutation: role integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": 12345
}
```

HTTP *

[Asserts]
status >= 400
status < 500

