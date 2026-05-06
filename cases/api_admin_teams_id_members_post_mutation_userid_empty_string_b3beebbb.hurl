# ── POST /api/admin/teams/{id}/members - mutation: userId empty string ──
# case_id=TC-b3beebbb
# case_name=POST /api/admin/teams/{id}/members - mutation: userId empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

