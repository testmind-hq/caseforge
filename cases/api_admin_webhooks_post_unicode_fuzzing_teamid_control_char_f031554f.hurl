# ── POST /api/admin/webhooks - [unicode_fuzzing] teamId control_char ──
# case_id=TC-f031554f
# case_name=POST /api/admin/webhooks - [unicode_fuzzing] teamId control_char
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
  "teamId": "hello\u0000world",
  "url": "https://www.futureb2c.biz/viral/vortals/clicks-and-mortar"
}
```

HTTP 400

