# ── PUT /api/admin/services/{serviceId}/team - [mass_assignment] financial probe ──
# case_id=TC-297a0e33
# case_name=PUT /api/admin/services/{serviceId}/team - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "price": 1,
  "teamId": "205575fc-05ed-461e-8bb1-47206ee3fe2a"
}
```

HTTP 400

