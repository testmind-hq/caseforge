# ── PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId zero_width ──
# case_id=TC-1c0a1d4a
# case_name=PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": "​hello"
}
```

HTTP 400

