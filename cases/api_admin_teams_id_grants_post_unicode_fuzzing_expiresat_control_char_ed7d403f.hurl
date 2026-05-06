# ── POST /api/admin/teams/{id}/grants - [unicode_fuzzing] expiresAt control_char ──
# case_id=TC-ed7d403f
# case_name=POST /api/admin/teams/{id}/grants - [unicode_fuzzing] expiresAt control_char
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
  "expiresAt": "hello\u0000world",
  "granteeTeamId": "2a53aa87-8c01-4d67-959b-fef4369d639b",
  "granteeUserId": "2c643fc3-d35d-409a-9c59-573ae43f37e6",
  "serviceId": "20b4459a-3c62-4892-9a64-1ee451cbd675"
}
```

HTTP 400

