# ── POST /api/admin/webhooks - [schema_violation] name_too_short ──
# case_id=TC-b49ea6fa
# case_name=POST /api/admin/webhooks - [schema_violation] name_too_short
# step_id=step-main
# step_type=test
# technique=schema_violation
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "hundred"
  ],
  "name": "",
  "providerType": "me",
  "teamId": "8afc12a7-a242-4e1f-b05b-4ade3fb01c0f",
  "url": "https://www.legacyincubate.io/seize"
}
```

HTTP 422

