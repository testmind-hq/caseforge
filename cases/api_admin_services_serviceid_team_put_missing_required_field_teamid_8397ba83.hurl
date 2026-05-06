# ── PUT /api/admin/services/{serviceId}/team - missing required field "teamId" ──
# case_id=TC-8397ba83
# case_name=PUT /api/admin/services/{serviceId}/team - missing required field "teamId"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{}
```

HTTP 422

