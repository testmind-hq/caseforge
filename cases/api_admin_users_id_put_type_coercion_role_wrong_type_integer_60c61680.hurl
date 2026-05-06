# ── PUT /api/admin/users/{id} - [type_coercion] role wrong_type_integer ──
# case_id=TC-60c61680
# case_name=PUT /api/admin/users/{id} - [type_coercion] role wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": 123
}
```

HTTP 422

