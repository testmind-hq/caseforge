# ── POST /api/admin/teams - mutation: displayName integer instead of string ──
# case_id=TC-c361779d
# case_name=POST /api/admin/teams - mutation: displayName integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "At this point the review, you want the number.",
  "displayName": 12345,
  "name": "Alysson Tucker"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

