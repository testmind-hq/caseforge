# ── POST /api/admin/teams/{id}/members - [mass_assignment] financial probe ──
# case_id=TC-31f44a55
# case_name=POST /api/admin/teams/{id}/members - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "price": 1,
  "role": "member",
  "userId": "b21cab01-ede4-49da-9080-18aced242f70"
}
```

HTTP 400

