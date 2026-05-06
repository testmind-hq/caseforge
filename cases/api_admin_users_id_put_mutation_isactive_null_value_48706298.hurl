# ── PUT /api/admin/users/{id} - mutation: isActive null value ──
# case_id=TC-48706298
# case_name=PUT /api/admin/users/{id} - mutation: isActive null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": null,
  "role": "super_admin"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

