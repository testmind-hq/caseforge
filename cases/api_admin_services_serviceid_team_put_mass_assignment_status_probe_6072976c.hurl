# ── PUT /api/admin/services/{serviceId}/team - [mass_assignment] status probe ──
# case_id=TC-6072976c
# case_name=PUT /api/admin/services/{serviceId}/team - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "teamId": "205575fc-05ed-461e-8bb1-47206ee3fe2a",
  "verified": true
}
```

HTTP 400

