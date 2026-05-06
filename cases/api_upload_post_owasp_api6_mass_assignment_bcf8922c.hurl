# ── [OWASP-API6] POST /api/upload — mass assignment ──
# case_id=TC-bcf8922c
# case_name=[OWASP-API6] POST /api/upload — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/upload
Content-Type: application/json
```json
{
  "branch": "second",
  "commitSha": "he",
  "createdAt": "2000-01-01T00:00:00Z",
  "id": 99999,
  "service": "his",
  "specContent": "of",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 201

[Asserts]
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"

