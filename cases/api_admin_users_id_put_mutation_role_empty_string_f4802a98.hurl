# ── PUT /api/admin/users/{id} - mutation: role empty string ──
# case_id=TC-f4802a98
# case_name=PUT /api/admin/users/{id} - mutation: role empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

