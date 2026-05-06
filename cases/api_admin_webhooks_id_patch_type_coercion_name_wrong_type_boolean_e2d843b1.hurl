# ── PATCH /api/admin/webhooks/:id - [type_coercion] name wrong_type_boolean ──
# case_id=TC-e2d843b1
# case_name=PATCH /api/admin/webhooks/:id - [type_coercion] name wrong_type_boolean
# step_id=step-main
# step_type=test
# technique=type_coercion
# priority=P2

PATCH {{base_url}}/api/admin/webhooks/:id
Content-Type: application/json
```json
{
  "events": [
    "some"
  ],
  "isActive": false,
  "name": true,
  "url": "https://www.financeoptimize.com/transform/cross-media/technologies"
}
```

HTTP 422

