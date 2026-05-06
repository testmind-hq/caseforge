# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName bidi_override ──
# case_id=TC-7c97c5e9
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Warm starts beat cold work.",
  "displayName": "‮hello"
}
```

HTTP 400

