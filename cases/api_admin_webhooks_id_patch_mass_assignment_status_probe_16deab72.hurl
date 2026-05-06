# ── PATCH /api/admin/webhooks/:id - [mass_assignment] status probe ──
# case_id=TC-16deab72
# case_name=PATCH /api/admin/webhooks/:id - [mass_assignment] status probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "approved": true,
  "banned": false,
  "disabled": false,
  "events": [
    "of"
  ],
  "isActive": false,
  "name": "Nathaniel Yang",
  "url": "https://www.forwardinteractive.com/architect/reintermediate/user-centric",
  "verified": true
}
```

HTTP 400

