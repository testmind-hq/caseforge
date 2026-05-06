# ── PATCH /api/admin/webhooks/:id - [mass_assignment] privilege probe ──
# case_id=TC-d0ddffec
# case_name=PATCH /api/admin/webhooks/:id - [mass_assignment] privilege probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "admin": true,
  "events": [
    "of"
  ],
  "isActive": false,
  "isAdmin": true,
  "is_admin": true,
  "name": "Nathaniel Yang",
  "role": "__probe__",
  "url": "https://www.forwardinteractive.com/architect/reintermediate/user-centric"
}
```

HTTP 400

