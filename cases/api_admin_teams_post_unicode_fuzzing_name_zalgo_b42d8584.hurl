# ── POST /api/admin/teams - [unicode_fuzzing] name zalgo ──
# case_id=TC-b42d8584
# case_name=POST /api/admin/teams - [unicode_fuzzing] name zalgo
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
  "name": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

