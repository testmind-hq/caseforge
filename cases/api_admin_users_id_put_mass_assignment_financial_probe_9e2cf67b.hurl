# ── PUT /api/admin/users/{id} - [mass_assignment] financial probe ──
# case_id=TC-9e2cf67b
# case_name=PUT /api/admin/users/{id} - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "isActive": true,
  "price": 1,
  "role": "super_admin"
}
```

HTTP 400

