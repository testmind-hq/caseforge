# ── PUT /api/admin/teams/{id}/members/{userId} - mutation: role empty string ──
# case_id=TC-9334c130
# case_name=PUT /api/admin/teams/{id}/members/{userId} - mutation: role empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

