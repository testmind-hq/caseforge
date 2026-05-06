# ── PUT /api/admin/users/{id} - [unicode_fuzzing] role bidi_override ──
# case_id=TC-a2217373
# case_name=PUT /api/admin/users/{id} - [unicode_fuzzing] role bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": false,
  "role": "‮hello"
}
```

HTTP 400

