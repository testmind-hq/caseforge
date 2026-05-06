# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] description zalgo ──
# case_id=TC-e354e0de
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] description zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "z̀́̂̃̄̅̆̇a",
  "displayName": "quarterly"
}
```

HTTP 400

