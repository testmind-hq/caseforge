# ── [OWASP-API6] POST /api/tokens — mass assignment ──
# case_id=TC-d9979992
# case_name=[OWASP-API6] POST /api/tokens — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/tokens
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "id": 99999,
  "name": "Marianne Nolan",
  "scope": "write",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 201

[Asserts]
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"

