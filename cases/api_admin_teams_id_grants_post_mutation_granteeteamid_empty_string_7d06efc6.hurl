# ── POST /api/admin/teams/{id}/grants - mutation: granteeTeamId empty string ──
# case_id=TC-7d06efc6
# case_name=POST /api/admin/teams/{id}/grants - mutation: granteeTeamId empty string
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "himself"
  ],
  "expiresAt": "2008-02-06T15:08:34Z",
  "granteeTeamId": "",
  "granteeUserId": "75344e9e-e685-43e6-adbc-d314f39096cb",
  "serviceId": "b5523acd-dea4-47c8-826c-23ff81d46ccf"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

