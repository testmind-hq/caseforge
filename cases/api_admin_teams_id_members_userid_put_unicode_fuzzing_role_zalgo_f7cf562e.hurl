# ── PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role zalgo ──
# case_id=TC-f7cf562e
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

