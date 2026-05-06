# ── POST /api/admin/teams - [mass_assignment] privilege probe ──
# case_id=TC-ed2bac60
# case_name=POST /api/admin/teams - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "admin": true,
  "description": "Prefer predictable group over surprising thing.",
  "displayName": "tensely",
  "isAdmin": true,
  "is_admin": true,
  "name": "Jalen Lyons",
  "role": "__probe__"
}
```

HTTP 400

