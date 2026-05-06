# ── PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId control_char ──
# case_id=TC-00caba6f
# case_name=PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": "hello\u0000world"
}
```

HTTP 400

