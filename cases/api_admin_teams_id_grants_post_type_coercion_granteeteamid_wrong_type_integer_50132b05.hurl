# ── POST /api/admin/teams/{id}/grants - [type_coercion] granteeTeamId wrong_type_integer ──
# case_id=TC-50132b05
# case_name=POST /api/admin/teams/{id}/grants - [type_coercion] granteeTeamId wrong_type_integer
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
  "granteeTeamId": 123,
  "granteeUserId": "a5671961-7609-4c30-865c-7e0cf0f5a413",
  "serviceId": "5d4bc090-b3d1-49ab-83a2-4855a8fbaee8"
}
```

HTTP 422

