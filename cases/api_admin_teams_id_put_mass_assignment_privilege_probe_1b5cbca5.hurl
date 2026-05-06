# ── PUT /api/admin/teams/{id} - [mass_assignment] privilege probe ──
# case_id=TC-1b5cbca5
# case_name=PUT /api/admin/teams/{id} - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "admin": true,
  "description": "Alert on way thresholds yesterday.",
  "displayName": "this",
  "isAdmin": true,
  "is_admin": true,
  "role": "__probe__"
}
```

HTTP 400

