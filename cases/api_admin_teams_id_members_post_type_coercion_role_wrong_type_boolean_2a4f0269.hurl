# ── POST /api/admin/teams/{id}/members - [type_coercion] role wrong_type_boolean ──
# case_id=TC-2a4f0269
# case_name=POST /api/admin/teams/{id}/members - [type_coercion] role wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": true,
  "userId": "8aa00d9d-7b81-42a4-830e-092302d2f2c4"
}
```

HTTP 422

