# ── PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role bidi_override ──
# case_id=TC-0b0faf09
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [unicode_fuzzing] role bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "‮hello"
}
```

HTTP 400

