# ── POST /api/admin/webhooks - [field_boundary] name valid_min ──
# case_id=TC-85b28596
# case_name=POST /api/admin/webhooks - [field_boundary] name valid_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P1

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "that"
  ],
  "name": "a",
  "providerType": "year",
  "teamId": "2078e75e-ac88-4a37-93b9-0aad2a57623c",
  "url": "http://www.principalinteractive.net/turn-key/redefine"
}
```

HTTP *

[Asserts]
status >= 200
status < 300

