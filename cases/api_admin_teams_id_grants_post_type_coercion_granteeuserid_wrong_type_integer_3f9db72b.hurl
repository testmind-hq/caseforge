# ── POST /api/admin/teams/{id}/grants - [type_coercion] granteeUserId wrong_type_integer ──
# case_id=TC-3f9db72b
# case_name=POST /api/admin/teams/{id}/grants - [type_coercion] granteeUserId wrong_type_integer
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
  "granteeUserId": 123,
  "serviceId": "5d4bc090-b3d1-49ab-83a2-4855a8fbaee8"
}
```

HTTP 422

