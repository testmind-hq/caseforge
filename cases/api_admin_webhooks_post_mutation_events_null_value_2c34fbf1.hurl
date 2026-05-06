# ── POST /api/admin/webhooks - mutation: events null value ──
# case_id=TC-2c34fbf1
# case_name=POST /api/admin/webhooks - mutation: events null value
# step_id=step-main
# step_type=test
# technique=mutation
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": null,
  "name": "Javier Bogan",
  "providerType": "regiment",
  "teamId": "e174fc1d-b8a7-4b7c-936d-1a4ea32d8bd1",
  "url": "http://www.groupembrace.net/engage/best-of-breed/scale"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

