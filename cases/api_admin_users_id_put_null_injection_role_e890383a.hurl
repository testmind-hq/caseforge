# ── PUT /api/admin/users/{id} - null injection: role ──
# case_id=TC-e890383a
# case_name=PUT /api/admin/users/{id} - null injection: role
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": null
}
```

HTTP 422

