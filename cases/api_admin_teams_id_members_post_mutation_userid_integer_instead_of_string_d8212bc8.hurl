# ── POST /api/admin/teams/{id}/members - mutation: userId integer instead of string ──
# case_id=TC-d8212bc8
# case_name=POST /api/admin/teams/{id}/members - mutation: userId integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": 12345
}
```

HTTP *

[Asserts]
status >= 400
status < 500

