# ── PUT /api/admin/teams/{id} - mutation: description integer instead of string ──
# case_id=TC-f0d62caa
# case_name=PUT /api/admin/teams/{id} - mutation: description integer instead of string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": 12345,
  "displayName": "shall"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

