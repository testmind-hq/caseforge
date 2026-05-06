# ── PUT /api/admin/users/{id} - [type_coercion] isActive wrong_type_integer ──
# case_id=TC-308337db
# case_name=PUT /api/admin/users/{id} - [type_coercion] isActive wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": 1,
  "role": "super_admin"
}
```

HTTP 422

