# ── POST /api/admin/webhooks - [field_boundary] name invalid_below_min ──
# case_id=TC-7b9e5b4d
# case_name=POST /api/admin/webhooks - [field_boundary] name invalid_below_min
# step_id=step-main
# step_type=test
# technique=field_boundary
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "regularly"
  ],
  "name": "",
  "providerType": "pen",
  "teamId": "8e786d80-b9b5-471b-8643-4dea8db9db45",
  "url": "http://www.seniorb2b.io/webservices/repurpose/mindshare"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

