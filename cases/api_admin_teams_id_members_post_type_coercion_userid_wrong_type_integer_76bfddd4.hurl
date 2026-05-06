# ── POST /api/admin/teams/{id}/members - [type_coercion] userId wrong_type_integer ──
# case_id=TC-76bfddd4
# case_name=POST /api/admin/teams/{id}/members - [type_coercion] userId wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": 123
}
```

HTTP 422

