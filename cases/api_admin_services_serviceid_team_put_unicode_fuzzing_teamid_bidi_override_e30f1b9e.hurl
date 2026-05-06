# ── PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId bidi_override ──
# case_id=TC-e30f1b9e
# case_name=PUT /api/admin/services/{serviceId}/team - [unicode_fuzzing] teamId bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "teamId": "‮hello"
}
```

HTTP 400

