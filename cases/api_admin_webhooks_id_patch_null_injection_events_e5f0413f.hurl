# ── PATCH /api/admin/webhooks/:id - null injection: events ──
# case_id=TC-e5f0413f
# case_name=PATCH /api/admin/webhooks/:id - null injection: events
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": null,
  "isActive": true,
  "name": "Opal Deckow",
  "url": "http://www.dynamicmarkets.net/vertical"
}
```

HTTP 422

