# ── PUT /api/admin/teams/{id} - [type_coercion] displayName wrong_type_integer ──
# case_id=TC-3ade9411
# case_name=PUT /api/admin/teams/{id} - [type_coercion] displayName wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Visualize hand for faster decisions.",
  "displayName": 123
}
```

HTTP 422

