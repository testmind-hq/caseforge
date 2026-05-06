# ── POST /api/admin/webhooks - [unicode_fuzzing] name zalgo ──
# case_id=TC-a7f8f480
# case_name=POST /api/admin/webhooks - [unicode_fuzzing] name zalgo
# step_id=step-main
# step_type=test
# technique=unicode_fuzzing
# priority=P3

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "you"
  ],
  "name": "z̀́̂̃̄̅̆̇a",
  "providerType": "anyway",
  "teamId": "70da9d0e-d12a-4d62-a990-03d8a490c890",
  "url": "https://www.futureb2c.biz/viral/vortals/clicks-and-mortar"
}
```

HTTP 400

