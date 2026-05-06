# ── PUT /api/admin/teams/{id}/members/{userId} - missing required field "role" ──
# case_id=TC-7f67bdd2
# case_name=PUT /api/admin/teams/{id}/members/{userId} - missing required field "role"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{}
```

HTTP 422

