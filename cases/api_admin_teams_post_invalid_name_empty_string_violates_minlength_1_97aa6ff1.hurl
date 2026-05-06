# ── POST /api/admin/teams - invalid name: empty string violates minLength 1 ──
# case_id=TC-97aa6ff1
# case_name=POST /api/admin/teams - invalid name: empty string violates minLength 1
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Track thing over time weekly.",
  "displayName": "everybody",
  "name": ""
}
```

HTTP 422

