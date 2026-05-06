# ── POST /api/admin/teams - [unicode_fuzzing] displayName control_char ──
# case_id=TC-7ead4ab7
# case_name=POST /api/admin/teams - [unicode_fuzzing] displayName control_char
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Review the woman every 2 weeks.",
  "displayName": "hello\u0000world",
  "name": "Thomas Castillo"
}
```

HTTP 400

