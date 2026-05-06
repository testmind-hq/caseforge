# ── POST /api/admin/teams - [mass_assignment] identity probe ──
# case_id=TC-82f380ef
# case_name=POST /api/admin/teams - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "description": "Prefer predictable group over surprising thing.",
  "displayName": "tensely",
  "name": "Jalen Lyons",
  "ownerId": "__probe__",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

