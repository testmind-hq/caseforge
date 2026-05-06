# ── [OWASP-API6] POST /auth/register — mass assignment ──
# case_id=TC-900b6a9f
# case_name=[OWASP-API6] POST /auth/register — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/register
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "email": "gennarogislason@newton.io",
  "id": 99999,
  "password": "did",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 201

[Asserts]
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"
jsonpath "$.id" != 99999

