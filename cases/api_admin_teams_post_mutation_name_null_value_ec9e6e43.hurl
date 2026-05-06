# ── POST /api/admin/teams - mutation: name null value ──
# case_id=TC-ec9e6e43
# case_name=POST /api/admin/teams - mutation: name null value
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
  "name": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

