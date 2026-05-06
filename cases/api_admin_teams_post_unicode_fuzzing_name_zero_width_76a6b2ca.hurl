# ── POST /api/admin/teams - [unicode_fuzzing] name zero_width ──
# case_id=TC-76a6b2ca
# case_name=POST /api/admin/teams - [unicode_fuzzing] name zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Review the woman every 2 weeks.",
  "displayName": "example",
  "name": "​hello"
}
```

HTTP 400

