# ── POST /api/admin/webhooks - [mass_assignment] financial probe ──
# case_id=TC-241955ee
# case_name=POST /api/admin/webhooks - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "events": [
    "actor"
  ],
  "name": "Agustina McKenzie",
  "price": 1,
  "providerType": "eye",
  "teamId": "304932c0-8102-4bb0-bd20-eb20fbf9ab2f",
  "url": "http://www.vicemethodologies.com/virtual/metrics"
}
```

HTTP 400

