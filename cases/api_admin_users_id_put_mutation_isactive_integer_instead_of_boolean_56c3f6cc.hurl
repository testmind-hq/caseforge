# ── PUT /api/admin/users/{id} - mutation: isActive integer instead of boolean ──
# case_id=TC-56c3f6cc
# case_name=PUT /api/admin/users/{id} - mutation: isActive integer instead of boolean
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": 1,
  "role": "super_admin"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

