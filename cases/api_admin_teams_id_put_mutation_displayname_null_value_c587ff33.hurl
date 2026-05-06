# ── PUT /api/admin/teams/{id} - mutation: displayName null value ──
# case_id=TC-c587ff33
# case_name=PUT /api/admin/teams/{id} - mutation: displayName null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "First of all, document the company and specify the rest.",
  "displayName": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

