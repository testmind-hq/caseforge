# ── PUT /api/admin/users/{id} - mutation: role integer instead of string ──
# case_id=TC-1d2d0cbd
# case_name=PUT /api/admin/users/{id} - mutation: role integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": 12345
}
```

HTTP *

[Asserts]
status >= 400
status < 500

