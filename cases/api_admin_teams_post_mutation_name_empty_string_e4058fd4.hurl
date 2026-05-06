# ── POST /api/admin/teams - mutation: name empty string ──
# case_id=TC-e4058fd4
# case_name=POST /api/admin/teams - mutation: name empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "At this point the review, you want the number.",
  "displayName": "his",
  "name": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

