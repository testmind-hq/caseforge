# ── PATCH /api/admin/webhooks/:id - null injection: isActive ──
# case_id=TC-f681cd0b
# case_name=PATCH /api/admin/webhooks/:id - null injection: isActive
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
  "isActive": null,
  "name": "Opal Deckow",
  "url": "http://www.dynamicmarkets.net/vertical"
}
```

HTTP 422

