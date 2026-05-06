# ── PATCH /api/admin/webhooks/:id - wrong content-type (text/plain) ──
# case_id=TC-94225ad6
# case_name=PATCH /api/admin/webhooks/:id - wrong content-type (text/plain)
# step_id=step-main
# step_type=test
# technique=constraint_mutation
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: text/plain
```json
{
  "events": [
    "aloof"
  ],
  "isActive": true,
  "name": "Opal Deckow",
  "url": "http://www.dynamicmarkets.net/vertical"
}
```

HTTP 415

