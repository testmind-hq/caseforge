# ── POST /api/admin/teams - [type_coercion] description wrong_type_integer ──
# case_id=TC-1aea557e
# case_name=POST /api/admin/teams - [type_coercion] description wrong_type_integer
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": 123,
  "displayName": "yet",
  "name": "Ardith Cole"
}
```

HTTP 422

