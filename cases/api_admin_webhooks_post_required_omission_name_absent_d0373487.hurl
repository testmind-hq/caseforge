# ── POST /api/admin/webhooks - [required_omission] name absent ──
# case_id=TC-d0373487
# case_name=POST /api/admin/webhooks - [required_omission] name absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "it"
  ],
  "providerType": "few",
  "teamId": "949cf797-62f1-45ef-9b37-71379d7223ec",
  "url": "http://www.regionalproactive.io/scalable"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

