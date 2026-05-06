# ── POST /api/admin/teams/{id}/grants - [mass_assignment] identity probe ──
# case_id=TC-74060ffe
# case_name=POST /api/admin/teams/{id}/grants - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "these"
  ],
  "createdBy": "__probe__",
  "expiresAt": "1935-06-17T15:07:26Z",
  "granteeTeamId": "02c4dc55-7e2a-4090-a2d0-b4fed5e1277e",
  "granteeUserId": "85fb4919-bc0a-470e-9fae-9fa164ef5b88",
  "ownerId": "__probe__",
  "serviceId": "b5371d8e-203f-403f-bbb6-ab0e4e8f8466",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

