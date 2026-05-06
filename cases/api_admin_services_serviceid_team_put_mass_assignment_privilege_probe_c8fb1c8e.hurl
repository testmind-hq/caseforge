# ── PUT /api/admin/services/{serviceId}/team - [mass_assignment] privilege probe ──
# case_id=TC-c8fb1c8e
# case_name=PUT /api/admin/services/{serviceId}/team - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "admin": true,
  "isAdmin": true,
  "is_admin": true,
  "role": "__probe__",
  "teamId": "205575fc-05ed-461e-8bb1-47206ee3fe2a"
}
```

HTTP 400

