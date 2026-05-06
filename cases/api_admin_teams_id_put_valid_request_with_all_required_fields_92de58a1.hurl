# ── PUT /api/admin/teams/{id} - valid request with all required fields ──
# case_id=TC-92de58a1
# case_name=PUT /api/admin/teams/{id} - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

PUT {{base_url}}/api/admin/teams/{id}
Content-Type: application/json
```json
{
  "description": "Optimize company for lovely clarity.",
  "displayName": "snore"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.isDeletable" exists
jsonpath "$.name" exists
jsonpath "$.createdAt" exists
jsonpath "$.description" exists
jsonpath "$.displayName" exists
jsonpath "$.id" exists
jsonpath "$.isDefault" exists

