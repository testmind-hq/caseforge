# ── PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName zalgo ──
# case_id=TC-5add01e6
# case_name=PUT /api/admin/teams/{id} - [unicode_fuzzing] displayName zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Warm starts beat cold work.",
  "displayName": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

