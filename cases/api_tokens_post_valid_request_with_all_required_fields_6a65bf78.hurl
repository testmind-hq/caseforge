# ── POST /api/tokens - valid request with all required fields ──
# case_id=TC-6a65bf78
# case_name=POST /api/tokens - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "name": "Allison Hunter",
  "scope": "read"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.name" exists
jsonpath "$.prefix" exists
jsonpath "$.scope" exists
jsonpath "$.token" exists
jsonpath "$.createdAt" exists
jsonpath "$.id" exists

