# ── POST /api/admin/teams - [required_omission] name absent ──
# case_id=TC-7a6a3b1a
# case_name=POST /api/admin/teams - [required_omission] name absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Sample week at 11s intervals.",
  "displayName": "annually"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

