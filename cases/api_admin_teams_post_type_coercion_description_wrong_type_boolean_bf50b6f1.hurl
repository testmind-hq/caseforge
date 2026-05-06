# ── POST /api/admin/teams - [type_coercion] description wrong_type_boolean ──
# case_id=TC-bf50b6f1
# case_name=POST /api/admin/teams - [type_coercion] description wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": true,
  "displayName": "yet",
  "name": "Ardith Cole"
}
```

HTTP 422

