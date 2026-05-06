# ── POST /api/admin/webhooks - missing required field "name" ──
# case_id=TC-6c83435b
# case_name=POST /api/admin/webhooks - missing required field "name"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "as"
  ],
  "providerType": "his",
  "teamId": "4c031d9f-941f-4af7-bf94-9bb5b7ae85a3",
  "url": "https://www.investormethodologies.net/maximize"
}
```

HTTP 422

