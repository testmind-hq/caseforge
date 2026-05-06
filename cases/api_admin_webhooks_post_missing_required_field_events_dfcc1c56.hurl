# ── POST /api/admin/webhooks - missing required field "events" ──
# case_id=TC-dfcc1c56
# case_name=POST /api/admin/webhooks - missing required field "events"
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P1

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "name": "Beulah Douglas",
  "providerType": "his",
  "teamId": "4c031d9f-941f-4af7-bf94-9bb5b7ae85a3",
  "url": "https://www.investormethodologies.net/maximize"
}
```

HTTP 422

