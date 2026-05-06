# ── POST /api/admin/webhooks - [required_omission] events absent ──
# case_id=TC-09946d4c
# case_name=POST /api/admin/webhooks - [required_omission] events absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "name": "Molly Hudson",
  "providerType": "next",
  "teamId": "6c927896-300a-4cc9-a530-93b2a15d5633",
  "url": "http://www.humanusers.name/engage"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

