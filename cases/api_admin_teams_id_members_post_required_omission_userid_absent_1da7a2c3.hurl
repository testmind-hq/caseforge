# ── POST /api/admin/teams/{id}/members - [required_omission] userId absent ──
# case_id=TC-1da7a2c3
# case_name=POST /api/admin/teams/{id}/members - [required_omission] userId absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "owner"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

