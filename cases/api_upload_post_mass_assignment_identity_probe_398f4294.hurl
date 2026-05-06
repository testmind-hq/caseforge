# ── POST /api/upload - [mass_assignment] identity probe ──
# case_id=TC-398f4294
# case_name=POST /api/upload - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "oops",
  "commitSha": "mustering",
  "createdBy": "__probe__",
  "ownerId": "__probe__",
  "service": "I",
  "specContent": "cut",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

