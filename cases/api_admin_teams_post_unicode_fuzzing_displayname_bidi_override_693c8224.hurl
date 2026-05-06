# ── POST /api/admin/teams - [unicode_fuzzing] displayName bidi_override ──
# case_id=TC-693c8224
# case_name=POST /api/admin/teams - [unicode_fuzzing] displayName bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Review the woman every 2 weeks.",
  "displayName": "‮hello",
  "name": "Thomas Castillo"
}
```

HTTP 400

