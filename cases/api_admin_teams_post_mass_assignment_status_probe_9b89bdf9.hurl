# ── POST /api/admin/teams - [mass_assignment] status probe ──
# case_id=TC-9b89bdf9
# case_name=POST /api/admin/teams - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "description": "Prefer predictable group over surprising thing.",
  "disabled": false,
  "displayName": "tensely",
  "name": "Jalen Lyons",
  "verified": true
}
```

HTTP 400

