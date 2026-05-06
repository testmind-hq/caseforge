# ── POST /api/admin/webhooks - [unicode_fuzzing] url control_char ──
# case_id=TC-c4479bd1
# case_name=POST /api/admin/webhooks - [unicode_fuzzing] url control_char
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
  "teamId": "70da9d0e-d12a-4d62-a990-03d8a490c890",
  "url": "hello\u0000world"
}
```

HTTP 400

