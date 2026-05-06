# ── PUT /api/admin/users/{id} - mutation: isActive string instead of boolean ──
# case_id=TC-c83a8b69
# case_name=PUT /api/admin/users/{id} - mutation: isActive string instead of boolean
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": "yes",
  "role": "super_admin"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

