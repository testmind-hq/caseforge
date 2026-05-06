# ── PUT /api/admin/teams/{id} - mutation: description null value ──
# case_id=TC-df8e9c3a
# case_name=PUT /api/admin/teams/{id} - mutation: description null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": null,
  "displayName": "shall"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

