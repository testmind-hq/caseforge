# ── POST /api/tokens - [mass_assignment] identity probe ──
# case_id=TC-b46880dc
# case_name=POST /api/tokens - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "name": "Jalen Phillips",
  "ownerId": "__probe__",
  "scope": "write",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

