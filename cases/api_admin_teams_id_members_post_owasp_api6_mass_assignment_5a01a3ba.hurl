# ── [OWASP-API6] POST /api/admin/teams/{id}/members — mass assignment ──
# case_id=TC-5a01a3ba
# case_name=[OWASP-API6] POST /api/admin/teams/{id}/members — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/{id}/members
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "id": 99999,
  "role": "owner",
  "updatedAt": "2000-01-01T00:00:00Z",
  "userId": "4409317f-6972-4069-8ed6-942e90d42ec2"
}
```

HTTP 201

[Asserts]
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"

