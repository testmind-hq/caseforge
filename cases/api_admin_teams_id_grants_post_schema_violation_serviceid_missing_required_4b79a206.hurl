# ── POST /api/admin/teams/{id}/grants - [schema_violation] serviceId_missing_required ──
# case_id=TC-4b79a206
# case_name=POST /api/admin/teams/{id}/grants - [schema_violation] serviceId_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "am"
  ],
  "expiresAt": "1970-08-02T20:53:06Z",
  "granteeTeamId": "7a8e7c06-efab-4a89-8471-23bbf2a20eea",
  "granteeUserId": "55b411ae-4ae9-4cf6-802a-a4a242203443"
}
```

HTTP 422

