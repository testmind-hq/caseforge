# ── POST /api/admin/webhooks - null injection: name ──
# case_id=TC-169dbf8c
# case_name=POST /api/admin/webhooks - null injection: name
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
  "name": null,
  "providerType": "patiently",
  "teamId": "19ccbd87-5161-4a81-beda-3e6a1d5aa25e",
  "url": "https://www.seniorsynergies.info/one-to-one"
}
```

HTTP 422

