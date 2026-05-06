# ── POST /api/admin/webhooks - [unicode_fuzzing] teamId zalgo ──
# case_id=TC-bba333a6
# case_name=POST /api/admin/webhooks - [unicode_fuzzing] teamId zalgo
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
  "name": "Anika Lane",
  "providerType": "anyway",
  "teamId": "z̀́̂̃̄̅̆̇a",
  "url": "https://www.futureb2c.biz/viral/vortals/clicks-and-mortar"
}
```

HTTP 400

