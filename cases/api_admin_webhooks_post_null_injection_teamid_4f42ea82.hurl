# ── POST /api/admin/webhooks - null injection: teamId ──
# case_id=TC-4f42ea82
# case_name=POST /api/admin/webhooks - null injection: teamId
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": [
    "from"
  ],
  "name": "Tanner Gardner",
  "providerType": "patiently",
  "teamId": null,
  "url": "https://www.seniorsynergies.info/one-to-one"
}
```

HTTP 422

