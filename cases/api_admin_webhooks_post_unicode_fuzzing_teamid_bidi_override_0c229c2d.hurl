# ── POST /api/admin/webhooks - [unicode_fuzzing] teamId bidi_override ──
# case_id=TC-0c229c2d
# case_name=POST /api/admin/webhooks - [unicode_fuzzing] teamId bidi_override
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
  "teamId": "‮hello",
  "url": "https://www.futureb2c.biz/viral/vortals/clicks-and-mortar"
}
```

HTTP 400

