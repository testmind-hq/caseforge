# ── POST /api/admin/webhooks - missing required field "url" ──
# case_id=TC-f322285b
# case_name=POST /api/admin/webhooks - missing required field "url"
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
  "name": "Beulah Douglas",
  "providerType": "his",
  "teamId": "4c031d9f-941f-4af7-bf94-9bb5b7ae85a3"
}
```

HTTP 422

