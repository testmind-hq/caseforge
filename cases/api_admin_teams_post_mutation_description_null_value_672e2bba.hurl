# ── POST /api/admin/teams - mutation: description null value ──
# case_id=TC-672e2bba
# case_name=POST /api/admin/teams - mutation: description null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": null,
  "displayName": "his",
  "name": "Alysson Tucker"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

