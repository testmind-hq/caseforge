# ── [OWASP-API6] POST /auth/login — mass assignment ──
# case_id=TC-09c747ae
# case_name=[OWASP-API6] POST /auth/login — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/auth/login
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "email": "eddhanson@thomas.net",
  "id": 99999,
  "password": "we",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 201

[Asserts]
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"

