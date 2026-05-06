# ── POST /api/admin/teams - mutation: displayName null value ──
# case_id=TC-782f4da8
# case_name=POST /api/admin/teams - mutation: displayName null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "At this point the review, you want the number.",
  "displayName": null,
  "name": "Alysson Tucker"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

