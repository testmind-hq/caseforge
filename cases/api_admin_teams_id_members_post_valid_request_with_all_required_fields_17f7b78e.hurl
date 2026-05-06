# ── POST /api/admin/teams/{id}/members - valid request with all required fields ──
# case_id=TC-17f7b78e
# case_name=POST /api/admin/teams/{id}/members - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "role": "member",
  "userId": "a3bd36d6-0660-42cd-82e2-4ffe231776bc"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.ok" exists

