# ── [OWASP-API6] PUT /api/admin/services/{serviceId}/team — mass assignment ──
# case_id=TC-29a92605
# case_name=[OWASP-API6] PUT /api/admin/services/{serviceId}/team — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/services/{serviceId}/team
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "id": 99999,
  "teamId": "d9bf3e10-6529-49aa-b714-03fd1a939f04",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 200

[Asserts]
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"

