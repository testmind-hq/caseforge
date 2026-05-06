# ── POST /api/admin/teams/{id}/members - [schema_violation] userId_missing_required ──
# case_id=TC-71efcd62
# case_name=POST /api/admin/teams/{id}/members - [schema_violation] userId_missing_required
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member"
}
```

HTTP 422

