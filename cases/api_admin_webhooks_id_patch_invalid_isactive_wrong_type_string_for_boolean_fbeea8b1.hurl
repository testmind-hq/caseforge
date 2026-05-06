# ── PATCH /api/admin/webhooks/:id - invalid isActive: wrong type (string for boolean) ──
# case_id=TC-fbeea8b1
# case_name=PATCH /api/admin/webhooks/:id - invalid isActive: wrong type (string for boolean)
# step_id=step-main
# step_type=test
# technique=isolated_negative
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "regularly"
  ],
  "isActive": "not_a_boolean",
  "name": "Halle Lewis",
  "url": "http://www.technicalschemas.com/web-enabled"
}
```

HTTP 422

