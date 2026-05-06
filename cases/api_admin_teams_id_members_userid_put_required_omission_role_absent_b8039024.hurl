# ── PUT /api/admin/teams/{id}/members/{userId} - [required_omission] role absent ──
# case_id=TC-b8039024
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [required_omission] role absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{}
```

HTTP *

[Asserts]
status >= 400
status < 500

