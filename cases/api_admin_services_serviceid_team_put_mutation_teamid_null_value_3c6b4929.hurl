# ── PUT /api/admin/services/{serviceId}/team - mutation: teamId null value ──
# case_id=TC-3c6b4929
# case_name=PUT /api/admin/services/{serviceId}/team - mutation: teamId null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": null
}
```

HTTP *

[Asserts]
status >= 400
status < 500

