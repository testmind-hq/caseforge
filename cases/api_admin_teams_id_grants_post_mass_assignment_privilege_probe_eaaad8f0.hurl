# ── POST /api/admin/teams/{id}/grants - [mass_assignment] privilege probe ──
# case_id=TC-eaaad8f0
# case_name=POST /api/admin/teams/{id}/grants - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "admin": true,
  "branches": [
    "these"
  ],
  "expiresAt": "1935-06-17T15:07:26Z",
  "granteeTeamId": "02c4dc55-7e2a-4090-a2d0-b4fed5e1277e",
  "granteeUserId": "85fb4919-bc0a-470e-9fae-9fa164ef5b88",
  "isAdmin": true,
  "is_admin": true,
  "role": "__probe__",
  "serviceId": "b5371d8e-203f-403f-bbb6-ab0e4e8f8466"
}
```

HTTP 400

