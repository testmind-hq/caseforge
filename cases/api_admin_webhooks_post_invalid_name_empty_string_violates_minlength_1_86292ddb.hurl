# ── POST /api/admin/webhooks - invalid name: empty string violates minLength 1 ──
# case_id=TC-86292ddb
# case_name=POST /api/admin/webhooks - invalid name: empty string violates minLength 1
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "as"
  ],
  "name": "",
  "providerType": "his",
  "teamId": "4c031d9f-941f-4af7-bf94-9bb5b7ae85a3",
  "url": "https://www.investormethodologies.net/maximize"
}
```

HTTP 422

