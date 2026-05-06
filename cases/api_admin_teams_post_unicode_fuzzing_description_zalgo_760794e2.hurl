# ── POST /api/admin/teams - [unicode_fuzzing] description zalgo ──
# case_id=TC-760794e2
# case_name=POST /api/admin/teams - [unicode_fuzzing] description zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "z̀́̂̃̄̅̆̇a",
  "displayName": "example",
  "name": "Thomas Castillo"
}
```

HTTP 400

