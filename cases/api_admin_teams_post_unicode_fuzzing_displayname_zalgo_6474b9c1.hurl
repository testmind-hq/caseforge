# ── POST /api/admin/teams - [unicode_fuzzing] displayName zalgo ──
# case_id=TC-6474b9c1
# case_name=POST /api/admin/teams - [unicode_fuzzing] displayName zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Review the woman every 2 weeks.",
  "displayName": "z̀́̂̃̄̅̆̇a",
  "name": "Thomas Castillo"
}
```

HTTP 400

