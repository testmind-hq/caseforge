# ── [OWASP-API6] POST /api/admin/teams/{id}/grants — mass assignment ──
# case_id=TC-e74b3c2c
# case_name=[OWASP-API6] POST /api/admin/teams/{id}/grants — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/teams/{id}/grants
Content-Type: application/json
```json
{
  "branches": [
    "might"
  ],
  "createdAt": "2000-01-01T00:00:00Z",
  "expiresAt": "1904-11-16T00:21:56Z",
  "granteeTeamId": "80cfeb39-de1f-4afc-b29b-dbf268b668eb",
  "granteeUserId": "af0ce4e0-f8fb-4c7c-b929-9d7dfc463d99",
  "id": 99999,
  "serviceId": "3751ed85-6162-4db7-8287-4b7491018fb0",
  "updatedAt": "2000-01-01T00:00:00Z"
}
```

HTTP 201

[Asserts]
jsonpath "$.id" != 99999
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"

