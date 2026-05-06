# ── PUT /api/admin/teams/{id} - [type_coercion] description wrong_type_boolean ──
# case_id=TC-6dd640a7
# case_name=PUT /api/admin/teams/{id} - [type_coercion] description wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": true,
  "displayName": "addition"
}
```

HTTP 422

