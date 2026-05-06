# ── POST /api/admin/teams/{id}/grants - null injection: serviceId ──
# case_id=TC-2571eb1b
# case_name=POST /api/admin/teams/{id}/grants - null injection: serviceId
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
  "granteeTeamId": "bcaeb7d9-6d53-4be0-8f2e-d1beacfc2fa1",
  "granteeUserId": "44099659-ceca-4310-b565-88e5257ae6f0",
  "serviceId": null
}
```

HTTP 422

