# ── POST /api/admin/teams - [unicode_fuzzing] description bidi_override ──
# case_id=TC-d96ca637
# case_name=POST /api/admin/teams - [unicode_fuzzing] description bidi_override
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "‮hello",
  "displayName": "example",
  "name": "Thomas Castillo"
}
```

HTTP 400

