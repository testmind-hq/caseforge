# ── POST /api/admin/teams/{id}/grants - null injection: granteeUserId ──
# case_id=TC-593b0773
# case_name=POST /api/admin/teams/{id}/grants - null injection: granteeUserId
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
  "granteeUserId": null,
  "serviceId": "4e8d3cff-ce68-4019-af70-67a1bb961ec8"
}
```

HTTP 422

