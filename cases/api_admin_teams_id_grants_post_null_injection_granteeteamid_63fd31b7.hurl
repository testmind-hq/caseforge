# ── POST /api/admin/teams/{id}/grants - null injection: granteeTeamId ──
# case_id=TC-63fd31b7
# case_name=POST /api/admin/teams/{id}/grants - null injection: granteeTeamId
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "sari"
  ],
  "expiresAt": "1914-05-11T22:00:14Z",
  "granteeTeamId": null,
  "granteeUserId": "44099659-ceca-4310-b565-88e5257ae6f0",
  "serviceId": "4e8d3cff-ce68-4019-af70-67a1bb961ec8"
}
```

HTTP 422

