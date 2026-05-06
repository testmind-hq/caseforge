# ── POST /api/admin/teams - [unicode_fuzzing] description control_char ──
# case_id=TC-8656dd0b
# case_name=POST /api/admin/teams - [unicode_fuzzing] description control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "hello\u0000world",
  "displayName": "example",
  "name": "Thomas Castillo"
}
```

HTTP 400

