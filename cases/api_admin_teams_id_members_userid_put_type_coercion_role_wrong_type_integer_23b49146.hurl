# ── PUT /api/admin/teams/{id}/members/{userId} - [type_coercion] role wrong_type_integer ──
# case_id=TC-23b49146
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [type_coercion] role wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": 123
}
```

HTTP 422

