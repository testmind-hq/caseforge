# ── PATCH /api/admin/webhooks/:id - [mass_assignment] identity probe ──
# case_id=TC-1274d148
# case_name=PATCH /api/admin/webhooks/:id - [mass_assignment] identity probe
# step_id=step-main
# step_type=test
# technique=mass_assignment
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "createdBy": "__probe__",
  "events": [
    "of"
  ],
  "isActive": false,
  "name": "Nathaniel Yang",
  "ownerId": "__probe__",
  "url": "https://www.forwardinteractive.com/architect/reintermediate/user-centric",
  "userId": "__probe__",
  "user_id": "__probe__"
}
```

HTTP 400

