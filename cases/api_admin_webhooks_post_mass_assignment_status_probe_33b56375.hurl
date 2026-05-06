# ── POST /api/admin/webhooks - [mass_assignment] status probe ──
# case_id=TC-33b56375
# case_name=POST /api/admin/webhooks - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "events": [
    "actor"
  ],
  "name": "Agustina McKenzie",
  "providerType": "eye",
  "teamId": "304932c0-8102-4bb0-bd20-eb20fbf9ab2f",
  "url": "http://www.vicemethodologies.com/virtual/metrics",
  "verified": true
}
```

HTTP 400

