# ── POST /api/upload - valid request with all required fields ──
# case_id=TC-e3da0de9
# case_name=POST /api/upload - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "my",
  "commitSha": "where",
  "service": "Asian",
  "specContent": "soon"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.endpointCount" exists
jsonpath "$.service" exists
jsonpath "$.unchanged" exists
jsonpath "$.warnings" exists
jsonpath "$.wasConverted" exists
jsonpath "$.branch" exists

