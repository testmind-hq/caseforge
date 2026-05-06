# ── PATCH /api/admin/webhooks/:id - [type_coercion] url wrong_type_boolean ──
# case_id=TC-d9bfd2d8
# case_name=PATCH /api/admin/webhooks/:id - [type_coercion] url wrong_type_boolean
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
  "name": "Emile Jones",
  "url": true
}
```

HTTP 422

