# ── PUT /api/admin/users/{id} - null injection: isActive ──
# case_id=TC-c8deaf48
# case_name=PUT /api/admin/users/{id} - null injection: isActive
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": null,
  "role": "super_admin"
}
```

HTTP 422

