# ── POST /api/admin/teams/{id}/members - [unicode_fuzzing] userId zero_width ──
# case_id=TC-bdeeed04
# case_name=POST /api/admin/teams/{id}/members - [unicode_fuzzing] userId zero_width
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "owner",
  "userId": "​hello"
}
```

HTTP 400

