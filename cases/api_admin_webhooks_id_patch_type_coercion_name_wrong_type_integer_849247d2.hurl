# ── PATCH /api/admin/webhooks/:id - [type_coercion] name wrong_type_integer ──
# case_id=TC-849247d2
# case_name=PATCH /api/admin/webhooks/:id - [type_coercion] name wrong_type_integer
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
  "name": 123,
  "url": "https://www.financeoptimize.com/transform/cross-media/technologies"
}
```

HTTP 422

