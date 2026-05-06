# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] description bidi_override ──
# case_id=TC-c42ef106
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] description bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "‮hello",
  "displayName": "quarterly"
}
```

HTTP 400

