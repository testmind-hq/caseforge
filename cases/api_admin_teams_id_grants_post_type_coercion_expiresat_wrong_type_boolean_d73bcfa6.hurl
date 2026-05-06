# ── POST /api/admin/teams/{id}/grants - [type_coercion] expiresAt wrong_type_boolean ──
# case_id=TC-d73bcfa6
# case_name=POST /api/admin/teams/{id}/grants - [type_coercion] expiresAt wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "bad"
  ],
  "expiresAt": true,
  "granteeTeamId": "caab12dc-a7c5-4aae-b5a0-a7bac0b2deab",
  "granteeUserId": "a5671961-7609-4c30-865c-7e0cf0f5a413",
  "serviceId": "5d4bc090-b3d1-49ab-83a2-4855a8fbaee8"
}
```

HTTP 422

