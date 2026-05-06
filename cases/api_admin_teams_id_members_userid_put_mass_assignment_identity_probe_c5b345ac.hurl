# ── PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] identity probe ──
# case_id=TC-c5b345ac
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "ownerId": "__probe__",
  "role": "member",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

