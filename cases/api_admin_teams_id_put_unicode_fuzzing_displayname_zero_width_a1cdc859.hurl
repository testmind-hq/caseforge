# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName zero_width ──
# case_id=TC-a1cdc859
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Warm starts beat cold work.",
  "displayName": "​hello"
}
```

HTTP 400

