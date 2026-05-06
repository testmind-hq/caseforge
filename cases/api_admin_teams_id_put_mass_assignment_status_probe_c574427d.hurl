# ── PUT /api/admin/teams/{id} - [mass_assignment] status probe ──
# case_id=TC-c574427d
# case_name=PUT /api/admin/teams/{id} - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "description": "Alert on way thresholds yesterday.",
  "disabled": false,
  "displayName": "this",
  "verified": true
}
```

HTTP 400

