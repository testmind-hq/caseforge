# ── POST /api/admin/teams - [mass_assignment] financial probe ──
# case_id=TC-3c2025cc
# case_name=POST /api/admin/teams - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "description": "Prefer predictable group over surprising thing.",
  "discount": 0,
  "displayName": "tensely",
  "name": "Jalen Lyons",
  "price": 1
}
```

HTTP 400

