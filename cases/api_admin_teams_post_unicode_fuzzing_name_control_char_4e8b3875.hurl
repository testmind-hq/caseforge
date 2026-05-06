# ── POST /api/admin/teams - [unicode_fuzzing] name control_char ──
# case_id=TC-4e8b3875
# case_name=POST /api/admin/teams - [unicode_fuzzing] name control_char
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
  "name": "hello\u0000world"
}
```

HTTP 400

