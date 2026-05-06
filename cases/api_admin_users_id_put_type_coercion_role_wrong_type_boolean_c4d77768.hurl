# ── PUT /api/admin/users/{id} - [type_coercion] role wrong_type_boolean ──
# case_id=TC-c4d77768
# case_name=PUT /api/admin/users/{id} - [type_coercion] role wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": true
}
```

HTTP 422

