# ── POST /api/admin/teams/{id}/members - [mass_assignment] privilege probe ──
# case_id=TC-850dd902
# case_name=POST /api/admin/teams/{id}/members - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "admin": true,
  "isAdmin": true,
  "is_admin": true,
  "role": "__probe__",
  "userId": "b21cab01-ede4-49da-9080-18aced242f70"
}
```

HTTP 400

