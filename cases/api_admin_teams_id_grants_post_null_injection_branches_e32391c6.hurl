# ── POST /api/admin/teams/{id}/grants - null injection: branches ──
# case_id=TC-e32391c6
# case_name=POST /api/admin/teams/{id}/grants - null injection: branches
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": null,
  "expiresAt": "1914-05-11T22:00:14Z",
  "granteeTeamId": "bcaeb7d9-6d53-4be0-8f2e-d1beacfc2fa1",
  "granteeUserId": "44099659-ceca-4310-b565-88e5257ae6f0",
  "serviceId": "4e8d3cff-ce68-4019-af70-67a1bb961ec8"
}
```

HTTP 422

