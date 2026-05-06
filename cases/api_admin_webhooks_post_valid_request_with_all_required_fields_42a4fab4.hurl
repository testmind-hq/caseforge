# ── POST /api/admin/webhooks - valid request with all required fields ──
# case_id=TC-42a4fab4
# case_name=POST /api/admin/webhooks - valid request with all required fields
# step_id=step-main
# step_type=test
# technique=equivalence_partitioning
# priority=P0

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "set"
  ],
  "name": "Fletcher Mendez",
  "providerType": "these",
  "teamId": "7b7e7d08-a4c7-4b59-a185-b2a7b8576f2e",
  "url": "http://www.nationalcross-platform.org/infomediaries/killer/technologies/frictionless"
}
```

HTTP 200

[Asserts]
duration < 2000
jsonpath "$.createdBy" exists
jsonpath "$.isActive" exists
jsonpath "$.providerType" exists
jsonpath "$.teamId" exists
jsonpath "$.name" exists
jsonpath "$.url" exists
jsonpath "$.createdAt" exists
jsonpath "$.id" exists
jsonpath "$.events" exists

