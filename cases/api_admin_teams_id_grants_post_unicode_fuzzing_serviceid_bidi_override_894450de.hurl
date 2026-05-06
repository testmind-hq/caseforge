# ── POST /api/admin/teams/{id}/grants - [unicode_fuzzing] serviceId bidi_override ──
# case_id=TC-894450de
# case_name=POST /api/admin/teams/{id}/grants - [unicode_fuzzing] serviceId bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "eye"
  ],
  "expiresAt": "1910-02-22T19:02:33Z",
  "granteeTeamId": "2a53aa87-8c01-4d67-959b-fef4369d639b",
  "granteeUserId": "2c643fc3-d35d-409a-9c59-573ae43f37e6",
  "serviceId": "‮hello"
}
```

HTTP 400

