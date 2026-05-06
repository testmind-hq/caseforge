# ── POST /api/admin/teams - mutation: description integer instead of string ──
# case_id=TC-4d295fcc
# case_name=POST /api/admin/teams - mutation: description integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": 12345,
  "displayName": "his",
  "name": "Alysson Tucker"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

