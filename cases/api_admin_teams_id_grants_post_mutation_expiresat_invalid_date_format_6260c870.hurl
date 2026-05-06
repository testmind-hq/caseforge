# ── POST /api/admin/teams/{id}/grants - mutation: expiresAt invalid date format ──
# case_id=TC-6260c870
# case_name=POST /api/admin/teams/{id}/grants - mutation: expiresAt invalid date format
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
  "expiresAt": "not-a-date",
  "granteeTeamId": "7147d6bf-cb22-4db1-b2c4-dfaecc9b5353",
  "granteeUserId": "75344e9e-e685-43e6-adbc-d314f39096cb",
  "serviceId": "b5523acd-dea4-47c8-826c-23ff81d46ccf"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

