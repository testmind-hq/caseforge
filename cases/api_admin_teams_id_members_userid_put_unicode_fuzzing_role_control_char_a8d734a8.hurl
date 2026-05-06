# ── PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role control_char ──
# case_id=TC-a8d734a8
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "hello\u0000world"
}
```

HTTP 400

