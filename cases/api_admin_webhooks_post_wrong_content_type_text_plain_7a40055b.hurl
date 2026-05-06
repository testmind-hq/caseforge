# ── POST /api/admin/webhooks - wrong content-type (text/plain) ──
# case_id=TC-7a40055b
# case_name=POST /api/admin/webhooks - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: text/plain
```json
{
  "events": [
    "from"
  ],
  "name": "Tanner Gardner",
  "providerType": "patiently",
  "teamId": "19ccbd87-5161-4a81-beda-3e6a1d5aa25e",
  "url": "https://www.seniorsynergies.info/one-to-one"
}
```

HTTP 415

