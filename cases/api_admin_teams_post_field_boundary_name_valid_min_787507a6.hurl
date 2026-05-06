# ── POST /api/admin/teams - [field_boundary] name valid_min ──
# case_id=TC-787507a6
# case_name=POST /api/admin/teams - [field_boundary] name valid_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P1

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Guard world with sensible limits.",
  "displayName": "those",
  "name": "a"
}
```

HTTP *

[Asserts]
status >= 200
status < 300

