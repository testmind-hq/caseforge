# ── PUT /api/admin/teams/{id}/members/{userId} - mutation: role null value ──
# case_id=TC-8380cf38
# case_name=PUT /api/admin/teams/{id}/members/{userId} - mutation: role null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

