# ── POST /api/admin/webhooks - [mass_assignment] identity probe ──
# case_id=TC-30b18c5f
# case_name=POST /api/admin/webhooks - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "events": [
    "actor"
  ],
  "name": "Agustina McKenzie",
  "ownerId": "__probe__",
  "providerType": "eye",
  "teamId": "304932c0-8102-4bb0-bd20-eb20fbf9ab2f",
  "url": "http://www.vicemethodologies.com/virtual/metrics",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

