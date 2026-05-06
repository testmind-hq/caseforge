# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] description control_char ──
# case_id=TC-d9200d81
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] description control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "hello\u0000world",
  "displayName": "quarterly"
}
```

HTTP 400

