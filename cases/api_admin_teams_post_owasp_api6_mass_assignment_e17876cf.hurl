# ── [OWASP-API6] POST /api/admin/teams — mass assignment ──
# case_id=TC-e17876cf
# case_name=[OWASP-API6] POST /api/admin/teams — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "description": "Prefer predictable government over surprising work.",
  "displayName": "can",
  "id": 99999,
  "name": "Dane Bates",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 201

[Asserts]
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"

