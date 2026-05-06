# ── PUT /api/admin/users/{id} - [schema_violation] role_invalid_enum ──
# case_id=TC-3765a2be
# case_name=PUT /api/admin/users/{id} - [schema_violation] role_invalid_enum
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": true,
  "role": "__invalid__"
}
```

HTTP 422

