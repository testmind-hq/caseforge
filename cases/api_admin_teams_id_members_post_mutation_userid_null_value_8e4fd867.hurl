# ── POST /api/admin/teams/{id}/members - mutation: userId null value ──
# case_id=TC-8e4fd867
# case_name=POST /api/admin/teams/{id}/members - mutation: userId null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

