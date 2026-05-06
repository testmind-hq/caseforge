# ── POST /api/admin/teams/{id}/members - [mass_assignment] identity probe ──
# case_id=TC-09f9b8eb
# case_name=POST /api/admin/teams/{id}/members - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
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

