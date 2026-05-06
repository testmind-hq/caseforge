# ── POST /api/admin/webhooks - [required_omission] url absent ──
# case_id=TC-6d3bc221
# case_name=POST /api/admin/webhooks - [required_omission] url absent
# step_id=step-main
# step_type=test
# technique=required_omission
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "last"
  ],
  "name": "Alvina Powell",
  "providerType": "itself",
  "teamId": "3652daaf-fcaf-461d-97f6-ccc7da39f569"
}
```

HTTP *

[Asserts]
status >= 400
status < 500

