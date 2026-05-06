# ── PUT /api/admin/users/{id} - mutation: role null value ──
# case_id=TC-091acd05
# case_name=PUT /api/admin/users/{id} - mutation: role null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

