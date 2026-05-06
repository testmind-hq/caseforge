# ── PUT /api/admin/services/{serviceId}/team - [mass_assignment] identity probe ──
# case_id=TC-c9fe2f6f
# case_name=PUT /api/admin/services/{serviceId}/team - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "ownerId": "__probe__",
  "teamId": "205575fc-05ed-461e-8bb1-47206ee3fe2a",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

