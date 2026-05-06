# ── PUT /api/admin/teams/{id} - mutation: displayName integer instead of string ──
# case_id=TC-05b44595
# case_name=PUT /api/admin/teams/{id} - mutation: displayName integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "First of all, document the company and specify the rest.",
  "displayName": 12345
}
```

HTTP *

[Asserts]
status >= 400
status < 500

