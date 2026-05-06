# ── PUT /api/admin/teams/{id}/members/{userId} - mutation: role oversized string (300 chars) ──
# case_id=TC-c4c6cb7f
# case_name=PUT /api/admin/teams/{id}/members/{userId} - mutation: role oversized string (300 chars)
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

