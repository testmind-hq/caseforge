# ── PUT /api/admin/users/{id} - [unicode_fuzzing] role control_char ──
# case_id=TC-be44c91e
# case_name=PUT /api/admin/users/{id} - [unicode_fuzzing] role control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "hello\u0000world"
}
```

HTTP 400

