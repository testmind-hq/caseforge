# ── PUT /api/admin/teams/{id} - [mass_assignment] financial probe ──
# case_id=TC-4c631268
# case_name=PUT /api/admin/teams/{id} - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "description": "Alert on way thresholds yesterday.",
  "discount": 0,
  "displayName": "this",
  "price": 1
}
```

HTTP 400

