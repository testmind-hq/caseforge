# ── POST /api/admin/teams/{id}/members - null injection: userId ──
# case_id=TC-1b45482b
# case_name=POST /api/admin/teams/{id}/members - null injection: userId
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "owner",
  "userId": null
}
```

HTTP 422

