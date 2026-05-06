# ── PATCH /api/admin/webhooks/:id - [mass_assignment] financial probe ──
# case_id=TC-ed85e04f
# case_name=PATCH /api/admin/webhooks/:id - [mass_assignment] financial probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "balance": 1,
  "credits": 1,
  "discount": 0,
  "events": [
    "of"
  ],
  "isActive": false,
  "name": "Nathaniel Yang",
  "price": 1,
  "url": "https://www.forwardinteractive.com/architect/reintermediate/user-centric"
}
```

HTTP 400

