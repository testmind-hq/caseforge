# ── PUT /api/admin/teams/{id}/members/{userId} - missing required field "role" ──
# case_id=TC-02cdac38
# case_name=PUT /api/admin/teams/{id}/members/{userId} - missing required field "role"
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P1

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{}
```

HTTP 422

