# ── PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] privilege probe ──
# case_id=TC-830ae193
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "admin": true,
  "isAdmin": true,
  "is_admin": true,
  "role": "__probe__"
}
```

HTTP 400

