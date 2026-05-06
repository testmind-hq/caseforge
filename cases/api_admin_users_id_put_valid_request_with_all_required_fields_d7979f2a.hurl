# ── PUT /api/admin/users/{id} - valid request with all required fields ──
# case_id=TC-d7979f2a
# case_name=PUT /api/admin/users/{id} - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "isActive": true,
  "role": "team_owner"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

