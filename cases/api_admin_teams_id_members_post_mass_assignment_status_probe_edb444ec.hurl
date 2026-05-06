# ── POST /api/admin/teams/{id}/members - [mass_assignment] status probe ──
# case_id=TC-edb444ec
# case_name=POST /api/admin/teams/{id}/members - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "role": "member",
  "userId": "b21cab01-ede4-49da-9080-18aced242f70",
  "verified": true
}
```

HTTP 400

