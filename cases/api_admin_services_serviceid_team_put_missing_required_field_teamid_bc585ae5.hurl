# ── PUT /api/admin/services/{serviceId}/team - missing required field "teamId" ──
# case_id=TC-bc585ae5
# case_name=PUT /api/admin/services/{serviceId}/team - missing required field "teamId"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{}
```

HTTP 422

