# ── POST /api/admin/teams/{id}/members - null injection: role ──
# case_id=TC-a2c2e196
# case_name=POST /api/admin/teams/{id}/members - null injection: role
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": null,
  "userId": "b6f51cc4-2389-42c5-a864-35545c08cda9"
}
```

HTTP 422

