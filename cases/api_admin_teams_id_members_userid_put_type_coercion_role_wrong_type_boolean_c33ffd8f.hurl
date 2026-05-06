# ── PUT /api/admin/teams/{id}/members/{userId} - [type_coercion] role wrong_type_boolean ──
# case_id=TC-c33ffd8f
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [type_coercion] role wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": true
}
```

HTTP 422

