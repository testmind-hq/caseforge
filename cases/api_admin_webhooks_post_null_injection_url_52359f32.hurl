# ── POST /api/admin/webhooks - null injection: url ──
# case_id=TC-52359f32
# case_name=POST /api/admin/webhooks - null injection: url
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
  "teamId": "19ccbd87-5161-4a81-beda-3e6a1d5aa25e",
  "url": null
}
```

HTTP 422

