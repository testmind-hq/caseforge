# ── [OWASP-API6] POST /api/admin/webhooks — mass assignment ──
# case_id=TC-1b59ba48
# case_name=[OWASP-API6] POST /api/admin/webhooks — mass assignment
# step_id=step-1
# step_type=test
# technique=owasp_api_top10
# priority=P0

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "createdAt": "2000-01-01T00:00:00Z",
  "events": [
    "Lebanese"
  ],
  "id": 99999,
  "name": "Rowan Bartell",
  "providerType": "Polish",
  "teamId": "5bfa6b50-a743-4866-b2b2-f649decc8c37",
  "updatedAt": "2000-01-01T00:00:00Z",
  "url": "https://www.regionalfacilitate.com/users/intuitive"
}
```

HTTP 201

[Asserts]
jsonpath "$.createdAt" != "2000-01-01T00:00:00Z"
jsonpath "$.updatedAt" != "2000-01-01T00:00:00Z"
jsonpath "$.id" != 99999

