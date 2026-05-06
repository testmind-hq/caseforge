# ── PUT /api/admin/users/{id} - [mass_assignment] status probe ──
# case_id=TC-1054f864
# case_name=PUT /api/admin/users/{id} - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "isActive": true,
  "role": "super_admin",
  "verified": true
}
```

HTTP 400

