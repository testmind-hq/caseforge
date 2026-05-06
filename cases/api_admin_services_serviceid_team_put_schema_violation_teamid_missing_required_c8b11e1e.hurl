# ── PUT /api/admin/services/{serviceId}/team - [schema_violation] teamId_missing_required ──
# case_id=TC-c8b11e1e
# case_name=PUT /api/admin/services/{serviceId}/team - [schema_violation] teamId_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{}
```

HTTP 422

