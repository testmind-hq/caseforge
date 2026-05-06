# ── POST /api/admin/webhooks - null injection: events ──
# case_id=TC-35254559
# case_name=POST /api/admin/webhooks - null injection: events
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

POST {{base_url}}/api/admin/webhooks
Content-Type: application/json
```json
{
  "events": null,
  "name": "Tanner Gardner",
  "providerType": "patiently",
  "teamId": "19ccbd87-5161-4a81-beda-3e6a1d5aa25e",
  "url": "https://www.seniorsynergies.info/one-to-one"
}
```

HTTP 422

