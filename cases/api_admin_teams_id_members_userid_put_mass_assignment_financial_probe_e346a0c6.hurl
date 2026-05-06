# ── PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] financial probe ──
# case_id=TC-e346a0c6
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "price": 1,
  "role": "member"
}
```

HTTP 400

