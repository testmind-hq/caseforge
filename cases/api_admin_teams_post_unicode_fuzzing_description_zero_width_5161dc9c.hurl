# ── POST /api/admin/teams - [unicode_fuzzing] description zero_width ──
# case_id=TC-5161dc9c
# case_name=POST /api/admin/teams - [unicode_fuzzing] description zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "​hello",
  "displayName": "example",
  "name": "Thomas Castillo"
}
```

HTTP 400

