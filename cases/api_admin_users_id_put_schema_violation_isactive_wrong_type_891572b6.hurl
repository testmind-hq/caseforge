# ── PUT /api/admin/users/{id} - [schema_violation] isActive_wrong_type ──
# case_id=TC-891572b6
# case_name=PUT /api/admin/users/{id} - [schema_violation] isActive_wrong_type
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": "not_a_boolean",
  "role": "team_owner"
}
```

HTTP 422

