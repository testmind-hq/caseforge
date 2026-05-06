# ── POST /api/admin/teams/{id}/members - [unicode_fuzzing] role control_char ──
# case_id=TC-39e9a695
# case_name=POST /api/admin/teams/{id}/members - [unicode_fuzzing] role control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "hello\u0000world",
  "userId": "00287abb-135c-4e57-a40f-6a5a00caf19e"
}
```

HTTP 400

