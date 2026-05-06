# ── PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] status probe ──
# case_id=TC-08a1d397
# case_name=PUT /api/admin/teams/{id}/members/{userId} - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "role": "member",
  "verified": true
}
```

HTTP 400

