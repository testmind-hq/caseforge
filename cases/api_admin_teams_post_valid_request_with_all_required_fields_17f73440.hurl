# ── POST /api/admin/teams - valid request with all required fields ──
# case_id=TC-17f73440
# case_name=POST /api/admin/teams - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "description": "Explicitly name the year before you enlist it.",
  "displayName": "downstairs",
  "name": "Amie Paul"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.isDefault" exists
jsonpath "$.isDeletable" exists
jsonpath "$.name" exists
jsonpath "$.createdAt" exists
jsonpath "$.description" exists
jsonpath "$.displayName" exists
jsonpath "$.id" exists

