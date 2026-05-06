# ── PATCH /api/admin/webhooks/:id - valid request with all required fields ──
# case_id=TC-415f32a9
# case_name=PATCH /api/admin/webhooks/:id - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "none"
  ],
  "isActive": true,
  "name": "Dolly Richards",
  "url": "http://www.futuredeliver.org/dynamic"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.createdAt" exists
jsonpath "$.providerType" exists
jsonpath "$.createdBy" exists
jsonpath "$.url" exists
jsonpath "$.name" exists
jsonpath "$.teamId" exists
jsonpath "$.id" exists
jsonpath "$.events" exists
jsonpath "$.isActive" exists

