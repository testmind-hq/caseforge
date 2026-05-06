# ── POST /api/admin/teams/{id}/grants - null injection: expiresAt ──
# case_id=TC-df39db3e
# case_name=POST /api/admin/teams/{id}/grants - null injection: expiresAt
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
  "expiresAt": null,
  "granteeTeamId": "bcaeb7d9-6d53-4be0-8f2e-d1beacfc2fa1",
  "granteeUserId": "44099659-ceca-4310-b565-88e5257ae6f0",
  "serviceId": "4e8d3cff-ce68-4019-af70-67a1bb961ec8"
}
```

HTTP 422

