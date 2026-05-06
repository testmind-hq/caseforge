# ── POST /api/admin/webhooks - [mass_assignment] privilege probe ──
# case_id=TC-f5c743f7
# case_name=POST /api/admin/webhooks - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "admin": true,
  "events": [
    "actor"
  ],
  "isAdmin": true,
  "is_admin": true,
  "name": "Agustina McKenzie",
  "providerType": "eye",
  "role": "__probe__",
  "teamId": "304932c0-8102-4bb0-bd20-eb20fbf9ab2f",
  "url": "http://www.vicemethodologies.com/virtual/metrics"
}
```

HTTP 400

