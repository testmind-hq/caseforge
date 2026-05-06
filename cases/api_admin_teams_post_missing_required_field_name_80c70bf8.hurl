# ── POST /api/admin/teams - missing required field "name" ──
# case_id=TC-80c70bf8
# case_name=POST /api/admin/teams - missing required field "name"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Track thing over time weekly.",
  "displayName": "everybody"
}
```

HTTP 422

