# ── PUT /api/admin/teams/{id} - mutation: description empty string ──
# case_id=TC-eb263846
# case_name=PUT /api/admin/teams/{id} - mutation: description empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "",
  "displayName": "shall"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

