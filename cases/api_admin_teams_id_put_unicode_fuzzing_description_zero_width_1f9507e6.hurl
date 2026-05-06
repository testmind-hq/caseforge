# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] description zero_width ──
# case_id=TC-1f9507e6
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] description zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "​hello",
  "displayName": "quarterly"
}
```

HTTP 400

