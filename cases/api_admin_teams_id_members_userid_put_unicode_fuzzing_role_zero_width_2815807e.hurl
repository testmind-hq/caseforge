# ── PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role zero_width ──
# case_id=TC-2815807e
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "​hello"
}
```

HTTP 400

