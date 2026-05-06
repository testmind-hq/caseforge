# ── POST /api/admin/teams - [unicode_fuzzing] displayName zero_width ──
# case_id=TC-8b028ce1
# case_name=POST /api/admin/teams - [unicode_fuzzing] displayName zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Review the woman every 2 weeks.",
  "displayName": "​hello",
  "name": "Thomas Castillo"
}
```

HTTP 400

