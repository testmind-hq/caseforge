# ── POST /api/admin/teams/{id}/members - [unicode_fuzzing] userId zalgo ──
# case_id=TC-9cd03a11
# case_name=POST /api/admin/teams/{id}/members - [unicode_fuzzing] userId zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "owner",
  "userId": "z̀́̂̃̄̅̆̇a"
}
```

HTTP 400

