# ── PUT /api/admin/users/{id} - [mass_assignment] privilege probe ──
# case_id=TC-a6a6cd31
# case_name=PUT /api/admin/users/{id} - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "admin": true,
  "isActive": true,
  "isAdmin": true,
  "is_admin": true,
  "role": "__probe__"
}
```

HTTP 400

