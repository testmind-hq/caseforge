# ── POST /api/admin/teams/{id}/grants - [unicode_fuzzing] granteeTeamId zalgo ──
# case_id=TC-603eeaa8
# case_name=POST /api/admin/teams/{id}/grants - [unicode_fuzzing] granteeTeamId zalgo
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
  "granteeTeamId": "z̀́̂̃̄̅̆̇a",
  "granteeUserId": "2c643fc3-d35d-409a-9c59-573ae43f37e6",
  "serviceId": "20b4459a-3c62-4892-9a64-1ee451cbd675"
}
```

HTTP 400

