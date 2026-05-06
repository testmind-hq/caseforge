# ── PUT /api/admin/teams/{id} - mutation: description oversized string (300 chars) ──
# case_id=TC-68ace4a3
# case_name=PUT /api/admin/teams/{id} - mutation: description oversized string (300 chars)
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "displayName": "shall"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

