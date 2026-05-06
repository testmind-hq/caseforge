# ── PUT /api/admin/teams/{id}/members/{userId} - null injection: role ──
# case_id=TC-92d17333
# case_name=PUT /api/admin/teams/{id}/members/{userId} - null injection: role
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": null
}
```

HTTP 422

