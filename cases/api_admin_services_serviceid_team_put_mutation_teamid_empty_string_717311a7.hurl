# ── PUT /api/admin/services/{serviceId}/team - mutation: teamId empty string ──
# case_id=TC-717311a7
# case_name=PUT /api/admin/services/{serviceId}/team - mutation: teamId empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": ""
}
```

HTTP *

[Asserts]
status >= 400
status < 500

