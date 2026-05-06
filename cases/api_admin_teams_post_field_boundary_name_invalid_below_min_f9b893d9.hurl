# ── POST /api/admin/teams - [field_boundary] name invalid_below_min ──
# case_id=TC-f9b893d9
# case_name=POST /api/admin/teams - [field_boundary] name invalid_below_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "The lingering fact been unexpectedly tensely.",
  "displayName": "yours",
  "name": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

