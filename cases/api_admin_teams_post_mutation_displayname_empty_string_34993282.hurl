# ── POST /api/admin/teams - mutation: displayName empty string ──
# case_id=TC-34993282
# case_name=POST /api/admin/teams - mutation: displayName empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "At this point the review, you want the number.",
  "displayName": "",
  "name": "Alysson Tucker"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

