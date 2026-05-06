# ── POST /api/admin/teams/{id}/grants - [unicode_fuzzing] granteeUserId bidi_override ──
# case_id=TC-57831769
# case_name=POST /api/admin/teams/{id}/grants - [unicode_fuzzing] granteeUserId bidi_override
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
  "granteeUserId": "‮hello",
  "serviceId": "20b4459a-3c62-4892-9a64-1ee451cbd675"
}
```

HTTP 400

