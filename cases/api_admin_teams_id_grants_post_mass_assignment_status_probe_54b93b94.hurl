# ── POST /api/admin/teams/{id}/grants - [mass_assignment] status probe ──
# case_id=TC-54b93b94
# case_name=POST /api/admin/teams/{id}/grants - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "branches": [
    "these"
  ],
  "disabled": false,
  "expiresAt": "1935-06-17T15:07:26Z",
  "granteeTeamId": "02c4dc55-7e2a-4090-a2d0-b4fed5e1277e",
  "granteeUserId": "85fb4919-bc0a-470e-9fae-9fa164ef5b88",
  "serviceId": "b5371d8e-203f-403f-bbb6-ab0e4e8f8466",
  "verified": true
}
```

HTTP 400

