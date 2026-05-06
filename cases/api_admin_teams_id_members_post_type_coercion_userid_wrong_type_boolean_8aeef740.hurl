# ── POST /api/admin/teams/{id}/members - [type_coercion] userId wrong_type_boolean ──
# case_id=TC-8aeef740
# case_name=POST /api/admin/teams/{id}/members - [type_coercion] userId wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": true
}
```

HTTP 422

