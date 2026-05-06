# ── PATCH /api/admin/webhooks/:id - [type_coercion] isActive wrong_type_integer ──
# case_id=TC-4c590e85
# case_name=PATCH /api/admin/webhooks/:id - [type_coercion] isActive wrong_type_integer
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
  "isActive": 1,
  "name": "Emile Jones",
  "url": "https://www.financeoptimize.com/transform/cross-media/technologies"
}
```

HTTP 422

