# ── PUT /api/admin/users/{id} - invalid role: value not in enum ──
# case_id=TC-be8b477d
# case_name=PUT /api/admin/users/{id} - invalid role: value not in enum
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": true,
  "role": "__invalid_enum__"
}
```

HTTP 422

