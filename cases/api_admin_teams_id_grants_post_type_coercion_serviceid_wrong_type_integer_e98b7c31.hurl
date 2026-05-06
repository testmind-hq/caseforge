# ── POST /api/admin/teams/{id}/grants - [type_coercion] serviceId wrong_type_integer ──
# case_id=TC-e98b7c31
# case_name=POST /api/admin/teams/{id}/grants - [type_coercion] serviceId wrong_type_integer
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
  "expiresAt": "2013-09-12T21:41:49Z",
  "granteeTeamId": "caab12dc-a7c5-4aae-b5a0-a7bac0b2deab",
  "granteeUserId": "a5671961-7609-4c30-865c-7e0cf0f5a413",
  "serviceId": 123
}
```

HTTP 422

