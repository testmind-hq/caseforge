# ── POST /api/admin/teams/{id}/members - missing required field "userId" ──
# case_id=TC-4eda623b
# case_name=POST /api/admin/teams/{id}/members - missing required field "userId"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member"
}
```

HTTP 422

