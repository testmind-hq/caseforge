# ── [OWASP-API6] PUT /api/admin/users/{id} — mass assignment ──
# case_id=TC-38dd166b
# case_name=[OWASP-API6] PUT /api/admin/users/{id} — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

PUT {{base_url}}/api/admin/users/{id}
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "id": 99999,
  "isActive": false,
  "role": "team_member",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 200

[Asserts]
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"

