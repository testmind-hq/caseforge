# ── PUT /api/admin/teams/{id} - mutation: displayName empty string ──
# case_id=TC-13a9f6ae
# case_name=PUT /api/admin/teams/{id} - mutation: displayName empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "First of all, document the company and specify the rest.",
  "displayName": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

