# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName control_char ──
# case_id=TC-39195267
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Warm starts beat cold work.",
  "displayName": "hello\u0000world"
}
```

HTTP 400

