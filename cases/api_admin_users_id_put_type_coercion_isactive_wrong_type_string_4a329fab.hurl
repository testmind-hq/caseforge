# ── PUT /api/admin/users/{id} - [type_coercion] isActive wrong_type_string ──
# case_id=TC-4a329fab
# case_name=PUT /api/admin/users/{id} - [type_coercion] isActive wrong_type_string
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": "not_a_boolean",
  "role": "super_admin"
}
```

HTTP 422

