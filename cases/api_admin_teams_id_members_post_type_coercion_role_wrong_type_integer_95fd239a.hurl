# ── POST /api/admin/teams/{id}/members - [type_coercion] role wrong_type_integer ──
# case_id=TC-95fd239a
# case_name=POST /api/admin/teams/{id}/members - [type_coercion] role wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": 123,
  "userId": "8aa00d9d-7b81-42a4-830e-092302d2f2c4"
}
```

HTTP 422

