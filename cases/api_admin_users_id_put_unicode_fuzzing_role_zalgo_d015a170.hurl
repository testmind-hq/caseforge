# ── PUT /api/admin/users/{id} - [unicode_fuzzing] role zalgo ──
# case_id=TC-d015a170
# case_name=PUT /api/admin/users/{id} - [unicode_fuzzing] role zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

