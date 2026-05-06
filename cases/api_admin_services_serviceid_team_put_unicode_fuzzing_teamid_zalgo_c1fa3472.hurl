# ── PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId zalgo ──
# case_id=TC-c1fa3472
# case_name=PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

