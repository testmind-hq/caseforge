# ── PUT /api/admin/teams/{id}/members/{userId} - [schema_violation] role_invalid_enum ──
# case_id=TC-128b22a3
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [schema_violation] role_invalid_enum
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "__invalid__"
}
```

HTTP 422

