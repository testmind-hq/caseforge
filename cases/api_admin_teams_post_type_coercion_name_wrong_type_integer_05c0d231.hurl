# ── POST /api/admin/teams - [type_coercion] name wrong_type_integer ──
# case_id=TC-05c0d231
# case_name=POST /api/admin/teams - [type_coercion] name wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Map the happy path through part.",
  "displayName": "yet",
  "name": 123
}
```

HTTP 422

