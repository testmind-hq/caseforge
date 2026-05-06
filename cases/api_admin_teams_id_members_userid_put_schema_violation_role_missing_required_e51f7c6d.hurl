# ── PUT /api/admin/teams/{id}/members/{userId} - [schema_violation] role_missing_required ──
# case_id=TC-e51f7c6d
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [schema_violation] role_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{}
```

HTTP 422

