# ── PATCH /api/admin/webhooks/:id - null injection: url ──
# case_id=TC-6597f138
# case_name=PATCH /api/admin/webhooks/:id - null injection: url
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "aloof"
  ],
  "isActive": true,
  "name": "Opal Deckow",
  "url": null
}
```

HTTP 422

