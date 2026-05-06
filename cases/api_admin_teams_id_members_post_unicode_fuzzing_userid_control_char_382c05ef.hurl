# ── POST /api/admin/teams/{id}/members - [unicode_fuzzing] userId control_char ──
# case_id=TC-382c05ef
# case_name=POST /api/admin/teams/{id}/members - [unicode_fuzzing] userId control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "owner",
  "userId": "hello\u0000world"
}
```

HTTP 400

