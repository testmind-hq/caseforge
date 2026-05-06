# ── POST /api/admin/teams - mutation: description empty string ──
# case_id=TC-569a3993
# case_name=POST /api/admin/teams - mutation: description empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "",
  "displayName": "his",
  "name": "Alysson Tucker"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

