# ── PUT /api/admin/teams/{id}/members/{userId} - valid request with all required fields ──
# case_id=TC-b950209e
# case_name=PUT /api/admin/teams/{id}/members/{userId} - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

PUT {{base_url}}/api/admin/teams/{id}/members/{userId}
Content-Type: application/json
```json
{
  "role": "member"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

