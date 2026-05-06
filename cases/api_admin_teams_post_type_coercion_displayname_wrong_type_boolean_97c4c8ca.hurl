# ── POST /api/admin/teams - [type_coercion] displayName wrong_type_boolean ──
# case_id=TC-97c4c8ca
# case_name=POST /api/admin/teams - [type_coercion] displayName wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Map the happy path through part.",
  "displayName": true,
  "name": "Ardith Cole"
}
```

HTTP 422

